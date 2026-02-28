package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type ChirpClaims struct {
	jwt.RegisteredClaims
}

func HashPassword(password string) (string, error) {
	return argon2id.CreateHash(password, argon2id.DefaultParams)
}

func CheckPasswordHash(password string, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(password, hash)
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	claims := ChirpClaims{
		jwt.RegisteredClaims{
			Issuer:    "chirpy-access",
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
			Subject:   userID.String(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &ChirpClaims{}, func(token *jwt.Token) (any, error) { return []byte(tokenSecret), nil })

	if err != nil {
		return uuid.UUID{}, err
	}

	id, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.UUID{}, err
	}
	if userID, err := uuid.Parse(id); err != nil {
		return uuid.UUID{}, err
	} else {
		return userID, nil
	}
}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if len(authHeader) == 0 {
		return "", fmt.Errorf("Invalid Bearer Token")
	}
	//validate header scheme
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", fmt.Errorf("Invalid Bearer Token")
	}
	if parts := strings.Split(authHeader, " "); len(parts) != 2 {
		return "", fmt.Errorf("Invalid Bearer Token")
	}

	token := strings.TrimSpace(strings.Replace(authHeader, "Bearer", "", 1))
	if len(token) == 0 {
		return "", fmt.Errorf("Invalid Bearer Token")
	}
	return token, nil
}

func GetAPIKey(headers http.Header) (string, error) {

	authHeader := headers.Get("Authorization")
	if len(authHeader) == 0 {
		return "", fmt.Errorf("Invalid API Key")
	}
	//validate header scheme
	if !strings.HasPrefix(authHeader, "ApiKey ") {
		return "", fmt.Errorf("Invalid API Key")
	}
	if parts := strings.Split(authHeader, " "); len(parts) != 2 {
		return "", fmt.Errorf("Invalid API Key")
	}

	token := strings.TrimSpace(strings.Replace(authHeader, "ApiKey", "", 1))
	if len(token) == 0 {
		return "", fmt.Errorf("Invalid API Key")
	}
	return token, nil
}

func MakeRefreshToken() string {
	data := make([]byte, 32)
	_, err := rand.Read(data)
	if err != nil {
		return ""
	}
	token := hex.EncodeToString(data)

	return token
}
