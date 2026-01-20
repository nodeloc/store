package config

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/nodeloc-faka/database"
	"golang.org/x/crypto/bcrypt"
)

// Config 应用配置结构
type Config struct {
	// NodeLoc OAuth 配置
	NodeLocURL          string
	NodeLocClientID     string
	NodeLocClientSecret string
	NodeLocRedirectURI  string

	// 服务器配置
	ServerPort    string
	SessionSecret string

	// 数据库配置
	Database *database.Config

	// 后台配置
	AdminPath     string
	AdminUsername string
	AdminPassword string
}

var AppConfig *Config

// LoadConfig 加载配置
func LoadConfig() (*Config, error) {
	// 尝试加载 .env 文件（如果存在）
	godotenv.Load()

	config := &Config{
		NodeLocURL:          getEnv("NODELOC_URL", "https://www.nodeloc.com"),
		NodeLocClientID:     getEnv("NODELOC_CLIENT_ID", ""),
		NodeLocClientSecret: getEnv("NODELOC_CLIENT_SECRET", ""),
		NodeLocRedirectURI:  getEnv("NODELOC_REDIRECT_URI", "http://localhost:8080/auth/callback"),
		ServerPort:          getEnv("PORT", "3000"),
		SessionSecret:       getEnv("SESSION_SECRET", ""),
		Database: &database.Config{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "faka"),
		},
		AdminPath:     getEnv("ADMIN_PATH", ""),
		AdminUsername: getEnv("ADMIN_USERNAME", ""),
		AdminPassword: getEnv("ADMIN_PASSWORD", ""),
	}

	// 如果没有设置SESSION_SECRET，则生成一个
	if config.SessionSecret == "" {
		config.SessionSecret = generateRandomString(32)
		log.Println("生成随机 SESSION_SECRET")
	}

	AppConfig = config
	return config, nil
}

// GenerateAdminCredentials 生成管理员凭证
func GenerateAdminCredentials() (path, username, password string) {
	path = "/" + generateRandomString(10)
	username = "admin_" + generateRandomString(6)
	password = generateRandomString(12)
	return
}

// HashPassword 哈希密码
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword 验证密码
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// generateRandomString 生成随机字符串
func generateRandomString(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "default_random_string"
	}
	str := hex.EncodeToString(bytes)
	// 只使用小写字母和数字
	str = strings.ToLower(str)
	if len(str) > length {
		str = str[:length]
	}
	return str
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// SaveToEnv 保存配置到 .env 文件
func SaveToEnv(key, value string) error {
	// 读取现有的 .env 文件
	content, err := os.ReadFile(".env")
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	lines := strings.Split(string(content), "\n")
	found := false
	for i, line := range lines {
		if strings.HasPrefix(line, key+"=") {
			lines[i] = fmt.Sprintf("%s=%s", key, value)
			found = true
			break
		}
	}

	if !found {
		lines = append(lines, fmt.Sprintf("%s=%s", key, value))
	}

	return os.WriteFile(".env", []byte(strings.Join(lines, "\n")), 0644)
}

// GetConfig 获取当前配置
func GetConfig() *Config {
	return AppConfig
}
