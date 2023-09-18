package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

const TOKEN_EXP = time.Hour * 24 * 365
const SECRET_KEY = "supersecretkey"

type claims struct {
	jwt.RegisteredClaims
	UserID string
}

type Auth struct {
	secretKey string
}

func NewAuth() *Auth {
	secretKey := SECRET_KEY
	if parsedSecretKey, exist := os.LookupEnv("SECRET_KEY"); exist {
		secretKey = parsedSecretKey
	}
	return &Auth{
		secretKey: secretKey,
	}
}

func (a Auth) GenerateToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TOKEN_EXP)),
		},
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(a.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken checks if token is valid and returns uuid
func (a Auth) ValidateToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(a.secretKey), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*claims); ok && token.Valid {
		return claims.UserID, nil
	} else {
		return "", fmt.Errorf("invalid token")
	}
}

func (a Auth) GenerateUUID() string {
	id := uuid.New().String()
	return id
}
