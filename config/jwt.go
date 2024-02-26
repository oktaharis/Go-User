package config

import (
	"github.com/golang-jwt/jwt/v4"
	_"github.com/golang-jwt/jwt/v5"

)


var JWT_KEY = []byte("asdfghjklsadasfassad12")

type JWTClaim struct {
	Email string
	Role string
	Product string
	jwt.RegisteredClaims
}