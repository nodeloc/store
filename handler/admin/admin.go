package admin

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nodeloc-faka/models"
	"github.com/nodeloc-faka/services"
)

// AdminHandler 管理员 API 处理器
type AdminHandler struct {
	categoryService *services.CategoryService
	productService  *services.ProductService
	cardKeyService  *services.CardKeyService
	orderService    *services.OrderService
	userService     *services.UserService
	settingService  *services.SettingService
}

// NewAdminHandler 创建管理员处理器
func NewAdminHandler() *AdminHandler {
	return &AdminHandler{
		categoryService: services.NewCategoryService(),
		productService:  services.NewProductService(),
		cardKeyService:  services.NewCardKeyService(),
		orderService:    services.NewOrderService(),
		userService:     services.NewUserService(),
		settingService:  services.NewSettingService(),
	}
}

// ============================================
// 商品分类管理
// ============================================

// GetCategories 获取所有分类
func (h *AdminHandler) GetCategories(c *gin.Context) {
	categories, err := h.categoryService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取分类失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"categories": categories})
}

// GetCategory 获取单个分类
func (h *AdminHandler) GetCategory(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	category, err := h.categoryService.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "分类不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"category": category})
}

// CreateCategory 创建分类
func (h *AdminHandler) CreateCategory(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
		Sort        int    `json:"sort"`
		IsActive    bool   `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	category := &models.Category{
		Name:        req.Name,
		Description: req.Description,
		Icon:        req.Icon,
		Sort:        req.Sort,
		IsActive:    req.IsActive,
	}

	if err := h.categoryService.Create(category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建分类失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"category": category})
}

// UpdateCategory 更新分类
func (h *AdminHandler) UpdateCategory(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	category, err := h.categoryService.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "分类不存在"})
		return
	}

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
		Sort        *int   `json:"sort"`
		IsActive    *bool  `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	if req.Name != "" {
		category.Name = req.Name
	}
	if req.Description != "" {
		category.Description = req.Description
	}
	if req.Icon != "" {
		category.Icon = req.Icon
	}
	if req.Sort != nil {
		category.Sort = *req.Sort
	}
	if req.IsActive != nil {
		category.IsActive = *req.IsActive
	}

	if err := h.categoryService.Update(category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新分类失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"category": category})
}

// DeleteCategory 删除分类
func (h *AdminHandler) DeleteCategory(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := h.categoryService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除分类失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// ============================================
// 商品管理
// ============================================

// GetProducts 获取所有商品（分页）
func (h *AdminHandler) GetProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	products, total, err := h.productService.GetWithPagination(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取商品失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"products": products,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// GetProduct 获取单个商品
func (h *AdminHandler) GetProduct(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	product, err := h.productService.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "商品不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"product": product})
}

// CreateProduct 创建商品
func (h *AdminHandler) CreateProduct(c *gin.Context) {
	var req struct {
		CategoryID  uint    `json:"category_id" binding:"required"`
		Name        string  `json:"name" binding:"required"`
		Description string  `json:"description"`
		Price       float64 `json:"price" binding:"required"`
		OrigPrice   float64 `json:"orig_price"`
		Image       string  `json:"image"`
		Sort        int     `json:"sort"`
		IsActive    bool    `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	product := &models.Product{
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		OrigPrice:   req.OrigPrice,
		Image:       req.Image,
		Sort:        req.Sort,
		IsActive:    req.IsActive,
		StockCount:  0,
		SalesCount:  0,
	}

	if err := h.productService.Create(product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建商品失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"product": product})
}

// UpdateProduct 更新商品
func (h *AdminHandler) UpdateProduct(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	product, err := h.productService.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "商品不存在"})
		return
	}

	var req struct {
		CategoryID  *uint    `json:"category_id"`
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Price       *float64 `json:"price"`
		OrigPrice   *float64 `json:"orig_price"`
		Image       string   `json:"image"`
		Sort        *int     `json:"sort"`
		IsActive    *bool    `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	if req.CategoryID != nil {
		product.CategoryID = *req.CategoryID
	}
	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Price != nil {
		product.Price = *req.Price
	}
	if req.OrigPrice != nil {
		product.OrigPrice = *req.OrigPrice
	}
	if req.Image != "" {
		product.Image = req.Image
	}
	if req.Sort != nil {
		product.Sort = *req.Sort
	}
	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}

	if err := h.productService.Update(product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新商品失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"product": product})
}

// DeleteProduct 删除商品
func (h *AdminHandler) DeleteProduct(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := h.productService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除商品失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// ============================================
// 卡密管理
// ============================================

// GetCardKeys 获取商品的卡密列表
func (h *AdminHandler) GetCardKeys(c *gin.Context) {
	productID, _ := strconv.ParseUint(c.Query("product_id"), 10, 32)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "50"))

	cardKeys, total, err := h.cardKeyService.GetWithPagination(uint(productID), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取卡密失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"card_keys": cardKeys,
		"total":     total,
		"page":      page,
		"pageSize":  pageSize,
	})
}

// AddCardKeys 批量添加卡密
func (h *AdminHandler) AddCardKeys(c *gin.Context) {
	var req struct {
		ProductID uint   `json:"product_id" binding:"required"`
		CardsText string `json:"cards_text" binding:"required"` // 格式：每行一个卡密，支持 "卡号----密码" 或只有卡号
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	count, err := h.cardKeyService.BatchCreate(req.ProductID, req.CardsText)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "添加卡密失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "添加成功", "count": count})
}

// DeleteCardKey 删除卡密
func (h *AdminHandler) DeleteCardKey(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := h.cardKeyService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除卡密失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// ============================================
// 订单管理
// ============================================

// GetOrders 获取所有订单
func (h *AdminHandler) GetOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status, _ := strconv.Atoi(c.DefaultQuery("status", "-1"))

	orders, total, err := h.orderService.GetWithPagination(page, pageSize, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取订单失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"orders":   orders,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// GetOrder 获取单个订单
func (h *AdminHandler) GetOrder(c *gin.Context) {
	orderNo := c.Param("orderNo")
	order, err := h.orderService.FindByOrderNo(orderNo)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "订单不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"order": order})
}

// UpdateOrderStatus 更新订单状态
func (h *AdminHandler) UpdateOrderStatus(c *gin.Context) {
	orderNo := c.Param("orderNo")
	var req struct {
		Status int `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	order, err := h.orderService.FindByOrderNo(orderNo)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "订单不存在"})
		return
	}

	// 根据状态调用不同方法
	switch req.Status {
	case models.OrderStatusPaid:
		err = h.orderService.MarkAsPaid(order.ID)
	case models.OrderStatusCompleted:
		err = h.orderService.Complete(order.ID)
	case models.OrderStatusCancelled:
		err = h.orderService.Cancel(order.ID)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的状态"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新订单状态失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// ============================================
// 用户管理
// ============================================

// GetUsers 获取所有用户
func (h *AdminHandler) GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	users, total, err := h.userService.GetWithPagination(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users":    users,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// GetUser 获取单个用户
func (h *AdminHandler) GetUser(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	user, err := h.userService.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}

// UpdateUser 更新用户信息
func (h *AdminHandler) UpdateUser(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var req struct {
		IsAdmin   *bool `json:"is_admin"`
		IsBlocked *bool `json:"is_blocked"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	if req.IsAdmin != nil {
		if err := h.userService.SetAdmin(uint(id), *req.IsAdmin); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "更新管理员状态失败"})
			return
		}
	}

	if req.IsBlocked != nil {
		if *req.IsBlocked {
			h.userService.Block(uint(id))
		} else {
			h.userService.Unblock(uint(id))
		}
	}

	user, _ := h.userService.FindByID(uint(id))
	c.JSON(http.StatusOK, gin.H{"user": user})
}

// ============================================
// 系统设置管理
// ============================================

// GetSettings 获取所有设置
func (h *AdminHandler) GetSettings(c *gin.Context) {
	settings := gin.H{
		"site_name":             h.settingService.Get(services.SettingSiteName),
		"site_description":      h.settingService.Get(services.SettingSiteDescription),
		"footer_text":           h.settingService.Get(services.SettingFooterText),
		"announcement":          h.settingService.Get(services.SettingAnnouncement),
		"nodeloc_client_id":     h.settingService.Get(services.SettingNodeLocClientID),
		"nodeloc_client_secret": h.settingService.Get(services.SettingNodeLocClientSecret),
		"nodeloc_redirect_uri":  h.settingService.Get(services.SettingNodeLocRedirectURI),
		"payment_enabled":       h.settingService.Get(services.SettingPaymentEnabled) == "true",
		"payment_id":            h.settingService.Get(services.SettingPaymentID),
		"payment_secret":        h.settingService.Get(services.SettingPaymentSecret),
		"payment_callback_uri":  h.settingService.Get(services.SettingPaymentCallback),
	}
	c.JSON(http.StatusOK, gin.H{"settings": settings})
}

// UpdateSettings 更新设置
func (h *AdminHandler) UpdateSettings(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	// 转换为 map[string]string
	settings := make(map[string]string)
	for key, value := range req {
		switch v := value.(type) {
		case string:
			settings[key] = v
		case bool:
			if v {
				settings[key] = "true"
			} else {
				settings[key] = "false"
			}
		case float64:
			settings[key] = strconv.FormatFloat(v, 'f', -1, 64)
		default:
			settings[key] = fmt.Sprintf("%v", v)
		}
	}

	if err := h.settingService.SetMultiple(settings); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新设置失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// ============================================
// 统计信息
// ============================================

// GetDashboard 获取仪表板统计信息
func (h *AdminHandler) GetDashboard(c *gin.Context) {
	userCount := h.userService.Count()
	productCount := h.productService.Count()
	orderCount := h.orderService.Count()
	categoryCount := h.categoryService.Count()

	c.JSON(http.StatusOK, gin.H{
		"stats": gin.H{
			"users":      userCount,
			"products":   productCount,
			"orders":     orderCount,
			"categories": categoryCount,
		},
	})
}
