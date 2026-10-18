package services

import (
	"math"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/nodeloc-faka/database"
	"github.com/nodeloc-faka/models"
)

// ProductService 商品服务
type ProductService struct{}

// NewProductService 创建商品服务
func NewProductService() *ProductService {
	return &ProductService{}
}

// Validate 校验商品字段，返回的 *ServiceError 可直接展示给用户。
// 价格必须是正整数：NodeLoc 能量支付按整数金额收款，小数价格会被截断，买家实际少付。
func (s *ProductService) Validate(p *models.Product) error {
	p.Name = strings.TrimSpace(p.Name)
	switch {
	case p.Name == "":
		return &ServiceError{Message: "商品名称不能为空"}
	case utf8.RuneCountInString(p.Name) > 200:
		return &ServiceError{Message: "商品名称不能超过 200 个字符"}
	case p.Price <= 0 || p.Price != math.Trunc(p.Price) || p.Price > 1e9:
		return &ServiceError{Message: "价格必须是正整数"}
	case p.PromoPrice < 0 || p.PromoPrice != math.Trunc(p.PromoPrice) || p.PromoPrice > 1e9:
		return &ServiceError{Message: "促销价必须是非负整数"}
	case p.OrigPrice < 0 || p.OrigPrice > 1e9:
		return &ServiceError{Message: "原价不能为负数"}
	case p.StockCount < 0:
		return &ServiceError{Message: "库存不能为负数"}
	case !models.ValidDeliveryType(p.DeliveryType):
		return &ServiceError{Message: "发货方式无效"}
	case p.DeliveryType == models.DeliveryTypeFixed && strings.TrimSpace(p.FixedContent) == "":
		return &ServiceError{Message: "固定内容发货需要填写发货内容"}
	case len(p.Image) > 500 || len(p.Images) > 2000:
		return &ServiceError{Message: "图片地址过长"}
	}
	return nil
}

// CategoryExists 分类是否存在
func (s *ProductService) CategoryExists(id uint) bool {
	var n int64
	database.GetDB().Model(&models.Category{}).Where("id = ?", id).Count(&n)
	return n > 0
}

// Create 创建商品
func (s *ProductService) Create(product *models.Product) error {
	return database.GetDB().Create(product).Error
}

// Update 更新商品
func (s *ProductService) Update(product *models.Product) error {
	return database.GetDB().Save(product).Error
}

// Delete 删除商品
func (s *ProductService) Delete(id uint) error {
	// 检查是否有未售出的卡密
	var count int64
	database.GetDB().Model(&models.CardKey{}).
		Where("product_id = ? AND status = ?", id, models.CardKeyStatusAvailable).
		Count(&count)
	if count > 0 {
		return ErrProductHasCards
	}
	return database.GetDB().Delete(&models.Product{}, id).Error
}

// FindByID 根据ID查找商品
func (s *ProductService) FindByID(id uint) (*models.Product, error) {
	var product models.Product
	if err := database.GetDB().Preload("Category").Preload("Shop").First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

// GetAll 获取所有商品
func (s *ProductService) GetAll() ([]models.Product, error) {
	var products []models.Product
	if err := database.GetDB().Preload("Category").Order("sort asc, id desc").Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

// GetActive 获取启用的商品（仅来自已批准的店铺）
func (s *ProductService) GetActive() ([]models.Product, error) {
	var products []models.Product
	if err := database.GetDB().
		Preload("Category").
		Preload("Shop").
		Joins("JOIN shops ON shops.id = products.shop_id AND shops.status = ?", models.ShopStatusApproved).
		Where("products.is_active = ?", true).
		Order("products.sort asc, products.id desc").
		Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

// GetByCategory 根据分类获取商品（仅已批准店铺）
func (s *ProductService) GetByCategory(categoryID uint) ([]models.Product, error) {
	var products []models.Product
	if err := database.GetDB().
		Preload("Category").
		Preload("Shop").
		Joins("JOIN shops ON shops.id = products.shop_id AND shops.status = ?", models.ShopStatusApproved).
		Where("products.category_id = ? AND products.is_active = ?", categoryID, true).
		Order("products.sort asc, products.id desc").
		Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

// GetByShop 根据店铺获取商品
func (s *ProductService) GetByShop(shopID uint, onlyActive bool) ([]models.Product, error) {
	var products []models.Product
	db := database.GetDB().Preload("Category").Preload("Shop").Where("shop_id = ?", shopID)
	if onlyActive {
		db = db.Where("is_active = ?", true)
	}
	if err := db.Order("sort asc, id desc").Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

// GetWithPagination 分页获取商品
func (s *ProductService) GetWithPagination(page, pageSize int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	db := database.GetDB().Model(&models.Product{})
	db.Count(&total)

	offset := (page - 1) * pageSize
	if err := db.Preload("Category").
		Order("sort asc, id desc").
		Offset(offset).
		Limit(pageSize).
		Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// UpdateStock 更新库存
func (s *ProductService) UpdateStock(id uint) error {
	var product models.Product
	if err := database.GetDB().Select("id", "delivery_type").First(&product, id).Error; err != nil {
		return err
	}
	// 手动/固定内容商品的库存由店主设定，不能被卡密数量覆盖
	if product.DeliveryType != models.DeliveryTypeCardKey {
		return nil
	}
	var count int64
	database.GetDB().Model(&models.CardKey{}).
		Where("product_id = ? AND status = ?", id, models.CardKeyStatusAvailable).
		Count(&count)

	return database.GetDB().Model(&models.Product{}).
		Where("id = ?", id).
		Update("stock_count", count).Error
}

// IncrementSales 增加销量
func (s *ProductService) IncrementSales(id uint, quantity int) error {
	return database.GetDB().Model(&models.Product{}).
		Where("id = ?", id).
		UpdateColumn("sales_count", database.GetDB().Raw("sales_count + ?", quantity)).
		Error
}

// SearchParams 商品搜索参数
type SearchParams struct {
	CategoryID uint
	Keyword    string
	MinPrice   float64
	MaxPrice   float64
	Sort       string // default | price_asc | price_desc | sales | newest
	PromoOnly  bool   // 仅返回促销进行中的商品
	Page       int
	PageSize   int
}

// SearchProducts 带过滤、排序、分页的商品搜索（仅已批准店铺）
func (s *ProductService) SearchProducts(p SearchParams) ([]models.Product, int64, error) {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PageSize < 1 || p.PageSize > 100 {
		p.PageSize = 20
	}

	db := database.GetDB().
		Model(&models.Product{}).
		Joins("JOIN shops ON shops.id = products.shop_id AND shops.status = ?", models.ShopStatusApproved).
		Where("products.is_active = ?", true)

	if p.CategoryID > 0 {
		db = db.Where("products.category_id = ?", p.CategoryID)
	}
	if p.Keyword != "" {
		db = db.Where("products.name LIKE ?", "%"+p.Keyword+"%")
	}
	if p.MinPrice > 0 {
		db = db.Where("products.price >= ?", p.MinPrice)
	}
	if p.MaxPrice > 0 {
		db = db.Where("products.price <= ?", p.MaxPrice)
	}
	if p.PromoOnly {
		now := time.Now()
		db = db.Where("products.promo_price > 0").
			Where("products.promo_start IS NULL OR products.promo_start <= ?", now).
			Where("products.promo_end IS NULL OR products.promo_end > ?", now)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	switch p.Sort {
	case "price_asc":
		db = db.Order("products.price asc")
	case "price_desc":
		db = db.Order("products.price desc")
	case "sales":
		db = db.Order("products.sales_count desc")
	case "newest":
		db = db.Order("products.id desc")
	default:
		db = db.Order("products.sort asc, products.id desc")
	}

	db = db.Offset((p.Page - 1) * p.PageSize).Limit(p.PageSize)

	var products []models.Product
	if err := db.Preload("Category").Preload("Shop").Find(&products).Error; err != nil {
		return nil, 0, err
	}
	return products, total, nil
}

// CountByShop 获取指定店铺的上架商品数量
func (s *ProductService) CountByShop(shopID uint, out *int64) {
	database.GetDB().Model(&models.Product{}).
		Where("shop_id = ? AND is_active = ?", shopID, true).
		Count(out)
}

// Count 获取商品数量
func (s *ProductService) Count() int64 {
	var count int64
	database.GetDB().Model(&models.Product{}).Count(&count)
	return count
}

// 错误定义
var ErrProductHasCards = &ServiceError{Message: "该商品下有未售出的卡密，无法删除"}
