/*
   OIUP-Backend Project is developed by KSkun and licensed under GPL-3.0.
   Copyright (c) KSkun, 2020
*/
package middleware

import (
	"OIUP-Backend/config"
	"OIUP-Backend/util"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func ValidateJWTToken(context *gin.Context, secret []byte) {
	token, err := request.ParseFromRequest(context.Request, request.OAuth2Extractor, func(token *jwt.Token) (interface{}, error) {
		if config.Config.JWT.JWTSigningMethod != token.Method.Alg() {
			return nil, errors.New("unexpected signing method")
		}
		return secret, nil
	})
	if err != nil {
		util.ErrorResponse(context, http.StatusUnauthorized, "invalid token", nil)
		context.Abort()
		return
	}

	if !token.Valid {
		util.ErrorResponse(context, http.StatusUnauthorized, "invalid token", nil)
		context.Abort()
		return
	}
	if time.Now().Unix() > int64(token.Claims.(jwt.MapClaims)["exp"].(float64)) {
		util.ErrorResponse(context, http.StatusUnauthorized, "token has expired", nil)
		context.Abort()
		return
	}

	context.Set("id", token.Claims.(jwt.MapClaims)["aud"])
}

func ValidateUserToken(context *gin.Context) {
	ValidateJWTToken(context, []byte(config.Config.JWT.JWTUserSecret))
}
