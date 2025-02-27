package test

import (
	"crypto/tls"
	"testing"

	"gopkg.in/gomail.v2"
)

func TestGoMail(t *testing.T) {
	// SMTP 服务器信息
	smtpHost := "smtp.gmail.com"
	smtpPort := 465
	smtpUser := "phpbest64@gmail.com"
	smtpPass := "kzkrpaxfpvclchbk"
	fromEmail := "notice@wsslink.com"
	// 创建一个新的邮件对象
	m := gomail.NewMessage()

	// 设置发件人地址（可以是别名）
	m.SetHeader("From", fromEmail)

	// 设置收件人地址
	m.SetHeader("To", "1214939285@qq.com")

	// 设置邮件主题
	m.SetHeader("Subject", "Test Email")

	// 设置邮件正文
	m.SetBody("text/plain", "This is a test email sent using Go and Gmail SMTP server.")

	// 使用 SMTP 服务器发送邮件
	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)

	// 设置 TLS 配置
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}

	println("Email sent successfully!")
}
