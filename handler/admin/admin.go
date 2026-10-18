package admin

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

// AdminHandler 管理员 API 处理器
type AdminHandler struct {
	categoryService *services.CategoryService
	productService  *services.ProductService
	cardKeyService  *services.CardKeyService
	orderService    *services.OrderService
	userService     *services.UserService
	settingService  *services.SettingService
	balanceService  *services.BalanceService
	shopService     *services.ShopService
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
		balanceService:  services.NewBalanceService(),
		shopService:     services.NewShopService(),
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
		Image       string `json:"image"`
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
		Image:       req.Image,
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
		Image       string `json:"image"`
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
	category.Image = req.Image
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

	out := make([]models.OwnerProduct, len(products))
	for i, p := range products {
		out[i] = models.ToOwnerProduct(p)
	}
	c.JSON(http.StatusOK, gin.H{
		"products": out,
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
	c.JSON(http.StatusOK, gin.H{"product": models.ToOwnerProduct(*product)})
}

// CreateProduct 创建商品
func (h *AdminHandler) CreateProduct(c *gin.Context) {
	var req struct {
		CategoryID  uint       `json:"category_id" binding:"required"`
		Name        string     `json:"name" binding:"required"`
		Description string     `json:"description"`
		Price       float64    `json:"price" binding:"required"`
		OrigPrice   float64    `json:"orig_price"`
		Image       string     `json:"image"`
		Images      string     `json:"images"`
		StockCount  int        `json:"stock_count"`
		Sort        int        `json:"sort"`
		IsActive    bool       `json:"is_active"`
		DeliveryType int       `json:"delivery_type"`
		FixedContent string    `json:"fixed_content"`
		PromoPrice  float64    `json:"promo_price"`
		PromoStart  *time.Time `json:"promo_start"`
		PromoEnd    *time.Time `json:"promo_end"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	if req.Images == "" {
		req.Images = "[]"
	}

	// 确保官方店存在，并获取其 ID
	officialShop, err := h.shopService.EnsureOfficialShop()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取官方店铺失败"})
		return
	}

	product := &models.Product{
		ShopID:       officialShop.ID,
		CategoryID:   req.CategoryID,
		Name:         req.Name,
		Description:  req.Description,
		Price:        req.Price,
		OrigPrice:    req.OrigPrice,
		Image:        req.Image,
		Images:       req.Images,
		StockCount:   req.StockCount,
		Sort:         req.Sort,
		IsActive:     req.IsActive,
		DeliveryType: req.DeliveryType,
		FixedContent: req.FixedContent,
		PromoPrice:   req.PromoPrice,
		PromoStart:   req.PromoStart,
		PromoEnd:     req.PromoEnd,
	}
	if err := h.productService.Validate(product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !h.productService.CategoryExists(product.CategoryID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "商品分类不存在"})
		return
	}
	if err := h.productService.Create(product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建商品失败"})
		return
	}
	product.ComputePromo()
	c.JSON(http.StatusOK, gin.H{"product": models.ToOwnerProduct(*product)})
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
		CategoryID   *uint      `json:"category_id"`
		Name         string     `json:"name"`
		Description  string     `json:"description"`
		Price        *float64   `json:"price"`
		OrigPrice    *float64   `json:"orig_price"`
		Image        string     `json:"image"`
		Images       *string    `json:"images"`
		StockCount   *int       `json:"stock_count"`
		Sort         *int       `json:"sort"`
		IsActive     *bool      `json:"is_active"`
		DeliveryType *int       `json:"delivery_type"`
		FixedContent *string    `json:"fixed_content"`
		PromoPrice   *float64   `json:"promo_price"`
		PromoStart   *time.Time `json:"promo_start"`
		PromoEnd     *time.Time `json:"promo_end"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	if req.CategoryID != nil { product.CategoryID = *req.CategoryID }
	if req.Name != "" { product.Name = req.Name }
	product.Description = req.Description
	if req.Price != nil { product.Price = *req.Price }
	if req.OrigPrice != nil { product.OrigPrice = *req.OrigPrice }
	product.Image = req.Image
	if req.Images != nil {
		if *req.Images == "" { product.Images = "[]" } else { product.Images = *req.Images }
	}
	if req.StockCount != nil { product.StockCount = *req.StockCount }
	if req.Sort != nil { product.Sort = *req.Sort }
	if req.IsActive != nil { product.IsActive = *req.IsActive }
	if req.DeliveryType != nil { product.DeliveryType = *req.DeliveryType }
	if req.FixedContent != nil { product.FixedContent = *req.FixedContent }
	if req.PromoPrice != nil { product.PromoPrice = *req.PromoPrice }
	product.PromoStart = req.PromoStart
	product.PromoEnd = req.PromoEnd
	if err := h.productService.Validate(product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.productService.Update(product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新商品失败"})
		return
	}
	product.ComputePromo()
	c.JSON(http.StatusOK, gin.H{"product": models.ToOwnerProduct(*product)})
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

	res, err := h.cardKeyService.BatchCreate(req.ProductID, req.CardsText)
	if err != nil {
		if se, ok := err.(*services.ServiceError); ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": se.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "添加卡密失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "添加成功", "count": res.Added, "duplicates": res.Duplicates, "invalid": res.Invalid, "lines": res.Lines})
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
	keyword := c.Query("keyword")

	db := database.GetDB().Model(&models.Order{})
	if status >= 0 {
		db = db.Where("orders.status = ?", status)
	}
	if keyword != "" {
		db = db.Joins("LEFT JOIN users ON users.id = orders.user_id").
			Where("orders.order_no LIKE ? OR users.username LIKE ? OR users.name LIKE ?",
				"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	var total int64
	db.Count(&total)

	var orders []models.Order
	offset := (page - 1) * pageSize
	if err := db.Preload("User").Preload("Product").
		Order("orders.id desc").
		Offset(offset).Limit(pageSize).
		Find(&orders).Error; err != nil {
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

// AdjustUserBalance 管理员调整用户余额
func (h *AdminHandler) AdjustUserBalance(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var req struct {
		Amount      float64 `json:"amount" binding:"required"`
		Description string  `json:"description" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误：金额和备注均必填"})
		return
	}
	if req.Amount == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "调整金额不能为 0"})
		return
	}

	// 使用数据库事务
	tx := database.GetDB().Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "开启事务失败"})
		return
	}

	err := h.balanceService.AddBalance(tx, uint(id), req.Amount, models.BalanceTxAdminAdjust, req.Description, "admin", 0)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tx.Commit()

	user, _ := h.userService.FindByID(uint(id))
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
		"site_name":           h.settingService.Get(services.SettingSiteName),
		"site_description":    h.settingService.Get(services.SettingSiteDescription),
		"site_logo":           h.settingService.Get(services.SettingSiteLogo),
		"contact_tel":         h.settingService.Get(services.SettingContactTel),
		"contact_tel_label":   h.settingService.Get(services.SettingContactTelLabel),
		"contact_email":       h.settingService.Get(services.SettingContactEmail),
		"contact_address":     h.settingService.Get(services.SettingContactAddress),
		"footer_text":         h.settingService.Get(services.SettingFooterText),
		"announcement":        h.settingService.Get(services.SettingAnnouncement),
		"home_banners":        h.settingService.Get(services.SettingHomeBanners),
		"home_promo_banners":  h.settingService.Get(services.SettingHomePromoBanners),
		"home_flash_banners":  h.settingService.Get(services.SettingHomeFlashBanners),
		"home_offer_cards":    h.settingService.Get(services.SettingHomeOfferCards),
		"home_bestsell_cta":   h.settingService.Get(services.SettingHomeBestsellCTA),
		"home_newsletter_img": h.settingService.Get(services.SettingHomeNewsletterImg),
		"deals_of_week_shop_id":    h.settingService.Get(services.SettingDealsOfWeekShopID),
		"deals_of_week_product_id": h.settingService.Get(services.SettingDealsOfWeekProductID),
		"shop_apply_enabled":       h.settingService.Get(services.SettingShopApplyEnabled),
		"home_features":            h.settingService.Get(services.SettingHomeFeatures),
		"home_hot_deal":            h.settingService.Get(services.SettingHomeHotDeal),
		"footer_config":            h.settingService.Get(services.SettingFooterConfig),
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
	// Keys managed via .env, not editable through admin UI
	envOnlyKeys := map[string]bool{
		"nodeloc_client_id": true, "nodeloc_client_secret": true, "nodeloc_redirect_uri": true,
		"payment_id": true, "payment_secret": true, "payment_callback_uri": true, "payment_enabled": true,
	}
	settings := make(map[string]string)
	for key, value := range req {
		if envOnlyKeys[key] {
			continue
		}
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

// ============================================
// 邮件订阅管理
// ============================================

// GetNewsletterSubscribers 获取订阅列表
func (h *AdminHandler) GetNewsletterSubscribers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "50"))
	keyword := c.Query("keyword")

	db := database.GetDB().Model(&models.NewsletterSubscriber{})
	if keyword != "" {
		db = db.Where("email LIKE ?", "%"+keyword+"%")
	}

	var total int64
	db.Count(&total)

	var subs []models.NewsletterSubscriber
	offset := (page - 1) * pageSize
	db.Order("id desc").Offset(offset).Limit(pageSize).Find(&subs)

	c.JSON(http.StatusOK, gin.H{"data": subs, "total": total, "page": page, "page_size": pageSize})
}

// DeleteNewsletterSubscriber 删除订阅
func (h *AdminHandler) DeleteNewsletterSubscriber(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	database.GetDB().Delete(&models.NewsletterSubscriber{}, id)
	c.JSON(http.StatusOK, gin.H{"success": true})
}
func (h *AdminHandler) GetDashboard(c *gin.Context) {
	db := database.GetDB()

	// 基础计数
	userCount := h.userService.Count()
	productCount := h.productService.Count()
	orderCount := h.orderService.Count()
	categoryCount := h.categoryService.Count()

	// 订单各状态数量
	pendingOrders := h.orderService.CountByStatus(models.OrderStatusPending)
	completedOrders := h.orderService.CountByStatus(models.OrderStatusCompleted)
	cancelledOrders := h.orderService.CountByStatus(models.OrderStatusCancelled)

	// 总销售额（已完成+已支付+已发货）
	totalSales := h.orderService.GetTotalSales()
	todaySales := h.orderService.GetTodaySales()

	// 今日新增用户
	var todayUsers int64
	today := time.Now().Format("2006-01-02")
	db.Model(&models.User{}).Where("DATE(created_at) = ?", today).Count(&todayUsers)

	// 今日新增订单
	var todayOrders int64
	db.Model(&models.Order{}).Where("DATE(created_at) = ?", today).Count(&todayOrders)

	// 待审核店铺数
	var pendingShops int64
	db.Model(&models.Shop{}).Where("status = ?", models.ShopStatusPending).Count(&pendingShops)

	// 待审核提现数
	var pendingWithdrawals int64
	db.Model(&models.WithdrawalRequest{}).Where("status = ?", models.WithdrawalStatusPending).Count(&pendingWithdrawals)

	// 最近 10 条订单
	var recentOrders []models.Order
	db.Preload("User").Preload("Product").
		Order("id desc").Limit(10).
		Find(&recentOrders)

	// 最近 7 天每日订单数和销售额
	type DayStat struct {
		Date   string  `json:"date"`
		Orders int64   `json:"orders"`
		Sales  float64 `json:"sales"`
	}
	var dailyStats []DayStat
	for i := 6; i >= 0; i-- {
		day := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		var cnt int64
		var sales float64
		db.Model(&models.Order{}).
			Where("DATE(created_at) = ? AND status IN ?", day,
				[]int{models.OrderStatusCompleted, models.OrderStatusPaid, models.OrderStatusShipped}).
			Count(&cnt)
		db.Model(&models.Order{}).
			Where("DATE(created_at) = ? AND status IN ?", day,
				[]int{models.OrderStatusCompleted, models.OrderStatusPaid, models.OrderStatusShipped}).
			Select("COALESCE(SUM(total_amount), 0)").Scan(&sales)
		dailyStats = append(dailyStats, DayStat{Date: day, Orders: cnt, Sales: sales})
	}

	// 热销商品 Top 5
	type TopProduct struct {
		ID         uint    `json:"id"`
		Name       string  `json:"name"`
		SalesCount int     `json:"sales_count"`
		Price      float64 `json:"price"`
		Image      string  `json:"image"`
	}
	var topProducts []TopProduct
	db.Model(&models.Product{}).
		Select("id, name, sales_count, price, image").
		Where("is_active = ?", true).
		Order("sales_count desc").
		Limit(5).Scan(&topProducts)

	c.JSON(http.StatusOK, gin.H{
		"stats": gin.H{
			"users":               userCount,
			"products":            productCount,
			"orders":              orderCount,
			"categories":          categoryCount,
			"pending_orders":      pendingOrders,
			"completed_orders":    completedOrders,
			"cancelled_orders":    cancelledOrders,
			"total_sales":         totalSales,
			"today_sales":         todaySales,
			"today_users":         todayUsers,
			"today_orders":        todayOrders,
			"pending_shops":       pendingShops,
			"pending_withdrawals": pendingWithdrawals,
		},
		"recent_orders": recentOrders,
		"daily_stats":   dailyStats,
		"top_products":  topProducts,
	})
}
