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

	"github.com/nodeloc-faka/database"
	"github.com/nodeloc-faka/models"
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
	Amount      int    `json:"amount"`      // 支付金额（积分）
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

	resp, err := http.PostForm(apiURL, formData)
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
func (s *PaymentService) VerifyCallback(callback *PaymentCallback) bool {
	cfg := s.GetConfig()
	if cfg.SecretKey == "" {
		fmt.Println("回调验证失败: SecretKey 为空")
		return false
	}

	// 准备签名参数（排除signature本身）
	params := map[string]string{
		"transaction_id":     callback.TransactionID,
		"external_reference": callback.ExternalReference,
		"amount":             fmt.Sprintf("%d", callback.Amount),
		"platform_fee":       fmt.Sprintf("%d", callback.PlatformFee),
		"merchant_points":    fmt.Sprintf("%d", callback.MerchantPoints),
		"status":             callback.Status,
		"paid_at":            callback.PaidAt,
	}

	// 生成签名并比对（回调验证也使用 token_hash，和发起支付一样）
	expectedSignature := s.generateSignatureForPayment(params, cfg.SecretKey)

	// 详细日志
	fmt.Println("=== 回调签名验证 ===")
	fmt.Printf("接收到的回调参数:\n")
	fmt.Printf("  transaction_id: %s\n", callback.TransactionID)
	fmt.Printf("  external_reference: %s\n", callback.ExternalReference)
	fmt.Printf("  amount: %d\n", callback.Amount)
	fmt.Printf("  platform_fee: %d\n", callback.PlatformFee)
	fmt.Printf("  merchant_points: %d\n", callback.MerchantPoints)
	fmt.Printf("  status: %s\n", callback.Status)
	fmt.Printf("  paid_at: %s\n", callback.PaidAt)
	fmt.Printf("  signature (接收): %s\n", callback.Signature)
	fmt.Printf("计算的签名: %s\n", expectedSignature)
	fmt.Printf("SecretKey: %s\n", cfg.SecretKey)

	isValid := hmac.Equal([]byte(expectedSignature), []byte(callback.Signature))
	fmt.Printf("签名验证结果: %v\n", isValid)
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

	resp, err := http.PostForm(apiURL, formData)
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

// ProcessPaymentCallback 处理支付回调
func (s *PaymentService) ProcessPaymentCallback(callback *PaymentCallback) error {
	// 1. 验证签名
	if !s.VerifyCallback(callback) {
		return fmt.Errorf("签名验证失败")
	}

	// 2. 查找订单
	orderService := NewOrderService()
	order, err := orderService.FindByOrderNo(callback.ExternalReference)
	if err != nil {
		return fmt.Errorf("订单不存在: %s", callback.ExternalReference)
	}

	// 3. 检查订单状态（防止重复处理）
	if order.Status != models.OrderStatusPending {
		// 订单已处理，直接返回成功（幂等）
		return nil
	}

	// 4. 验证金额
	expectedAmount := int(order.TotalAmount) // 假设1积分=1元
	if callback.Amount != expectedAmount {
		return fmt.Errorf("金额不匹配: 期望 %d, 实际 %d", expectedAmount, callback.Amount)
	}

	// 5. 检查支付状态
	if callback.Status != "completed" {
		return fmt.Errorf("支付未完成: %s", callback.Status)
	}

	// 6. 更新订单状态
	now := time.Now()
	order.Status = models.OrderStatusPaid
	order.PaidAt = &now
	order.PayMethod = "nodeloc"
	order.TransactionID = callback.TransactionID

	if err := database.GetDB().Save(order).Error; err != nil {
		return fmt.Errorf("更新订单失败: %w", err)
	}

	// 7. 分配卡密
	cardKeyService := NewCardKeyService()
	availableCards, err := cardKeyService.GetAvailableByProduct(order.ProductID, order.Quantity)
	if err != nil || len(availableCards) < order.Quantity {
		// 库存不足，记录日志但不影响支付状态
		// 实际项目中应该发送告警
		return nil
	}

	cardIDs := make([]uint, len(availableCards))
	for i, card := range availableCards {
		cardIDs[i] = card.ID
	}
	cardKeyService.MarkAsSold(cardIDs, order.ID)

	// 8. 更新商品库存和销量
	productService := NewProductService()
	productService.UpdateStock(order.ProductID)
	productService.IncrementSales(order.ProductID, order.Quantity)

	// 9. 标记订单完成
	order.Status = models.OrderStatusCompleted
	database.GetDB().Save(order)

	return nil
}
