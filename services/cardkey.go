package services

import (
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/nodeloc-faka/database"
	"github.com/nodeloc-faka/models"
	"gorm.io/gorm"
)

// CardKeyService 卡密服务
type CardKeyService struct{}

// NewCardKeyService 创建卡密服务
func NewCardKeyService() *CardKeyService {
	return &CardKeyService{}
}

// Create 创建单个卡密
func (s *CardKeyService) Create(cardKey *models.CardKey) error {
	err := database.GetDB().Create(cardKey).Error
	if err == nil {
		// 更新商品库存
		NewProductService().UpdateStock(cardKey.ProductID)
	}
	return err
}

// 批量导入限制
const (
	cardKeyMaxLen       = 500  // 与 card_no/card_pwd 列宽一致
	cardKeyMaxBatchRows = 5000 // 单次导入的最大行数
)

var (
	ErrCardKeysTooMany       = &ServiceError{Message: fmt.Sprintf("单次最多导入 %d 条卡密，请分批添加", cardKeyMaxBatchRows)}
	ErrCardKeysLookLikeCode  = &ServiceError{Message: "内容看起来像源码或文档，而不是一行一条的卡密。整段内容请把商品发货方式改为「固定内容自动发货」。"}
	ErrCardKeysFixedProduct  = &ServiceError{Message: "该商品为固定内容发货，不需要添加卡密"}
	ErrCardKeysManualProduct = &ServiceError{Message: "该商品为手动发货，卡密不会自动发出。如需自动发货，请先把发货方式改为「卡密自动发货」。"}
)

// CardKeyBatchResult 批量导入结果
type CardKeyBatchResult struct {
	Added      int `json:"count"`
	Duplicates int `json:"duplicates"`
	Invalid    int `json:"invalid"`
	Lines      int `json:"lines"`
}

var (
	codeOnlyPunct   = regexp.MustCompile(`^[\[\]{}()<>;,:]+$`)
	codeKeywordHead = regexp.MustCompile(`^(import|from|def|class|return|if|elif|else|for|while|try|except|finally|with|print|function|const|let|var|public|private|package|async|await|pass|break|continue|#include)\b`)
	codeShellHead   = regexp.MustCompile(`^(pip3?|npm|yarn|pnpm|python3?|node|sudo|git|curl|wget|chmod|docker|apt|yum)\s`)
	codeCommentHead = regexp.MustCompile("^(#|//|/\\*|\\*/|<!--|```)")
	codeTail        = regexp.MustCompile(`[{:;(,\\：]$`)
	codeAssign      = regexp.MustCompile(`\s(=|==|\+=|-=|=>)\s`)
)

// looksLikeCode 判断一行是否更像源码/配置而不是卡密
func looksLikeCode(line string) bool {
	return codeOnlyPunct.MatchString(line) ||
		codeKeywordHead.MatchString(line) ||
		codeShellHead.MatchString(line) ||
		codeCommentHead.MatchString(line) ||
		codeTail.MatchString(line) ||
		codeAssign.MatchString(line)
}

// batchLooksLikeCode 至少 5 行且 30% 以上的行像代码时，认为整段是源码/文档
func batchLooksLikeCode(lines []string) bool {
	codeLike := 0
	for _, line := range lines {
		if looksLikeCode(line) {
			codeLike++
		}
	}
	return codeLike >= 5 && codeLike*100 >= len(lines)*30
}

func parseCardLine(line string) (cardNo, cardPwd string) {
	if strings.Contains(line, "----") {
		parts := strings.SplitN(line, "----", 2)
		return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
	}
	return line, ""
}

// BatchCreate 批量创建卡密：每行一条，支持「卡号」或「卡号----密码」。
// 超长/空卡号计为无效，批内及库内已存在的计为重复，整体像源码时拒绝导入。
func (s *CardKeyService) BatchCreate(productID uint, cardsText string) (*CardKeyBatchResult, error) {
	db := database.GetDB()

	var product models.Product
	if err := db.Select("id", "delivery_type").First(&product, productID).Error; err != nil {
		return nil, ErrProductNotFound
	}
	switch product.DeliveryType {
	case models.DeliveryTypeFixed:
		return nil, ErrCardKeysFixedProduct
	case models.DeliveryTypeManual:
		return nil, ErrCardKeysManualProduct
	}

	var lines []string
	for _, raw := range strings.Split(strings.ReplaceAll(cardsText, "\r\n", "\n"), "\n") {
		if line := strings.TrimSpace(raw); line != "" {
			lines = append(lines, line)
		}
	}
	res := &CardKeyBatchResult{Lines: len(lines)}
	if len(lines) == 0 {
		return res, nil
	}
	if len(lines) > cardKeyMaxBatchRows {
		return nil, ErrCardKeysTooMany
	}

	if batchLooksLikeCode(lines) {
		return nil, ErrCardKeysLookLikeCode
	}

	type pair struct{ No, Pwd string }
	var existing []pair
	if err := db.Model(&models.CardKey{}).Select("card_no AS no, card_pwd AS pwd").
		Where("product_id = ?", productID).Scan(&existing).Error; err != nil {
		return nil, err
	}
	seen := make(map[pair]struct{}, len(existing)+len(lines))
	for _, e := range existing {
		seen[e] = struct{}{}
	}

	toCreate := make([]models.CardKey, 0, len(lines))
	for _, line := range lines {
		no, pwd := parseCardLine(line)
		if no == "" || utf8.RuneCountInString(no) > cardKeyMaxLen || utf8.RuneCountInString(pwd) > cardKeyMaxLen {
			res.Invalid++
			continue
		}
		k := pair{no, pwd}
		if _, dup := seen[k]; dup {
			res.Duplicates++
			continue
		}
		seen[k] = struct{}{}
		toCreate = append(toCreate, models.CardKey{
			ProductID: productID,
			CardNo:    no,
			CardPwd:   pwd,
			Status:    models.CardKeyStatusAvailable,
		})
	}

	if len(toCreate) > 0 {
		if err := db.CreateInBatches(toCreate, 200).Error; err != nil {
			return nil, err
		}
	}
	res.Added = len(toCreate)

	NewProductService().UpdateStock(productID)
	return res, nil
}

// Update 更新卡密
func (s *CardKeyService) Update(cardKey *models.CardKey) error {
	return database.GetDB().Save(cardKey).Error
}

// Delete 删除卡密
func (s *CardKeyService) Delete(id uint) error {
	var cardKey models.CardKey
	if err := database.GetDB().First(&cardKey, id).Error; err != nil {
		return err
	}

	if cardKey.Status == models.CardKeyStatusSold {
		return ErrCardKeySold
	}

	err := database.GetDB().Delete(&models.CardKey{}, id).Error
	if err == nil {
		// 更新商品库存
		NewProductService().UpdateStock(cardKey.ProductID)
	}
	return err
}

// BatchDelete 批量删除未售出的卡密
func (s *CardKeyService) BatchDelete(ids []uint) (int, error) {
	result := database.GetDB().
		Where("id IN ? AND status = ?", ids, models.CardKeyStatusAvailable).
		Delete(&models.CardKey{})
	return int(result.RowsAffected), result.Error
}

// FindByID 根据ID查找卡密
func (s *CardKeyService) FindByID(id uint) (*models.CardKey, error) {
	var cardKey models.CardKey
	if err := database.GetDB().Preload("Product").First(&cardKey, id).Error; err != nil {
		return nil, err
	}
	return &cardKey, nil
}

// GetByProduct 根据商品获取卡密
func (s *CardKeyService) GetByProduct(productID uint) ([]models.CardKey, error) {
	var cardKeys []models.CardKey
	if err := database.GetDB().
		Where("product_id = ?", productID).
		Order("status asc, id desc").
		Find(&cardKeys).Error; err != nil {
		return nil, err
	}
	return cardKeys, nil
}

// GetAvailableByProduct 获取商品的可售卡密
func (s *CardKeyService) GetAvailableByProduct(productID uint, limit int) ([]models.CardKey, error) {
	var cardKeys []models.CardKey
	if err := database.GetDB().
		Where("product_id = ? AND status = ?", productID, models.CardKeyStatusAvailable).
		Order("id asc").
		Limit(limit).
		Find(&cardKeys).Error; err != nil {
		return nil, err
	}
	return cardKeys, nil
}

// MarkAsSold 标记卡密为已售出
func (s *CardKeyService) MarkAsSold(ids []uint, orderID uint) error {
	return s.MarkAsSoldTx(database.GetDB(), ids, orderID)
}

// MarkAsSoldTx 在指定事务中标记卡密为已售出
func (s *CardKeyService) MarkAsSoldTx(tx *gorm.DB, ids []uint, orderID uint) error {
	now := time.Now()
	return tx.Model(&models.CardKey{}).
		Where("id IN ?", ids).
		Updates(map[string]interface{}{
			"status":   models.CardKeyStatusSold,
			"order_id": orderID,
			"sold_at":  now,
		}).Error
}

// GetWithPagination 分页获取卡密
func (s *CardKeyService) GetWithPagination(productID uint, page, pageSize int) ([]models.CardKey, int64, error) {
	var cardKeys []models.CardKey
	var total int64

	db := database.GetDB().Model(&models.CardKey{})
	if productID > 0 {
		db = db.Where("product_id = ?", productID)
	}
	db.Count(&total)

	offset := (page - 1) * pageSize
	if err := db.Preload("Product").
		Order("status asc, id desc").
		Offset(offset).
		Limit(pageSize).
		Find(&cardKeys).Error; err != nil {
		return nil, 0, err
	}

	return cardKeys, total, nil
}

// CountByProduct 获取商品的卡密数量
func (s *CardKeyService) CountByProduct(productID uint, status int) int64 {
	var count int64
	db := database.GetDB().Model(&models.CardKey{}).Where("product_id = ?", productID)
	if status >= 0 {
		db = db.Where("status = ?", status)
	}
	db.Count(&count)
	return count
}

// CountAll 获取所有卡密数量
func (s *CardKeyService) CountAll() int64 {
	var count int64
	database.GetDB().Model(&models.CardKey{}).Count(&count)
	return count
}

// CountAvailable 获取可售卡密数量
func (s *CardKeyService) CountAvailable() int64 {
	var count int64
	database.GetDB().Model(&models.CardKey{}).Where("status = ?", models.CardKeyStatusAvailable).Count(&count)
	return count
}

// 错误定义
var ErrCardKeySold = &ServiceError{Message: "该卡密已售出，无法删除"}
