package test

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"testing"
)

func TestEmailSend(t *testing.T) {
	// SMTP 服务器配置
	smtpHost := "smtp.gmail.com"
	smtpPort := "465"
	senderEmail := "phpbest64@gmail.com"
	password := "kzkrpaxfpvclchbk"
	fromEmail := "notice@wsslink.com"
	// 收件人
	to := []string{"1214939285@qq.com"}

	// 邮件内容
	subject := "Subject: Test Email\n"
	body := "This is a test email from Go."
	message := []byte("From: " + fromEmail + "\"n" + subject + "\n" + body)

	// 认证
	auth := smtp.PlainAuth("", senderEmail, password, smtpHost)

	// TLS 配置
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpHost,
	}

	// 建立连接
	conn, err := tls.Dial("tcp", smtpHost+":"+smtpPort, tlsconfig)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer client.Quit()

	// 认证
	if err = client.Auth(auth); err != nil {
		fmt.Println("Error:", err)
		return
	}
	// 设置发件人和收件人
	if err = client.Mail(fromEmail); err != nil {
		fmt.Println("Error:", err)
		return
	}
	for _, addr := range to {
		if err = client.Rcpt(addr); err != nil {
			fmt.Println("Error:", err)
			return
		}
	}

	// 发送邮件数据
	w, err := client.Data()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	_, err = w.Write(message)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	err = w.Close()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Email sent successfully!")
}
