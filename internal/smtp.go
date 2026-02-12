package mailer

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
)

type SMTPMailer struct {
	Host string
	Port string
	User string
	Pass string
	From string
}

func (m *SMTPMailer) Send(to, subject, body string) error {
	addr := fmt.Sprintf("%s:%s", m.Host, m.Port)

	// connect
	c, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer c.Close()

	// STARTTLS
	if err := c.StartTLS(&tls.Config{ServerName: m.Host}); err != nil {
		return err
	}

	// auth
	auth := smtp.PlainAuth("", m.User, m.Pass, m.Host)
	if err := c.Auth(auth); err != nil {
		return err
	}

	if err := c.Mail(m.From); err != nil {
		return err
	}
	if err := c.Rcpt(to); err != nil {
		return err
	}

	w, err := c.Data()
	if err != nil {
		return err
	}
	defer w.Close()

	msg := "" +
		"From: " + m.From + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/plain; charset=\"utf-8\"\r\n" +
		"\r\n" +
		body + "\r\n"

	_, err = w.Write([]byte(msg))
	return err
}
