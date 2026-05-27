package service

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strings"

	"github.com/ajianaz/gofin-full/api/internal/config"
)

// MailService sends plain text emails via SMTP.
type MailService struct {
	host   string
	port   int
	user   string
	pass   string
	from   string
	useTLS bool
}

// NewMailService creates a new MailService from config.
// If SMTP is not configured, the service can still be created but
// SendEmail will return an error.
func NewMailService(cfg *config.Config) *MailService {
	return &MailService{
		host:   cfg.SMTPHost,
		port:   cfg.SMTPPort,
		user:   cfg.SMTPUser,
		pass:   cfg.SMTPPass,
		from:   cfg.SMTPFrom,
		useTLS: cfg.SMTPUseTLS,
	}
}

// Configured returns true if SMTP host is set.
func (s *MailService) Configured() bool {
	return s.host != ""
}

// SendEmail sends a plain text email.
// Returns an error if SMTP is not configured or the send fails.
func (s *MailService) SendEmail(to, subject, body string) error {
	if !s.Configured() {
		return fmt.Errorf("SMTP not configured")
	}

	from := s.from
	if from == "" {
		from = "noreply@gofin.local"
	}

	addr := net.JoinHostPort(s.host, fmt.Sprintf("%d", s.port))
	msg := "From: " + from + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/plain; charset=\"utf-8\"\r\n" +
		"\r\n" +
		body + "\r\n"

	// Port 465: direct TLS connection
	if s.port == 465 {
		tlsConfig := &tls.Config{ServerName: s.host}
		if !s.useTLS {
			tlsConfig = nil
		}
		conn, err := tls.Dial("tcp", addr, tlsConfig)
		if err != nil {
			return fmt.Errorf("smtp tls dial: %w", err)
		}
		client, err := smtp.NewClient(conn, s.host)
		if err != nil {
			return fmt.Errorf("smtp new client: %w", err)
		}
		defer client.Close()

		if err = client.Auth(smtp.PlainAuth("", s.user, s.pass, s.host)); err != nil {
			return fmt.Errorf("smtp auth: %w", err)
		}
		if err = client.Mail(from); err != nil {
			return fmt.Errorf("smtp mail from: %w", err)
		}
		if err = client.Rcpt(to); err != nil {
			return fmt.Errorf("smtp rcpt to: %w", err)
		}
		w, err := client.Data()
		if err != nil {
			return fmt.Errorf("smtp data: %w", err)
		}
		if _, err = w.Write([]byte(msg)); err != nil {
			return fmt.Errorf("smtp write: %w", err)
		}
		if err = w.Close(); err != nil {
			return fmt.Errorf("smtp close: %w", err)
		}
		return client.Quit()
	}

	// Port 587 (or any other): STARTTLS
	auth := smtp.PlainAuth("", s.user, s.pass, s.host)
	if s.useTLS {
		return smtp.SendMail(addr, auth, from, []string{to}, []byte(msg))
	}

	// Plain (no TLS)
	client, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("smtp dial: %w", err)
	}
	defer client.Close()

	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("smtp auth: %w", err)
	}
	if err = client.Mail(from); err != nil {
		return fmt.Errorf("smtp mail from: %w", err)
	}
	if err = client.Rcpt(to); err != nil {
		return fmt.Errorf("smtp rcpt to: %w", err)
	}
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("smtp data: %w", err)
	}
	if _, err = w.Write([]byte(msg)); err != nil {
		return fmt.Errorf("smtp write: %w", err)
	}
	if err = w.Close(); err != nil {
		return fmt.Errorf("smtp close: %w", err)
	}
	return client.Quit()
}

// sanitizeEmail prevents header injection in email addresses.
func sanitizeEmail(email string) string {
	return strings.ReplaceAll(email, "\r\n", "")
}
