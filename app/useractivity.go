package app

import (
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
)

func GetUserRole(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method")
		}
		return []byte("secret"), nil
	})
	var role string
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		role = claims["role"].(string)
		return role, nil
	}

	return role, err

}

func VerifyUser(tokenString string) (int64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method")
		}
		return []byte("secret"), nil
	})
	var userid int64
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userid = int64(claims["uid"].(float64))
		return userid, nil
	}

	return userid, err

}
