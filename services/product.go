package services

import (
	"github.com/nodeloc-faka/database"
	"github.com/nodeloc-faka/models"
)

// ProductService 商品服务
type ProductService struct{}

// NewProductService 创建商品服务
func NewProductService() *ProductService {
	return &ProductService{}
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
	if err := database.GetDB().Preload("Category").First(&product, id).Error; err != nil {
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

// GetActive 获取启用的商品
func (s *ProductService) GetActive() ([]models.Product, error) {
	var products []models.Product
	if err := database.GetDB().Preload("Category").
		Where("is_active = ?", true).
		Order("sort asc, id desc").
		Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

// GetByCategory 根据分类获取商品
func (s *ProductService) GetByCategory(categoryID uint) ([]models.Product, error) {
	var products []models.Product
	if err := database.GetDB().Preload("Category").
		Where("category_id = ? AND is_active = ?", categoryID, true).
		Order("sort asc, id desc").
		Find(&products).Error; err != nil {
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

// Count 获取商品数量
func (s *ProductService) Count() int64 {
	var count int64
	database.GetDB().Model(&models.Product{}).Count(&count)
	return count
}

// 错误定义
var ErrProductHasCards = &ServiceError{Message: "该商品下有未售出的卡密，无法删除"}
