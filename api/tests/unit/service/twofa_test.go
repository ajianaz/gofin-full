package service_test

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/azfirazka/gofin-full/api/internal/repository"
	"github.com/azfirazka/gofin-full/api/internal/service"
)

func newTwoFAService() *service.TwoFAService {
	return service.NewTwoFAService(&repository.UserRepository{})
}

func TestValidateTOTP_RejectsHardcodedBypass(t *testing.T) {
	svc := newTwoFAService()

	// The old hardcoded bypass "123456" must be rejected
	secret := base64.StdEncoding.EncodeToString([]byte("anything"))
	valid := svc.ValidateTOTP(secret, "123456")
	assert.False(t, valid, "hardcoded bypass code '123456' must be rejected")
}

func TestValidateTOTP_RejectsInvalidSecret(t *testing.T) {
	svc := newTwoFAService()

	// Invalid base64 secret should not crash, just return false
	valid := svc.ValidateTOTP("!!!invalid-base64!!!", "654321")
	assert.False(t, valid, "invalid secret should return false")
}

func TestValidateTOTP_RejectsWrongCode(t *testing.T) {
	svc := newTwoFAService()

	// Use a dummy base64 secret — wrong code should be rejected
	dummySecret := base64.StdEncoding.EncodeToString([]byte("test-secret-key"))
	valid := svc.ValidateTOTP(dummySecret, "000000")
	assert.False(t, valid, "wrong TOTP code should be rejected")
}

func TestGenerateBackupCodes(t *testing.T) {
	svc := newTwoFAService()

	codes, err := svc.GenerateBackupCodes()
	require.NoError(t, err)
	assert.Len(t, codes, 10, "should generate exactly 10 backup codes")

	for _, code := range codes {
		assert.NotEmpty(t, code)
		assert.Len(t, code, 16, "each backup code should be 16 chars")
	}

	// All codes should be unique
	seen := make(map[string]bool)
	for _, code := range codes {
		assert.False(t, seen[code], "backup codes must be unique")
		seen[code] = true
	}
}

func TestGenerateRandomBytes(t *testing.T) {
	b1, err := service.GenerateRandomBytes(32)
	require.NoError(t, err)
	assert.Len(t, b1, 32)

	b2, err := service.GenerateRandomBytes(32)
	require.NoError(t, err)
	assert.Len(t, b2, 32)

	// Two calls should produce different values (extremely unlikely to collide)
	assert.NotEqual(t, b1, b2, "random bytes should differ between calls")
}

func TestGenerateRefreshToken(t *testing.T) {
	token, err := service.GenerateRefreshToken()
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	token2, err := service.GenerateRefreshToken()
	require.NoError(t, err)
	assert.NotEqual(t, token, token2, "refresh tokens should be unique")
}

func TestGenerateSessionID(t *testing.T) {
	session, err := service.GenerateSessionID()
	require.NoError(t, err)
	assert.NotEmpty(t, session)
}
