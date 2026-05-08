package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/nodeloc-faka/database"
	"github.com/nodeloc-faka/models"
	"gorm.io/gorm"
)

// OrderService 订单服务
type OrderService struct{}

// NewOrderService 创建订单服务
func NewOrderService() *OrderService {
	return &OrderService{}
}

// Create 创建订单
func (s *OrderService) Create(order *models.Order) error {
	// 生成订单号
	order.OrderNo = s.generateOrderNo()
	return database.GetDB().Create(order).Error
}

// generateOrderNo 生成订单号
func (s *OrderService) generateOrderNo() string {
	return fmt.Sprintf("%s%d", time.Now().Format("20060102150405"), time.Now().UnixNano()%10000)
}

// Update 更新订单
func (s *OrderService) Update(order *models.Order) error {
	return database.GetDB().Save(order).Error
}

// FindByID 根据ID查找订单
func (s *OrderService) FindByID(id uint) (*models.Order, error) {
	var order models.Order
	if err := database.GetDB().
		Preload("User").
		Preload("Product").
		Preload("CardKeys").
		First(&order, id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

// FindByOrderNo 根据订单号查找订单
func (s *OrderService) FindByOrderNo(orderNo string) (*models.Order, error) {
	var order models.Order
	if err := database.GetDB().
		Preload("User").
		Preload("Product").
		Preload("CardKeys").
		Where("order_no = ?", orderNo).
		First(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

// GetByUser 获取用户的订单
func (s *OrderService) GetByUser(userID uint) ([]models.Order, error) {
	var orders []models.Order
	if err := database.GetDB().
		Preload("Product").
		Preload("CardKeys").
		Where("user_id = ?", userID).
		Order("id desc").
		Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

// GetByShop 分页获取店铺的订单
func (s *OrderService) GetByShop(shopID uint, page, pageSize, status int) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64

	db := database.GetDB().Model(&models.Order{}).Where("shop_id = ?", shopID)
	if status >= 0 {
		db = db.Where("status = ?", status)
	}
	db.Count(&total)

	offset := (page - 1) * pageSize
	if err := db.Preload("Product").Preload("User").
		Order("id desc").
		Offset(offset).Limit(pageSize).
		Find(&orders).Error; err != nil {
		return nil, 0, err
	}
	return orders, total, nil
}

// GetAll 获取所有订单
func (s *OrderService) GetAll() ([]models.Order, error) {
	var orders []models.Order
	if err := database.GetDB().
		Preload("User").
		Preload("Product").
		Order("id desc").
		Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

// GetWithPagination 分页获取订单
func (s *OrderService) GetWithPagination(page, pageSize int, status int) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64

	db := database.GetDB().Model(&models.Order{})
	if status >= 0 {
		db = db.Where("status = ?", status)
	}
	db.Count(&total)

	offset := (page - 1) * pageSize
	if err := db.Preload("User").
		Preload("Product").
		Order("id desc").
		Offset(offset).
		Limit(pageSize).
		Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

// MarkAsPaid 标记订单为已支付
func (s *OrderService) MarkAsPaid(id uint) error {
	now := time.Now()
	return database.GetDB().Model(&models.Order{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":  models.OrderStatusPaid,
			"paid_at": now,
		}).Error
}

// Complete 完成订单
func (s *OrderService) Complete(id uint) error {
	return database.GetDB().Model(&models.Order{}).
		Where("id = ?", id).
		Update("status", models.OrderStatusCompleted).Error
}

// Cancel 取消订单
func (s *OrderService) Cancel(id uint) error {
	return database.GetDB().Model(&models.Order{}).
		Where("id = ?", id).
		Update("status", models.OrderStatusCancelled).Error
}

// Count 获取订单数量
func (s *OrderService) Count() int64 {
	var count int64
	database.GetDB().Model(&models.Order{}).Count(&count)
	return count
}

// CountByStatus 根据状态获取订单数量
func (s *OrderService) CountByStatus(status int) int64 {
	var count int64
	database.GetDB().Model(&models.Order{}).Where("status = ?", status).Count(&count)
	return count
}

// GetTotalSales 获取总销售额
func (s *OrderService) GetTotalSales() float64 {
	var total float64
	database.GetDB().Model(&models.Order{}).
		Where("status IN ?", []int{models.OrderStatusPaid, models.OrderStatusCompleted}).
		Select("COALESCE(SUM(total_amount), 0)").
		Scan(&total)
	return total
}

// GetTodaySales 获取今日销售额
func (s *OrderService) GetTodaySales() float64 {
	var total float64
	today := time.Now().Format("2006-01-02")
	database.GetDB().Model(&models.Order{}).
		Where("status IN ? AND DATE(created_at) = ?", []int{models.OrderStatusPaid, models.OrderStatusCompleted}, today).
		Select("COALESCE(SUM(total_amount), 0)").
		Scan(&total)
	return total
}

// CreatePendingOrder 创建待支付订单
func (s *OrderService) CreatePendingOrder(userID, productID uint, quantity int, contact, remark string) (*models.Order, error) {
	productService := NewProductService()
	product, err := productService.FindByID(productID)
	if err != nil {
		return nil, ErrProductNotFound
	}

	// 检查库存
	if product.DeliveryType == models.DeliveryTypeManual {
		if product.StockCount < quantity {
			return nil, ErrInsufficientStock
		}
	} else {
		cardKeyService := NewCardKeyService()
		count := cardKeyService.CountByProduct(productID, models.CardKeyStatusAvailable)
		if int(count) < quantity {
			return nil, ErrInsufficientStock
		}
	}

	expiredAt := time.Now().Add(30 * time.Minute)
	order := &models.Order{
		UserID:      userID,
		ProductID:   productID,
		ShopID:      product.ShopID,
		Quantity:    quantity,
		TotalAmount: product.GetEffectivePrice() * float64(quantity),
		Status:      models.OrderStatusPending,
		PayMethod:   models.PayMethodNodeLoc,
		Contact:     contact,
		Remark:      remark,
		ExpiredAt:   &expiredAt,
	}

	if err := s.Create(order); err != nil {
		return nil, err
	}

	return order, nil
}

// SetPaymentInfo 设置支付信息
func (s *OrderService) SetPaymentInfo(orderID uint, transactionID, paymentURL string) error {
	return database.GetDB().Model(&models.Order{}).
		Where("id = ?", orderID).
		Updates(map[string]interface{}{
			"transaction_id": transactionID,
			"payment_url":    paymentURL,
		}).Error
}

// FindByTransactionID 根据交易ID查找订单
func (s *OrderService) FindByTransactionID(transactionID string) (*models.Order, error) {
	var order models.Order
	if err := database.GetDB().
		Preload("User").
		Preload("Product").
		Preload("CardKeys").
		Where("transaction_id = ?", transactionID).
		First(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

// ProcessPaymentCallback 处理支付回调（NodeLoc 支付完成后调用）
func (s *OrderService) ProcessPaymentCallback(transactionID string, amount, platformFee, merchantPoints int) (*models.Order, error) {
	order, err := s.FindByTransactionID(transactionID)
	if err != nil {
		return nil, ErrOrderNotFound
	}

	// 幂等：订单已处理则直接返回
	if order.Status != models.OrderStatusPending {
		return order, nil
	}

	expectedAmount := int(order.TotalAmount)
	if amount != expectedAmount {
		return nil, ErrAmountMismatch
	}

	now := time.Now()
	order.PaidAt = &now
	order.PlatformFee = platformFee
	order.MerchantPoints = merchantPoints
	if order.PayMethod == "" {
		order.PayMethod = models.PayMethodNodeLoc
	}

	// 使用数据库乐观锁：只更新 pending 状态的订单，防止并发重复处理
	updated := database.GetDB().Model(&models.Order{}).
		Where("id = ? AND status = ?", order.ID, models.OrderStatusPending).
		Updates(map[string]interface{}{
			"paid_at":         now,
			"platform_fee":    platformFee,
			"merchant_points": merchantPoints,
			"pay_method":      order.PayMethod,
		})
	if updated.Error != nil {
		return nil, updated.Error
	}
	if updated.RowsAffected == 0 {
		// 已被其他请求处理（幂等）
		return s.FindByID(order.ID)
	}

	// 手动发货：仅标记已支付，等待卖家发货
	if order.Product != nil && order.Product.DeliveryType == models.DeliveryTypeManual {
		if err := database.GetDB().Transaction(func(tx *gorm.DB) error {
			if err := tx.Model(&models.Order{}).Where("id = ?", order.ID).
				Update("status", models.OrderStatusPaid).Error; err != nil {
				return err
			}
			// 扣减手动库存
			return tx.Model(&models.Product{}).Where("id = ?", order.ProductID).
				UpdateColumn("stock_count", gorm.Expr("stock_count - ?", order.Quantity)).Error
		}); err != nil {
			return nil, err
		}
		return s.FindByID(order.ID)
	}

	// 卡密自动发货：获取卡密并完成订单
	cardKeyService := NewCardKeyService()
	availableCards, err := cardKeyService.GetAvailableByProduct(order.ProductID, order.Quantity)
	if err != nil {
		return nil, err
	}
	if len(availableCards) < order.Quantity {
		return nil, ErrInsufficientStock
	}

	order.Status = models.OrderStatusCompleted
	if err := s.completeOrderTx(order, availableCards); err != nil {
		return nil, err
	}

	return s.FindByID(order.ID)
}

// CreateAndProcess 创建并处理订单（免费模式）
func (s *OrderService) CreateAndProcess(userID, productID uint, quantity int, contact, remark string) (*models.Order, error) {
	productService := NewProductService()
	product, err := productService.FindByID(productID)
	if err != nil {
		return nil, ErrProductNotFound
	}

	now := time.Now()
	order := &models.Order{
		UserID:      userID,
		ProductID:   productID,
		ShopID:      product.ShopID,
		Quantity:    quantity,
		TotalAmount: product.GetEffectivePrice() * float64(quantity),
		PayMethod:   models.PayMethodFree,
		Contact:     contact,
		Remark:      remark,
		PaidAt:      &now,
	}

	// 手动发货：付款后待发货
	if product.DeliveryType == models.DeliveryTypeManual {
		if product.StockCount < quantity {
			return nil, ErrInsufficientStock
		}
		order.Status = models.OrderStatusPaid
		if err := s.Create(order); err != nil {
			return nil, err
		}
		database.GetDB().Model(&models.Product{}).Where("id = ?", productID).
			UpdateColumn("stock_count", gorm.Expr("stock_count - ?", quantity))
		return s.FindByID(order.ID)
	}

	// 卡密自动发货
	cardKeyService := NewCardKeyService()
	availableCards, err := cardKeyService.GetAvailableByProduct(productID, quantity)
	if err != nil {
		return nil, err
	}
	if len(availableCards) < quantity {
		return nil, ErrInsufficientStock
	}
	order.Status = models.OrderStatusCompleted
	if err := s.Create(order); err != nil {
		return nil, err
	}
	if err := s.completeOrderTx(order, availableCards); err != nil {
		return nil, err
	}
	return s.FindByID(order.ID)
}

// CreateAndProcessByBalance 使用买家余额支付下单
func (s *OrderService) CreateAndProcessByBalance(userID, productID uint, quantity int, contact, remark string) (*models.Order, error) {
	productService := NewProductService()
	product, err := productService.FindByID(productID)
	if err != nil {
		return nil, ErrProductNotFound
	}

	totalAmount := product.GetEffectivePrice() * float64(quantity)
	balanceService := NewBalanceService()
	now := time.Now()

	order := &models.Order{
		UserID:      userID,
		ProductID:   productID,
		ShopID:      product.ShopID,
		Quantity:    quantity,
		TotalAmount: totalAmount,
		Status:      models.OrderStatusPending,
		PayMethod:   models.PayMethodBalance,
		Contact:     contact,
		Remark:      remark,
	}
	if err := s.Create(order); err != nil {
		return nil, err
	}

	// 手动发货：扣余额 + 标记已支付，不发卡不结算
	if product.DeliveryType == models.DeliveryTypeManual {
		if product.StockCount < quantity {
			s.Cancel(order.ID)
			return nil, ErrInsufficientStock
		}
		if err := database.GetDB().Transaction(func(tx *gorm.DB) error {
			if err := balanceService.DeductBalance(tx, userID, totalAmount,
				models.BalanceTxPurchase,
				fmt.Sprintf("余额支付订单 %s", order.OrderNo),
				"order", order.ID,
			); err != nil {
				return err
			}
			order.Status = models.OrderStatusPaid
			order.PaidAt = &now
			if err := tx.Save(order).Error; err != nil {
				return err
			}
			return tx.Model(&models.Product{}).Where("id = ?", productID).
				UpdateColumn("stock_count", gorm.Expr("stock_count - ?", quantity)).Error
		}); err != nil {
			s.Cancel(order.ID)
			return nil, err
		}
		return s.FindByID(order.ID)
	}

	// 卡密自动发货：扣余额 + 发卡 + 结算
	cardKeyService := NewCardKeyService()
	availableCards, err := cardKeyService.GetAvailableByProduct(productID, quantity)
	if err != nil || len(availableCards) < quantity {
		s.Cancel(order.ID)
		return nil, ErrInsufficientStock
	}

	if err := database.GetDB().Transaction(func(tx *gorm.DB) error {
		if err := balanceService.DeductBalance(tx, userID, totalAmount,
			models.BalanceTxPurchase,
			fmt.Sprintf("余额支付订单 %s", order.OrderNo),
			"order", order.ID,
		); err != nil {
			return err
		}
		order.Status = models.OrderStatusCompleted
		order.PaidAt = &now
		return s.completeOrderTxWithDB(tx, order, availableCards)
	}); err != nil {
		s.Cancel(order.ID)
		return nil, err
	}

	return s.FindByID(order.ID)
}

// ShipOrder 卖家发货（手动发货商品）
func (s *OrderService) ShipOrder(orderID uint, shopUserID uint, content, note string) error {
	order, err := s.FindByID(orderID)
	if err != nil {
		return ErrOrderNotFound
	}
	if order.Status != models.OrderStatusPaid {
		return &ServiceError{Message: "订单状态不允许发货"}
	}
	// 验证归属
	var shop models.Shop
	if err := database.GetDB().First(&shop, order.ShopID).Error; err != nil || shop.UserID != shopUserID {
		return &ServiceError{Message: "无权操作"}
	}
	if err := database.GetDB().Model(&models.Order{}).Where("id = ?", orderID).Updates(map[string]interface{}{
		"status":          models.OrderStatusShipped,
		"deliver_content": content,
		"deliver_note":    note,
	}).Error; err != nil {
		return err
	}

	// 如果订单关联了卡密（卡密商品走了手动发货流程），发货时同步标记卡密为已售出
	cardKeyService := NewCardKeyService()
	var cardKeys []models.CardKey
	database.GetDB().Where("product_id = ? AND status = ? AND order_id IS NULL", order.ProductID, models.CardKeyStatusAvailable).
		Limit(order.Quantity).Find(&cardKeys)
	if len(cardKeys) > 0 {
		ids := make([]uint, len(cardKeys))
		for i, ck := range cardKeys {
			ids[i] = ck.ID
		}
		cardKeyService.MarkAsSold(ids, orderID)
	}

	// 发货通知邮件
	order.DeliverContent = content
	order.DeliverNote = note
	NewEmailService().SendOrderShipped(order)
	return nil
}

// ConfirmReceipt 买家确认收货，触发结算
func (s *OrderService) ConfirmReceipt(orderNo string, userID uint) (*models.Order, error) {
	order, err := s.FindByOrderNo(orderNo)
	if err != nil {
		return nil, ErrOrderNotFound
	}
	if order.UserID != userID {
		return nil, &ServiceError{Message: "无权操作"}
	}
	if order.Status != models.OrderStatusShipped {
		return nil, &ServiceError{Message: "订单状态不允许确认收货"}
	}

	balanceService := NewBalanceService()
	if err := database.GetDB().Transaction(func(tx *gorm.DB) error {
		// 标记完成
		if err := tx.Model(&models.Order{}).Where("id = ?", order.ID).Updates(map[string]interface{}{
			"status": models.OrderStatusCompleted,
		}).Error; err != nil {
			return err
		}
		// 结算给卖家
		if order.ShopID > 0 && !order.ShopSettled {
			var shop models.Shop
			if err := tx.First(&shop, order.ShopID).Error; err != nil {
				return err
			}
			if !shop.IsOfficial && shop.UserID > 0 {
				income := order.TotalAmount
				if err := balanceService.AddBalance(tx, shop.UserID, income,
					models.BalanceTxSaleIncome,
					fmt.Sprintf("订单 %s 销售收入（确认收货）", order.OrderNo),
					"order", order.ID,
				); err != nil {
					return err
				}
				return tx.Model(&models.Order{}).Where("id = ?", order.ID).Updates(map[string]interface{}{
					"shop_settled": true,
					"shop_income":  income,
				}).Error
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return s.FindByOrderNo(orderNo)
}

// completeOrderTx 完成订单：开启事务 → 发卡 + 库存 + 销量 + 店主余额结算
func (s *OrderService) completeOrderTx(order *models.Order, availableCards []models.CardKey) error {
	return database.GetDB().Transaction(func(tx *gorm.DB) error {
		return s.completeOrderTxWithDB(tx, order, availableCards)
	})
}

// completeOrderTxWithDB 在指定事务中完成订单
func (s *OrderService) completeOrderTxWithDB(tx *gorm.DB, order *models.Order, availableCards []models.CardKey) error {
	cardKeyService := NewCardKeyService()
	balanceService := NewBalanceService()
	emailService := NewEmailService()

	// 保存订单
	if err := tx.Save(order).Error; err != nil {
		return err
	}

	// 分配卡密
	cardIDs := make([]uint, len(availableCards))
	for i, card := range availableCards {
		cardIDs[i] = card.ID
	}
	if err := cardKeyService.MarkAsSoldTx(tx, cardIDs, order.ID); err != nil {
		return err
	}

	// 更新商品库存和销量（直接 SQL，无需重新查询）
	if err := tx.Model(&models.Product{}).Where("id = ?", order.ProductID).
		UpdateColumn("sales_count", gorm.Expr("sales_count + ?", order.Quantity)).Error; err != nil {
		return err
	}
	// 重新计算库存（必须在同一事务内，否则新连接会被当前事务的行锁阻塞 50s）
	var availableCount int64
	tx.Model(&models.CardKey{}).
		Where("product_id = ? AND status = ?", order.ProductID, models.CardKeyStatusAvailable).
		Count(&availableCount)
	if err := tx.Model(&models.Product{}).Where("id = ?", order.ProductID).
		Update("stock_count", availableCount).Error; err != nil {
		return err
	}

	// 店主结算：ShopID > 0 且未结算
	if order.ShopID > 0 && !order.ShopSettled {
		var shop models.Shop
		err := tx.First(&shop, order.ShopID).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err // 真正的数据库错误才回滚
		}
		// 非官方店且有店主用户 → 结算余额
		if err == nil && !shop.IsOfficial && shop.UserID > 0 {
			income := order.TotalAmount
			if err := balanceService.AddBalance(tx, shop.UserID, income,
				models.BalanceTxSaleIncome,
				fmt.Sprintf("订单 %s 销售收入", order.OrderNo),
				"order", order.ID,
			); err != nil {
				return err
			}
			if err := tx.Model(&models.Order{}).Where("id = ?", order.ID).Updates(map[string]interface{}{
				"shop_settled": true,
				"shop_income":  income,
			}).Error; err != nil {
				return err
			}

			// 通知店主有新订单（异步，不影响事务）
			var shopOwner models.User
			if err := tx.First(&shopOwner, shop.UserID).Error; err == nil && shopOwner.Email != "" {
				go emailService.SendNewOrderToShop(order, shopOwner.Email)
			}
		}
	}

	// 通知买家订单完成（卡密已发放）
	if order.User != nil && order.User.Email != "" {
		// 重新加载卡密（事务内已分配）
		var fullOrder models.Order
		if err := tx.Preload("User").Preload("Product").Preload("CardKeys").First(&fullOrder, order.ID).Error; err == nil {
			go emailService.SendOrderCompleted(&fullOrder)
		}
	}

	return nil
}

// CancelExpiredOrders 取消过期订单
func (s *OrderService) CancelExpiredOrders() (int64, error) {
	result := database.GetDB().Model(&models.Order{}).
		Where("status = ? AND expired_at < ?", models.OrderStatusPending, time.Now()).
		Update("status", models.OrderStatusCancelled)
	return result.RowsAffected, result.Error
}

// 错误定义
var (
	ErrProductNotFound   = &ServiceError{Message: "商品不存在"}
	ErrInsufficientStock = &ServiceError{Message: "库存不足"}
	ErrOrderNotFound     = &ServiceError{Message: "订单不存在"}
	ErrAmountMismatch    = &ServiceError{Message: "支付金额不匹配"}
	ErrOrderExpired      = &ServiceError{Message: "订单已过期"}
)

// ShopDashboardStats 店铺仪表盘统计
type ShopDashboardStats struct {
	TotalOrders     int64   `json:"total_orders"`
	CompletedOrders int64   `json:"completed_orders"`
	PendingOrders   int64   `json:"pending_orders"`
	CancelledOrders int64   `json:"cancelled_orders"`
	TotalIncome     float64 `json:"total_income"`
	TodayOrders     int64   `json:"today_orders"`
	TodayIncome     float64 `json:"today_income"`
	WeekOrders      int64   `json:"week_orders"`
	WeekIncome      float64 `json:"week_income"`
	MonthOrders     int64   `json:"month_orders"`
	MonthIncome     float64 `json:"month_income"`
	RecentOrders    []models.Order `json:"recent_orders"`
}

// GetShopDashboard 获取店铺仪表盘统计数据
func (s *OrderService) GetShopDashboard(shopID uint) (*ShopDashboardStats, error) {
	db := database.GetDB()
	stats := &ShopDashboardStats{}

	// 总订单数
	db.Model(&models.Order{}).Where("shop_id = ?", shopID).Count(&stats.TotalOrders)

	// 各状态订单数
	db.Model(&models.Order{}).Where("shop_id = ? AND status IN ?", shopID,
		[]int{models.OrderStatusCompleted, models.OrderStatusPaid, models.OrderStatusShipped}).
		Count(&stats.CompletedOrders)
	db.Model(&models.Order{}).Where("shop_id = ? AND status = ?", shopID, models.OrderStatusPending).
		Count(&stats.PendingOrders)
	db.Model(&models.Order{}).Where("shop_id = ? AND status = ?", shopID, models.OrderStatusCancelled).
		Count(&stats.CancelledOrders)

	// 总收入（已完成+已支付+已发货的订单）
	db.Model(&models.Order{}).
		Where("shop_id = ? AND status IN ?", shopID,
			[]int{models.OrderStatusCompleted, models.OrderStatusPaid, models.OrderStatusShipped}).
		Select("COALESCE(SUM(total_amount), 0)").Scan(&stats.TotalIncome)

	// 今日
	today := time.Now().Format("2006-01-02")
	db.Model(&models.Order{}).
		Where("shop_id = ? AND status IN ? AND DATE(created_at) = ?", shopID,
			[]int{models.OrderStatusCompleted, models.OrderStatusPaid, models.OrderStatusShipped}, today).
		Count(&stats.TodayOrders)
	db.Model(&models.Order{}).
		Where("shop_id = ? AND status IN ? AND DATE(created_at) = ?", shopID,
			[]int{models.OrderStatusCompleted, models.OrderStatusPaid, models.OrderStatusShipped}, today).
		Select("COALESCE(SUM(total_amount), 0)").Scan(&stats.TodayIncome)

	// 本周
	weekStart := time.Now().AddDate(0, 0, -int(time.Now().Weekday()))
	weekStartStr := weekStart.Format("2006-01-02")
	db.Model(&models.Order{}).
		Where("shop_id = ? AND status IN ? AND DATE(created_at) >= ?", shopID,
			[]int{models.OrderStatusCompleted, models.OrderStatusPaid, models.OrderStatusShipped}, weekStartStr).
		Count(&stats.WeekOrders)
	db.Model(&models.Order{}).
		Where("shop_id = ? AND status IN ? AND DATE(created_at) >= ?", shopID,
			[]int{models.OrderStatusCompleted, models.OrderStatusPaid, models.OrderStatusShipped}, weekStartStr).
		Select("COALESCE(SUM(total_amount), 0)").Scan(&stats.WeekIncome)

	// 本月
	monthStart := time.Now().Format("2006-01") + "-01"
	db.Model(&models.Order{}).
		Where("shop_id = ? AND status IN ? AND DATE(created_at) >= ?", shopID,
			[]int{models.OrderStatusCompleted, models.OrderStatusPaid, models.OrderStatusShipped}, monthStart).
		Count(&stats.MonthOrders)
	db.Model(&models.Order{}).
		Where("shop_id = ? AND status IN ? AND DATE(created_at) >= ?", shopID,
			[]int{models.OrderStatusCompleted, models.OrderStatusPaid, models.OrderStatusShipped}, monthStart).
		Select("COALESCE(SUM(total_amount), 0)").Scan(&stats.MonthIncome)

	// 最近 5 条订单
	db.Preload("Product").Preload("User").
		Where("shop_id = ?", shopID).
		Order("id desc").Limit(5).
		Find(&stats.RecentOrders)

	return stats, nil
}
