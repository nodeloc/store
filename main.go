package main

import (
	"log"
	"net/http"
	"os"

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
	uploadHandler := handler.NewUploadHandler()

	// 设置 Gin
	router := gin.Default()
	log.Println("✓ API 模式启动")

	// 应用 Session 中间件
	router.Use(middleware.SessionMiddleware(sessionStore))

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
		apiGroup.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})
	}

	// 需要认证的 API
	apiAuthGroup := apiGroup.Group("", middleware.AuthRequired())
	{
		apiAuthGroup.GET("/user/info", apiHandler.GetUserInfo)
		apiAuthGroup.GET("/orders", apiHandler.GetOrders)
		apiAuthGroup.GET("/orders/:orderNo", apiHandler.GetOrder)
		apiAuthGroup.POST("/orders/:orderNo/repay", apiHandler.RepayOrder)
		apiAuthGroup.POST("/orders/create", apiHandler.CreateOrder)
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

		// 系统设置
		adminAPIGroup.GET("/settings", adminHandler.GetSettings)
		adminAPIGroup.PUT("/settings", adminHandler.UpdateSettings)
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
	router.POST("/api/order/create", paymentHandler.CreateOrder)
	router.GET("/payment/callback", paymentHandler.PaymentCallback)
	router.GET("/api/order/:order_no/status", paymentHandler.QueryOrder)
	router.POST("/api/order/:order_no/cancel", paymentHandler.CancelOrder)

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

	// 检查是否已初始化
	initialized := settingService.Get(services.SettingInitialized)
	if initialized == "true" {
		// 已初始化，但需要确保支付配置同步到数据库
		syncPaymentConfig(settingService)
		return
	}

	log.Println("首次运行，正在初始化系统...")

	// 保存基础设置
	err := settingService.SetMultiple(map[string]string{
		services.SettingSiteName:        "NodeLoc 社区发卡",
		services.SettingSiteDescription: "基于 NodeLoc OAuth 的社区发卡系统",
		services.SettingFooterText:      "© 2026 NodeLoc 社区发卡系统",
		services.SettingInitialized:     "true",
	})
	if err != nil {
		log.Fatalf("保存系统设置失败: %v", err)
	}

	// 同步支付配置
	syncPaymentConfig(settingService)

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
