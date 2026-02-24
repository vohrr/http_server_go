package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeJWT(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"
	expiresIn := time.Hour

	token, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT returned error: %v", err)
	}

	if token == "" {
		t.Fatal("MakeJWT returned empty token")
	}
}

func TestValidateJWT_ValidToken(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"
	expiresIn := time.Hour

	token, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT returned error: %v", err)
	}

	validatedUserID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("ValidateJWT returned error: %v", err)
	}

	if validatedUserID != userID {
		t.Errorf("ValidateJWT returned wrong user ID: got %s, want %s", validatedUserID, userID)
	}
}

func TestValidateJWT_WrongSecret(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"
	wrongSecret := "wrong-secret"
	expiresIn := time.Hour

	token, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT returned error: %v", err)
	}

	_, err = ValidateJWT(token, wrongSecret)
	if err == nil {
		t.Fatal("ValidateJWT should have returned error for wrong secret")
	}
}

func TestValidateJWT_InvalidToken(t *testing.T) {
	secret := "test-secret"

	_, err := ValidateJWT("invalid-token", secret)
	if err == nil {
		t.Fatal("ValidateJWT should have returned error for invalid token")
	}
}

func TestValidateJWT_ExpiredToken(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"
	expiresIn := -time.Hour

	token, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT returned error: %v", err)
	}

	_, err = ValidateJWT(token, secret)
	if err == nil {
		t.Fatal("ValidateJWT should have returned error for expired token")
	}
}

func TestGetBearerToken_Valid(t *testing.T) {
	headers := http.Header{}
	headers.Set("Authorization", "Bearer mytoken123")

	token, err := GetBearerToken(headers)
	if err != nil {
		t.Fatalf("GetBearerToken returned error: %v", err)
	}

	if token != "mytoken123" {
		t.Errorf("GetBearerToken returned wrong token: got %s, want mytoken123", token)
	}
}

func TestGetBearerToken_MissingHeader(t *testing.T) {
	headers := http.Header{}

	_, err := GetBearerToken(headers)
	if err == nil {
		t.Fatal("GetBearerToken should have returned error for missing header")
	}
}

func TestGetBearerToken_InvalidScheme(t *testing.T) {
	headers := http.Header{}
	headers.Set("Authorization", "Basic somer token")

	_, err := GetBearerToken(headers)
	if err == nil {
		t.Fatal("GetBearerToken should have returned error for invalid scheme")
	}
}

func TestGetBearerToken_LowercaseScheme(t *testing.T) {
	headers := http.Header{}
	headers.Set("Authorization", "bearer mytoken123")

	_, err := GetBearerToken(headers)
	if err == nil {
		t.Fatal("GetBearerToken should have returned error for lowercase scheme")
	}
}

func TestGetBearerToken_EmptyAfterTrim(t *testing.T) {
	headers := http.Header{}
	headers.Set("Authorization", "Bearer ")

	_, err := GetBearerToken(headers)
	if err == nil {
		t.Fatal("GetBearerToken should have returned error for empty token")
	}
}
