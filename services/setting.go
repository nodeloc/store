package services

import (
	"github.com/nodeloc-faka/database"
	"github.com/nodeloc-faka/models"
)

// SettingService 设置服务
type SettingService struct{}

// NewSettingService 创建设置服务
func NewSettingService() *SettingService {
	return &SettingService{}
}

// Get 获取设置值
func (s *SettingService) Get(key string) string {
	var setting models.Setting
	// 使用原生 SQL 查询，避免保留字问题
	result := database.GetDB().Raw("SELECT id, `key`, value, created_at, updated_at FROM settings WHERE `key` = ? LIMIT 1", key).Scan(&setting)
	if result.Error != nil || result.RowsAffected == 0 {
		return ""
	}
	return setting.Value
}

// Set 设置值
func (s *SettingService) Set(key, value string) error {
	var setting models.Setting
	// 先查询是否存在
	result := database.GetDB().Raw("SELECT id, `key`, value, created_at, updated_at FROM settings WHERE `key` = ? LIMIT 1", key).Scan(&setting)
	
	if result.Error != nil || result.RowsAffected == 0 {
		// 不存在，创建新的
		insertResult := database.GetDB().Exec("INSERT INTO settings (`key`, value, created_at, updated_at) VALUES (?, ?, NOW(), NOW())", key, value)
		return insertResult.Error
	}
	
	// 存在，更新
	updateResult := database.GetDB().Exec("UPDATE settings SET value = ?, updated_at = NOW() WHERE `key` = ?", value, key)
	return updateResult.Error
}

// GetAll 获取所有设置
func (s *SettingService) GetAll() map[string]string {
	var settings []models.Setting
	database.GetDB().Raw("SELECT id, `key`, value, created_at, updated_at FROM settings").Scan(&settings)
	
	result := make(map[string]string)
	for _, setting := range settings {
		result[setting.Key] = setting.Value
	}
	return result
}

// SetMultiple 批量设置
func (s *SettingService) SetMultiple(settings map[string]string) error {
	for key, value := range settings {
		if err := s.Set(key, value); err != nil {
			return err
		}
	}
	return nil
}

// 常用设置键
const (
	SettingSiteName        = "site_name"
	SettingSiteDescription = "site_description"
	SettingSiteLogo        = "site_logo"
	SettingSiteKeywords    = "site_keywords"
	SettingAdminPath       = "admin_path"
	SettingNodeLocClientID     = "nodeloc_client_id"
	SettingNodeLocClientSecret = "nodeloc_client_secret"
	SettingNodeLocRedirectURI  = "nodeloc_redirect_uri"
	SettingSessionSecret   = "session_secret"
	SettingContactEmail    = "contact_email"
	SettingContactQQ       = "contact_qq"
	SettingAnnouncement    = "announcement"
	SettingFooterText      = "footer_text"
	SettingInitialized     = "initialized"
	// 支付相关设置
	SettingPaymentID       = "payment_id"
	SettingPaymentSecret   = "payment_secret"
	SettingPaymentEnabled  = "payment_enabled"
	SettingPaymentCallback = "payment_callback"
)

// GetSiteSettings 获取网站设置
func (s *SettingService) GetSiteSettings() map[string]string {
	keys := []string{
		SettingSiteName,
		SettingSiteDescription,
		SettingSiteLogo,
		SettingSiteKeywords,
		SettingContactEmail,
		SettingContactQQ,
		SettingAnnouncement,
		SettingFooterText,
	}
	
	result := make(map[string]string)
	for _, key := range keys {
		result[key] = s.Get(key)
	}
	return result
}
