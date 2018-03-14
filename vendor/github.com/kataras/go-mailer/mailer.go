// Package mailer is a simple e-mail sender for Go Programming Language
package mailer

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/valyala/bytebufferpool"
	"net/mail"
	"net/smtp"
	"os/exec"
	"strings"
)

const (
	// Version current version number
	Version = "0.0.2"
)

var buf bytebufferpool.Pool

type (
	// Service is the interface which mail sender(mailer) should implement
	Service interface {
		// Send sends a mail to recipients
		// the body can be html also
		//
		// Note: you can change the UseCommand in runtime
		Send(string, string, ...string) error
		// UpdateConfig replaces the current configuration with the receiver
		UpdateConfig(Config)
	}

	mailer struct {
		config        Config
		fromAddr      mail.Address
		auth          smtp.Auth
		authenticated bool
	}
)

// New creates and returns a new mail service
func New(cfg Config) Service {
	m := &mailer{config: cfg}
	addr := cfg.FromAddr
	if addr == "" {
		addr = cfg.Username
	}

	if cfg.FromAlias == "" {
		if !cfg.UseCommand && cfg.Username != "" && strings.Contains(cfg.Username, "@") {
			m.fromAddr = mail.Address{Name: cfg.Username[0:strings.IndexByte(cfg.Username, '@')], Address: addr}
		}
	} else {
		m.fromAddr = mail.Address{Name: cfg.FromAlias, Address: addr}
	}
	return m
}

func (m *mailer) UpdateConfig(cfg Config) {
	m.config = cfg
}

// Send sends a mail to recipients
// the body can be html also
//
// Note: you can change the UseCommand in runtime
func (m *mailer) Send(subject string, body string, to ...string) error {
	if m.config.UseCommand {
		return m.sendCmd(subject, body, to)
	}

	return m.sendSMTP(subject, body, to)
}

func (m *mailer) sendSMTP(subject string, body string, to []string) error {
	buffer := buf.Get()
	defer buf.Put(buffer)

	if !m.authenticated {
		cfg := m.config
		if cfg.Username == "" || cfg.Password == "" || cfg.Host == "" || cfg.Port <= 0 {
			return fmt.Errorf("Username, Password, Host & Port cannot be empty when using SMTP!")
		}
		m.auth = smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)
		m.authenticated = true
	}

	fullhost := fmt.Sprintf("%s:%d", m.config.Host, m.config.Port)

	header := make(map[string]string)
	header["From"] = m.fromAddr.String()
	header["To"] = strings.Join(to, ",")
	header["Subject"] = subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(body))

	return smtp.SendMail(
		fmt.Sprintf(fullhost),
		m.auth,
		m.config.Username,
		to,
		[]byte(message),
	)
}

func (m *mailer) sendCmd(subject string, body string, to []string) error {
	buffer := buf.Get()
	defer buf.Put(buffer)

	header := make(map[string]string)
	header["To"] = strings.Join(to, ",")
	header["Subject"] = subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(body))
	buffer.WriteString(message)

	cmd := exec.Command("sendmail", "-F", m.fromAddr.Name, "-f", m.fromAddr.Address, "-t")
	cmd.Stdin = bytes.NewBuffer(buffer.Bytes())
	_, err := cmd.CombinedOutput()
	return err
}
