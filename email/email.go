package email

import (
	"github.com/XC-Zero/zero_common/config"
	"github.com/jordan-wright/email"
	"github.com/pkg/errors"

	"net/smtp"
	"strings"
)

// SendEmail 发送邮件 ，仅支持SMTP服务
func SendEmail(cfg config.EmailConfig, subject, text, content string, to ...string) error {
	e := email.NewEmail()
	e.From = cfg.SenderEmail
	e.To = to
	e.Subject = subject
	e.Text = []byte(text)
	e.HTML = []byte(content)
	err := e.Send(cfg.EmailServerAddr, smtp.PlainAuth("", cfg.SenderEmail,
		cfg.EmailSecret, strings.Split(cfg.EmailServerAddr, ":")[0]))
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
