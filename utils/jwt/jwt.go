package jwt

import (
	"apps/config"
	"apps/utils/helper"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	TokenExpiredTime = 5 * 60 * 60 // 5 hours
)

func GenerateToken(payload map[string]interface{}) (string, error) {
	cfg, _ := config.LoadConfig()
	tokenClaims := jwt.MapClaims{
		"payload": payload,
		"exp":     time.Now().Add(time.Second * TokenExpiredTime).Unix(),
	}
	jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tokenClaims)
	token, err := jwtToken.SignedString([]byte(cfg.Jwt.SecretKey))
	if err != nil {
		return "", err
	}

	return token, nil
}

func ValidateToken(jwtToken string) (map[string]string, error) {
	cfg, _ := config.LoadConfig()
	cleanJWT := strings.Replace(jwtToken, "Bearer ", "", -1)
	tokenData := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(cleanJWT, tokenData, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.Jwt.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrInvalidKey
	}

	var data map[string]string
	helper.Copy(&data, tokenData["payload"])

	return data, nil
}