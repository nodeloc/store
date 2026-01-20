package services

import (
	"errors"
	"time"

	"github.com/nodeloc-faka/config"
	"github.com/nodeloc-faka/database"
	"github.com/nodeloc-faka/models"
)

// AdminService 管理员服务
type AdminService struct{}

// NewAdminService 创建管理员服务
func NewAdminService() *AdminService {
	return &AdminService{}
}

// Create 创建管理员
func (s *AdminService) Create(username, password string) (*models.Admin, error) {
	hash, err := config.HashPassword(password)
	if err != nil {
		return nil, err
	}

	admin := &models.Admin{
		Username:     username,
		PasswordHash: hash,
	}

	if err := database.GetDB().Create(admin).Error; err != nil {
		return nil, err
	}

	return admin, nil
}

// FindByUsername 根据用户名查找管理员
func (s *AdminService) FindByUsername(username string) (*models.Admin, error) {
	var admin models.Admin
	if err := database.GetDB().Where("username = ?", username).First(&admin).Error; err != nil {
		return nil, err
	}
	return &admin, nil
}

// FindByID 根据ID查找管理员
func (s *AdminService) FindByID(id uint) (*models.Admin, error) {
	var admin models.Admin
	if err := database.GetDB().First(&admin, id).Error; err != nil {
		return nil, err
	}
	return &admin, nil
}

// Verify 验证管理员登录
func (s *AdminService) Verify(username, password string) (*models.Admin, error) {
	admin, err := s.FindByUsername(username)
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	if !config.CheckPassword(password, admin.PasswordHash) {
		return nil, errors.New("用户名或密码错误")
	}

	// 更新最后登录时间
	now := time.Now()
	admin.LastLoginAt = &now
	database.GetDB().Save(admin)

	return admin, nil
}

// UpdatePassword 更新密码
func (s *AdminService) UpdatePassword(id uint, newPassword string) error {
	hash, err := config.HashPassword(newPassword)
	if err != nil {
		return err
	}

	return database.GetDB().Model(&models.Admin{}).Where("id = ?", id).Update("password_hash", hash).Error
}

// GetAll 获取所有管理员
func (s *AdminService) GetAll() ([]models.Admin, error) {
	var admins []models.Admin
	if err := database.GetDB().Find(&admins).Error; err != nil {
		return nil, err
	}
	return admins, nil
}

// Delete 删除管理员
func (s *AdminService) Delete(id uint) error {
	return database.GetDB().Delete(&models.Admin{}, id).Error
}

// Count 获取管理员数量
func (s *AdminService) Count() int64 {
	var count int64
	database.GetDB().Model(&models.Admin{}).Count(&count)
	return count
}
