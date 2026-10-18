package services

import (
	"crypto/tls"
	"fmt"
	"html"
	"mime"
	"net"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/nodeloc-faka/models"
)

// EmailService 邮件服务
type EmailService struct {
	host     string
	port     int
	user     string
	pass     string
	from     string
	admin    string
	enabled  bool
}

var globalEmailService *EmailService

// NewEmailService 创建邮件服务（从环境变量读取配置）
func NewEmailService() *EmailService {
	if globalEmailService != nil {
		return globalEmailService
	}
	host := os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	user := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASS")
	from := os.Getenv("SMTP_FROM")
	admin := os.Getenv("ADMIN_EMAIL")

	port := 465
	if portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil {
			port = p
		}
	}
	if from == "" && user != "" {
		from = user
	}

	globalEmailService = &EmailService{
		host:    host,
		port:    port,
		user:    user,
		pass:    pass,
		from:    from,
		admin:   admin,
		enabled: host != "" && user != "" && pass != "",
	}
	return globalEmailService
}

// IsEnabled 是否已配置邮件
func (s *EmailService) IsEnabled() bool {
	return s.enabled
}

// send 发送邮件（内部方法，支持 SSL/TLS 465 和 STARTTLS 587）
func (s *EmailService) send(to []string, subject, body string) error {
	if !s.enabled {
		return nil
	}

	// 构建 MIME 邮件
	header := make(map[string]string)
	header["From"] = s.from
	header["To"] = strings.Join(to, ", ")
	header["Subject"] = mime.QEncoding.Encode("UTF-8", subject)
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=UTF-8"
	header["Date"] = time.Now().Format(time.RFC1123Z)

	var msg strings.Builder
	for k, v := range header {
		msg.WriteString(k + ": " + v + "\r\n")
	}
	msg.WriteString("\r\n" + body)

	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	auth := smtp.PlainAuth("", s.user, s.pass, s.host)

	// 465 端口用 SSL/TLS，587 用 STARTTLS
	if s.port == 465 {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: false,
			ServerName:         s.host,
		}
		conn, err := tls.Dial("tcp", addr, tlsConfig)
		if err != nil {
			return fmt.Errorf("smtp tls dial: %w", err)
		}
		defer conn.Close()

		client, err := smtp.NewClient(conn, s.host)
		if err != nil {
			return fmt.Errorf("smtp new client: %w", err)
		}
		defer client.Close()

		if err = client.Auth(auth); err != nil {
			return fmt.Errorf("smtp auth: %w", err)
		}
		if err = client.Mail(s.user); err != nil {
			return err
		}
		for _, addr := range to {
			if err = client.Rcpt(addr); err != nil {
				return err
			}
		}
		w, err := client.Data()
		if err != nil {
			return err
		}
		_, err = w.Write([]byte(msg.String()))
		if err != nil {
			return err
		}
		return w.Close()
	}

	// STARTTLS (587 / 25)
	host, _, _ := net.SplitHostPort(addr)
	_ = host
	return smtp.SendMail(addr, auth, s.user, to, []byte(msg.String()))
}

// sendAsync 异步发送（不阻塞主流程）
func (s *EmailService) sendAsync(to []string, subject, body string) {
	if !s.enabled {
		return
	}
	go func() {
		if err := s.send(to, subject, body); err != nil {
			fmt.Printf("[Email] 发送失败 to=%v subject=%s err=%v\n", to, subject, err)
		}
	}()
}

// ============================================================
// 业务邮件模板
// ============================================================

// SendOrderCompleted 订单完成 → 发送卡密给买家
func (s *EmailService) SendOrderCompleted(order *models.Order) {
	if !s.enabled || order.User == nil || order.User.Email == "" {
		return
	}

	var cardLines strings.Builder
	for _, ck := range order.CardKeys {
		if ck.CardPwd != "" {
			cardLines.WriteString(fmt.Sprintf("<li>卡号：<b>%s</b> &nbsp; 密码：<b>%s</b></li>", html.EscapeString(ck.CardNo), html.EscapeString(ck.CardPwd)))
		} else {
			cardLines.WriteString(fmt.Sprintf("<li>卡号：<b>%s</b></li>", html.EscapeString(ck.CardNo)))
		}
	}

	productName := ""
	if order.Product != nil {
		productName = order.Product.Name
	}

	intro, heading := "以下是您购买的卡密：", "卡密信息"
	deliverBlock := fmt.Sprintf(`<ul style="background:#f9f9f9;padding:16px;border-radius:8px">%s</ul>`, cardLines.String())
	if len(order.CardKeys) == 0 && order.DeliverContent != "" {
		intro, heading = "以下是您购买的商品内容：", "商品内容"
		deliverBlock = fmt.Sprintf(`<pre style="background:#f9f9f9;padding:16px;border-radius:8px;white-space:pre-wrap;word-break:break-all;font-family:monospace">%s</pre>`, html.EscapeString(order.DeliverContent))
	}

	body := fmt.Sprintf(`
<div style="font-family:sans-serif;max-width:600px;margin:0 auto">
  <h2 style="color:#e6162d">订单完成通知</h2>
  <p>您好，您的订单已完成，%s</p>
  <table style="width:100%%;border-collapse:collapse;margin:16px 0">
    <tr><td style="padding:8px;background:#f5f5f5;font-weight:bold">订单号</td><td style="padding:8px">%s</td></tr>
    <tr><td style="padding:8px;background:#f5f5f5;font-weight:bold">商品</td><td style="padding:8px">%s</td></tr>
    <tr><td style="padding:8px;background:#f5f5f5;font-weight:bold">数量</td><td style="padding:8px">%d</td></tr>
    <tr><td style="padding:8px;background:#f5f5f5;font-weight:bold">金额</td><td style="padding:8px">%g 能量</td></tr>
  </table>
  <h3>%s</h3>
  %s
  <p style="color:#888;font-size:12px">请妥善保管，如有问题请联系店铺客服。</p>
</div>`, intro, html.EscapeString(order.OrderNo), html.EscapeString(productName), order.Quantity, order.TotalAmount, heading, deliverBlock)

	s.sendAsync([]string{order.User.Email}, "订单完成 - 商品已发放", body)
}

// SendOrderShipped 手动发货通知 → 发给买家
func (s *EmailService) SendOrderShipped(order *models.Order) {
	if !s.enabled || order.User == nil || order.User.Email == "" {
		return
	}

	productName := ""
	if order.Product != nil {
		productName = order.Product.Name
	}

	body := fmt.Sprintf(`
<div style="font-family:sans-serif;max-width:600px;margin:0 auto">
  <h2 style="color:#e6162d">发货通知</h2>
  <p>您好，您的订单已发货，请查收：</p>
  <table style="width:100%%;border-collapse:collapse;margin:16px 0">
    <tr><td style="padding:8px;background:#f5f5f5;font-weight:bold">订单号</td><td style="padding:8px">%s</td></tr>
    <tr><td style="padding:8px;background:#f5f5f5;font-weight:bold">商品</td><td style="padding:8px">%s</td></tr>
  </table>
  <h3>发货内容</h3>
  <div style="background:#f9f9f9;padding:16px;border-radius:8px;white-space:pre-wrap">%s</div>
  %s
  <p style="margin-top:16px">收到商品后请在订单页面确认收货，确认后款项将结算给卖家。</p>
</div>`,
		html.EscapeString(order.OrderNo), html.EscapeString(productName), html.EscapeString(order.DeliverContent),
		func() string {
			if order.DeliverNote != "" {
				return fmt.Sprintf(`<p><b>备注：</b>%s</p>`, html.EscapeString(order.DeliverNote))
			}
			return ""
		}())

	s.sendAsync([]string{order.User.Email}, "发货通知 - "+productName, body)
}

// SendNewOrderToShop 新订单通知 → 发给店主
func (s *EmailService) SendNewOrderToShop(order *models.Order, shopOwnerEmail string) {
	if !s.enabled || shopOwnerEmail == "" {
		return
	}

	productName := ""
	if order.Product != nil {
		productName = order.Product.Name
	}
	buyerName := ""
	if order.User != nil {
		buyerName = order.User.Username
	}

	body := fmt.Sprintf(`
<div style="font-family:sans-serif;max-width:600px;margin:0 auto">
  <h2 style="color:#e6162d">新订单通知</h2>
  <p>您的店铺收到了一笔新订单：</p>
  <table style="width:100%%;border-collapse:collapse;margin:16px 0">
    <tr><td style="padding:8px;background:#f5f5f5;font-weight:bold">订单号</td><td style="padding:8px">%s</td></tr>
    <tr><td style="padding:8px;background:#f5f5f5;font-weight:bold">商品</td><td style="padding:8px">%s</td></tr>
    <tr><td style="padding:8px;background:#f5f5f5;font-weight:bold">数量</td><td style="padding:8px">%d</td></tr>
    <tr><td style="padding:8px;background:#f5f5f5;font-weight:bold">金额</td><td style="padding:8px">%g 能量</td></tr>
    <tr><td style="padding:8px;background:#f5f5f5;font-weight:bold">买家</td><td style="padding:8px">%s</td></tr>
    <tr><td style="padding:8px;background:#f5f5f5;font-weight:bold">联系方式</td><td style="padding:8px">%s</td></tr>
  </table>
  <p>请登录后台处理订单。</p>
</div>`, html.EscapeString(order.OrderNo), html.EscapeString(productName), order.Quantity, order.TotalAmount,
		html.EscapeString(buyerName), html.EscapeString(order.Contact))

	s.sendAsync([]string{shopOwnerEmail}, "新订单通知 - "+productName, body)
}

// SendNewsletterConfirm 订阅确认邮件
func (s *EmailService) SendNewsletterConfirm(email, siteName string) {
	if !s.enabled || email == "" {
		return
	}
	body := fmt.Sprintf(`
<div style="font-family:sans-serif;max-width:600px;margin:0 auto">
  <h2 style="color:#e6162d">订阅成功</h2>
  <p>感谢您订阅 <b>%s</b> 的最新资讯！</p>
  <p>我们将不定期向您发送最新商品、优惠活动等信息。</p>
  <p style="color:#888;font-size:12px">如果您没有订阅过我们的邮件，请忽略此邮件。</p>
</div>`, html.EscapeString(siteName))

	s.sendAsync([]string{email}, "订阅成功 - "+siteName, body)
}

// SendWithdrawalNotify 提现申请通知 → 发给管理员
func (s *EmailService) SendWithdrawalNotify(req *models.WithdrawalRequest) {
	if !s.enabled || s.admin == "" {
		return
	}

	username := ""
	if req.User != nil {
		username = req.User.Username
	}

	body := fmt.Sprintf(`
<div style="font-family:sans-serif;max-width:600px;margin:0 auto">
  <h2 style="color:#e6162d">提现申请通知</h2>
  <p>有新的提现申请待审核：</p>
  <table style="width:100%%;border-collapse:collapse;margin:16px 0">
    <tr><td style="padding:8px;background:#f5f5f5;font-weight:bold">申请人</td><td style="padding:8px">%s</td></tr>
    <tr><td style="padding:8px;background:#f5f5f5;font-weight:bold">申请金额</td><td style="padding:8px">%g 能量</td></tr>
    <tr><td style="padding:8px;background:#f5f5f5;font-weight:bold">实际到账</td><td style="padding:8px">%g 能量</td></tr>
    <tr><td style="padding:8px;background:#f5f5f5;font-weight:bold">备注</td><td style="padding:8px">%s</td></tr>
  </table>
  <p>请登录后台审核。</p>
</div>`, html.EscapeString(username), req.Amount, req.ActualAmount, html.EscapeString(req.Remark))

	s.sendAsync([]string{s.admin}, "提现申请待审核", body)
}
