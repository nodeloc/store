package services

import (
	"strings"
	"time"

	"github.com/nodeloc-faka/database"
	"github.com/nodeloc-faka/models"
)

// CardKeyService 卡密服务
type CardKeyService struct{}

// NewCardKeyService 创建卡密服务
func NewCardKeyService() *CardKeyService {
	return &CardKeyService{}
}

// Create 创建单个卡密
func (s *CardKeyService) Create(cardKey *models.CardKey) error {
	err := database.GetDB().Create(cardKey).Error
	if err == nil {
		// 更新商品库存
		NewProductService().UpdateStock(cardKey.ProductID)
	}
	return err
}

// BatchCreate 批量创建卡密
func (s *CardKeyService) BatchCreate(productID uint, cardsText string) (int, error) {
	lines := strings.Split(cardsText, "\n")
	count := 0

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// 支持两种格式：
		// 1. 卡号----密码
		// 2. 只有卡号
		var cardNo, cardPwd string
		if strings.Contains(line, "----") {
			parts := strings.SplitN(line, "----", 2)
			cardNo = strings.TrimSpace(parts[0])
			if len(parts) > 1 {
				cardPwd = strings.TrimSpace(parts[1])
			}
		} else {
			cardNo = line
		}

		cardKey := &models.CardKey{
			ProductID: productID,
			CardNo:    cardNo,
			CardPwd:   cardPwd,
			Status:    models.CardKeyStatusAvailable,
		}

		if err := database.GetDB().Create(cardKey).Error; err != nil {
			continue
		}
		count++
	}

	// 更新商品库存
	NewProductService().UpdateStock(productID)

	return count, nil
}

// Update 更新卡密
func (s *CardKeyService) Update(cardKey *models.CardKey) error {
	return database.GetDB().Save(cardKey).Error
}

// Delete 删除卡密
func (s *CardKeyService) Delete(id uint) error {
	var cardKey models.CardKey
	if err := database.GetDB().First(&cardKey, id).Error; err != nil {
		return err
	}

	if cardKey.Status == models.CardKeyStatusSold {
		return ErrCardKeySold
	}

	err := database.GetDB().Delete(&models.CardKey{}, id).Error
	if err == nil {
		// 更新商品库存
		NewProductService().UpdateStock(cardKey.ProductID)
	}
	return err
}

// BatchDelete 批量删除未售出的卡密
func (s *CardKeyService) BatchDelete(ids []uint) (int, error) {
	result := database.GetDB().
		Where("id IN ? AND status = ?", ids, models.CardKeyStatusAvailable).
		Delete(&models.CardKey{})
	return int(result.RowsAffected), result.Error
}

// FindByID 根据ID查找卡密
func (s *CardKeyService) FindByID(id uint) (*models.CardKey, error) {
	var cardKey models.CardKey
	if err := database.GetDB().Preload("Product").First(&cardKey, id).Error; err != nil {
		return nil, err
	}
	return &cardKey, nil
}

// GetByProduct 根据商品获取卡密
func (s *CardKeyService) GetByProduct(productID uint) ([]models.CardKey, error) {
	var cardKeys []models.CardKey
	if err := database.GetDB().
		Where("product_id = ?", productID).
		Order("status asc, id desc").
		Find(&cardKeys).Error; err != nil {
		return nil, err
	}
	return cardKeys, nil
}

// GetAvailableByProduct 获取商品的可售卡密
func (s *CardKeyService) GetAvailableByProduct(productID uint, limit int) ([]models.CardKey, error) {
	var cardKeys []models.CardKey
	if err := database.GetDB().
		Where("product_id = ? AND status = ?", productID, models.CardKeyStatusAvailable).
		Order("id asc").
		Limit(limit).
		Find(&cardKeys).Error; err != nil {
		return nil, err
	}
	return cardKeys, nil
}

// MarkAsSold 标记卡密为已售出
func (s *CardKeyService) MarkAsSold(ids []uint, orderID uint) error {
	now := time.Now()
	return database.GetDB().Model(&models.CardKey{}).
		Where("id IN ?", ids).
		Updates(map[string]interface{}{
			"status":   models.CardKeyStatusSold,
			"order_id": orderID,
			"sold_at":  now,
		}).Error
}

// GetWithPagination 分页获取卡密
func (s *CardKeyService) GetWithPagination(productID uint, page, pageSize int) ([]models.CardKey, int64, error) {
	var cardKeys []models.CardKey
	var total int64

	db := database.GetDB().Model(&models.CardKey{})
	if productID > 0 {
		db = db.Where("product_id = ?", productID)
	}
	db.Count(&total)

	offset := (page - 1) * pageSize
	if err := db.Preload("Product").
		Order("status asc, id desc").
		Offset(offset).
		Limit(pageSize).
		Find(&cardKeys).Error; err != nil {
		return nil, 0, err
	}

	return cardKeys, total, nil
}

// CountByProduct 获取商品的卡密数量
func (s *CardKeyService) CountByProduct(productID uint, status int) int64 {
	var count int64
	db := database.GetDB().Model(&models.CardKey{}).Where("product_id = ?", productID)
	if status >= 0 {
		db = db.Where("status = ?", status)
	}
	db.Count(&count)
	return count
}

// CountAll 获取所有卡密数量
func (s *CardKeyService) CountAll() int64 {
	var count int64
	database.GetDB().Model(&models.CardKey{}).Count(&count)
	return count
}

// CountAvailable 获取可售卡密数量
func (s *CardKeyService) CountAvailable() int64 {
	var count int64
	database.GetDB().Model(&models.CardKey{}).Where("status = ?", models.CardKeyStatusAvailable).Count(&count)
	return count
}

// 错误定义
var ErrCardKeySold = &ServiceError{Message: "该卡密已售出，无法删除"}
