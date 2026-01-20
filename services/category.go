package services

import (
	"github.com/nodeloc-faka/database"
	"github.com/nodeloc-faka/models"
)

// CategoryService 分类服务
type CategoryService struct{}

// NewCategoryService 创建分类服务
func NewCategoryService() *CategoryService {
	return &CategoryService{}
}

// Create 创建分类
func (s *CategoryService) Create(category *models.Category) error {
	return database.GetDB().Create(category).Error
}

// Update 更新分类
func (s *CategoryService) Update(category *models.Category) error {
	return database.GetDB().Save(category).Error
}

// Delete 删除分类
func (s *CategoryService) Delete(id uint) error {
	// 检查是否有关联商品
	var count int64
	database.GetDB().Model(&models.Product{}).Where("category_id = ?", id).Count(&count)
	if count > 0 {
		return ErrCategoryHasProducts
	}
	return database.GetDB().Delete(&models.Category{}, id).Error
}

// FindByID 根据ID查找分类
func (s *CategoryService) FindByID(id uint) (*models.Category, error) {
	var category models.Category
	if err := database.GetDB().First(&category, id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

// GetAll 获取所有分类
func (s *CategoryService) GetAll() ([]models.Category, error) {
	var categories []models.Category
	if err := database.GetDB().Order("sort asc, id asc").Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// GetActive 获取启用的分类
func (s *CategoryService) GetActive() ([]models.Category, error) {
	var categories []models.Category
	if err := database.GetDB().Where("is_active = ?", true).Order("sort asc, id asc").Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// GetWithProducts 获取分类及其商品
func (s *CategoryService) GetWithProducts() ([]models.Category, error) {
	var categories []models.Category
	if err := database.GetDB().Preload("Products", "is_active = ?", true).
		Where("is_active = ?", true).
		Order("sort asc, id asc").
		Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// Count 获取分类数量
func (s *CategoryService) Count() int64 {
	var count int64
	database.GetDB().Model(&models.Category{}).Count(&count)
	return count
}

// 错误定义
var ErrCategoryHasProducts = &ServiceError{Message: "该分类下有商品，无法删除"}
