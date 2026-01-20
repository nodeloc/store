package public

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nodeloc-faka/models"
	"github.com/nodeloc-faka/services"
)

// Handler 前台处理器
type Handler struct {
	settingService  *services.SettingService
	categoryService *services.CategoryService
	productService  *services.ProductService
	orderService    *services.OrderService
	userService     *services.UserService
}

// NewHandler 创建前台处理器
func NewHandler() *Handler {
	return &Handler{
		settingService:  services.NewSettingService(),
		categoryService: services.NewCategoryService(),
		productService:  services.NewProductService(),
		orderService:    services.NewOrderService(),
		userService:     services.NewUserService(),
	}
}

// getSiteData 获取网站公共数据
func (h *Handler) getSiteData(c *gin.Context) gin.H {
	settings := h.settingService.GetSiteSettings()
	categories, _ := h.categoryService.GetActive()

	// 获取当前用户
	var currentUser *models.User
	if userInterface, exists := c.Get("user"); exists && userInterface != nil {
		currentUser = userInterface.(*models.User)
	}

	// 检查支付是否启用
	paymentEnabled := h.settingService.Get(services.SettingPaymentEnabled) == "true"

	return gin.H{
		"site":           settings,
		"categories":     categories,
		"user":           currentUser,
		"paymentEnabled": paymentEnabled,
	}
}

// Index 首页
func (h *Handler) Index(c *gin.Context) {
	data := h.getSiteData(c)

	// 获取带商品的分类
	categoriesWithProducts, _ := h.categoryService.GetWithProducts()

	data["title"] = "首页"
	data["categoriesWithProducts"] = categoriesWithProducts

	c.HTML(http.StatusOK, "public/index", data)
}

// Category 分类页面
func (h *Handler) Category(c *gin.Context) {
	data := h.getSiteData(c)

	categoryID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	category, err := h.categoryService.FindByID(uint(categoryID))
	if err != nil {
		c.HTML(http.StatusNotFound, "public/404", data)
		return
	}

	products, _ := h.productService.GetByCategory(uint(categoryID))

	data["title"] = category.Name
	data["category"] = category
	data["products"] = products

	c.HTML(http.StatusOK, "public/category", data)
}

// Product 商品详情
func (h *Handler) Product(c *gin.Context) {
	data := h.getSiteData(c)

	productID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	product, err := h.productService.FindByID(uint(productID))
	
	// 商品不存在或未激活，显示 404
	if err != nil || product == nil || !product.IsActive {
		data["title"] = "商品不存在"
		c.HTML(http.StatusNotFound, "public/404", data)
		return
	}

	data["title"] = product.Name
	data["product"] = product

	c.HTML(http.StatusOK, "public/product", data)
}

// Purchase 购买商品
func (h *Handler) Purchase(c *gin.Context) {
	data := h.getSiteData(c)

	// 检查是否登录
	userInterface, exists := c.Get("user")
	if !exists || userInterface == nil {
		c.Redirect(http.StatusFound, "/login?redirect="+c.Request.URL.String())
		return
	}
	user := userInterface.(*models.User)

	// 检查用户是否被封禁
	if user.IsBlocked {
		data["title"] = "错误"
		data["error"] = "您的账号已被封禁，无法购买"
		c.HTML(http.StatusForbidden, "public/error", data)
		return
	}

	productID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	product, err := h.productService.FindByID(uint(productID))
	
	// 商品不存在或未激活，显示 404
	if err != nil || product == nil || !product.IsActive {
		data["title"] = "商品不存在"
		c.HTML(http.StatusNotFound, "public/404", data)
		return
	}
	
	// 检查库存
	if product.StockCount <= 0 {
		data["title"] = "库存不足"
		data["error"] = "该商品暂时缺货"
		c.HTML(http.StatusOK, "public/error", data)
		return
	}

	if c.Request.Method == "GET" {
		data["title"] = "购买 - " + product.Name
		data["product"] = product
		c.HTML(http.StatusOK, "public/purchase", data)
		return
	}

	// POST 请求已改为通过 API 处理，这里重定向回 GET
	c.Redirect(http.StatusFound, c.Request.URL.Path)
}

// Orders 我的订单
func (h *Handler) Orders(c *gin.Context) {
	data := h.getSiteData(c)

	// 检查是否登录
	userInterface, exists := c.Get("user")
	if !exists || userInterface == nil {
		c.Redirect(http.StatusFound, "/login?redirect=/orders")
		return
	}
	user := userInterface.(*models.User)

	orders, _ := h.orderService.GetByUser(user.ID)

	data["title"] = "我的订单"
	data["orders"] = orders

	c.HTML(http.StatusOK, "public/orders", data)
}

// OrderDetail 订单详情
func (h *Handler) OrderDetail(c *gin.Context) {
	data := h.getSiteData(c)

	orderNo := c.Param("order_no")
	order, err := h.orderService.FindByOrderNo(orderNo)
	if err != nil {
		c.HTML(http.StatusNotFound, "public/404", data)
		return
	}

	// 检查是否是自己的订单
	userInterface, exists := c.Get("user")
	if !exists || userInterface == nil {
		c.Redirect(http.StatusFound, "/login?redirect=/order/"+orderNo)
		return
	}
	user := userInterface.(*models.User)

	if order.UserID != user.ID {
		c.HTML(http.StatusForbidden, "public/403", data)
		return
	}

	data["title"] = "订单详情 - " + order.OrderNo
	data["order"] = order

	// 如果订单是待支付状态，显示支付页面
	if order.Status == models.OrderStatusPending {
		c.HTML(http.StatusOK, "public/order_pending", data)
		return
	}

	c.HTML(http.StatusOK, "public/order_detail", data)
}

// Profile 个人中心
func (h *Handler) Profile(c *gin.Context) {
	data := h.getSiteData(c)

	// 检查是否登录
	userInterface, exists := c.Get("user")
	if !exists || userInterface == nil {
		c.Redirect(http.StatusFound, "/login?redirect=/profile")
		return
	}
	user := userInterface.(*models.User)

	orders, _ := h.orderService.GetByUser(user.ID)

	data["title"] = "个人中心"
	data["orders"] = orders

	c.HTML(http.StatusOK, "public/profile", data)
}
