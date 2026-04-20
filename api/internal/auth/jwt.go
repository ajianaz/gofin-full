package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TokenPair holds access and refresh tokens.
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

// Claims represents the JWT claims for access tokens.
type Claims struct {
	jwt.RegisteredClaims
	UserID   int64  `json:"user_id"`
	Email    string `json:"email"`
	GroupID  *int64 `json:"group_id,omitempty"`
	DemoUser bool   `json:"demo_user,omitempty"`
}

// refreshClaims represents the JWT claims for refresh tokens.
type refreshClaims struct {
	jwt.RegisteredClaims
	UserID   int64  `json:"user_id"`
	Email    string `json:"email"`
	GroupID  *int64 `json:"group_id,omitempty"`
	DemoUser bool   `json:"demo_user,omitempty"`
	TokenID  string `json:"tid"`
}

// JWTManager handles JWT token creation and validation.
type JWTManager struct {
	secret        []byte
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

// NewJWTManager creates a new JWT manager.
func NewJWTManager(secret string, accessExpiryMin, refreshExpiryDays int) *JWTManager {
	return &JWTManager{
		secret:        []byte(secret),
		accessExpiry:  time.Duration(accessExpiryMin) * time.Minute,
		refreshExpiry: time.Duration(refreshExpiryDays) * 24 * time.Hour,
	}
}

// GenerateTokenPair creates an access + refresh token pair.
func (m *JWTManager) GenerateTokenPair(identity *UserIdentity, groupID *int64) (*TokenPair, error) {
	now := time.Now()

	// Access token
	accessClaims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        generateTokenID(),
			Subject:   fmt.Sprintf("%d", identity.ID),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.accessExpiry)),
		},
		UserID:   identity.ID,
		Email:    identity.Email,
		GroupID:  groupID,
		DemoUser: identity.DemoUser,
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessStr, err := accessToken.SignedString(m.secret)
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	// Refresh token (signed JWT with longer expiry)
	tid := generateTokenID()
	refreshClaims := refreshClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprintf("%d", identity.ID),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.refreshExpiry)),
		},
		UserID:   identity.ID,
		Email:    identity.Email,
		GroupID:  groupID,
		DemoUser: identity.DemoUser,
		TokenID:  tid,
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshStr, err := refreshToken.SignedString(m.secret)
	if err != nil {
		return nil, fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessStr,
		RefreshToken: refreshStr,
		ExpiresIn:    int64(m.accessExpiry.Seconds()),
		TokenType:    "Bearer",
	}, nil
}

// ValidateAccessToken parses and validates an access token, returning claims.
func (m *JWTManager) ValidateAccessToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return m.secret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

// ValidateRefreshToken parses and validates a refresh token, returning claims.
func (m *JWTManager) ValidateRefreshToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &refreshClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return m.secret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	rc, ok := token.Claims.(*refreshClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid refresh token claims")
	}

	return &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: rc.Subject,
		},
		UserID:   rc.UserID,
		Email:    rc.Email,
		GroupID:  rc.GroupID,
		DemoUser: rc.DemoUser,
	}, nil
}

// HashRefreshToken hashes a refresh token for storage.
func HashRefreshToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return fmt.Sprintf("%x", h)
}

func generateTokenID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("%x", time.Now().UnixNano())
	}
	return fmt.Sprintf("%x", b)
}
