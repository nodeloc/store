package services

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

// PaymentService 支付服务
type PaymentService struct {
	settingService *SettingService
}

// NewPaymentService 创建支付服务
func NewPaymentService() *PaymentService {
	return &PaymentService{
		settingService: NewSettingService(),
	}
}

// PaymentConfig 支付配置
type PaymentConfig struct {
	PaymentID  string // Payment ID (pay_xxx)
	SecretKey  string // Secret Key
	PaymentURL string // 支付API地址
}

// GetConfig 获取支付配置
func (s *PaymentService) GetConfig() *PaymentConfig {
	return &PaymentConfig{
		PaymentID:  s.settingService.Get(SettingPaymentID),
		SecretKey:  s.settingService.Get(SettingPaymentSecret),
		PaymentURL: "https://www.nodeloc.com",
	}
}

// IsConfigured 检查是否已配置支付
func (s *PaymentService) IsConfigured() bool {
	cfg := s.GetConfig()
	return cfg.PaymentID != "" && cfg.SecretKey != ""
}

// CreatePaymentRequest 发起支付请求
type CreatePaymentRequest struct {
	Amount      int    `json:"amount"`      // 支付金额（能量）
	Description string `json:"description"` // 交易描述
	OrderID     string `json:"order_id"`    // 订单号
}

// CreatePaymentResponse 发起支付响应
type CreatePaymentResponse struct {
	PaymentURL    string `json:"payment_url"`
	TransactionID string `json:"transaction_id"`
	Status        string `json:"status"`
	Amount        int    `json:"amount"`
}

// CreatePayment 发起支付
func (s *PaymentService) CreatePayment(req *CreatePaymentRequest) (*CreatePaymentResponse, error) {
	cfg := s.GetConfig()
	if cfg.PaymentID == "" || cfg.SecretKey == "" {
		fmt.Printf("支付配置检查 - PaymentID: %s, SecretKey: %s\n", cfg.PaymentID, cfg.SecretKey)
		return nil, fmt.Errorf("支付未配置 (PaymentID=%s, SecretKey=%s)", cfg.PaymentID, cfg.SecretKey)
	}

	// 准备签名参数
	params := map[string]string{
		"amount":      fmt.Sprintf("%d", req.Amount),
		"description": req.Description,
		"order_id":    req.OrderID,
	}

	// 生成签名（发起支付使用 token_hash）
	signature := s.generateSignatureForPayment(params, cfg.SecretKey)
	params["signature"] = signature

	// 发送请求
	apiURL := fmt.Sprintf("%s/payment/pay/%s/process", cfg.PaymentURL, cfg.PaymentID)
	fmt.Printf("调用支付API: %s\n", apiURL)
	fmt.Printf("支付参数: amount=%d, description=%s, order_id=%s\n", req.Amount, req.Description, req.OrderID)

	// 构建表单数据
	formData := url.Values{}
	for k, v := range params {
		formData.Set(k, v)
	}

	client := &http.Client{Timeout: 8 * time.Second}
	resp, err := client.PostForm(apiURL, formData)
	if err != nil {
		return nil, fmt.Errorf("请求支付API失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("支付API返回错误: %s", string(body))
	}

	var result CreatePaymentResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &result, nil
}

// PaymentCallback 支付回调参数
type PaymentCallback struct {
	TransactionID     string `form:"transaction_id"`
	ExternalReference string `form:"external_reference"` // 订单号
	Amount            int    `form:"amount"`
	PlatformFee       int    `form:"platform_fee"`
	MerchantPoints    int    `form:"merchant_points"`
	Status            string `form:"status"`
	PaidAt            string `form:"paid_at"`
	Signature         string `form:"signature"`
}

// VerifyCallback 验证回调签名
// NodeLoc 回调签名使用直接 secret_key 作为 HMAC 密钥（不经过 token_hash）
func (s *PaymentService) VerifyCallback(callback *PaymentCallback) bool {
	cfg := s.GetConfig()
	if cfg.SecretKey == "" {
		fmt.Println("回调验证失败: SecretKey 为空")
		return false
	}

	// 准备签名参数（排除 signature 本身）
	params := map[string]string{
		"transaction_id":     callback.TransactionID,
		"external_reference": callback.ExternalReference,
		"amount":             fmt.Sprintf("%d", callback.Amount),
		"platform_fee":       fmt.Sprintf("%d", callback.PlatformFee),
		"merchant_points":    fmt.Sprintf("%d", callback.MerchantPoints),
		"status":             callback.Status,
		"paid_at":            callback.PaidAt,
	}

	// 回调验证：直接用 secret_key 做 HMAC（不经过 token_hash）
	expectedSignature := s.generateSignatureForCallback(params, cfg.SecretKey)

	fmt.Println("=== 回调签名验证 ===")
	fmt.Printf("  transaction_id: %s\n", callback.TransactionID)
	fmt.Printf("  external_reference: %s\n", callback.ExternalReference)
	fmt.Printf("  signature (接收): %s\n", callback.Signature)
	fmt.Printf("  signature (计算): %s\n", expectedSignature)

	isValid := hmac.Equal([]byte(expectedSignature), []byte(callback.Signature))
	fmt.Printf("  验证结果: %v\n", isValid)

	// 兼容：如果直接 secret_key 方式失败，再尝试 token_hash 方式
	if !isValid {
		expectedAlt := s.generateSignatureForPayment(params, cfg.SecretKey)
		isValid = hmac.Equal([]byte(expectedAlt), []byte(callback.Signature))
		if isValid {
			fmt.Println("  (使用 token_hash 方式验证通过)")
		}
	}
	fmt.Println("==================")
	return isValid
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
func (s *PaymentService) QueryPayment(transactionID string) (*QueryPaymentResponse, error) {
	cfg := s.GetConfig()
	if cfg.PaymentID == "" || cfg.SecretKey == "" {
		return nil, fmt.Errorf("支付未配置")
	}

	// 准备签名参数
	params := map[string]string{
		"transaction_id": transactionID,
	}

	// 生成签名（查询也使用 token_hash）
	signature := s.generateSignatureForPayment(params, cfg.SecretKey)
	params["signature"] = signature

	// 发送请求
	apiURL := fmt.Sprintf("%s/payment/query/%s", cfg.PaymentURL, cfg.PaymentID)

	formData := url.Values{}
	for k, v := range params {
		formData.Set(k, v)
	}

	client := &http.Client{Timeout: 8 * time.Second}
	resp, err := client.PostForm(apiURL, formData)
	if err != nil {
		return nil, fmt.Errorf("请求查询API失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("查询API返回错误: %s", string(body))
	}

	var result QueryPaymentResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &result, nil
}

// generateSignatureForPayment 生成发起支付的签名（使用 token_hash）
func (s *PaymentService) generateSignatureForPayment(params map[string]string, token string) string {
	// 1. 按键名排序
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 2. 拼接字符串
	var pairs []string
	for _, k := range keys {
		pairs = append(pairs, fmt.Sprintf("%s=%s", k, params[k]))
	}
	paramString := strings.Join(pairs, "&")

	// 3. 计算 token_hash = SHA256(token) - 根据文档第 69 行
	tokenHashBytes := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(tokenHashBytes[:])

	// 4. 计算 HMAC-SHA256(token_hash, paramString) - 根据文档第 70 行
	h := hmac.New(sha256.New, []byte(tokenHash))
	h.Write([]byte(paramString))
	signature := hex.EncodeToString(h.Sum(nil))

	fmt.Printf("签名调试 - 参数字符串: %s\n", paramString)
	fmt.Printf("签名调试 - token: %s\n", token)
	fmt.Printf("签名调试 - token_hash (SHA256): %s\n", tokenHash)
	fmt.Printf("签名调试 - signature (HMAC-SHA256): %s\n", signature)

	return signature
}

// generateSignatureForCallback 生成回调验证的签名（直接使用 secret_key）
func (s *PaymentService) generateSignatureForCallback(params map[string]string, secretKey string) string {
	// 1. 按键名排序
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 2. 拼接字符串
	var pairs []string
	for _, k := range keys {
		pairs = append(pairs, fmt.Sprintf("%s=%s", k, params[k]))
	}
	paramString := strings.Join(pairs, "&")

	// 3. 计算 HMAC-SHA256(secret_key, paramString) - 根据文档第 157 行
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(paramString))
	signature := hex.EncodeToString(h.Sum(nil))

	fmt.Printf("回调签名生成 - 参数字符串: %s\n", paramString)
	fmt.Printf("回调签名生成 - secret_key: %s\n", secretKey)
	fmt.Printf("回调签名生成 - signature: %s\n", signature)

	return signature
}

// ProcessPaymentCallback 处理支付回调 — 委托给 OrderService 统一处理
// （包含手动发货判断、事务保护、店主余额结算）
func (s *PaymentService) ProcessPaymentCallback(callback *PaymentCallback) error {
	// 1. 验证签名
	if !s.VerifyCallback(callback) {
		return fmt.Errorf("签名验证失败")
	}

	// 2. 检查支付状态
	if callback.Status != "completed" {
		return fmt.Errorf("支付未完成: %s", callback.Status)
	}

	// 3. 委托 OrderService 处理（含事务、手动发货、店主结算）
	orderService := NewOrderService()
	_, err := orderService.ProcessPaymentCallback(
		callback.TransactionID,
		callback.Amount,
		callback.PlatformFee,
		callback.MerchantPoints,
	)
	return err
}
