package handler

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nodeloc-faka/models"
	"github.com/nodeloc-faka/services"
)

// PaymentHandler 支付处理器
type PaymentHandler struct {
	paymentService *services.PaymentService
	orderService   *services.OrderService
}

// NewPaymentHandler 创建支付处理器
func NewPaymentHandler() *PaymentHandler {
	return &PaymentHandler{
		paymentService: services.NewPaymentService(),
		orderService:   services.NewOrderService(),
	}
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

	// 获取前端基础 URL（从 NODELOC_REDIRECT_URI 推断，或使用 Host header）
	frontendBase := ""
	redirectURI := os.Getenv("NODELOC_REDIRECT_URI")
	if redirectURI != "" {
		// 从 https://domain.com/auth/callback 提取 https://domain.com
		if idx := len(redirectURI) - len("/auth/callback"); idx > 0 && redirectURI[idx:] == "/auth/callback" {
			frontendBase = redirectURI[:idx]
		}
	}
	if frontendBase == "" {
		// fallback: 用请求的 Host
		scheme := "https"
		if c.Request.TLS == nil && c.GetHeader("X-Forwarded-Proto") != "https" {
			scheme = "http"
		}
		frontendBase = scheme + "://" + c.Request.Host
	}
	orderURL := frontendBase + "/order/" + url.PathEscape(callback.ExternalReference)

	// 处理回调
	if err := h.paymentService.ProcessPaymentCallback(callback); err != nil {
		log.Printf("支付回调处理失败 tx=%s order=%s: %v", callback.TransactionID, callback.ExternalReference, err)
		c.Redirect(http.StatusFound, orderURL+"?error=payment_failed")
		return
	}

	// 重定向到订单详情页（使用完整 URL 避免代理跟随问题）
	c.Redirect(http.StatusFound, orderURL+"?success=1")
}

// QueryOrder 主动向 NodeLoc 查询待支付订单的支付状态（浏览器回调丢失时的补偿）
func (h *PaymentHandler) QueryOrder(c *gin.Context) {
	user := GetCurrentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "请先登录"})
		return
	}

	orderNo := c.Param("order_no")
	order, err := h.orderService.FindByOrderNo(orderNo)
	if err != nil || order.UserID != user.ID {
		c.JSON(http.StatusNotFound, gin.H{"error": "订单不存在"})
		return
	}

	if order.Status == models.OrderStatusPending && order.TransactionID != "" {
		queryResp, err := h.paymentService.QueryPayment(order.TransactionID)
		if err == nil && queryResp.Status == "completed" {
			// 主动查询是服务端直连 NodeLoc，结果可信，无需签名
			if _, processErr := h.orderService.ProcessPaymentCallback(
				queryResp.TransactionID,
				queryResp.Amount,
				queryResp.PlatformFee,
				queryResp.MerchantPoints,
			); processErr != nil {
				log.Printf("QueryOrder 处理订单 %s 失败: %v", orderNo, processErr)
			}
			if refreshed, err := h.orderService.FindByOrderNo(orderNo); err == nil {
				order = refreshed
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"order_no": order.OrderNo,
		"status":   order.Status,
		"paid":     order.Status != models.OrderStatusPending && order.Status != models.OrderStatusCancelled,
	})
}

// CancelOrder 取消订单
func (h *PaymentHandler) CancelOrder(c *gin.Context) {
	user := GetCurrentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "请先登录"})
		return
	}

	order, err := h.orderService.FindByOrderNo(c.Param("order_no"))
	if err != nil || order.UserID != user.ID {
		c.JSON(http.StatusNotFound, gin.H{"error": "订单不存在"})
		return
	}

	if err := h.orderService.CancelPending(order.ID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "订单状态不允许取消"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}
