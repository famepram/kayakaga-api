package tokenizer

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func GenerateAccessToken(userID uint, email string) (string, error) {
	secret := getEnv("JWT_SECRET", "finai-jwt-secret-change-in-production")
	expiry := getEnv("JWT_ACCESS_EXPIRY", "15m")

	expiryDuration, err := time.ParseDuration(expiry)
	if err != nil {
		expiryDuration = 15 * time.Minute
	}

	claims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiryDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			Issuer:    "kayakaga-api",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ValidateAccessToken(tokenString string) (*Claims, error) {
	secret := getEnv("JWT_SECRET", "finai-jwt-secret-change-in-production")

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func GenerateRefreshToken() string {
	return uuid.New().String()
}

func GetRefreshTokenExpiry() time.Time {
	expiry := getEnv("JWT_REFRESH_EXPIRY", "720h")
	expiryDuration, err := time.ParseDuration(expiry)
	if err != nil {
		expiryDuration = 720 * time.Hour
	}
	return time.Now().UTC().Add(expiryDuration)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
