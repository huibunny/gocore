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
func ParseToken(tokenStr string, secret string) (string, int64, map[string]interface{}, error) {
	if len(secret) == 0 {
		secret = DEFAULT_SECRET
	}
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (i interface{}, e error) {
		return []byte(secret), nil
	})
	if err != nil {
		return "", -1, nil, err
	}
	userInfo := token.Claims.(jwt.MapClaims)
	userId := userInfo["user_id"].(string)
	expireTime := int64(userInfo["expire_time"].(float64))
	var extra map[string]interface{}
	if userInfo["extra"] != nil {
		extra = userInfo["extra"].(map[string]interface{})
	} else {
		//
	}

	return userId, expireTime, extra, nil
}
