package services

import (
	"errors"
	"fmt"
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
	return s.FindByOrderNo(orderNo)
}
