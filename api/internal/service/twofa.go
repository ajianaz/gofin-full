package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/google/uuid"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"

	"github.com/ajianaz/gofin-full/api/internal/repository"
)

type TwoFAService struct {
	repo *repository.UserRepository
}

func NewTwoFAService(repo *repository.UserRepository) *TwoFAService {
	return &TwoFAService{repo: repo}
}

// TOTPSecret represents a TOTP secret for a user.
type TOTPSecret struct {
	Secret   string `json:"secret"`
	QRCodeURL string `json:"qr_code_url"`
}

// GenerateTOTP creates a new TOTP secret for a user.
func (s *TwoFAService) GenerateTOTP(ctx context.Context, userID uuid.UUID, email string) (*TOTPSecret, error) {
	key, err := GenerateRandomBytes(32)
	if err != nil {
		return nil, fmt.Errorf("failed to generate secret: %w", err)
	}

	secret := base64.StdEncoding.EncodeToString(key)
	otpOpts := totp.GenerateOpts{
		Issuer:      "Gofin",
		AccountName: email,
		Secret:      key,
		Period:      30,
		Digits:      otp.DigitsSix,
		Algorithm:   otp.AlgorithmSHA256,
	}

	k, err := totp.Generate(otpOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to create TOTP: %w", err)
	}

	url := k.URL()

	return &TOTPSecret{
		Secret:   secret,
		QRCodeURL: url,
	}, nil
}

// ValidateTOTP validates a TOTP code against a secret.
func (s *TwoFAService) ValidateTOTP(secret, code string) bool {
	decoded, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		return false
	}
	return totp.Validate(code, string(decoded))
}

// GenerateBackupCodes generates 10 backup codes for 2FA recovery.
func (s *TwoFAService) GenerateBackupCodes() ([]string, error) {
	codes := make([]string, 10)
	for i := 0; i < 10; i++ {
		b, err := GenerateRandomBytes(8)
		if err != nil {
			return nil, err
		}
		codes[i] = base32Encode(b)
	}
	return codes, nil
}

// GenerateRandomBytes generates cryptographically secure random bytes.
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	return b, err
}

func base32Encode(data []byte) string {
	chars := "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	result := make([]byte, len(data)*2)
	for i, b := range data {
		result[i*2] = chars[b%32]
		result[i*2+1] = chars[(b>>3)%32]
	}
	return string(result)
}

// GenerateRefreshToken creates a secure random refresh token.
func GenerateRefreshToken() (string, error) {
	b, err := GenerateRandomBytes(32)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// GenerateSessionID creates a secure random session ID.
func GenerateSessionID() (string, error) {
	b, err := GenerateRandomBytes(16)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
