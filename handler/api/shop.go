package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nodeloc-faka/database"
	"github.com/nodeloc-faka/models"
	"github.com/nodeloc-faka/services"
)

// ShopHandler 店铺相关 API 处理器
type ShopHandler struct {
	shopService       *services.ShopService
	productService    *services.ProductService
	cardKeyService    *services.CardKeyService
	orderService      *services.OrderService
	balanceService    *services.BalanceService
	withdrawalService *services.WithdrawalService
	settingService    *services.SettingService
	reviewService     *services.ReviewService
}

// NewShopHandler 创建店铺处理器
func NewShopHandler() *ShopHandler {
	return &ShopHandler{
		shopService:       services.NewShopService(),
		productService:    services.NewProductService(),
		cardKeyService:    services.NewCardKeyService(),
		orderService:      services.NewOrderService(),
		balanceService:    services.NewBalanceService(),
		withdrawalService: services.NewWithdrawalService(),
		settingService:    services.NewSettingService(),
		reviewService:     services.NewReviewService(),
	}
}

func parsePage(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return page, pageSize
}

func currentUser(c *gin.Context) (*models.User, bool) {
	v, ok := c.Get("user")
	if !ok || v == nil {
		return nil, false
	}
	u, ok := v.(*models.User)
	if !ok || u == nil {
		return nil, false
	}
	return u, true
}

// ============================== 公开接口 ==============================

// ListShops 公开店铺列表（仅已批准）
func (h *ShopHandler) ListShops(c *gin.Context) {
	page, pageSize := parsePage(c)
	keyword := c.Query("keyword")

	db := database.GetDB().Model(&models.Shop{}).Where("status = ?", models.ShopStatusApproved)
	if keyword != "" {
		db = db.Where("name LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	var total int64
	db.Count(&total)

	var shops []models.Shop
	offset := (page - 1) * pageSize
	if err := db.Order("is_official desc, id asc").Offset(offset).Limit(pageSize).Find(&shops).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 填充各店铺的商品数量
	type shopWithCount struct {
		models.Shop
		ProductCount int64 `json:"product_count"`
	}
	result := make([]shopWithCount, len(shops))
	for i, s := range shops {
		var cnt int64
		h.productService.CountByShop(s.ID, &cnt)
		result[i] = shopWithCount{Shop: s, ProductCount: cnt}
	}
	c.JSON(http.StatusOK, gin.H{"data": result, "total": total, "page": page, "page_size": pageSize})
}

// GetShop 获取店铺详情
func (h *ShopHandler) GetShop(c *gin.Context) {
	id := ParseUint(c.Param("id"))
	shop, err := h.shopService.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "店铺不存在"})
		return
	}
	if shop.Status != models.ShopStatusApproved {
		c.JSON(http.StatusNotFound, gin.H{"error": "店铺不存在"})
		return
	}
	c.JSON(http.StatusOK, shop)
}

// GetShopProducts 获取店铺的商品
func (h *ShopHandler) GetShopProducts(c *gin.Context) {
	id := ParseUint(c.Param("id"))
	shop, err := h.shopService.FindByID(id)
	if err != nil || shop.Status != models.ShopStatusApproved {
		c.JSON(http.StatusNotFound, gin.H{"error": "店铺不存在"})
		return
	}
	products, err := h.productService.GetByShop(id, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, products)
}

// ============================== 用户：开店申请 ==============================

// Apply 申请开店
func (h *ShopHandler) Apply(c *gin.Context) {
	u, ok := currentUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "请先登录"})
		return
	}
	if h.settingService.Get(services.SettingShopApplyEnabled) == "false" {
		c.JSON(http.StatusForbidden, gin.H{"error": "开店申请已关闭"})
		return
	}
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Logo        string `json:"logo"`
		Contact     string `json:"contact"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	shop := &models.Shop{
		UserID:      u.ID,
		Name:        req.Name,
		Description: req.Description,
		Logo:        req.Logo,
		Contact:     req.Contact,
	}
	// 没有设置 logo 时，使用用户头像
	if shop.Logo == "" && u.AvatarURL != "" {
		shop.Logo = u.AvatarURL
	}
	if err := h.shopService.Apply(shop); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, shop)
}

// MyShop 我的店铺
func (h *ShopHandler) MyShop(c *gin.Context) {
	u, ok := currentUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "请先登录"})
		return
	}
	shop, err := h.shopService.GetByUser(u.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "未开店"})
		return
	}
	c.JSON(http.StatusOK, shop)
}

// UpdateMyShop 更新我的店铺信息（审核通过后修改将触发重新审核）
func (h *ShopHandler) UpdateMyShop(c *gin.Context) {
	u, ok := currentUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "请先登录"})
		return
	}
	shop, err := h.shopService.GetByUser(u.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "店铺不存在"})
		return
	}
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Logo        string `json:"logo"`
		Contact     string `json:"contact"`
		Features    string `json:"features"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	if req.Name != "" {
		shop.Name = req.Name
	}
	shop.Description = req.Description
	shop.Logo = req.Logo
	shop.Contact = req.Contact
	shop.Features = req.Features
	// 修改信息后需重新审核
	if shop.Status == models.ShopStatusApproved {
		shop.Status = models.ShopStatusPending
		shop.RejectReason = ""
	}
	if err := h.shopService.Update(shop); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, shop)
}

// ============================== 余额 / 提现 ==============================

// Balance 我的余额 + 提现配置
func (h *ShopHandler) Balance(c *gin.Context) {
	u, ok := currentUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "请先登录"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"balance":            u.Balance,
		"withdrawal_enabled": h.withdrawalService.IsEnabled(),
		"fee_rate":           h.withdrawalService.GetFeeRate(),
		"min_amount":         h.withdrawalService.GetMinAmount(),
	})
}

// BalanceTxs 余额流水
func (h *ShopHandler) BalanceTxs(c *gin.Context) {
	u, ok := currentUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "请先登录"})
		return
	}
	page, pageSize := parsePage(c)
	list, total, err := h.balanceService.GetUserTxs(u.ID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list, "total": total, "page": page, "page_size": pageSize})
}

// ApplyWithdrawal 提现申请
func (h *ShopHandler) ApplyWithdrawal(c *gin.Context) {
	u, ok := currentUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "请先登录"})
		return
	}
	shop, err := h.shopService.GetByUser(u.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请先开通店铺"})
		return
	}
	if shop.Status != models.ShopStatusApproved {
		c.JSON(http.StatusBadRequest, gin.H{"error": "店铺未通过审核"})
		return
	}
	var req struct {
		Amount float64 `json:"amount" binding:"required,gt=0"`
		Remark string  `json:"remark"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	wreq, err := h.withdrawalService.Apply(u.ID, shop.ID, req.Amount, req.Remark)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, wreq)
}

// MyWithdrawals 我的提现记录
func (h *ShopHandler) MyWithdrawals(c *gin.Context) {
	u, ok := currentUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "请先登录"})
		return
	}
	page, pageSize := parsePage(c)
	list, total, err := h.withdrawalService.GetByUser(u.ID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list, "total": total, "page": page, "page_size": pageSize})
}

// ============================== 店主：商品/订单管理 ==============================

// MyShopProducts 店铺自己的商品（包括下架的）
func (h *ShopHandler) MyShopProducts(c *gin.Context) {
	u, ok := currentUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "请先登录"})
		return
	}
	shop, err := h.shopService.GetByUser(u.ID)
	if err != nil || shop.Status != models.ShopStatusApproved {
		c.JSON(http.StatusForbidden, gin.H{"error": "店铺未通过审核"})
		return
	}
	products, err := h.productService.GetByShop(shop.ID, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	out := make([]models.OwnerProduct, len(products))
	for i, p := range products {
		out[i] = models.ToOwnerProduct(p)
	}
	c.JSON(http.StatusOK, out)
}

// myApprovedShop 获取当前用户的已审核店铺（内部辅助）
func (h *ShopHandler) myApprovedShop(c *gin.Context) (*models.User, *models.Shop, bool) {
	u, ok := currentUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "请先登录"})
		return nil, nil, false
	}
	shop, err := h.shopService.GetByUser(u.ID)
	if err != nil || shop.Status != models.ShopStatusApproved {
		c.JSON(http.StatusForbidden, gin.H{"error": "店铺未通过审核"})
		return nil, nil, false
	}
	return u, shop, true
}

// MyShopCreateProduct 店主创建商品
func (h *ShopHandler) MyShopCreateProduct(c *gin.Context) {
	_, shop, ok := h.myApprovedShop(c)
	if !ok {
		return
	}
	var req struct {
		CategoryID   uint       `json:"category_id" binding:"required"`
		Name         string     `json:"name" binding:"required"`
		Description  string     `json:"description"`
		Price        float64    `json:"price" binding:"required"`
		OrigPrice    float64    `json:"orig_price"`
		Image        string     `json:"image"`
		Images       string     `json:"images"`
		IsActive     bool       `json:"is_active"`
		DeliveryType int        `json:"delivery_type"`
		FixedContent string     `json:"fixed_content"`
		StockCount   int        `json:"stock_count"`
		PromoPrice   float64    `json:"promo_price"`
		PromoStart   *time.Time `json:"promo_start"`
		PromoEnd     *time.Time `json:"promo_end"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	if req.Images == "" {
		req.Images = "[]"
	}
	product := &models.Product{
		FixedContent: req.FixedContent,
		ShopID:       shop.ID,
		CategoryID:   req.CategoryID,
		Name:         req.Name,
		Description:  req.Description,
		Price:        req.Price,
		OrigPrice:    req.OrigPrice,
		Image:        req.Image,
		Images:       req.Images,
		IsActive:     req.IsActive,
		DeliveryType: req.DeliveryType,
		StockCount:   req.StockCount,
		PromoPrice:   req.PromoPrice,
		PromoStart:   req.PromoStart,
		PromoEnd:     req.PromoEnd,
	}
	if !h.validateProduct(c, product, true) {
		return
	}
	if err := h.productService.Create(product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建商品失败"})
		return
	}
	product.ComputePromo()
	c.JSON(http.StatusOK, models.ToOwnerProduct(*product))
}

// MyShopUpdateProduct 店主更新商品
func (h *ShopHandler) MyShopUpdateProduct(c *gin.Context) {
	_, shop, ok := h.myApprovedShop(c)
	if !ok {
		return
	}
	pid := ParseUint(c.Param("pid"))
	product, err := h.productService.FindByID(pid)
	if err != nil || product.ShopID != shop.ID {
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
		IsActive     *bool      `json:"is_active"`
		DeliveryType *int       `json:"delivery_type"`
		FixedContent *string    `json:"fixed_content"`
		StockCount   *int       `json:"stock_count"`
		PromoPrice   *float64   `json:"promo_price"`
		PromoStart   *time.Time `json:"promo_start"`
		PromoEnd     *time.Time `json:"promo_end"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	categoryChanged := req.CategoryID != nil && *req.CategoryID != product.CategoryID
	if req.CategoryID != nil { product.CategoryID = *req.CategoryID }
	if req.Name != "" { product.Name = req.Name }
	product.Description = req.Description
	if req.Price != nil { product.Price = *req.Price }
	if req.OrigPrice != nil { product.OrigPrice = *req.OrigPrice }
	product.Image = req.Image
	if req.Images != nil {
		if *req.Images == "" { product.Images = "[]" } else { product.Images = *req.Images }
	}
	if req.IsActive != nil { product.IsActive = *req.IsActive }
	if req.DeliveryType != nil { product.DeliveryType = *req.DeliveryType }
	if req.FixedContent != nil { product.FixedContent = *req.FixedContent }
	if req.StockCount != nil { product.StockCount = *req.StockCount }
	if req.PromoPrice != nil { product.PromoPrice = *req.PromoPrice }
	product.PromoStart = req.PromoStart
	product.PromoEnd = req.PromoEnd
	if !h.validateProduct(c, product, categoryChanged) {
		return
	}

	if err := h.productService.Update(product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新商品失败"})
		return
	}
	product.ComputePromo()
	c.JSON(http.StatusOK, models.ToOwnerProduct(*product))
}

// validateProduct 校验商品字段，失败时直接写 400 响应
func (h *ShopHandler) validateProduct(c *gin.Context, product *models.Product, checkCategory bool) bool {
	if err := h.productService.Validate(product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return false
	}
	if checkCategory && !h.productService.CategoryExists(product.CategoryID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "商品分类不存在"})
		return false
	}
	return true
}

// MyShopDeleteProduct 店主删除商品
func (h *ShopHandler) MyShopDeleteProduct(c *gin.Context) {
	_, shop, ok := h.myApprovedShop(c)
	if !ok {
		return
	}
	pid := ParseUint(c.Param("pid"))
	product, err := h.productService.FindByID(pid)
	if err != nil || product.ShopID != shop.ID {
		c.JSON(http.StatusNotFound, gin.H{"error": "商品不存在"})
		return
	}
	if err := h.productService.Delete(pid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// MyShopGetCardKeys 获取商品卡密列表
func (h *ShopHandler) MyShopGetCardKeys(c *gin.Context) {
	_, shop, ok := h.myApprovedShop(c)
	if !ok {
		return
	}
	pid := ParseUint(c.Param("pid"))
	product, err := h.productService.FindByID(pid)
	if err != nil || product.ShopID != shop.ID {
		c.JSON(http.StatusNotFound, gin.H{"error": "商品不存在"})
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "50"))
	if pageSize < 1 || pageSize > 200 {
		pageSize = 50
	}
	cardKeys, total, err := h.cardKeyService.GetWithPagination(pid, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取卡密失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"card_keys": cardKeys, "total": total})
}

// MyShopAddCardKeys 批量添加卡密
func (h *ShopHandler) MyShopAddCardKeys(c *gin.Context) {
	_, shop, ok := h.myApprovedShop(c)
	if !ok {
		return
	}
	pid := ParseUint(c.Param("pid"))
	product, err := h.productService.FindByID(pid)
	if err != nil || product.ShopID != shop.ID {
		c.JSON(http.StatusNotFound, gin.H{"error": "商品不存在"})
		return
	}
	var req struct {
		CardsText string `json:"cards_text" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	res, err := h.cardKeyService.BatchCreate(pid, req.CardsText)
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

// MyShopDeleteCardKey 删除卡密
func (h *ShopHandler) MyShopDeleteCardKey(c *gin.Context) {
	_, shop, ok := h.myApprovedShop(c)
	if !ok {
		return
	}
	kid := ParseUint(c.Param("kid"))
	cardKey, err := h.cardKeyService.FindByID(kid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "卡密不存在"})
		return
	}
	product, err := h.productService.FindByID(cardKey.ProductID)
	if err != nil || product.ShopID != shop.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权操作"})
		return
	}
	if err := h.cardKeyService.Delete(kid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// MyShopOrders 店铺订单列表
func (h *ShopHandler) MyShopOrders(c *gin.Context) {
	_, shop, ok := h.myApprovedShop(c)
	if !ok {
		return
	}
	page, pageSize := parsePage(c)
	status, _ := strconv.Atoi(c.DefaultQuery("status", "-1"))
	orders, total, err := h.orderService.GetByShop(shop.ID, page, pageSize, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": orders, "total": total, "page": page, "page_size": pageSize})
}

// MyShopShipOrder 卖家发货
func (h *ShopHandler) MyShopShipOrder(c *gin.Context) {
	u, _, ok := h.myApprovedShop(c)
	if !ok {
		return
	}
	oid := ParseUint(c.Param("oid"))
	var req struct {
		Content string `json:"content" binding:"required"`
		Note    string `json:"note"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请填写发货内容"})
		return
	}
	if err := h.orderService.ShipOrder(oid, u.ID, req.Content, req.Note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// ============================== 店主：Dashboard ==============================

// MyShopDashboard 店铺仪表盘统计数据
func (h *ShopHandler) MyShopDashboard(c *gin.Context) {
	_, shop, ok := h.myApprovedShop(c)
	if !ok {
		return
	}

	stats, err := h.orderService.GetShopDashboard(shop.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取统计数据失败"})
		return
	}

	// 商品数量
	var productCount int64
	h.productService.CountByShop(shop.ID, &productCount)

	// 评价统计
	reviewStats, _ := h.reviewService.GetShopReviewStats(shop.ID)

	c.JSON(http.StatusOK, gin.H{
		"order_stats":   stats,
		"product_count": productCount,
		"review_stats":  reviewStats,
	})
}

// ============================== 店主：评价管理 ==============================

// MyShopReviews 店铺所有商品的评价列表
func (h *ShopHandler) MyShopReviews(c *gin.Context) {
	_, shop, ok := h.myApprovedShop(c)
	if !ok {
		return
	}

	page, pageSize := parsePage(c)
	ratingFilter, _ := strconv.Atoi(c.DefaultQuery("rating", "0"))

	reviews, total, err := h.reviewService.GetByShop(shop.ID, page, pageSize, ratingFilter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取评价失败"})
		return
	}

	// 评价统计
	reviewStats, _ := h.reviewService.GetShopReviewStats(shop.ID)

	c.JSON(http.StatusOK, gin.H{
		"data":      reviews,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
		"stats":     reviewStats,
	})
}
