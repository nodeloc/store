package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/oauth2"
)

// NodeLocUser NodeLoc 用户信息结构
type NodeLocUser struct {
	ID         int    `json:"id"`
	Username   string `json:"username"`
	Name       string `json:"name"`
	AvatarURL  string `json:"avatar_url"`
	TrustLevel int    `json:"trust_level"`
	Email      string `json:"email,omitempty"`
}

// Client OAuth 客户端
type Client struct {
	config  *oauth2.Config
	baseURL string
}

// NewClient 创建新的 OAuth 客户端
func NewClient(baseURL, clientID, clientSecret, redirectURI string) *Client {
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURI,
		Scopes:       []string{"openid", "profile", "email"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  fmt.Sprintf("%s/oauth-provider/authorize", baseURL),
			TokenURL: fmt.Sprintf("%s/oauth-provider/token", baseURL),
		},
	}

	return &Client{
		config:  config,
		baseURL: baseURL,
	}
}

// UpdateConfig 更新配置
func (c *Client) UpdateConfig(clientID, clientSecret, redirectURI string) {
	if clientID != "" {
		c.config.ClientID = clientID
	}
	if clientSecret != "" {
		c.config.ClientSecret = clientSecret
	}
	if redirectURI != "" {
		c.config.RedirectURL = redirectURI
	}
}

// GetAuthURL 获取授权 URL
func (c *Client) GetAuthURL(state string) string {
	return c.config.AuthCodeURL(state, oauth2.AccessTypeOnline)
}

// ExchangeToken 使用授权码交换 Access Token
func (c *Client) ExchangeToken(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := c.config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange token: %w", err)
	}
	return token, nil
}

// GetUserInfo 获取用户信息
func (c *Client) GetUserInfo(ctx context.Context, token *oauth2.Token) (*NodeLocUser, error) {
	// 创建带有 token 的 HTTP 客户端
	client := c.config.Client(ctx, token)

	// 请求用户信息
	userInfoURL := fmt.Sprintf("%s/oauth-provider/userinfo", c.baseURL)
	resp, err := client.Get(userInfoURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get user info: status=%d, body=%s", resp.StatusCode, string(body))
	}

	// 解析响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var user NodeLocUser
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, fmt.Errorf("failed to parse user info: %w", err)
	}

	return &user, nil
}

// RefreshToken 刷新 Access Token
func (c *Client) RefreshToken(ctx context.Context, refreshToken string) (*oauth2.Token, error) {
	tokenSource := c.config.TokenSource(ctx, &oauth2.Token{
		RefreshToken: refreshToken,
	})

	token, err := tokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	return token, nil
}

// IsConfigured 检查是否已配置
func (c *Client) IsConfigured() bool {
	return c.config.ClientID != "" && c.config.ClientSecret != ""
}
