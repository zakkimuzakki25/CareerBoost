package entity

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type AdminLogin struct {
	Username string `json:"username" gorm:"NOT NULL"`
	Password string `json:"password" gorm:"NOT NULL"`
}

type AdminClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func NewAdminClaims(username string, exp time.Duration) AdminClaims {
	return AdminClaims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(exp)),
		},
	}
}
