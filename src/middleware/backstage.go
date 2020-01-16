/*
   OIUP-Backend Project is developed by KSkun and licensed under GPL-3.0.
   Copyright (c) KSkun, 2020
*/
package middleware

import (
    "OIUP-Backend/config"
    "OIUP-Backend/util"
    "errors"
    "github.com/gin-gonic/gin"
    "net/http"
    "strings"
)

func getKey(headerStr string) (string, error) {
    if len(headerStr) > 6 && strings.ToLower(headerStr[0:7]) == "bearer " {
        return headerStr[7:], nil
    }
    return "", errors.New("invalid key")
}

func ValidateBackstageKey(context *gin.Context) {
    key, err := getKey(context.GetHeader("Authorization"))
    if err != nil {
        util.ErrorResponse(context, http.StatusUnauthorized, err.Error(), nil)
        context.Abort()
        return
    }

    if key != config.Config.HTTP.BackstageKey {
        util.ErrorResponse(context, http.StatusUnauthorized, "invalid key", nil)
        context.Abort()
        return
    }
}
