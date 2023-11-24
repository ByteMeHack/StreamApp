package auth

import (
	"fmt"
	"os"
	"strconv"

	"github.com/golang-jwt/jwt"
)

var secret []byte

func init() {
	secret = []byte(os.Getenv("AUTH_SECRET"))
}

func CreateToken(userID int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": strconv.FormatInt(userID, 10),
	})
	return token.SignedString(secret)
}

func Auth(tokenString string) (int64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		strId := claims["user"]
		if strId == nil {
			return 0, fmt.Errorf("incorrect auth token")
		}
		id, err := strconv.ParseInt(strId.(string), 10, 64)

		if err != nil {
			return 0, err
		}
		return id, nil
	}

	return 0, err
}
