package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

type Claims struct {
	UserID int    `json:"user_id"`
	Role   int    `json:"role"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// GenerateToken 生成 JWT
func GenerateToken(userID int, role int, email string) (string, error) {
	// 从配置读取 Secret 和 过期时间
	secret := viper.GetString("jwt.secret")
	expirationSeconds := viper.GetInt64("jwt.expiration_seconds")

	claims := Claims{
		UserID: userID,
		Role:   role,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expirationSeconds) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "job-platform-go",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ParseToken 解析 JWT
func ParseToken(tokenString string) (*Claims, error) {
	secret := viper.GetString("jwt.secret")

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
