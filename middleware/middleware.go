package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nodeloc-faka/database"
	"github.com/nodeloc-faka/models"
	"github.com/nodeloc-faka/services"
)

const (
	sessionCookieName   = "session_id"
	sessionTTL          = 7 * 24 * time.Hour
	sessionDestroyedKey = "session_destroyed"
)

// SessionStore 数据库持久化 session 存储
type SessionStore struct{}

// NewSessionStore 创建 session 存储
func NewSessionStore() *SessionStore {
	return &SessionStore{}
}

// Get 从数据库获取 session
func (s *SessionStore) Get(sessionID string) map[string]interface{} {
	var sess models.Session
	if err := database.GetDB().First(&sess, "id = ? AND expires_at > ?", sessionID, time.Now()).Error; err != nil {
		return make(map[string]interface{})
	}
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(sess.Data), &data); err != nil || data == nil {
		return make(map[string]interface{})
	}
	return data
}

// Set 保存 session 到数据库
func (s *SessionStore) Set(sessionID string, data map[string]interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		return
	}
	sess := models.Session{
		ID:        sessionID,
		Data:      string(b),
		ExpiresAt: time.Now().Add(sessionTTL),
		UpdatedAt: time.Now(),
	}
	database.GetDB().Save(&sess)
}

// Delete 删除 session
func (s *SessionStore) Delete(sessionID string) {
	database.GetDB().Delete(&models.Session{}, "id = ?", sessionID)
}

// PurgeExpired 分批删除过期 session，避免一次性大删除长时间锁表
func (s *SessionStore) PurgeExpired() int64 {
	var total int64
	for {
		res := database.GetDB().Exec("DELETE FROM sessions WHERE expires_at < ? LIMIT 5000", time.Now())
		if res.Error != nil {
			log.Printf("清理过期 session 失败: %v", res.Error)
			return total
		}
		total += res.RowsAffected
		if res.RowsAffected < 5000 {
			return total
		}
	}
}

// StartJanitor 启动后台协程定期清理过期 session
func (s *SessionStore) StartJanitor(interval time.Duration) {
	go func() {
		for {
			if n := s.PurgeExpired(); n > 0 {
				log.Printf("已清理过期 session %d 条", n)
			}
			time.Sleep(interval)
		}
	}()
}

// secureCookies 站点走 HTTPS 时给 cookie 加 Secure；本地开发是 http，不能加
func secureCookies() bool {
	return strings.HasPrefix(os.Getenv("NODELOC_REDIRECT_URI"), "https://")
}

func setSessionCookie(c *gin.Context, sessionID string, maxAge int) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(sessionCookieName, sessionID, maxAge, "/", "", secureCookies(), true)
}

// RegenerateSession 登录成功后更换 session ID（数据保留），防止会话固定攻击
func RegenerateSession(c *gin.Context) {
	if oldID := c.GetString("session_id"); oldID != "" {
		database.GetDB().Delete(&models.Session{}, "id = ?", oldID)
	}
	newID := generateSessionID()
	c.Set("session_id", newID)
	setSessionCookie(c, newID, int(sessionTTL.Seconds()))
}

// DestroySession 退出登录：删除 session 记录并让浏览器丢弃 cookie
func DestroySession(c *gin.Context) {
	if id := c.GetString("session_id"); id != "" {
		database.GetDB().Delete(&models.Session{}, "id = ?", id)
	}
	c.Set(sessionDestroyedKey, true)
	setSessionCookie(c, "", -1)
}

// SessionMiddleware Session 中间件
func SessionMiddleware(store *SessionStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, err := c.Cookie(sessionCookieName)
		if err != nil || sessionID == "" {
			sessionID = generateSessionID()
			setSessionCookie(c, sessionID, int(sessionTTL.Seconds()))
		}

		session := store.Get(sessionID)
		before, _ := json.Marshal(session)
		c.Set("session", session)
		c.Set("session_id", sessionID)

		// 如果 session 中有 user_id，从数据库加载用户；被封禁的用户视为未登录
		if rawID, exists := session["user_id"]; exists && rawID != nil {
			var userID uint
			switch v := rawID.(type) {
			case float64:
				userID = uint(v)
			case uint:
				userID = v
			}
			if userID > 0 {
				var user models.User
				if err := database.GetDB().First(&user, userID).Error; err == nil && !user.IsBlocked {
					c.Set("user", &user)
				}
			}
		}

		c.Next()

		if c.GetBool(sessionDestroyedKey) {
			return
		}
		sessionData := c.MustGet("session").(map[string]interface{})
		currentID := c.GetString("session_id")
		after, _ := json.Marshal(sessionData)
		// 只在内容变化（或 ID 轮换）时写库：匿名访客和爬虫的每个请求都会带来一个空 session，
		// 全部落库会让 sessions 表无限膨胀
		if currentID == sessionID && string(after) == string(before) {
			return
		}
		if len(sessionData) == 0 {
			store.Delete(currentID)
			return
		}
		store.Set(currentID, sessionData)
	}
}

// generateSessionID 生成 session ID
func generateSessionID() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		panic("crypto/rand unavailable: " + err.Error())
	}
	return base64.URLEncoding.EncodeToString(b)
}

// SecurityHeaders 基础安全响应头
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Next()
	}
}

// AdminAuthMiddleware 后台认证中间件
func AdminAuthMiddleware() gin.HandlerFunc {
	settingService := services.NewSettingService()

	return func(c *gin.Context) {
		session := c.MustGet("session").(map[string]interface{})
		adminPath := settingService.Get(services.SettingAdminPath)

		// 检查是否是登录页面
		if c.Request.URL.Path == adminPath+"/login" {
			c.Next()
			return
		}

		// 检查是否已登录
		adminID, exists := session["admin_id"]
		if !exists || adminID == nil {
			c.Redirect(http.StatusFound, adminPath+"/login")
			c.Abort()
			return
		}

		c.Set("admin_id", adminID)
		c.Set("admin_username", session["admin_username"])
		c.Next()
	}
}

// UserAuthMiddleware 用户认证中间件（可选认证）
func UserAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 不强制登录，只是检查是否已登录
		c.Next()
	}
}

// RequireUserMiddleware 需要用户登录中间件
func RequireUserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userInterface, exists := c.Get("user")
		if !exists || userInterface == nil {
			redirect := c.Request.URL.String()
			c.Redirect(http.StatusFound, "/login?redirect="+redirect)
			c.Abort()
			return
		}

		c.Next()
	}
}

// AuthRequired API 认证中间件（返回 JSON）
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		userInterface, exists := c.Get("user")
		if !exists || userInterface == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// AdminRequired API 管理员认证中间件（返回 JSON）
func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 首先检查是否已登录
		userInterface, exists := c.Get("user")
		if !exists || userInterface == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
			c.Abort()
			return
		}

		// 检查是否是管理员
		user, ok := userInterface.(*models.User)
		if !ok || !user.IsAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "需要管理员权限"})
			c.Abort()
			return
		}

		c.Next()
	}
}
