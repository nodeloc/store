package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nodeloc-faka/database"
	"github.com/nodeloc-faka/models"
	"github.com/nodeloc-faka/services"
)

// PaymentHandler 支付处理器
type PaymentHandler struct {
	paymentService  *services.PaymentService
	orderService    *services.OrderService
	productService  *services.ProductService
	settingService  *services.SettingService
}

// NewPaymentHandler 创建支付处理器
func NewPaymentHandler() *PaymentHandler {
	return &PaymentHandler{
		paymentService:  services.NewPaymentService(),
		orderService:    services.NewOrderService(),
		productService:  services.NewProductService(),
		settingService:  services.NewSettingService(),
	}
}

// CreateOrder 创建订单并发起支付
func (h *PaymentHandler) CreateOrder(c *gin.Context) {
	// 检查用户是否登录
	userInterface, exists := c.Get("user")
	if !exists || userInterface == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "请先登录"})
		return
	}
	user := userInterface.(*models.User)

	// 检查用户是否被封禁
	if user.IsBlocked {
		c.JSON(http.StatusForbidden, gin.H{"error": "您的账号已被封禁"})
		return
	}

	// 获取请求参数
	productID, _ := strconv.ParseUint(c.PostForm("product_id"), 10, 64)
	quantity, _ := strconv.Atoi(c.PostForm("quantity"))
	contact := c.PostForm("contact")
	remark := c.PostForm("remark")

	if quantity <= 0 {
		quantity = 1
	}

	// 获取商品
	product, err := h.productService.FindByID(uint(productID))
	if err != nil || !product.IsActive {
		c.JSON(http.StatusBadRequest, gin.H{"error": "商品不存在或已下架"})
		return
	}

	// 检查库存
	if product.StockCount < quantity {
		c.JSON(http.StatusBadRequest, gin.H{"error": "库存不足"})
		return
	}

	// 计算金额
	totalAmount := product.Price * float64(quantity)

	// 创建订单
	order := &models.Order{
		UserID:      user.ID,
		ProductID:   uint(productID),
		Quantity:    quantity,
		TotalAmount: totalAmount,
		Status:      models.OrderStatusPending,
		Contact:     contact,
		Remark:      remark,
	}

	// 设置过期时间（30分钟）
	expiredAt := time.Now().Add(30 * time.Minute)
	order.ExpiredAt = &expiredAt

	if err := h.orderService.Create(order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建订单失败"})
		return
	}

	// 检查是否配置了支付
	if !h.paymentService.IsConfigured() {
		// 未配置支付，直接完成订单（免费模式）
		h.completeOrderFree(order)
		c.JSON(http.StatusOK, gin.H{
			"success":  true,
			"order_no": order.OrderNo,
			"message":  "订单创建成功",
			"redirect": "/order/" + order.OrderNo,
		})
		return
	}

	// 发起支付请求
	payResp, err := h.paymentService.CreatePayment(&services.CreatePaymentRequest{
		Amount:      int(totalAmount), // 假设1积分=1元
		Description: fmt.Sprintf("购买 %s x%d", product.Name, quantity),
		OrderID:     order.OrderNo,
	})

	if err != nil {
		// 支付请求失败，记录错误但不影响订单
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":    "发起支付失败: " + err.Error(),
			"order_no": order.OrderNo,
		})
		return
	}

	// 保存支付信息
	order.TransactionID = payResp.TransactionID
	order.PaymentURL = payResp.PaymentURL
	database.GetDB().Save(order)

	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"order_no":     order.OrderNo,
		"payment_url":  payResp.PaymentURL,
		"redirect":     payResp.PaymentURL,
	})
}

// completeOrderFree 免费模式完成订单
func (h *PaymentHandler) completeOrderFree(order *models.Order) {
	cardKeyService := services.NewCardKeyService()
	
	// 分配卡密
	availableCards, _ := cardKeyService.GetAvailableByProduct(order.ProductID, order.Quantity)
	if len(availableCards) >= order.Quantity {
		cardIDs := make([]uint, order.Quantity)
		for i := 0; i < order.Quantity; i++ {
			cardIDs[i] = availableCards[i].ID
		}
		cardKeyService.MarkAsSold(cardIDs, order.ID)
	}

	// 更新订单状态
	now := time.Now()
	order.Status = models.OrderStatusCompleted
	order.PayMethod = "free"
	order.PaidAt = &now
	database.GetDB().Save(order)

	// 更新商品库存和销量
	h.productService.UpdateStock(order.ProductID)
	h.productService.IncrementSales(order.ProductID, order.Quantity)
}

// PaymentCallback 支付回调
func (h *PaymentHandler) PaymentCallback(c *gin.Context) {
	// 解析回调参数
	callback := &services.PaymentCallback{
		TransactionID:     c.Query("transaction_id"),
		ExternalReference: c.Query("external_reference"),
		Status:            c.Query("status"),
		PaidAt:            c.Query("paid_at"),
		Signature:         c.Query("signature"),
	}
	callback.Amount, _ = strconv.Atoi(c.Query("amount"))
	callback.PlatformFee, _ = strconv.Atoi(c.Query("platform_fee"))
	callback.MerchantPoints, _ = strconv.Atoi(c.Query("merchant_points"))

	// 处理回调
	err := h.paymentService.ProcessPaymentCallback(callback)
	if err != nil {
		// 记录错误日志
		fmt.Printf("支付回调处理失败: %v\n", err)
		// 仍然重定向到订单页面，但显示错误
		c.Redirect(http.StatusFound, "/order/"+callback.ExternalReference+"?error="+err.Error())
		return
	}

	// 重定向到订单详情页
	c.Redirect(http.StatusFound, "/order/"+callback.ExternalReference+"?success=支付成功")
}

// QueryOrder 查询订单支付状态
func (h *PaymentHandler) QueryOrder(c *gin.Context) {
	orderNo := c.Param("order_no")

	order, err := h.orderService.FindByOrderNo(orderNo)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "订单不存在"})
		return
	}

	// 如果订单已完成或已取消，直接返回状态
	if order.Status != models.OrderStatusPending {
		c.JSON(http.StatusOK, gin.H{
			"order_no": order.OrderNo,
			"status":   order.Status,
			"paid":     order.Status == models.OrderStatusCompleted || order.Status == models.OrderStatusPaid,
		})
		return
	}

	// 如果有交易ID，查询支付状态
	if order.TransactionID != "" {
		queryResp, err := h.paymentService.QueryPayment(order.TransactionID)
		if err == nil && queryResp.Status == "completed" {
			// 支付已完成，处理订单
			callback := &services.PaymentCallback{
				TransactionID:     queryResp.TransactionID,
				ExternalReference: order.OrderNo,
				Amount:            queryResp.Amount,
				PlatformFee:       queryResp.PlatformFee,
				MerchantPoints:    queryResp.MerchantPoints,
				Status:            queryResp.Status,
			}
			if queryResp.PaidAt != nil {
				callback.PaidAt = *queryResp.PaidAt
			}
			h.paymentService.ProcessPaymentCallback(callback)
			
			// 重新查询订单
			order, _ = h.orderService.FindByOrderNo(orderNo)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"order_no": order.OrderNo,
		"status":   order.Status,
		"paid":     order.Status == models.OrderStatusCompleted || order.Status == models.OrderStatusPaid,
	})
}

// CancelOrder 取消订单
func (h *PaymentHandler) CancelOrder(c *gin.Context) {
	orderNo := c.Param("order_no")

	// 检查用户
	userInterface, exists := c.Get("user")
	if !exists || userInterface == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "请先登录"})
		return
	}
	user := userInterface.(*models.User)

	order, err := h.orderService.FindByOrderNo(orderNo)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "订单不存在"})
		return
	}

	// 检查是否是自己的订单
	if order.UserID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权操作此订单"})
		return
	}

	// 只能取消待支付的订单
	if order.Status != models.OrderStatusPending {
		c.JSON(http.StatusBadRequest, gin.H{"error": "订单状态不允许取消"})
		return
	}

	// 取消订单
	h.orderService.Cancel(order.ID)

	c.JSON(http.StatusOK, gin.H{"success": true})
}
