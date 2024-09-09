package smtp

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/redis/go-redis/v9"
	gomail "gopkg.in/mail.v2"
	"io"
	"net/smtp"
	"test-auth/internal/config"
	"test-auth/pkg/token_manager"
	"time"
)

var (
	numbers = [10]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
)

type Smtp struct {
	cfg          config.SMTP
	cacheClient  *redis.Client
	tokenManager *token_manager.TokenManager
}

func NewSmtp(cfg config.SMTP, cacheClient *redis.Client, tokenManager *token_manager.TokenManager) *Smtp {
	return &Smtp{
		cfg:          cfg,
		cacheClient:  cacheClient,
		tokenManager: tokenManager,
	}
}

func generateCode() (string, error) {
	b := make([]byte, 5)
	n, err := io.ReadAtLeast(rand.Reader, b, 5)
	if err != nil {
		return "", err
	}
	if n != 5 {
		return "", errors.New("can't create code")
	}

	for i := 0; i < len(b); i++ {
		b[i] = numbers[int(b[i])%len(numbers)]
	}
	code := string(b)

	return code, nil
}

func (s *Smtp) SendCode(ctx context.Context, receiver, userId string) error {
	code, err := generateCode()
	if err != nil {
		return err
	}

	message := gomail.NewMessage()
	message.SetHeader("From", s.cfg.Username)
	message.SetHeader("To", receiver)
	message.SetHeader("Subject", "Your one time password")
	message.SetBody("text/plain", code)

	auth := smtp.PlainAuth("", s.cfg.Username, s.cfg.Password, s.cfg.Host)
	dialer := gomail.NewDialer(s.cfg.Host, s.cfg.Port, s.cfg.Username, s.cfg.Password)

	dialer.Auth = auth
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: false, ServerName: dialer.Host}

	err = dialer.DialAndSend(message)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("user_%s", userId)

	val, err := jsoniter.Marshal(code)
	if err != nil {
		return err
	}

	status := s.cacheClient.Set(ctx, key, val, 3*time.Minute)
	if _, err = status.Result(); err != nil {
		return err
	}
	return nil
}
