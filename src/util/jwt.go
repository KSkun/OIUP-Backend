/*
   OIUP-Backend Project is developed by KSkun and licensed under GPL-3.0.
   Copyright (c) KSkun, 2020
 */
package util

import (
	"OIUP-Backend/config"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"time"
)

func GetIDFromContext(context *gin.Context) string {
	id, _ := context.Get("id")
	return id.(string)
}

func NewJWTToken(sub string, aud string, secret []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Minute * time.Duration(config.Config.JWT.JWTTokenLife)).Unix(),
		"iss": "OIUP",
		"sub": sub,
		"aud": aud,
	})
	return token.SignedString(secret)
}

func NewUserJWTToken(id string) (string, error) {
	return NewJWTToken("OIUP User Token", id, []byte(config.Config.JWT.JWTSecret))
}
