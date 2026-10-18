package services

import (
	"github.com/nodeloc-faka/database"
	"github.com/nodeloc-faka/models"
)

// ReviewService 评价服务
type ReviewService struct{}

// NewReviewService 创建评价服务
func NewReviewService() *ReviewService {
	return &ReviewService{}
}

// ReviewStats 评价统计
type ReviewStats struct {
	Total            int                `json:"total"`
	AvgRating        float64            `json:"avg_rating"`
	RatingDistribution map[int]int      `json:"rating_distribution"`
}

// ReviewResponse 评价列表响应
type ReviewResponse struct {
	Reviews []models.ProductReview `json:"reviews"`
	Stats   ReviewStats            `json:"stats"`
}

// GetByProduct 获取商品评价列表（含用户信息）
func (s *ReviewService) GetByProduct(productID uint) (ReviewResponse, error) {
	db := database.GetDB()

	var reviews []models.ProductReview
	if err := db.Preload("User", PublicUserColumns).
		Where("product_id = ?", productID).
		Order("created_at DESC").
		Find(&reviews).Error; err != nil {
		return ReviewResponse{}, err
	}

	// 计算统计
	stats := ReviewStats{
		Total:              len(reviews),
		RatingDistribution: map[int]int{1: 0, 2: 0, 3: 0, 4: 0, 5: 0},
	}
	if len(reviews) > 0 {
		sum := 0
		for _, r := range reviews {
			sum += r.Rating
			if r.Rating >= 1 && r.Rating <= 5 {
				stats.RatingDistribution[r.Rating]++
			}
		}
		stats.AvgRating = float64(sum) / float64(len(reviews))
	}

	return ReviewResponse{Reviews: reviews, Stats: stats}, nil
}

// CanReview 检查用户是否有权评价（已完成订单）
func (s *ReviewService) CanReview(userID, productID uint) bool {
	db := database.GetDB()
	var count int64
	db.Model(&models.Order{}).
		Where("user_id = ? AND product_id = ? AND status = ?", userID, productID, models.OrderStatusCompleted).
		Count(&count)
	return count > 0
}

// HasReviewed 检查用户是否已评价
func (s *ReviewService) HasReviewed(userID, productID uint) bool {
	db := database.GetDB()
	var count int64
	db.Model(&models.ProductReview{}).
		Where("user_id = ? AND product_id = ?", userID, productID).
		Count(&count)
	return count > 0
}

// Create 创建评价
func (s *ReviewService) Create(review *models.ProductReview) error {
	db := database.GetDB()
	return db.Create(review).Error
}

// RatingStats 单个商品评分统计
type RatingStats struct {
	AvgRating float64
	Total     int
}

// GetAvgRatingByProducts 批量获取多个商品的评分统计
func (s *ReviewService) GetAvgRatingByProducts(productIDs []uint) map[uint]RatingStats {
	result := make(map[uint]RatingStats)
	if len(productIDs) == 0 {
		return result
	}

	type row struct {
		ProductID uint
		Avg       float64
		Total     int
	}
	var rows []row
	database.GetDB().Model(&models.ProductReview{}).
		Select("product_id, COALESCE(AVG(rating), 0) as avg, COUNT(*) as total").
		Where("product_id IN ?", productIDs).
		Group("product_id").
		Scan(&rows)

	for _, r := range rows {
		result[r.ProductID] = RatingStats{AvgRating: r.Avg, Total: r.Total}
	}
	return result
}
type ShopReviewStats struct {
	Total              int            `json:"total"`
	AvgRating          float64        `json:"avg_rating"`
	RatingDistribution map[int]int    `json:"rating_distribution"`
}

// GetShopReviewStats 获取店铺所有商品的评价统计
func (s *ReviewService) GetShopReviewStats(shopID uint) (*ShopReviewStats, error) {
	db := database.GetDB()

	stats := &ShopReviewStats{
		RatingDistribution: map[int]int{1: 0, 2: 0, 3: 0, 4: 0, 5: 0},
	}

	// 获取该店铺所有商品 ID
	var productIDs []uint
	if err := db.Model(&models.Product{}).Where("shop_id = ?", shopID).Pluck("id", &productIDs).Error; err != nil {
		return stats, err
	}
	if len(productIDs) == 0 {
		return stats, nil
	}

	// 总数和平均分
	var result struct {
		Total int
		Avg   float64
	}
	db.Model(&models.ProductReview{}).
		Where("product_id IN ?", productIDs).
		Select("COUNT(*) as total, COALESCE(AVG(rating), 0) as avg").
		Scan(&result)
	stats.Total = result.Total
	stats.AvgRating = result.Avg

	// 评分分布
	type ratingCount struct {
		Rating int
		Count  int
	}
	var counts []ratingCount
	db.Model(&models.ProductReview{}).
		Where("product_id IN ?", productIDs).
		Select("rating, COUNT(*) as count").
		Group("rating").
		Scan(&counts)
	for _, rc := range counts {
		if rc.Rating >= 1 && rc.Rating <= 5 {
			stats.RatingDistribution[rc.Rating] = rc.Count
		}
	}

	return stats, nil
}

// GetByShop 分页获取店铺所有商品的评价
func (s *ReviewService) GetByShop(shopID uint, page, pageSize, ratingFilter int) ([]models.ProductReview, int64, error) {
	db := database.GetDB()

	// 获取该店铺所有商品 ID
	var productIDs []uint
	if err := db.Model(&models.Product{}).Where("shop_id = ?", shopID).Pluck("id", &productIDs).Error; err != nil {
		return nil, 0, err
	}
	if len(productIDs) == 0 {
		return []models.ProductReview{}, 0, nil
	}

	query := db.Model(&models.ProductReview{}).Where("product_id IN ?", productIDs)
	if ratingFilter >= 1 && ratingFilter <= 5 {
		query = query.Where("rating = ?", ratingFilter)
	}

	var total int64
	query.Count(&total)

	var reviews []models.ProductReview
	offset := (page - 1) * pageSize
	if err := query.Preload("User", PublicUserColumns).Preload("Product").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&reviews).Error; err != nil {
		return nil, 0, err
	}

	return reviews, total, nil
}
