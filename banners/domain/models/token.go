package models

import jwt "github.com/dgrijalva/jwt-go"

type Token struct {
	UserID int64
	Email  string
	Role   string
	*jwt.StandardClaims
}
