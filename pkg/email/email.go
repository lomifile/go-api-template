// Package email provides interface for email client
package email

import (
	"github.com/lomifile/api/config"
	"github.com/lomifile/api/pkg/logger"
	gmail "github.com/wneessen/go-mail"
)

// SendEmailConfig Config what needs to be passed to send HTML email
type SendEmailConfig struct {
	To                string
	Subject           string
	AlternativeString string
	HTML              string
}

// Client Main email client wiht logger and sender
type Client struct {
	e    *gmail.Client
	l    *logger.Logger
	from string
}

// NewEmailClient Creates new instace of email client
func NewEmailClient(l *logger.Logger, c *config.Config) *Client {
	client, err := gmail.NewClient(
		c.Email.Host,
		gmail.WithPort(587),
		gmail.WithUsername(c.Email.Username),
		gmail.WithPassword(c.Email.Password),
		gmail.WithSMTPAuth(gmail.SMTPAuthPlain),
		gmail.WithTLSPolicy(gmail.TLSMandatory),
	)
	if err != nil {
		l.Fatal(err.Error())
	}

	l.Info("Email client ready")
	return &Client{e: client, l: l, from: c.Email.Username}
}

// SendHTMLEmail Sends HTML type email
func (em *Client) SendHTMLEmail(cfg *SendEmailConfig) error {
	msg := gmail.NewMsg()

	err := msg.From(em.from)
	if err != nil {
		return err
	}

	err = msg.To(cfg.To)
	if err != nil {
		return err
	}

	msg.Subject(cfg.Subject)
	msg.AddAlternativeString(gmail.TypeTextPlain, cfg.AlternativeString)
	msg.SetBodyString(gmail.TypeTextHTML, cfg.HTML)

	err = em.e.DialAndSend(msg)
	if err != nil {
		return err
	}

	return nil
}
