package utils

import (
	"github.com/golang-jwt/jwt"
)

const DEFAULT_SECRET = "gocoresec"

func CreateToken(obj map[string]interface{}, secret string) (string, error) {
	if len(secret) == 0 {
		secret = DEFAULT_SECRET
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(obj))
	tokenString, err := token.SignedString([]byte(secret))
	return tokenString, err
}
func ParseToken(tokenStr string, secret string) (map[string]interface{}, error) {
	if len(secret) == 0 {
		secret = DEFAULT_SECRET
	}
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (i interface{}, e error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	finToken := token.Claims.(jwt.MapClaims)
	return finToken, nil
}
