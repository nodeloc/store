package handler

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nodeloc-faka/models"
	"github.com/nodeloc-faka/oauth"
	"github.com/nodeloc-faka/services"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	oauthClient    *oauth.Client
	userService    *services.UserService
	settingService *services.SettingService
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(oauthClient *oauth.Client) *AuthHandler {
	return &AuthHandler{
		oauthClient:    oauthClient,
		userService:    services.NewUserService(),
		settingService: services.NewSettingService(),
	}
}

// Login 发起 OAuth 登录
func (h *AuthHandler) Login(c *gin.Context) {
	// 生成随机 state 用于防止 CSRF 攻击
	state, err := generateState()
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to generate state")
		return
	}

	// 保存重定向地址
	redirect := c.Query("redirect")
	if redirect == "" {
		redirect = "/"
	}

	// 将 state 和 redirect 保存到 session
	session := c.MustGet("session").(map[string]interface{})
	session["oauth_state"] = state
	session["oauth_redirect"] = redirect

	// 重定向到授权页面
	authURL := h.oauthClient.GetAuthURL(state)
	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

// Callback OAuth 回调处理
func (h *AuthHandler) Callback(c *gin.Context) {
	// 获取参数
	code := c.Query("code")
	state := c.Query("state")
	errorCode := c.Query("error")

	// 检查是否有错误
	if errorCode != "" {
		// OAuth 授权失败，重定向到登录页并显示错误
		c.Redirect(http.StatusTemporaryRedirect, "/?error="+errorCode)
		return
	}

	// 验证 state
	session := c.MustGet("session").(map[string]interface{})
	savedState, exists := session["oauth_state"].(string)
	if !exists || savedState != state {
		// State 验证失败，重定向到登录页
		c.Redirect(http.StatusTemporaryRedirect, "/?error=invalid_state")
		return
	}

	// 清除 state
	delete(session, "oauth_state")

	// 交换 token
	token, err := h.oauthClient.ExchangeToken(c.Request.Context(), code)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, "/?error=token_exchange_failed")
		return
	}

	// 获取用户信息
	nodeLocUser, err := h.oauthClient.GetUserInfo(c.Request.Context(), token)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, "/?error=get_user_info_failed")
		return
	}

	// 创建或更新用户
	user, err := h.userService.CreateOrUpdate(
		nodeLocUser.ID,
		nodeLocUser.Username,
		nodeLocUser.Name,
		nodeLocUser.Email,
		nodeLocUser.AvatarURL,
		nodeLocUser.TrustLevel,
	)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, "/?error=save_user_failed")
		return
	}

	// 检查用户是否被封禁
	if user.IsBlocked {
		c.Redirect(http.StatusTemporaryRedirect, "/?error=account_blocked")
		return
	}

	// 保存用户信息到 session
	session["user"] = user
	session["token"] = token

	// 获取重定向地址
	redirect := "/"
	if savedRedirect, ok := session["oauth_redirect"].(string); ok {
		redirect = savedRedirect
		delete(session, "oauth_redirect")
	}

	// 重定向
	c.Redirect(http.StatusTemporaryRedirect, redirect)
}

// Logout 退出登录
func (h *AuthHandler) Logout(c *gin.Context) {
	// 清除 session
	session := c.MustGet("session").(map[string]interface{})
	delete(session, "user")
	delete(session, "token")

	c.Redirect(http.StatusTemporaryRedirect, "/")
}

// GetCurrentUser 获取当前用户中间件辅助函数
func GetCurrentUser(c *gin.Context) *models.User {
	userInterface, exists := c.Get("user")
	if !exists || userInterface == nil {
		return nil
	}
	return userInterface.(*models.User)
}

// generateState 生成随机 state 字符串
func generateState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
