package services

import (
	"errors"
	"os"
	"time"

	"github.com/nodeloc-faka/database"
	"github.com/nodeloc-faka/models"
	"gorm.io/gorm"
)

// ShopService 店铺服务
type ShopService struct{}

// NewShopService 创建店铺服务
func NewShopService() *ShopService {
	return &ShopService{}
}

var (
	ErrShopAlreadyApplied = &ServiceError{Message: "您已经申请过店铺，请勿重复申请"}
	ErrShopNotFound       = &ServiceError{Message: "店铺不存在"}
	ErrShopNotApproved    = &ServiceError{Message: "店铺未通过审核"}
)

// Apply 用户申请开店
func (s *ShopService) Apply(shop *models.Shop) error {
	db := database.GetDB()

	// 检查是否已经申请
	var existing models.Shop
	err := db.Where("user_id = ?", shop.UserID).First(&existing).Error
	if err == nil {
		return ErrShopAlreadyApplied
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// AUTO_APPROVE_SHOP=true 时自动审核通过
	if os.Getenv("AUTO_APPROVE_SHOP") == "true" {
		now := time.Now()
		shop.Status = models.ShopStatusApproved
		shop.ReviewedAt = &now
	} else {
		shop.Status = models.ShopStatusPending
	}
	return db.Create(shop).Error
}

// GetByUser 根据用户ID查询店铺
func (s *ShopService) GetByUser(userID uint) (*models.Shop, error) {
	var shop models.Shop
	if err := database.GetDB().Where("user_id = ?", userID).First(&shop).Error; err != nil {
		return nil, err
	}
	return &shop, nil
}

// FindByID 根据ID查询
func (s *ShopService) FindByID(id uint) (*models.Shop, error) {
	var shop models.Shop
	if err := database.GetDB().Preload("User").First(&shop, id).Error; err != nil {
		return nil, err
	}
	return &shop, nil
}

// Update 更新店铺信息（店主可调用，仅可改名称/描述/联系方式/Logo）
func (s *ShopService) Update(shop *models.Shop) error {
	return database.GetDB().Model(&models.Shop{}).Where("id = ?", shop.ID).Updates(map[string]interface{}{
		"name":          shop.Name,
		"description":   shop.Description,
		"contact":       shop.Contact,
		"logo":          shop.Logo,
		"features":      shop.Features,
		"status":        shop.Status,
		"reject_reason": shop.RejectReason,
	}).Error
}

// Approve 批准店铺
func (s *ShopService) Approve(shopID uint) error {
	now := time.Now()
	return database.GetDB().Model(&models.Shop{}).Where("id = ?", shopID).Updates(map[string]interface{}{
		"status":        models.ShopStatusApproved,
		"reject_reason": "",
		"reviewed_at":   &now,
	}).Error
}

// Reject 拒绝店铺
func (s *ShopService) Reject(shopID uint, reason string) error {
	now := time.Now()
	return database.GetDB().Model(&models.Shop{}).Where("id = ?", shopID).Updates(map[string]interface{}{
		"status":        models.ShopStatusRejected,
		"reject_reason": reason,
		"reviewed_at":   &now,
	}).Error
}

// Block 封禁店铺
func (s *ShopService) Block(shopID uint, reason string) error {
	return database.GetDB().Model(&models.Shop{}).Where("id = ?", shopID).Updates(map[string]interface{}{
		"status":        models.ShopStatusBlocked,
		"reject_reason": reason,
	}).Error
}

// GetAll 获取所有店铺（分页 + 状态过滤），admin 用
func (s *ShopService) GetAll(page, pageSize, status int) ([]models.Shop, int64, error) {
	var shops []models.Shop
	var total int64

	db := database.GetDB().Model(&models.Shop{})
	if status >= 0 {
		db = db.Where("status = ?", status)
	}
	db.Count(&total)

	offset := (page - 1) * pageSize
	if err := db.Preload("User").Order("id desc").Offset(offset).Limit(pageSize).Find(&shops).Error; err != nil {
		return nil, 0, err
	}
	return shops, total, nil
}

// GetApproved 获取所有已批准的店铺（前台展示用，分页）
func (s *ShopService) GetApproved(page, pageSize int) ([]models.Shop, int64, error) {
	var shops []models.Shop
	var total int64
	db := database.GetDB().Model(&models.Shop{}).Where("status = ?", models.ShopStatusApproved)
	db.Count(&total)
	offset := (page - 1) * pageSize
	if err := db.Order("is_official desc, id asc").Offset(offset).Limit(pageSize).Find(&shops).Error; err != nil {
		return nil, 0, err
	}
	return shops, total, nil
}

// EnsureOfficialShop 确保平台官方店存在，并将无 ShopID 的存量商品挂到该店铺下
func (s *ShopService) EnsureOfficialShop() (*models.Shop, error) {
	db := database.GetDB()
	settingService := NewSettingService()

	// 查找或创建官方店
	var shop models.Shop
	err := db.Where("is_official = ?", true).First(&shop).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 创建官方店
		now := time.Now()
		shop = models.Shop{
			UserID:      0,
			Name:        "平台官方店",
			Description: "由平台直营的官方店铺",
			Status:      models.ShopStatusApproved,
			IsOfficial:  true,
			ReviewedAt:  &now,
		}
		if err := db.Create(&shop).Error; err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	// 每次都迁移 shop_id=0 的商品和订单到官方店（幂等操作）
	db.Model(&models.Product{}).
		Where("shop_id = 0 OR shop_id IS NULL").
		Update("shop_id", shop.ID)
	db.Model(&models.Order{}).
		Where("shop_id = 0 OR shop_id IS NULL").
		Update("shop_id", shop.ID)

	// 写入 setting
	settingService.Set(SettingOfficialShopID, IntToString(int(shop.ID)))
	return &shop, nil
}
