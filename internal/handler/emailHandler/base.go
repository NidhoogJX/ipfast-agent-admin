package emailhandler

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"sync"
)

// EmailConfig 包含发送邮件所需的配置
type EmailConfig struct {
	FromEmail string
	To        []string
	Subject   string
	Body      string
}

// SMTPClientPool 连接池
type SMTPClientPool struct {
	mu      sync.Mutex
	clients []*smtp.Client
	config  *tls.Config
	auth    smtp.Auth
	host    string
	port    string
}

var emailPool *SMTPClientPool

// NewSMTPClientPool 创建一个新的连接池
func NewSMTPClientPool(host, port, senderEmail, password string, poolSize int) error {
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	auth := smtp.PlainAuth("", senderEmail, password, host)

	emailPool = &SMTPClientPool{
		clients: make([]*smtp.Client, 0, poolSize),
		config:  tlsconfig,
		auth:    auth,
		host:    host,
		port:    port,
	}

	for i := 0; i < poolSize; i++ {
		client, err := emailPool.newClient()
		if err != nil {
			return err
		}
		emailPool.clients = append(emailPool.clients, client)
	}

	return nil
}

// newClient 创建一个新的 SMTP 客户端
func (p *SMTPClientPool) newClient() (*smtp.Client, error) {
	conn, err := tls.Dial("tcp", p.host+":"+p.port, p.config)
	if err != nil {
		return nil, fmt.Errorf("failed to dial TLS: %v", err)
	}

	client, err := smtp.NewClient(conn, p.host)
	if err != nil {
		return nil, fmt.Errorf("failed to create SMTP client: %v", err)
	}

	if err = client.Auth(p.auth); err != nil {
		return nil, fmt.Errorf("failed to authenticate: %v", err)
	}

	return client, nil
}

// GetClient 从连接池中获取一个 SMTP 客户端
func (p *SMTPClientPool) GetClient() (*smtp.Client, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for len(p.clients) > 0 {
		client := p.clients[len(p.clients)-1]
		p.clients = p.clients[:len(p.clients)-1]

		// 健康检查
		if err := client.Noop(); err == nil {
			return client, nil
		}

		// 如果客户端不可用，关闭并丢弃
		client.Close()
	}

	// 如果没有可用的客户端，创建一个新的
	return p.newClient()
}

// ReturnClient 将 SMTP 客户端归还到连接池
func (p *SMTPClientPool) ReturnClient(client *smtp.Client) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.clients = append(p.clients, client)
}

// SendEmail 发送邮件
func SendEmail(config EmailConfig) error {
	client, err := emailPool.GetClient()
	if err != nil {
		return err
	}
	defer emailPool.ReturnClient(client)

	// 邮件内容
	message := []byte("From: " + config.FromEmail + "\n" +
		"Subject: " + config.Subject + "\n" +
		"\n" + config.Body)

	// 设置发件人和收件人
	if err = client.Mail(config.FromEmail); err != nil {
		return fmt.Errorf("failed to set sender: %v", err)
	}
	for _, addr := range config.To {
		if err = client.Rcpt(addr); err != nil {
			return fmt.Errorf("failed to set recipient: %v", err)
		}
	}

	// 发送邮件数据
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to get data writer: %v", err)
	}

	_, err = w.Write(message)
	if err != nil {
		return fmt.Errorf("failed to write message: %v", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("failed to close writer: %v", err)
	}

	fmt.Println("Email sent successfully!")
	return nil
}

// AsyncSendEmail 异步发送邮件
func AsyncSendEmail(config EmailConfig) {
	go func() {
		if err := SendEmail(config); err != nil {
			fmt.Printf("Failed to send email: %v\n", err)
		}
	}()
}
