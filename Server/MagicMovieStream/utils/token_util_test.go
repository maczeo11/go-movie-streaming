package utils

import (
	"testing"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

func TestGenerateAndValidateToken(t *testing.T) {
	token, _, err := GenerateAllTokens("jane@example.com", "Jane", "Doe", "USER", "user_123")
	if err != nil {
		t.Fatal(err)
	}

	claims, err := ValidateToken(token)
	if err != nil {
		t.Fatalf("token should validate: %v", err)
	}
	if claims.Email != "jane@example.com" {
		t.Errorf("email mismatch: %s", claims.Email)
	}
	if claims.Role != "USER" {
		t.Errorf("role mismatch: %s", claims.Role)
	}
	if claims.UserId != "user_123" {
		t.Errorf("user id mismatch: %s", claims.UserId)
	}
}

func TestValidateTokenWrongKey(t *testing.T) {
	t.Setenv("SECRET_KEY", "key-a")

	token, _, err := GenerateAllTokens("a@b.c", "A", "B", "USER", "u1")
	if err != nil {
		t.Fatal(err)
	}

	t.Setenv("SECRET_KEY", "key-b")
	if _, err := ValidateToken(token); err == nil {
		t.Fatal("expected validation to fail with a different signing key")
	}
}

func TestValidateTokenGarbage(t *testing.T) {
	if _, err := ValidateToken("not.a.jwt"); err == nil {
		t.Fatal("expected error for a garbage token")
	}
}

func TestRefreshTokenValidation(t *testing.T) {
	_, refresh, err := GenerateAllTokens("a@b.c", "A", "B", "USER", "u1")
	if err != nil {
		t.Fatal(err)
	}

	claims, err := ValidateRefreshToken(refresh)
	if err != nil {
		t.Fatalf("refresh token should validate: %v", err)
	}
	if claims.UserId != "u1" {
		t.Errorf("user id mismatch: %s", claims.UserId)
	}
}

func TestExpiredTokenRejected(t *testing.T) {
	token := createExpiredToken(t, "u1", "USER")
	if _, err := ValidateToken(token); err == nil {
		t.Fatal("expired token should be rejected")
	}
}

// createExpiredToken signs a token that expired 10 minutes ago. It calls
// the real signing path so the test doesn't drift from production code.
func createExpiredToken(t *testing.T, userId, role string) string {
	t.Helper()
	t.Setenv("SECRET_KEY", "expiry-test-key")

	claims := &SignedDetails{
		Email:     "expired@example.com",
		FirstName: "Old",
		LastName:  "Token",
		Role:      role,
		UserId:    userId,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "MagicStream",
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-10 * time.Minute)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(getSecretKey()))
	if err != nil {
		t.Fatal(err)
	}
	return signed
}
