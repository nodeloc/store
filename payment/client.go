package payment

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

// Client NodeLoc 支付客户端
type Client struct {
	baseURL   string
	paymentID string
	secretKey string
	httpClient *http.Client
}

// NewClient 创建支付客户端
func NewClient(baseURL, paymentID, secretKey string) *Client {
	return &Client{
		baseURL:   baseURL,
		paymentID: paymentID,
		secretKey: secretKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// UpdateConfig 更新配置
func (c *Client) UpdateConfig(paymentID, secretKey string) {
	if paymentID != "" {
		c.paymentID = paymentID
	}
	if secretKey != "" {
		c.secretKey = secretKey
	}
}

// IsConfigured 检查是否已配置
func (c *Client) IsConfigured() bool {
	return c.paymentID != "" && c.secretKey != ""
}

// CreatePaymentRequest 发起支付请求
type CreatePaymentRequest struct {
	Amount      int    `json:"amount"`
	Description string `json:"description"`
	OrderID     string `json:"order_id"`
}

// CreatePaymentResponse 发起支付响应
type CreatePaymentResponse struct {
	PaymentURL    string `json:"payment_url"`
	TransactionID string `json:"transaction_id"`
	Status        string `json:"status"`
	Amount        int    `json:"amount"`
}

// CreatePayment 发起支付
func (c *Client) CreatePayment(req *CreatePaymentRequest) (*CreatePaymentResponse, error) {
	// 准备签名参数
	params := map[string]string{
		"amount":      fmt.Sprintf("%d", req.Amount),
		"description": req.Description,
		"order_id":    req.OrderID,
	}

	// 生成签名
	signature := c.generateSignature(params)

	// 构建请求
	formData := url.Values{}
	for k, v := range params {
		formData.Set(k, v)
	}
	formData.Set("signature", signature)

	// 发送请求
	apiURL := fmt.Sprintf("%s/payment/pay/%s/process", c.baseURL, c.paymentID)
	resp, err := c.httpClient.PostForm(apiURL, formData)
	if err != nil {
		return nil, fmt.Errorf("发起支付请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("支付请求失败: status=%d, body=%s", resp.StatusCode, string(body))
	}

	var result CreatePaymentResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &result, nil
}

// QueryPaymentRequest 查询支付请求
type QueryPaymentRequest struct {
	TransactionID string `json:"transaction_id"`
}

// QueryPaymentResponse 查询支付响应
type QueryPaymentResponse struct {
	TransactionID     string  `json:"transaction_id"`
	Status            string  `json:"status"`
	Amount            int     `json:"amount"`
	PlatformFee       int     `json:"platform_fee"`
	MerchantPoints    int     `json:"merchant_points"`
	Description       string  `json:"description"`
	ExternalReference string  `json:"external_reference"`
	CreatedAt         string  `json:"created_at"`
	UpdatedAt         string  `json:"updated_at"`
	PaidAt            *string `json:"paid_at"`
	ExpiresAt         string  `json:"expires_at"`
	Expired           bool    `json:"expired"`
	ErrorMessage      *string `json:"error_message"`
}

// QueryPayment 查询支付状态
func (c *Client) QueryPayment(transactionID string) (*QueryPaymentResponse, error) {
	// 准备签名参数
	params := map[string]string{
		"transaction_id": transactionID,
	}

	// 生成签名
	signature := c.generateSignature(params)

	// 构建请求
	formData := url.Values{}
	formData.Set("transaction_id", transactionID)
	formData.Set("signature", signature)

	// 发送请求
	apiURL := fmt.Sprintf("%s/payment/query/%s", c.baseURL, c.paymentID)
	resp, err := c.httpClient.PostForm(apiURL, formData)
	if err != nil {
		return nil, fmt.Errorf("查询支付请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("查询支付失败: status=%d, body=%s", resp.StatusCode, string(body))
	}

	var result QueryPaymentResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &result, nil
}

// CallbackParams 回调参数
type CallbackParams struct {
	TransactionID     string `json:"transaction_id"`
	ExternalReference string `json:"external_reference"`
	Amount            int    `json:"amount"`
	PlatformFee       int    `json:"platform_fee"`
	MerchantPoints    int    `json:"merchant_points"`
	Status            string `json:"status"`
	PaidAt            string `json:"paid_at"`
	Signature         string `json:"signature"`
}

// VerifyCallback 验证回调签名
func (c *Client) VerifyCallback(params *CallbackParams) bool {
	// 准备验证参数（排除 signature）
	verifyParams := map[string]string{
		"amount":             fmt.Sprintf("%d", params.Amount),
		"external_reference": params.ExternalReference,
		"merchant_points":    fmt.Sprintf("%d", params.MerchantPoints),
		"paid_at":            params.PaidAt,
		"platform_fee":       fmt.Sprintf("%d", params.PlatformFee),
		"status":             params.Status,
		"transaction_id":     params.TransactionID,
	}

	// 使用 secret_key 直接计算签名（回调验证）
	expectedSignature := c.generateCallbackSignature(verifyParams)

	return hmac.Equal([]byte(expectedSignature), []byte(params.Signature))
}

// generateSignature 生成签名（用于发起支付和查询）
// 使用 token_hash = SHA256(token), 然后 HMAC-SHA256(token_hash, params)
func (c *Client) generateSignature(params map[string]string) string {
	// 1. 对参数键排序
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 2. 拼接参数字符串
	var pairs []string
	for _, k := range keys {
		pairs = append(pairs, fmt.Sprintf("%s=%s", k, params[k]))
	}
	paramString := strings.Join(pairs, "&")

	// 3. 计算 token_hash = SHA256(secret_key)
	tokenHash := sha256.Sum256([]byte(c.secretKey))
	tokenHashHex := hex.EncodeToString(tokenHash[:])

	// 4. 生成 HMAC-SHA256 签名
	h := hmac.New(sha256.New, []byte(tokenHashHex))
	h.Write([]byte(paramString))
	return hex.EncodeToString(h.Sum(nil))
}

// generateCallbackSignature 生成回调签名（用于验证回调）
// 直接使用 secret_key 作为 HMAC 密钥
func (c *Client) generateCallbackSignature(params map[string]string) string {
	// 1. 对参数键排序
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 2. 拼接参数字符串
	var pairs []string
	for _, k := range keys {
		pairs = append(pairs, fmt.Sprintf("%s=%s", k, params[k]))
	}
	paramString := strings.Join(pairs, "&")

	// 3. 生成 HMAC-SHA256 签名
	h := hmac.New(sha256.New, []byte(c.secretKey))
	h.Write([]byte(paramString))
	return hex.EncodeToString(h.Sum(nil))
}

// PaymentStatus 支付状态常量
const (
	PaymentStatusPending    = "pending"
	PaymentStatusProcessing = "processing"
	PaymentStatusCompleted  = "completed"
	PaymentStatusFailed     = "failed"
	PaymentStatusCancelled  = "cancelled"
	PaymentStatusRefunded   = "refunded"
)
