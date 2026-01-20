package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/nodeloc-faka/models"
	"github.com/nodeloc-faka/services"
)

// SessionStore 简单的内存 session 存储
type SessionStore struct {
	sessions map[string]map[string]interface{}
	mu       sync.RWMutex
}

// NewSessionStore 创建新的 session 存储
func NewSessionStore() *SessionStore {
	return &SessionStore{
		sessions: make(map[string]map[string]interface{}),
	}
}

// Get 获取 session
func (s *SessionStore) Get(sessionID string) map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, exists := s.sessions[sessionID]
	if !exists {
		return make(map[string]interface{})
	}
	return session
}

// Set 设置 session
func (s *SessionStore) Set(sessionID string, data map[string]interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.sessions[sessionID] = data
}

// SessionMiddleware Session 中间件
func SessionMiddleware(store *SessionStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 cookie 获取或生成 session ID
		sessionID, err := c.Cookie("session_id")
		if err != nil {
			// 生成新的 session ID
			sessionID = generateSessionID()
			c.SetCookie("session_id", sessionID, 86400*7, "/", "", false, true)
		}

		// 获取 session 数据
		session := store.Get(sessionID)
		c.Set("session", session)
		c.Set("session_id", sessionID)

		// 如果 session 中有用户信息，设置到 context
		if userInterface, exists := session["user"]; exists {
			if user, ok := userInterface.(*models.User); ok {
				c.Set("user", user)
			}
		}

		c.Next()

		// 请求处理完成后保存 session
		sessionData := c.MustGet("session").(map[string]interface{})
		store.Set(sessionID, sessionData)
	}
}

// generateSessionID 生成 session ID
func generateSessionID() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "session_fallback"
	}
	return base64.URLEncoding.EncodeToString(b)
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
