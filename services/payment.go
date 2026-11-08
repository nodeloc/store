package services

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/nodeloc-faka/payment"
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
	cfg := &PaymentConfig{
		PaymentID:  s.settingService.Get(SettingPaymentID),
		SecretKey:  s.settingService.Get(SettingPaymentSecret),
		PaymentURL: "https://www.nodeloc.com",
	}
	// 与提现一致：后台未配置时回退到环境变量
	if cfg.PaymentID == "" || cfg.SecretKey == "" {
		env := getAppConfig()
		if cfg.PaymentID == "" {
			cfg.PaymentID = env.PaymentID
		}
		if cfg.SecretKey == "" {
			cfg.SecretKey = env.PaymentSecret
		}
	}
	return cfg
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
		return nil, fmt.Errorf("支付未配置")
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
		log.Println("回调验证失败: 支付密钥未配置")
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
	isValid := hmac.Equal([]byte(expectedSignature), []byte(callback.Signature))

	// 兼容：如果直接 secret_key 方式失败，再尝试 token_hash 方式
	if !isValid {
		expectedAlt := s.generateSignatureForPayment(params, cfg.SecretKey)
		isValid = hmac.Equal([]byte(expectedAlt), []byte(callback.Signature))
	}
	if !isValid {
		log.Printf("回调签名验证失败 tx=%s order=%s", callback.TransactionID, callback.ExternalReference)
	}
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

// ErrRefundUnsupported NodeLoc 尚未提供原路退款接口（路由返回 404）
var ErrRefundUnsupported = errors.New("NodeLoc 原路退款接口不可用")

// RefundPaymentResponse 原路退款响应
type RefundPaymentResponse struct {
	Success           bool   `json:"success"`
	TransactionID     string `json:"transaction_id"`
	ExternalReference string `json:"external_reference"`
	Status            string `json:"status"`
	Amount            int    `json:"amount"`
	MerchantPoints    int    `json:"merchant_points"`
	Error             string `json:"error"`
}

// RefundPayment 原路全额退款：买家拿回全部积分，商家退回实收积分（无手续费损失）。
// amount 必须等于交易积分；接口幂等，已退款的交易再次调用仍返回成功。
func (s *PaymentService) RefundPayment(transactionID string, amount int) (*RefundPaymentResponse, error) {
	cfg := s.GetConfig()
	if cfg.PaymentID == "" || cfg.SecretKey == "" {
		return nil, fmt.Errorf("支付未配置")
	}

	params := map[string]string{
		"amount":         fmt.Sprintf("%d", amount),
		"transaction_id": transactionID,
	}
	formData := url.Values{}
	for k, v := range params {
		formData.Set(k, v)
	}
	formData.Set("signature", s.generateSignatureForPayment(params, cfg.SecretKey))

	apiURL := fmt.Sprintf("%s/payment/refund/%s", cfg.PaymentURL, cfg.PaymentID)
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.PostForm(apiURL, formData)
	if err != nil {
		return nil, fmt.Errorf("请求退款API失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	// 交易已由 query 确认存在，此时空 body 的 404 只能是退款路由尚未部署
	if resp.StatusCode == http.StatusNotFound && len(strings.TrimSpace(string(body))) == 0 {
		return nil, ErrRefundUnsupported
	}

	var result RefundPaymentResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("退款API返回异常: status=%d body=%s", resp.StatusCode, string(body))
	}
	if resp.StatusCode != http.StatusOK || !result.Success {
		msg := result.Error
		if msg == "" {
			msg = fmt.Sprintf("status=%d", resp.StatusCode)
		}
		return nil, fmt.Errorf("退款失败: %s", msg)
	}
	return &result, nil
}

// TransferRefund 原路退款不可用时，从商户账户转账退回（商户承担转账手续费）。
// order_id 固定为 refund_<订单号>，NodeLoc 侧幂等，重复调用不会重复转账。
func (s *PaymentService) TransferRefund(orderNo string, toNodeLocID int, toUsername string, amount int) (*payment.TransferResponse, error) {
	cfg := s.GetConfig()
	if cfg.PaymentID == "" || cfg.SecretKey == "" {
		return nil, fmt.Errorf("支付未配置")
	}
	client := payment.NewClient(cfg.PaymentURL, cfg.PaymentID, cfg.SecretKey)
	return client.Transfer(&payment.TransferRequest{
		ToUserID:   toNodeLocID,
		ToUsername: toUsername,
		Amount:     amount,
		OrderID:    "refund_" + orderNo,
	})
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
	return hex.EncodeToString(h.Sum(nil))
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
	return hex.EncodeToString(h.Sum(nil))
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
