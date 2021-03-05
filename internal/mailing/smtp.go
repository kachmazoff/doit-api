package mailing

import (
	"net"
	"net/mail"
	"net/smtp"
)

// SMTPSender sends emails using an SMTP service.
type SMTPSender struct {
	From mail.Address
	Addr string
	Auth smtp.Auth
}

type SMTPConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

// NewSMTPSender implementation using an SMTP server.
func NewSMTPSender(config SMTPConfig) *SMTPSender {
	return &SMTPSender{
		From: mail.Address{Address: config.From},
		Addr: net.JoinHostPort(config.Host, config.Port),
		Auth: smtp.PlainAuth("", config.Username, config.Password, config.Host),
	}
}

// Send an email to the given email address.
func (s *SMTPSender) Send(to, subject, body string) error {
	toAddr := mail.Address{Address: to}
	msg := message(s.From, toAddr, subject, body)

	return smtp.SendMail(
		s.Addr,
		s.Auth,
		s.From.Address,
		[]string{toAddr.Address},
		[]byte(msg),
	)
}
