package config

import (
	"CareerBoost/entity"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// type JWTClaim struct {
// 	Username string
// 	jwt.RegisteredClaims
// }

func GenerateToken(payload entity.UserLogin) (string, error) {

	expStr := os.Getenv("JWT_EXP")
	var exp time.Duration
	exp, err := time.ParseDuration(expStr)
	if expStr == "" || err != nil {
		exp = time.Hour * 1
	}

	tokenJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, entity.NewUserClaims(payload.Email, exp))
	tokenJwtReal, err := tokenJWT.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return "", err
	}
	return tokenJwtReal, nil
}

func DecodeToken(signedToken string, ptrClaims jwt.Claims, KEY string) error {

	token, err := jwt.ParseWithClaims(signedToken, ptrClaims, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return "", errors.New("wrong signing method")
		}
		return []byte(KEY), nil
	})

	if err != nil {
		return fmt.Errorf("token has been tampered with")
	}

	if !token.Valid {
		// token is not valid
		return fmt.Errorf("invalid token")
	}

	return nil
}
