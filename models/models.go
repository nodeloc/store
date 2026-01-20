package models

import (
	"time"

	"gorm.io/gorm"
)

// Setting 系统设置
type Setting struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Key       string    `gorm:"column:key;type:varchar(100);uniqueIndex" json:"key"`
	Value     string    `gorm:"type:text" json:"value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName 指定表名
func (Setting) TableName() string {
	return "settings"
}

// BeforeCreate GORM hook - 在创建前执行
func (s *Setting) BeforeCreate(tx *gorm.DB) error {
	// GORM 会自动处理列名，这里不需要特殊处理
	return nil
}

// Admin 管理员
type Admin struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	Username     string     `gorm:"uniqueIndex;size:50" json:"username"`
	PasswordHash string     `gorm:"size:255" json:"-"`
	LastLoginAt  *time.Time `json:"last_login_at"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// User 用户（通过OAuth登录）
type User struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	NodeLocID   int        `gorm:"uniqueIndex" json:"nodeloc_id"`
	Username    string     `gorm:"size:100" json:"username"`
	Name        string     `gorm:"size:100" json:"name"`
	Email       string     `gorm:"size:200" json:"email"`
	AvatarURL   string     `gorm:"size:500" json:"avatar_url"`
	TrustLevel  int        `json:"trust_level"`
	Balance     float64    `gorm:"default:0" json:"balance"`
	IsAdmin     bool       `gorm:"default:false;index" json:"is_admin"`
	IsBlocked   bool       `gorm:"default:false" json:"is_blocked"`
	LastLoginAt *time.Time `json:"last_login_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Orders      []Order    `gorm:"foreignKey:UserID" json:"orders,omitempty"`
}

// Category 商品分类
type Category struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:100" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Icon        string    `gorm:"size:50" json:"icon"`
	Sort        int       `gorm:"default:0" json:"sort"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Products    []Product `gorm:"foreignKey:CategoryID" json:"products,omitempty"`
}

// Product 商品
type Product struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	CategoryID  uint      `gorm:"index" json:"category_id"`
	Category    *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Name        string    `gorm:"size:200" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Price       float64   `json:"price"`
	OrigPrice   float64   `json:"orig_price"`
	Image       string    `gorm:"size:500" json:"image"`
	StockCount  int       `gorm:"default:0" json:"stock_count"`
	SalesCount  int       `gorm:"default:0" json:"sales_count"`
	Sort        int       `gorm:"default:0" json:"sort"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CardKeys    []CardKey `gorm:"foreignKey:ProductID" json:"card_keys,omitempty"`
}

// CardKey 卡密
type CardKey struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	ProductID uint       `gorm:"index" json:"product_id"`
	Product   *Product   `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	CardNo    string     `gorm:"size:500" json:"card_no"`
	CardPwd   string     `gorm:"size:500" json:"card_pwd"`
	Status    int        `gorm:"default:0" json:"status"` // 0: 未售出, 1: 已售出, 2: 已锁定
	OrderID   *uint      `gorm:"index" json:"order_id"`
	Order     *Order     `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	SoldAt    *time.Time `json:"sold_at"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// CardKeyStatus 卡密状态
const (
	CardKeyStatusAvailable = 0 // 可售
	CardKeyStatusSold      = 1 // 已售出
	CardKeyStatusLocked    = 2 // 已锁定
)

// Order 订单
type Order struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	OrderNo   string    `gorm:"uniqueIndex;size:50" json:"order_no"`
	UserID    uint      `gorm:"index" json:"user_id"`
	User      *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ProductID uint      `gorm:"index" json:"product_id"`
	Product   *Product  `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Quantity  int       `json:"quantity"`
	TotalAmount float64 `json:"total_amount"`
	Status    int       `gorm:"default:0" json:"status"` // 0: 待支付, 1: 已支付, 2: 已完成, 3: 已取消
	PayMethod string    `gorm:"size:50" json:"pay_method"`
	PaidAt    *time.Time `json:"paid_at"`
	Contact   string    `gorm:"size:200" json:"contact"`
	Remark    string    `gorm:"type:text" json:"remark"`
	
	// NodeLoc Payment 支付字段
	TransactionID  string     `gorm:"size:100;index" json:"transaction_id"`  // 支付交易ID
	PaymentURL     string     `gorm:"size:500" json:"payment_url"`           // 支付链接
	PlatformFee    int        `gorm:"default:0" json:"platform_fee"`         // 平台手续费
	MerchantPoints int        `gorm:"default:0" json:"merchant_points"`      // 商家实收积分
	ExpiredAt      *time.Time `json:"expired_at"`                            // 订单过期时间
	
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	CardKeys  []CardKey  `gorm:"foreignKey:OrderID" json:"card_keys,omitempty"`
}

// OrderStatus 订单状态
const (
	OrderStatusPending   = 0 // 待支付
	OrderStatusPaid      = 1 // 已支付
	OrderStatusCompleted = 2 // 已完成
	OrderStatusCancelled = 3 // 已取消
)

// AutoMigrate 自动迁移数据库
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&Setting{},
		&Admin{},
		&User{},
		&Category{},
		&Product{},
		&CardKey{},
		&Order{},
	)
}
