package services

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/nodeloc-faka/database"
	"github.com/nodeloc-faka/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// AdminRefund 后台全额退款，退款后订单置为已取消。
//
// 能量支付优先走 NodeLoc 原路退款（无手续费损失）；原路退款接口未上线时改为从商户账户转账退回
// （商户承担转账手续费）。余额支付退回买家余额。已结算给店主的收入会一并扣回。
//
// 整个过程在锁住订单行的事务里完成：并发的重复点击会串行化，第二次看到 refunded_at 已设置直接拒绝；
// 远程退款/转账本身也是幂等的，事务提交失败后重试不会重复退钱。
func (s *OrderService) AdminRefund(orderNo string) (*models.Order, error) {
	paymentService := NewPaymentService()
	balanceService := NewBalanceService()

	err := database.GetDB().Transaction(func(tx *gorm.DB) error {
		var order models.Order
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("order_no = ?", orderNo).First(&order).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return &ServiceError{Message: "订单不存在"}
			}
			return err
		}
		if order.RefundedAt != nil {
			return &ServiceError{Message: "订单已退款"}
		}
		if order.Status == models.OrderStatusPending || order.PaidAt == nil {
			return &ServiceError{Message: "订单未付款，无需退款"}
		}
		if order.PayMethod == "free" || order.TotalAmount <= 0 {
			return &ServiceError{Message: "免费订单无需退款"}
		}

		// 先扣回已结算给店主的收入：余额不足就中止，此时还没动任何钱
		if order.ShopSettled && order.ShopIncome > 0 {
			var shop models.Shop
			err := tx.First(&shop, order.ShopID).Error
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			if err == nil && !shop.IsOfficial && shop.UserID > 0 {
				if err := balanceService.DeductBalance(tx, shop.UserID, order.ShopIncome,
					models.BalanceTxRefundReversal,
					fmt.Sprintf("订单 %s 退款，扣回销售收入", order.OrderNo),
					"order", order.ID,
				); err != nil {
					if errors.Is(err, ErrInsufficientBalance) {
						return &ServiceError{Message: fmt.Sprintf(
							"店主余额不足以扣回该订单收入 %.0f（可能已提现），请先处理店主余额", order.ShopIncome)}
					}
					return err
				}
			}
		}

		amount := order.TotalAmount
		var method, txID string
		switch order.PayMethod {
		case "balance":
			if err := balanceService.AddBalance(tx, order.UserID, amount,
				models.BalanceTxOrderRefund,
				fmt.Sprintf("订单 %s 退款", order.OrderNo),
				"order", order.ID,
			); err != nil {
				return err
			}
			method = "balance"

		case "nodeloc":
			if order.TransactionID == "" {
				return &ServiceError{Message: "订单缺少支付交易号，无法退款"}
			}
			q, err := paymentService.QueryPayment(order.TransactionID)
			if err != nil {
				return &ServiceError{Message: "查询支付状态失败: " + err.Error()}
			}
			if q.ExternalReference != order.OrderNo {
				return &ServiceError{Message: "支付交易与订单不匹配，已中止退款"}
			}
			amount = float64(q.Amount)

			switch q.Status {
			case "refunded": // 之前已原路退款（例如上次退款后事务未提交），只需补记
				method, txID = "native", q.TransactionID
			case "completed":
				r, err := paymentService.RefundPayment(order.TransactionID, q.Amount)
				switch {
				case err == nil:
					method, txID = "native", r.TransactionID
				case errors.Is(err, ErrRefundUnsupported):
					var buyer models.User
					if err := tx.First(&buyer, order.UserID).Error; err != nil {
						return &ServiceError{Message: "找不到买家信息"}
					}
					if buyer.NodeLocID == 0 {
						return &ServiceError{Message: "买家未绑定 NodeLoc 账号，无法转账退款"}
					}
					tr, err := paymentService.TransferRefund(order.OrderNo, buyer.NodeLocID, buyer.Username, q.Amount)
					if err != nil {
						return &ServiceError{Message: "转账退款失败: " + err.Error()}
					}
					method, txID = "transfer", tr.TransactionID
				default:
					return &ServiceError{Message: err.Error()}
				}
			default:
				return &ServiceError{Message: fmt.Sprintf("支付状态为 %s，无法退款", q.Status)}
			}

		default:
			return &ServiceError{Message: "不支持的支付方式: " + order.PayMethod}
		}

		// 手动发货商品付款时已扣库存，未发货就退款则补回
		if order.Status == models.OrderStatusPaid {
			var product models.Product
			if err := tx.First(&product, order.ProductID).Error; err == nil && product.DeliveryType == models.DeliveryTypeManual {
				if err := tx.Model(&models.Product{}).Where("id = ?", product.ID).
					UpdateColumn("stock_count", gorm.Expr("stock_count + ?", order.Quantity)).Error; err != nil {
					return err
				}
			}
		}

		updates := map[string]interface{}{
			"status":        models.OrderStatusCancelled,
			"refunded_at":   time.Now(),
			"refund_amount": amount,
			"refund_method": method,
			"refund_tx_id":  txID,
		}
		if order.ShopSettled {
			updates["shop_income"] = 0 // 已扣回；保持 shop_settled=true 防止再次结算
		}
		return tx.Model(&models.Order{}).Where("id = ?", order.ID).Updates(updates).Error
	})
	if err != nil {
		return nil, err
	}
	order, err := s.FindByOrderNo(orderNo)
	if err == nil {
		NewEmailService().SendOrderRefunded(order)
	}
	return order, err
}

// AdminSyncPayment 后台向 NodeLoc 查询未付款订单的支付状态，已付款则走与支付回调完全相同的处理
// （待支付订单发卡/结算；已取消订单记下付款等待退款）。替代原来只改状态的「标记已支付」。
func (s *OrderService) AdminSyncPayment(orderNo string) (*models.Order, string, error) {
	order, err := s.FindByOrderNo(orderNo)
	if err != nil {
		return nil, "", ErrOrderNotFound
	}
	if order.RefundedAt != nil {
		return nil, "", &ServiceError{Message: "订单已退款"}
	}
	if order.PaidAt != nil {
		return nil, "", &ServiceError{Message: "订单已记录付款"}
	}
	if order.PayMethod != models.PayMethodNodeLoc || order.TransactionID == "" {
		return nil, "", &ServiceError{Message: "该订单没有 NodeLoc 支付交易"}
	}
	if order.Status != models.OrderStatusPending && order.Status != models.OrderStatusCancelled {
		return nil, "", &ServiceError{Message: "当前订单状态无需同步支付"}
	}

	q, err := NewPaymentService().QueryPayment(order.TransactionID)
	if err != nil {
		return nil, "", &ServiceError{Message: "查询支付状态失败: " + err.Error()}
	}
	if q.ExternalReference != order.OrderNo {
		return nil, "", &ServiceError{Message: "支付交易与订单不匹配"}
	}
	if q.Status != "completed" {
		return order, fmt.Sprintf("NodeLoc 支付状态：%s，订单未付款", q.Status), nil
	}

	if _, err := s.ProcessPaymentCallback(q.TransactionID, q.Amount, q.PlatformFee, q.MerchantPoints); err != nil {
		order, _ = s.FindByOrderNo(orderNo)
		return order, "已确认付款，但处理失败：" + err.Error(), nil
	}
	order, err = s.FindByOrderNo(orderNo)
	if err != nil {
		return nil, "", err
	}
	if order.Status == models.OrderStatusCancelled {
		return order, "已确认付款，但订单已取消，请退款", nil
	}
	return order, "已确认付款并完成处理", nil
}

// recordPaymentOnCancelled 已取消的订单收到了付款（例如取消后买家仍在旧支付页付款）：
// 不复活订单，只记下付款信息让它出现在后台「待退款」里，并通知管理员。
func (s *OrderService) recordPaymentOnCancelled(order *models.Order, platformFee, merchantPoints int) {
	if order.PaidAt != nil || order.RefundedAt != nil {
		return
	}
	res := database.GetDB().Model(&models.Order{}).
		Where("id = ? AND status = ? AND paid_at IS NULL AND refunded_at IS NULL", order.ID, models.OrderStatusCancelled).
		Updates(map[string]interface{}{
			"paid_at":         time.Now(),
			"platform_fee":    platformFee,
			"merchant_points": merchantPoints,
		})
	if res.Error != nil {
		log.Printf("订单 %s 已取消但收到付款，记录失败: %v", order.OrderNo, res.Error)
		return
	}
	if res.RowsAffected == 1 {
		log.Printf("订单 %s 已取消但收到付款 %.0f，需退款", order.OrderNo, order.TotalAmount)
		notifyOrderNeedsAction(order, "订单已取消，但买家随后完成了付款。请在后台订单管理中退款。", false)
	}
}

// settleShopTx 把订单收入结算给店主（非官方店、未结算过），与订单状态更新在同一事务内
func settleShopTx(tx *gorm.DB, order *models.Order, desc string) error {
	if order.ShopID == 0 || order.ShopSettled {
		return nil
	}
	var shop models.Shop
	if err := tx.First(&shop, order.ShopID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	if shop.IsOfficial || shop.UserID == 0 {
		return nil
	}
	income := order.TotalAmount
	if err := NewBalanceService().AddBalance(tx, shop.UserID, income,
		models.BalanceTxSaleIncome, desc, "order", order.ID); err != nil {
		return err
	}
	return tx.Model(&models.Order{}).Where("id = ?", order.ID).Updates(map[string]interface{}{
		"shop_settled": true,
		"shop_income":  income,
	}).Error
}

// shopOwnerEmail 非官方店店主邮箱，没有则返回空
func shopOwnerEmail(shopID uint) string {
	if shopID == 0 {
		return ""
	}
	var shop models.Shop
	if err := database.GetDB().First(&shop, shopID).Error; err != nil || shop.IsOfficial || shop.UserID == 0 {
		return ""
	}
	var owner models.User
	if err := database.GetDB().First(&owner, shop.UserID).Error; err != nil {
		return ""
	}
	return owner.Email
}

// notifyShopNewOrder 手动发货订单付款后通知店主发货
func notifyShopNewOrder(order *models.Order) {
	if email := shopOwnerEmail(order.ShopID); email != "" {
		NewEmailService().SendNewOrderToShop(order, email)
	}
}

// notifyOrderNeedsAction 已付款订单需要人工处理时通知管理员（以及店主）
func notifyOrderNeedsAction(order *models.Order, reason string, includeShop bool) {
	var extra []string
	if includeShop {
		if email := shopOwnerEmail(order.ShopID); email != "" {
			extra = append(extra, email)
		}
	}
	NewEmailService().SendOrderNeedsAction(order, reason, extra)
}
