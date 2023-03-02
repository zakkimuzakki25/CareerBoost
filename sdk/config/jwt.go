package config

import (
	"github.com/golang-jwt/jwt/v4"
)

var JWT_KEY = []byte("ncjkanasdnmsaojdsaoijsocndinidSNWNWNDWJncinsjdff")

type JWTClaim struct {
	Username string
	jwt.RegisteredClaims
}
