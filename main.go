package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nodeloc-faka/config"
	"github.com/nodeloc-faka/database"
	"github.com/nodeloc-faka/handler"
	"github.com/nodeloc-faka/handler/admin"
	"github.com/nodeloc-faka/handler/api"
	"github.com/nodeloc-faka/middleware"
	"github.com/nodeloc-faka/models"
	"github.com/nodeloc-faka/oauth"
	"github.com/nodeloc-faka/services"
)

func main() {
	log.Println("==========================================")
	log.Println("NodeLoc 社区发卡系统 API 启动中...")
	log.Println("==========================================")

	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 连接数据库
	log.Println("正在连接数据库...")
	db, err := database.Connect(cfg.Database)
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}
	_ = db
	log.Println("✓ 数据库连接成功")

	// 数据库迁移
	log.Println("正在执行数据库迁移...")
	if err := models.AutoMigrate(database.GetDB()); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}
	log.Println("✓ 数据库迁移完成")

	// 初始化系统（简化版 - 只初始化基础设置）
	initSystemSimple()

	// 插入演示数据（仅在数据库为空时，且未禁用时）
	if os.Getenv("DISABLE_DEMO_DATA") != "true" {
		services.SeedDemoData()
	}

	// 创建 OAuth 客户端
	oauthClient := oauth.NewClient(
		cfg.NodeLocURL,
		cfg.NodeLocClientID,
		cfg.NodeLocClientSecret,
		cfg.NodeLocRedirectURI,
	)

	// 创建 session 存储
	sessionStore := middleware.NewSessionStore()

	// 创建处理器
	authHandler := handler.NewAuthHandler(oauthClient)
	paymentHandler := handler.NewPaymentHandler()
	apiHandler := api.NewAPIHandler()
	adminHandler := admin.NewAdminHandler()
	shopHandler := api.NewShopHandler()
	shopAdminHandler := admin.NewShopAdminHandler()
	uploadHandler := handler.NewUploadHandler()
	blogHandler := api.NewBlogHandler()
	blogAdminHandler := admin.NewBlogAdminHandler()
	couponHandler := api.NewCouponHandler()
	couponAdminHandler := admin.NewCouponAdminHandler()

	// 设置 Gin
	router := gin.Default()
	log.Println("✓ API 模式启动")

	router.Use(middleware.SecurityHeaders())
	router.Use(middleware.SessionMiddleware(sessionStore))
	sessionStore.StartJanitor(time.Hour)

	// ========================================
	// API 路由 (JSON 响应)
	// ========================================

	// 公开 API
	apiGroup := router.Group("/api")
	{
		apiGroup.GET("/settings", apiHandler.GetSettings)
		apiGroup.GET("/categories/with-products", apiHandler.GetCategoriesWithProducts)
		apiGroup.GET("/categories/:id", apiHandler.GetCategory)
		apiGroup.GET("/products", apiHandler.GetProducts)
		apiGroup.GET("/products/:id", apiHandler.GetProduct)
		apiGroup.GET("/products/:id/reviews", apiHandler.GetProductReviews)
		apiGroup.GET("/shops", shopHandler.ListShops)
		apiGroup.GET("/shops/:id", shopHandler.GetShop)
		apiGroup.GET("/shops/:id/products", shopHandler.GetShopProducts)
		apiGroup.GET("/deals-of-week", apiHandler.GetDealsOfWeek)
		apiGroup.GET("/stats", apiHandler.GetPublicStats)
		apiGroup.GET("/recent-purchases", apiHandler.GetRecentPurchases)
		apiGroup.POST("/newsletter/subscribe", apiHandler.SubscribeNewsletter)
		apiGroup.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})
		// 博客
		apiGroup.GET("/blog/posts", blogHandler.ListPosts)
		apiGroup.GET("/blog/posts/:slug", blogHandler.GetPost)
		apiGroup.GET("/blog/categories", blogHandler.GetCategories)
		apiGroup.GET("/blog/recent", blogHandler.RecentPosts)
	}

	// 需要认证的 API
	apiAuthGroup := apiGroup.Group("", middleware.AuthRequired())
	{
		apiAuthGroup.GET("/user/info", apiHandler.GetUserInfo)
		apiAuthGroup.GET("/orders", apiHandler.GetOrders)
		apiAuthGroup.GET("/orders/:orderNo", apiHandler.GetOrder)
		apiAuthGroup.POST("/orders/:orderNo/repay", apiHandler.RepayOrder)
		apiAuthGroup.POST("/orders/create", apiHandler.CreateOrder)
		apiAuthGroup.GET("/order/:order_no/status", paymentHandler.QueryOrder)
		apiAuthGroup.POST("/order/:order_no/cancel", paymentHandler.CancelOrder)

		// 评价
		apiAuthGroup.POST("/products/:id/reviews", apiHandler.CreateProductReview)

		// 已登录用户通用上传
		apiAuthGroup.POST("/upload/image", uploadHandler.UploadImage)

		// 店铺：申请 / 管理自己的店铺
		apiAuthGroup.POST("/shop/apply", shopHandler.Apply)
		apiAuthGroup.GET("/shop/me", shopHandler.MyShop)
		apiAuthGroup.PUT("/shop/me", shopHandler.UpdateMyShop)
		apiAuthGroup.GET("/shop/me/products", shopHandler.MyShopProducts)
		apiAuthGroup.POST("/shop/me/products", shopHandler.MyShopCreateProduct)
		apiAuthGroup.PUT("/shop/me/products/:pid", shopHandler.MyShopUpdateProduct)
		apiAuthGroup.DELETE("/shop/me/products/:pid", shopHandler.MyShopDeleteProduct)
		apiAuthGroup.GET("/shop/me/products/:pid/cardkeys", shopHandler.MyShopGetCardKeys)
		apiAuthGroup.POST("/shop/me/products/:pid/cardkeys", shopHandler.MyShopAddCardKeys)
		apiAuthGroup.DELETE("/shop/me/cardkeys/:kid", shopHandler.MyShopDeleteCardKey)
		apiAuthGroup.GET("/shop/me/orders", shopHandler.MyShopOrders)
		apiAuthGroup.POST("/shop/me/orders/:oid/ship", shopHandler.MyShopShipOrder)
		apiAuthGroup.GET("/shop/me/dashboard", shopHandler.MyShopDashboard)
		apiAuthGroup.GET("/shop/me/reviews", shopHandler.MyShopReviews)
		apiAuthGroup.POST("/orders/:orderNo/confirm", apiHandler.ConfirmReceipt)

		// 优惠券：校验（需登录）
		apiAuthGroup.POST("/coupons/validate", couponHandler.ValidateCoupon)
		// 店主：管理自己的优惠券
		apiAuthGroup.GET("/shop/me/coupons", couponHandler.MyShopCoupons)
		apiAuthGroup.POST("/shop/me/coupons", couponHandler.MyShopCreateCoupon)
		apiAuthGroup.PUT("/shop/me/coupons/:id", couponHandler.MyShopUpdateCoupon)
		apiAuthGroup.DELETE("/shop/me/coupons/:id", couponHandler.MyShopDeleteCoupon)

		// 余额 & 提现
		apiAuthGroup.GET("/balance", shopHandler.Balance)
		apiAuthGroup.GET("/balance/txs", shopHandler.BalanceTxs)
		apiAuthGroup.POST("/withdrawals/apply", shopHandler.ApplyWithdrawal)
		apiAuthGroup.GET("/withdrawals", shopHandler.MyWithdrawals)
	}

	// 管理员 API
	adminAPIGroup := apiGroup.Group("/admin", middleware.AdminRequired())
	{
		// 仪表板
		adminAPIGroup.GET("/dashboard", adminHandler.GetDashboard)

		// 文件上传
		adminAPIGroup.POST("/upload/image", uploadHandler.UploadImage)

		// 分类管理
		adminAPIGroup.GET("/categories", adminHandler.GetCategories)
		adminAPIGroup.GET("/categories/:id", adminHandler.GetCategory)
		adminAPIGroup.POST("/categories", adminHandler.CreateCategory)
		adminAPIGroup.PUT("/categories/:id", adminHandler.UpdateCategory)
		adminAPIGroup.DELETE("/categories/:id", adminHandler.DeleteCategory)

		// 商品管理
		adminAPIGroup.GET("/products", adminHandler.GetProducts)
		adminAPIGroup.GET("/products/:id", adminHandler.GetProduct)
		adminAPIGroup.POST("/products", adminHandler.CreateProduct)
		adminAPIGroup.PUT("/products/:id", adminHandler.UpdateProduct)
		adminAPIGroup.DELETE("/products/:id", adminHandler.DeleteProduct)

		// 卡密管理
		adminAPIGroup.GET("/card-keys", adminHandler.GetCardKeys)
		adminAPIGroup.POST("/card-keys", adminHandler.AddCardKeys)
		adminAPIGroup.DELETE("/card-keys/:id", adminHandler.DeleteCardKey)

		// 订单管理
		adminAPIGroup.GET("/orders", adminHandler.GetOrders)
		adminAPIGroup.GET("/orders/:orderNo", adminHandler.GetOrder)
		adminAPIGroup.PUT("/orders/:orderNo/status", adminHandler.UpdateOrderStatus)

		// 用户管理
		adminAPIGroup.GET("/users", adminHandler.GetUsers)
		adminAPIGroup.GET("/users/:id", adminHandler.GetUser)
		adminAPIGroup.PUT("/users/:id", adminHandler.UpdateUser)
		adminAPIGroup.POST("/users/:id/balance", adminHandler.AdjustUserBalance)

		// 系统设置
		adminAPIGroup.GET("/settings", adminHandler.GetSettings)
		adminAPIGroup.PUT("/settings", adminHandler.UpdateSettings)

		// 店铺审核
		adminAPIGroup.GET("/shops", shopAdminHandler.ListShops)
		adminAPIGroup.POST("/shops/:id/approve", shopAdminHandler.ApproveShop)
		adminAPIGroup.POST("/shops/:id/reject", shopAdminHandler.RejectShop)
		adminAPIGroup.POST("/shops/:id/block", shopAdminHandler.BlockShop)
		adminAPIGroup.PUT("/shops/:id", shopAdminHandler.UpdateShop)

		// 提现审核
		adminAPIGroup.GET("/withdrawals", shopAdminHandler.ListWithdrawals)
		adminAPIGroup.POST("/withdrawals/:id/approve", shopAdminHandler.ApproveWithdrawal)
		adminAPIGroup.POST("/withdrawals/:id/reject", shopAdminHandler.RejectWithdrawal)

		// 博客管理
		adminAPIGroup.GET("/blog/posts", blogAdminHandler.List)
		adminAPIGroup.POST("/blog/posts", blogAdminHandler.Create)
		adminAPIGroup.PUT("/blog/posts/:id", blogAdminHandler.Update)
		adminAPIGroup.DELETE("/blog/posts/:id", blogAdminHandler.Delete)

		// 邮件订阅管理
		adminAPIGroup.GET("/newsletter/subscribers", adminHandler.GetNewsletterSubscribers)
		adminAPIGroup.DELETE("/newsletter/subscribers/:id", adminHandler.DeleteNewsletterSubscriber)

		// 优惠券管理
		adminAPIGroup.GET("/coupons", couponAdminHandler.List)
		adminAPIGroup.POST("/coupons", couponAdminHandler.Create)
		adminAPIGroup.PUT("/coupons/:id", couponAdminHandler.Update)
		adminAPIGroup.DELETE("/coupons/:id", couponAdminHandler.Delete)
	}

	// ========================================
	// OAuth 认证路由（后端处理）
	// ========================================
	router.GET("/auth/login", authHandler.Login)
	router.GET("/auth/callback", authHandler.Callback)
	router.GET("/auth/logout", authHandler.Logout)

	// ========================================
	// 支付回调路由（后端处理）
	// ========================================
	router.GET("/payment/callback", paymentHandler.PaymentCallback)

	// 静态文件服务（上传的图片）
	router.Static("/uploads", "./uploads")

	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// 启动服务器
	addr := ":8080" // 固定使用 8080 端口（Docker 内部端口）
	log.Println("")
	log.Println("╔════════════════════════════════════════════════════════════╗")
	log.Println("║            NodeLoc 社区发卡系统 API v1.0                   ║")
	log.Println("╠════════════════════════════════════════════════════════════╣")
	log.Printf("║ API 服务: http://0.0.0.0%s                                ", addr)
	log.Println("║ 前端: 由 Vue + Nginx 提供                                  ║")
	log.Println("╚════════════════════════════════════════════════════════════╝")
	log.Println("")
	log.Printf("API 服务器启动在端口 8080，等待请求...")

	// 使用 0.0.0.0 监听所有网络接口（Docker 需要）
	if err := router.Run("0.0.0.0" + addr); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}

// initSystemSimple 简化的系统初始化（前后端分离版本）
func initSystemSimple() {
	settingService := services.NewSettingService()
	shopService := services.NewShopService()

	// 检查是否已初始化
	initialized := settingService.Get(services.SettingInitialized)
	if initialized == "true" {
		// 已初始化，但需要确保支付配置同步到数据库
		syncPaymentConfig(settingService)
		// 确保官方店铺存在
		if _, err := shopService.EnsureOfficialShop(); err != nil {
			log.Printf("⚠️  确保官方店铺失败: %v", err)
		}
		return
	}

	log.Println("首次运行，正在初始化系统...")

	// 保存基础设置
	err := settingService.SetMultiple(map[string]string{
		services.SettingSiteName:               "NodeLoc 社区发卡",
		services.SettingSiteDescription:        "基于 NodeLoc OAuth 的社区发卡系统",
		services.SettingFooterText:             "© 2026 NodeLoc 社区发卡系统",
		services.SettingInitialized:            "true",
		services.SettingShopApplyEnabled:       "true",
		services.SettingWithdrawalEnabled:      "true",
		services.SettingWithdrawalFeeRate:      "2",
		services.SettingWithdrawalMinAmount:    "10",
		services.SettingWithdrawalAutoMode:     "false",
	})
	if err != nil {
		log.Fatalf("保存系统设置失败: %v", err)
	}

	// 同步支付配置
	syncPaymentConfig(settingService)

	// 创建官方店铺并迁移现有商品
	if _, err := shopService.EnsureOfficialShop(); err != nil {
		log.Printf("⚠️  创建官方店铺失败: %v", err)
	} else {
		log.Println("✓ 官方店铺已就绪")
	}

	log.Println("✓ 系统初始化完成")
	log.Println("提示: 首个登录的 NodeLoc 用户将自动成为管理员")
}

// syncPaymentConfig 同步支付配置到数据库
func syncPaymentConfig(settingService *services.SettingService) {
	paymentID := os.Getenv("PAYMENT_ID")
	paymentSecret := os.Getenv("PAYMENT_SECRET")
	paymentCallback := os.Getenv("PAYMENT_CALLBACK")

	if paymentID != "" {
		settingService.Set(services.SettingPaymentID, paymentID)
		log.Printf("✓ 同步 PAYMENT_ID: %s", paymentID)
	}
	if paymentSecret != "" {
		settingService.Set(services.SettingPaymentSecret, paymentSecret)
		log.Printf("✓ 同步 PAYMENT_SECRET: %s", maskSecret(paymentSecret))
	}
	if paymentCallback != "" {
		settingService.Set(services.SettingPaymentCallback, paymentCallback)
		log.Printf("✓ 同步 PAYMENT_CALLBACK: %s", paymentCallback)
	}

	if paymentID != "" && paymentSecret != "" {
		settingService.Set(services.SettingPaymentEnabled, "true")
		log.Println("✓ 支付功能已启用")
	} else {
		log.Println("⚠️  支付未配置，请在 .env 文件中设置 PAYMENT_ID 和 PAYMENT_SECRET")
	}
}

// maskSecret 隐藏密钥（只显示前后几位）
func maskSecret(secret string) string {
	if len(secret) <= 8 {
		return "****"
	}
	return secret[:4] + "****" + secret[len(secret)-4:]
}
