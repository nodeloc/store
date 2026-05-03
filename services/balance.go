package services

import (
	"fmt"

	"github.com/nodeloc-faka/database"
	"github.com/nodeloc-faka/models"
	"gorm.io/gorm"
)

// BalanceService 余额服务（事务性）
type BalanceService struct{}

// NewBalanceService 创建余额服务
func NewBalanceService() *BalanceService {
	return &BalanceService{}
}

var ErrInsufficientBalance = &ServiceError{Message: "余额不足"}

// AddBalance 在事务中给用户增加余额并写入流水
func (s *BalanceService) AddBalance(tx *gorm.DB, userID uint, amount float64, txType, description, refType string, refID uint) error {
	if amount == 0 {
		return nil
	}
	if userID == 0 {
		return nil // 平台官方店等无用户的店铺，不做余额累计
	}

	// 锁定用户行（使用标准 GORM clause）
	var user models.User
	if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil // 用户不存在（如店主账号未注册），跳过结算
		}
		return err
	}

	newBalance := user.Balance + amount
	if newBalance < 0 {
		return ErrInsufficientBalance
	}

	if err := tx.Model(&models.User{}).Where("id = ?", userID).
		Update("balance", newBalance).Error; err != nil {
		return err
	}

	record := &models.BalanceTx{
		UserID:       userID,
		Type:         txType,
		Amount:       amount,
		BalanceAfter: newBalance,
		RefType:      refType,
		RefID:        refID,
		Description:  description,
	}
	if err := tx.Create(record).Error; err != nil {
		return err
	}
	return nil
}

// DeductBalance 扣减用户余额（amount 为正数），失败返回 ErrInsufficientBalance
func (s *BalanceService) DeductBalance(tx *gorm.DB, userID uint, amount float64, txType, description, refType string, refID uint) error {
	if amount <= 0 {
		return fmt.Errorf("扣减金额必须为正数")
	}
	return s.AddBalance(tx, userID, -amount, txType, description, refType, refID)
}

// GetUserTxs 获取用户余额流水（分页）
func (s *BalanceService) GetUserTxs(userID uint, page, pageSize int) ([]models.BalanceTx, int64, error) {
	var txs []models.BalanceTx
	var total int64
	db := database.GetDB().Model(&models.BalanceTx{}).Where("user_id = ?", userID)
	db.Count(&total)
	offset := (page - 1) * pageSize
	if err := db.Order("id desc").Offset(offset).Limit(pageSize).Find(&txs).Error; err != nil {
		return nil, 0, err
	}
	return txs, total, nil
}
