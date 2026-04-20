package auth_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ajianaz/gofin-full/api/internal/auth"
)

func TestJWTManager_GenerateAndValidate(t *testing.T) {
	mgr := auth.NewJWTManager("test-secret-key-32-chars-min!", 60, 30)

	identity := &auth.UserIdentity{
		ID:       1,
		Email:    "test@example.com",
		Blocked:  false,
		DemoUser: false,
	}

	tokens, err := mgr.GenerateTokenPair(identity, nil)
	require.NoError(t, err)

	assert.NotEmpty(t, tokens.AccessToken)
	assert.NotEmpty(t, tokens.RefreshToken)
	assert.Equal(t, "Bearer", tokens.TokenType)
	assert.Equal(t, int64(3600), tokens.ExpiresIn)

	// Validate the access token
	claims, err := mgr.ValidateAccessToken(tokens.AccessToken)
	require.NoError(t, err)
	assert.Equal(t, int64(1), claims.UserID)
	assert.Equal(t, "test@example.com", claims.Email)
}

func TestJWTManager_InvalidToken(t *testing.T) {
	mgr := auth.NewJWTManager("test-secret-key-32-chars-min!", 60, 30)

	_, err := mgr.ValidateAccessToken("invalid-token")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid token")
}

func TestJWTManager_WrongSecret(t *testing.T) {
	mgr1 := auth.NewJWTManager("secret-one-32-chars-minimum", 60, 30)
	mgr2 := auth.NewJWTManager("secret-two-32-chars-minimum", 60, 30)

	identity := &auth.UserIdentity{ID: 1, Email: "test@example.com"}
	tokens, err := mgr1.GenerateTokenPair(identity, nil)
	require.NoError(t, err)

	_, err = mgr2.ValidateAccessToken(tokens.AccessToken)
	assert.Error(t, err)
}

func TestJWTManager_ExpiredToken(t *testing.T) {
	// Create a manager with 0-minute expiry for testing
	mgr := auth.NewJWTManager("test-secret-key-32-chars-min!", 0, 30)

	identity := &auth.UserIdentity{ID: 1, Email: "test@example.com"}
	tokens, err := mgr.GenerateTokenPair(identity, nil)
	require.NoError(t, err)

	// Wait for the token to expire (should be instant with 0 min expiry)
	time.Sleep(10 * time.Millisecond)

	_, err = mgr.ValidateAccessToken(tokens.AccessToken)
	assert.Error(t, err)
}

func TestJWTManager_GroupID(t *testing.T) {
	mgr := auth.NewJWTManager("test-secret-key-32-chars-min!", 60, 30)

	identity := &auth.UserIdentity{ID: 1, Email: "test@example.com"}
	groupID := int64(42)
	tokens, err := mgr.GenerateTokenPair(identity, &groupID)
	require.NoError(t, err)

	claims, err := mgr.ValidateAccessToken(tokens.AccessToken)
	require.NoError(t, err)
	require.NotNil(t, claims.GroupID)
	assert.Equal(t, int64(42), *claims.GroupID)
}

func TestJWTManager_DemoUser(t *testing.T) {
	mgr := auth.NewJWTManager("test-secret-key-32-chars-min!", 60, 30)

	identity := &auth.UserIdentity{ID: 1, Email: "demo@example.com", DemoUser: true}
	tokens, err := mgr.GenerateTokenPair(identity, nil)
	require.NoError(t, err)

	claims, err := mgr.ValidateAccessToken(tokens.AccessToken)
	require.NoError(t, err)
	assert.True(t, claims.DemoUser)
}

func TestHashRefreshToken(t *testing.T) {
	token := "my-refresh-token-value"
	hash := auth.HashRefreshToken(token)

	assert.NotEmpty(t, hash)
	assert.Len(t, hash, 64) // SHA256 hex = 64 chars
	assert.NotEqual(t, token, hash)

	// Same input should produce same hash
	hash2 := auth.HashRefreshToken(token)
	assert.Equal(t, hash, hash2)
}

func TestAuthProvider_LocalDefault(t *testing.T) {
provider := auth.NewDisabledProvider()
	assert.Equal(t, "disabled", provider.Name())
}

func TestAuthProvider_DisabledAuth(t *testing.T) {
	provider := auth.NewDisabledProvider()
	identity, err := provider.Authenticate(nil, auth.Credentials{})
	require.NoError(t, err)
	assert.Equal(t, int64(1), identity.ID)
	assert.Equal(t, "admin@local", identity.Email)
	assert.False(t, identity.Blocked)
	assert.False(t, identity.DemoUser)
}

func TestHashPassword(t *testing.T) {
	hash, err := auth.HashPassword("my-secure-password-123456")
	require.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.True(t, auth.CheckPassword(hash, "my-secure-password-123456"))
	assert.False(t, auth.CheckPassword(hash, "wrong-password"))
}

func TestHashPassword_MinLength(t *testing.T) {
	hash, err := auth.HashPassword("short")
	require.NoError(t, err)
	assert.NotEmpty(t, hash)
}
