package daemon

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/logpacker/mailer/pkg/shared"
	"net/mail"
	"net/smtp"
	"os"
	"time"
)

// SMTPClient struct
type SMTPClient struct {
	Addr string
	Host string
}

// SMTPEmail struct
type SMTPEmail struct {
	From    mail.Address
	To      mail.Address
	Headers map[string]string
	Body    string
}

// BuildSMTPEmail func
func BuildSMTPEmail(email *shared.Email) *SMTPEmail {
	e := new(SMTPEmail)
	e.From = mail.Address{
		Address: email.From.Email,
		Name:    email.From.Name,
	}
	e.To = mail.Address{
		Address: email.To.Email,
		Name:    email.To.Name,
	}
	e.Headers = make(map[string]string)
	e.Headers["From"] = e.To.String()
	e.Headers["Reply-To"] = e.To.String()
	e.Headers["Return-Path"] = e.To.String()
	e.Headers["To"] = e.To.String()
	e.Headers["Subject"] = email.Subject
	e.Headers["MIME-Version"] = "1.0"
	e.Headers["Content-Type"] = "text/html; charset=\"utf-8\""
	e.Headers["Content-Transfer-Encoding"] = "base64"
	e.Headers["Date"] = time.Now().Format(time.RFC1123)
	e.Headers["Delivery-Date"] = time.Now().Format(time.RFC1123)
	e.Headers["Received"], _ = os.Hostname()
	e.Headers["Message-Id"] = fmt.Sprintf("%d", email.ID)
	e.Headers["List-Unsubscribe"] = "<" + email.URLUnsubscribe + ">"

	// 1px tracker image
	e.Body += "<img src=\"\" title=\"tracker-image\" width=1 height=1>"

	return e
}

// Init func
func (s *SMTPClient) Init(addr string) {
	s.Addr = addr
}

// Send func
func (s *SMTPClient) Send(smtpEmail *SMTPEmail) error {
	body := ""
	for k, v := range smtpEmail.Headers {
		body += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	body += fmt.Sprintf("\r\n%s", base64.StdEncoding.EncodeToString([]byte(smtpEmail.Body)))

	c, dialErr := smtp.Dial(s.Addr)
	if dialErr != nil {
		return dialErr
	}
	defer c.Close()

	c.Hello(s.Addr)
	c.Mail(smtpEmail.From.Address)
	c.Rcpt(smtpEmail.To.Address)

	wc, dataErr := c.Data()
	if dataErr != nil {
		return dataErr
	}
	defer wc.Close()

	buf := bytes.NewBufferString(body)
	if _, bufErr := buf.WriteTo(wc); bufErr != nil {
		return bufErr
	}
	c.Quit()

	return nil
}
