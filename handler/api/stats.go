package api

import (
	"net/http"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/nodeloc-faka/database"
	"github.com/nodeloc-faka/models"
)

type publicStats struct {
	SoldKeys    int64 `json:"sold_keys"`
	PaidOrders  int64 `json:"paid_orders"`
	Products    int64 `json:"products"`
	Shops       int64 `json:"shops"`
}

type recentPurchase struct {
	Buyer        string     `json:"buyer"`
	ProductID    uint       `json:"product_id"`
	ProductName  string     `json:"product_name"`
	CategoryIcon string     `json:"category_icon"`
	PaidAt       *time.Time `json:"paid_at"`
}

// Homepage endpoints are hit on every visit, so results are cached briefly.
var (
	statsMu       sync.Mutex
	statsCache    *publicStats
	statsCachedAt time.Time

	recentMu       sync.Mutex
	recentCache    []recentPurchase
	recentCachedAt time.Time
)

var paidStatuses = []int{models.OrderStatusPaid, models.OrderStatusCompleted, models.OrderStatusShipped}

// GetPublicStats 首页公开统计
func (h *APIHandler) GetPublicStats(c *gin.Context) {
	statsMu.Lock()
	defer statsMu.Unlock()
	if statsCache == nil || time.Since(statsCachedAt) > time.Minute {
		db := database.GetDB()
		s := &publicStats{}
		db.Model(&models.CardKey{}).Where("status = ?", models.CardKeyStatusSold).Count(&s.SoldKeys)
		db.Model(&models.Order{}).Where("status IN ?", paidStatuses).Count(&s.PaidOrders)
		db.Model(&models.Product{}).
			Joins("JOIN shops ON shops.id = products.shop_id AND shops.status = ?", models.ShopStatusApproved).
			Where("products.is_active = ?", true).Count(&s.Products)
		db.Model(&models.Shop{}).Where("status = ?", models.ShopStatusApproved).Count(&s.Shops)
		statsCache = s
		statsCachedAt = time.Now()
	}
	c.JSON(http.StatusOK, statsCache)
}

// GetRecentPurchases 最近成交（用户名脱敏，不含金额与订单号）
func (h *APIHandler) GetRecentPurchases(c *gin.Context) {
	recentMu.Lock()
	defer recentMu.Unlock()
	if recentCache == nil || time.Since(recentCachedAt) > 30*time.Second {
		var orders []models.Order
		database.GetDB().
			Where("status IN ? AND paid_at IS NOT NULL", paidStatuses).
			Preload("User").Preload("Product").Preload("Product.Category").
			Order("paid_at desc").Limit(12).Find(&orders)

		list := make([]recentPurchase, 0, len(orders))
		for _, o := range orders {
			if o.Product == nil {
				continue
			}
			item := recentPurchase{
				ProductID:   o.Product.ID,
				ProductName: o.Product.Name,
				PaidAt:      o.PaidAt,
			}
			if o.User != nil {
				item.Buyer = maskName(o.User.Username)
			} else {
				item.Buyer = "匿名"
			}
			if o.Product.Category != nil {
				item.CategoryIcon = o.Product.Category.Icon
			}
			list = append(list, item)
		}
		recentCache = list
		recentCachedAt = time.Now()
	}
	c.JSON(http.StatusOK, gin.H{"items": recentCache})
}

func maskName(name string) string {
	n := utf8.RuneCountInString(name)
	if n == 0 {
		return "匿名"
	}
	r := []rune(name)
	if n <= 2 {
		return string(r[0]) + "***"
	}
	return string(r[0]) + "***" + string(r[n-1])
}
