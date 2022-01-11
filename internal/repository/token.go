package repository

import "github.com/dgrijalva/jwt-go/v4"

type Token struct {
	jwt.StandardClaims
	Data map[string]string `jwt:"data"`
}