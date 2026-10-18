package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nodeloc-faka/database"
	"github.com/nodeloc-faka/models"
	"github.com/nodeloc-faka/services"
)

// APIHandler API处理器
type APIHandler struct {
	settingService  *services.SettingService
	categoryService *services.CategoryService
	productService  *services.ProductService
	orderService    *services.OrderService
	userService     *services.UserService
	paymentService  *services.PaymentService
	reviewService   *services.ReviewService
	shopService     *services.ShopService
}

// NewAPIHandler 创建API处理器
func NewAPIHandler() *APIHandler {
	return &APIHandler{
		settingService:  services.NewSettingService(),
		categoryService: services.NewCategoryService(),
		productService:  services.NewProductService(),
		orderService:    services.NewOrderService(),
		userService:     services.NewUserService(),
		paymentService:  services.NewPaymentService(),
		reviewService:   services.NewReviewService(),
		shopService:     services.NewShopService(),
	}
}

// GetSettings 获取网站设置
func (h *APIHandler) GetSettings(c *gin.Context) {
	settings := h.settingService.GetSiteSettings()
	c.JSON(http.StatusOK, settings)
}

// GetCategoriesWithProducts 获取带商品的分类列表
func (h *APIHandler) GetCategoriesWithProducts(c *gin.Context) {
	categories, err := h.categoryService.GetWithProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取分类失败"})
		return
	}
	c.JSON(http.StatusOK, categories)
}

// GetCategory 获取单个分类
func (h *APIHandler) GetCategory(c *gin.Context) {
	id := c.Param("id")
	category, err := h.categoryService.FindByID(ParseUint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "分类不存在"})
		return
	}
	c.JSON(http.StatusOK, category)
}

// GetProducts 获取商品列表（支持分类、关键词、价格区间、排序、分页）
func (h *APIHandler) GetProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	minPrice, _ := strconv.ParseFloat(c.Query("min_price"), 64)
	maxPrice, _ := strconv.ParseFloat(c.Query("max_price"), 64)

	params := services.SearchParams{
		CategoryID: ParseUint(c.Query("category_id")),
		Keyword:    c.Query("keyword"),
		MinPrice:   minPrice,
		MaxPrice:   maxPrice,
		Sort:       c.DefaultQuery("sort", "default"),
		PromoOnly:  c.Query("promo") == "1",
		Page:       page,
		PageSize:   pageSize,
	}

	products, total, err := h.productService.SearchProducts(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取商品失败"})
		return
	}
	for i := range products {
		products[i].ComputePromo()
	}

	// 批量查询评分统计
	productIDs := make([]uint, len(products))
	for i, p := range products {
		productIDs[i] = p.ID
	}
	ratingMap := h.reviewService.GetAvgRatingByProducts(productIDs)

	// 构建带评分的响应
	type productWithRating struct {
		models.Product
		AvgRating   float64 `json:"avg_rating"`
		ReviewCount int     `json:"review_count"`
	}
	result := make([]productWithRating, len(products))
	for i, p := range products {
		stats := ratingMap[p.ID]
		result[i] = productWithRating{
			Product:     p,
			AvgRating:   stats.AvgRating,
			ReviewCount: stats.Total,
		}
	}

	totalPages := (total + int64(pageSize) - 1) / int64(pageSize)
	c.JSON(http.StatusOK, gin.H{
		"products":    result,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": totalPages,
	})
}

// GetProduct 获取单个商品
func (h *APIHandler) GetProduct(c *gin.Context) {
	id := c.Param("id")
	product, err := h.productService.FindByID(ParseUint(id))
	if err != nil || !product.IsActive {
		c.JSON(http.StatusNotFound, gin.H{"error": "商品不存在"})
		return
	}
	product.ComputePromo()
	c.JSON(http.StatusOK, product)
}

// GetUserInfo 获取用户信息
func (h *APIHandler) GetUserInfo(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}

// GetOrders 获取用户订单列表
func (h *APIHandler) GetOrders(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	u := user.(*models.User)
	orders, err := h.orderService.GetByUser(u.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取订单失败"})
		return
	}
	c.JSON(http.StatusOK, orders)
}

// GetOrder 获取单个订单
func (h *APIHandler) GetOrder(c *gin.Context) {
	orderNo := c.Param("orderNo")
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	order, err := h.orderService.FindByOrderNo(orderNo)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "订单不存在"})
		return
	}

	// 检查是否是自己的订单
	u := user.(*models.User)
	if order.UserID != u.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权访问"})
		return
	}

	c.JSON(http.StatusOK, order)
}

// ConfirmReceipt 买家确认收货
func (h *APIHandler) ConfirmReceipt(c *gin.Context) {
	userInterface, exists := c.Get("user")
	if !exists || userInterface == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "请先登录"})
		return
	}
	u := userInterface.(*models.User)
	orderNo := c.Param("orderNo")
	order, err := h.orderService.ConfirmReceipt(orderNo, u.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, order)
}

// RepayOrder 重新支付订单
func (h *APIHandler) RepayOrder(c *gin.Context) {
	userInterface, exists := c.Get("user")
	if !exists || userInterface == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "请先登录"})
		return
	}
	u := userInterface.(*models.User)

	orderNo := c.Param("orderNo")
	order, err := h.orderService.FindByOrderNo(orderNo)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "订单不存在"})
		return
	}

	// 验证订单归属
	if order.UserID != u.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权访问此订单"})
		return
	}

	// 只有待支付订单才能重新支付
	if order.Status != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "订单状态不允许重新支付"})
		return
	}

	// 重新发起支付
	product, err := h.productService.FindByID(order.ProductID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取商品信息失败"})
		return
	}

	paymentReq := &services.CreatePaymentRequest{
		Amount:      int(order.TotalAmount),
		Description: fmt.Sprintf("购买商品: %s x%d", product.Name, order.Quantity),
		OrderID:     order.OrderNo,
	}

	paymentResp, err := h.paymentService.CreatePayment(paymentReq)
	if err != nil {
		log.Printf("重新支付失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "支付接口调用失败: " + err.Error()})
		return
	}

	// 更新订单支付信息
	order.TransactionID = paymentResp.TransactionID
	order.PaymentURL = paymentResp.PaymentURL
	h.orderService.Update(order)

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"payment_url": paymentResp.PaymentURL,
		"order":       order,
	})
}

// CreateOrder 创建订单
func (h *APIHandler) CreateOrder(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	u := user.(*models.User)

	// 检查用户是否被封禁
	if u.IsBlocked {
		c.JSON(http.StatusForbidden, gin.H{"error": "您的账号已被封禁"})
		return
	}

	var req struct {
		ProductID uint   `json:"product_id" binding:"required"`
		Quantity  int    `json:"quantity" binding:"required,min=1"`
		Contact   string `json:"contact"`
		Remark    string `json:"remark"`
		PayMethod string `json:"pay_method"` // nodeloc / balance；默认 nodeloc
		CouponID  uint   `json:"coupon_id"`  // 优惠券ID（可选）
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	// 验证商品
	product, err := h.productService.FindByID(req.ProductID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "商品不存在"})
		return
	}

	if !product.IsActive {
		c.JSON(http.StatusBadRequest, gin.H{"error": "商品已下架"})
		return
	}

	// 应用优惠券（如果有）
	couponSvc := services.NewCouponService()
	var couponDiscount float64
	var appliedCouponID uint
	if req.CouponID > 0 {
		coupon, err := couponSvc.FindByID(req.CouponID)
		if err == nil {
			amount := product.GetEffectivePrice() * float64(req.Quantity)
			result, err := couponSvc.Validate(coupon.Code, product.ShopID, amount)
			if err == nil {
				couponDiscount = result.Discount
				appliedCouponID = coupon.ID
			}
		}
	}
	_ = couponDiscount
	_ = appliedCouponID

	// 余额支付分支
	if req.PayMethod == models.PayMethodBalance {
		order, err := h.orderService.CreateAndProcessByBalance(u.ID, req.ProductID, req.Quantity, req.Contact, req.Remark)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success":   true,
			"order_no":  order.OrderNo,
			"order":     order,
			"completed": true,
		})
		return
	}

	// 预检查支付是否已配置，避免创建脏订单
	if !h.paymentService.IsConfigured() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "NodeLoc 能量支付暂未开启，请使用余额支付"})
		return
	}

	// 创建订单（使用 CreatePendingOrder 创建待支付订单）
	order, err := h.orderService.CreatePendingOrder(u.ID, req.ProductID, req.Quantity, req.Contact, req.Remark)
	if err != nil {
		if err.Error() == "insufficient stock" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "库存不足"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "创建订单失败"})
		}
		return
	}

	// 发起支付（NodeLoc Payment）
	paymentReq := &services.CreatePaymentRequest{
		Amount:      int(order.TotalAmount), // 能量 = 人民币（1:1）
		Description: fmt.Sprintf("购买商品: %s x%d", product.Name, req.Quantity),
		OrderID:     order.OrderNo,
	}

	paymentResp, err := h.paymentService.CreatePayment(paymentReq)
	if err != nil {
		// 支付接口调用失败，取消待支付订单并返回错误
		fmt.Printf("支付接口调用失败: %v\n", err)
		h.orderService.Cancel(order.ID)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "支付接口调用失败，请稍后重试或使用余额支付",
		})
		return
	}

	// 更新订单的支付信息
	order.TransactionID = paymentResp.TransactionID
	order.PaymentURL = paymentResp.PaymentURL
	h.orderService.Update(order)

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"order_no":    order.OrderNo,
		"payment_url": paymentResp.PaymentURL, // 前端需要跳转到这个URL
		"order":       order,
	})
}

// GetProductReviews 获取商品评价列表
func (h *APIHandler) GetProductReviews(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.reviewService.GetByProduct(ParseUint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取评价失败"})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// CreateProductReview 提交商品评价
func (h *APIHandler) CreateProductReview(c *gin.Context) {
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "请先登录"})
		return
	}
	u := userInterface.(*models.User)

	productID := ParseUint(c.Param("id"))

	// 检查是否已评价
	if h.reviewService.HasReviewed(u.ID, productID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "您已评价过该商品"})
		return
	}

	// 检查是否有权评价
	if !h.reviewService.CanReview(u.ID, productID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "只有购买并完成订单的用户才能评价"})
		return
	}

	var body struct {
		Rating  int    `json:"rating" binding:"required,min=1,max=5"`
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误：评分须在 1-5 之间"})
		return
	}

	review := &models.ProductReview{
		ProductID: productID,
		UserID:    u.ID,
		Rating:    body.Rating,
		Title:     body.Title,
		Content:   body.Content,
	}
	if err := h.reviewService.Create(review); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "提交评价失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "评价已提交"})
}

// ParseUint 解析 uint
func ParseUint(s string) uint {
	id, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0
	}
	return uint(id)
}

// SubscribeNewsletter 邮件订阅
func (h *APIHandler) SubscribeNewsletter(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请输入有效的邮箱地址"})
		return
	}

	// 存入数据库（重复订阅忽略）
	sub := &models.NewsletterSubscriber{Email: req.Email}
	database.GetDB().Where("email = ?", req.Email).FirstOrCreate(sub)

	siteName := h.settingService.Get(services.SettingSiteName)
	if siteName == "" {
		siteName = "NodeLoc发卡系统"
	}

	emailSvc := services.NewEmailService()
	if emailSvc.IsEnabled() {
		emailSvc.SendNewsletterConfirm(req.Email, siteName)
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "订阅成功"})
}
func (h *APIHandler) GetDealsOfWeek(c *gin.Context) {
	productIDStr := h.settingService.Get(services.SettingDealsOfWeekProductID)
	if productIDStr == "" {
		c.JSON(http.StatusOK, gin.H{"product": nil})
		return
	}

	productID, err := strconv.ParseUint(productIDStr, 10, 32)
	if err != nil || productID == 0 {
		c.JSON(http.StatusOK, gin.H{"product": nil})
		return
	}

	product, err := h.productService.FindByID(uint(productID))
	if err != nil || !product.IsActive {
		c.JSON(http.StatusOK, gin.H{"product": nil})
		return
	}
	product.ComputePromo()

	// 附带评分统计
	ratingMap := h.reviewService.GetAvgRatingByProducts([]uint{product.ID})
	stats := ratingMap[product.ID]

	c.JSON(http.StatusOK, gin.H{
		"product":      product,
		"avg_rating":   stats.AvgRating,
		"review_count": stats.Total,
	})
}
