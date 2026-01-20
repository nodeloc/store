package services

import (
	"time"

	"github.com/nodeloc-faka/database"
	"github.com/nodeloc-faka/models"
)

// UserService 用户服务
type UserService struct{}

// NewUserService 创建用户服务
func NewUserService() *UserService {
	return &UserService{}
}

// CreateOrUpdate 创建或更新用户（用于OAuth登录）
func (s *UserService) CreateOrUpdate(nodeLocID int, username, name, email, avatarURL string, trustLevel int) (*models.User, error) {
	var user models.User
	result := database.GetDB().Where("node_loc_id = ?", nodeLocID).First(&user)

	now := time.Now()
	if result.Error != nil {
		// 不存在，创建新用户
		// 检查是否是第一个用户，如果是则自动设为管理员
		var count int64
		database.GetDB().Model(&models.User{}).Count(&count)
		isFirstUser := count == 0
		
		user = models.User{
			NodeLocID:   nodeLocID,
			Username:    username,
			Name:        name,
			Email:       email,
			AvatarURL:   avatarURL,
			TrustLevel:  trustLevel,
			IsAdmin:     isFirstUser, // 第一个用户自动成为管理员
			LastLoginAt: &now,
		}
		if err := database.GetDB().Create(&user).Error; err != nil {
			return nil, err
		}
	} else {
		// 存在，更新信息（不修改 IsAdmin 状态）
		user.Username = username
		user.Name = name
		user.Email = email
		user.AvatarURL = avatarURL
		user.TrustLevel = trustLevel
		user.LastLoginAt = &now
		if err := database.GetDB().Save(&user).Error; err != nil {
			return nil, err
		}
	}

	return &user, nil
}

// FindByID 根据ID查找用户
func (s *UserService) FindByID(id uint) (*models.User, error) {
	var user models.User
	if err := database.GetDB().First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByNodeLocID 根据NodeLoc ID查找用户
func (s *UserService) FindByNodeLocID(nodeLocID int) (*models.User, error) {
	var user models.User
	if err := database.GetDB().Where("node_loc_id = ?", nodeLocID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetAll 获取所有用户
func (s *UserService) GetAll() ([]models.User, error) {
	var users []models.User
	if err := database.GetDB().Order("id desc").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// GetWithPagination 分页获取用户
func (s *UserService) GetWithPagination(page, pageSize int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	db := database.GetDB().Model(&models.User{})
	db.Count(&total)

	offset := (page - 1) * pageSize
	if err := db.Order("id desc").
		Offset(offset).
		Limit(pageSize).
		Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// Block 封禁用户
func (s *UserService) Block(id uint) error {
	return database.GetDB().Model(&models.User{}).Where("id = ?", id).Update("is_blocked", true).Error
}

// Unblock 解封用户
func (s *UserService) Unblock(id uint) error {
	return database.GetDB().Model(&models.User{}).Where("id = ?", id).Update("is_blocked", false).Error
}

// UpdateBalance 更新余额
func (s *UserService) UpdateBalance(id uint, amount float64) error {
	return database.GetDB().Model(&models.User{}).
		Where("id = ?", id).
		UpdateColumn("balance", database.GetDB().Raw("balance + ?", amount)).
		Error
}

// Count 获取用户数量
func (s *UserService) Count() int64 {
	var count int64
	database.GetDB().Model(&models.User{}).Count(&count)
	return count
}

// IsBlocked 检查用户是否被封禁
func (s *UserService) IsBlocked(id uint) bool {
	var user models.User
	if err := database.GetDB().First(&user, id).Error; err != nil {
		return false
	}
	return user.IsBlocked
}

// SetAdmin 设置用户为管理员
func (s *UserService) SetAdmin(id uint, isAdmin bool) error {
	return database.GetDB().Model(&models.User{}).Where("id = ?", id).Update("is_admin", isAdmin).Error
}

// IsAdmin 检查用户是否是管理员
func (s *UserService) IsAdmin(id uint) bool {
	var user models.User
	if err := database.GetDB().First(&user, id).Error; err != nil {
		return false
	}
	return user.IsAdmin
}
