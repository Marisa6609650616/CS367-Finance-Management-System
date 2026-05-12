package tests

import (
	"os"
	"testing"
	"time"

	"CS367-Finance-Management-System/pkg/utils"

	"github.com/golang-jwt/jwt/v5"
)

func setup() {
	os.Setenv("JWT_SECRET", "test_secret_key_for_unit_tests")
}

// GenerateToken Tests
func TestGenerateToken_Success(t *testing.T) {
	setup()
	token, err := utils.GenerateToken(1, "test@dome.tu.ac.th.com", "user")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if token == "" {
		t.Fatal("expected token string, got empty string")
	}
}

func TestGenerateToken_NoSecret(t *testing.T) {
	os.Setenv("JWT_SECRET", "")
	_, err := utils.GenerateToken(1, "test@dome.tu.ac.th.com", "user")
	if err == nil {
		t.Fatal("expected error when JWT_SECRET is empty, got nil")
	}
}

func TestGenerateToken_ContainsCorrectClaims(t *testing.T) {
	setup()
	token, _ := utils.GenerateToken(42, "MeMyself@gmail.com", "admin")
	claims, err := utils.ParseToken(token)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	if int(claims["user_id"].(float64)) != 42 {
		t.Errorf("expected user_id=42, got %v", claims["user_id"])
	}
	if claims["email"] != "MeMyself@gmail.com" {
		t.Errorf("expected email=MeMyself@gmail.com, got %v", claims["email"])
	}
	if claims["role"] != "admin" {
		t.Errorf("expected role=admin, got %v", claims["role"])
	}
}

// ParseToken Tests
func TestParseToken_ValidToken(t *testing.T) {
	setup()
	token, _ := utils.GenerateToken(1, "test@dome.tu.ac.th.com", "user")
	claims, err := utils.ParseToken(token)
	if err != nil {
		t.Fatalf("expected valid parse, got error: %v", err)
	}
	if claims == nil {
		t.Fatal("expected claims, got nil")
	}
}

func TestParseToken_EmptyString(t *testing.T) {
	setup()
	_, err := utils.ParseToken("")
	if err == nil {
		t.Fatal("expected error on empty token, got nil")
	}
}

func TestParseToken_RandomString(t *testing.T) {
	setup()
	_, err := utils.ParseToken("this.is.not.a.jwt")
	if err == nil {
		t.Fatal("expected error on fake token, got nil")
	}
}

func TestParseToken_TamperedPayload(t *testing.T) {
	setup()
	tampered := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9" +
		".eyJ1c2VyX2lkIjo5OTksImVtYWlsIjoiaGFja2VyQGV2aWwuY29tIiwicm9sZSI6ImFkbWluIn0" +
		".invalidsignature"
	_, err := utils.ParseToken(tampered)
	if err == nil {
		t.Fatal("expected error on tampered token, got nil")
	}
}

func TestParseToken_WrongAlgorithm(t *testing.T) {
	setup()
	token := jwt.NewWithClaims(jwt.SigningMethodNone, jwt.MapClaims{
		"user_id": 1,
		"email":   "hacker@inwza.com",
		"role":    "admin",
		"exp":     time.Now().Add(time.Hour).Unix(),
	})
	tokenStr, _ := token.SignedString(jwt.UnsafeAllowNoneSignatureType)

	_, err := utils.ParseToken(tokenStr)
	if err == nil {
		t.Fatal("expected error on 'none' algorithm token, got nil")
	}
}

func TestParseToken_ExpiredToken(t *testing.T) {
	setup()
	secret := []byte("test_secret_key_for_unit_tests")
	claims := jwt.MapClaims{
		"user_id": 1,
		"email":   "test@dome.tu.ac.th.com",
		"role":    "user",
		"exp":     time.Now().Add(-1 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, _ := token.SignedString(secret)

	_, err := utils.ParseToken(tokenStr)
	if err == nil {
		t.Fatal("expected error on expired token, got nil")
	}
}

func TestParseToken_WrongSecret(t *testing.T) {
	os.Setenv("JWT_SECRET", "secret_A")
	token, _ := utils.GenerateToken(1, "test@dome.tu.ac.th.com", "user")

	os.Setenv("JWT_SECRET", "secret_B")
	_, err := utils.ParseToken(token)
	if err == nil {
		t.Fatal("expected error when parsing with wrong secret, got nil")
	}
}
