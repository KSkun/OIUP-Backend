/*
   OIUP-Backend Project is developed by KSkun and licensed under GPL-3.0.
   Copyright (c) KSkun, 2020
*/
package middleware

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

func CrossOriginMiddleware(context *gin.Context) {
    context.Header("Access-Control-Allow-Origin", "*")
    context.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
    context.Header("Access-Control-Allow-Credentials", "true")
    context.Header("Access-Control-Allow-Headers", "Authorization,DNT,User-Agent,Keep-Alive,Content-Type,Accept,Origin,X-Requested-With")

    if context.Request.Method == "OPTIONS" {
        context.AbortWithStatus(http.StatusNoContent)
        return
    }

    context.Next()
}
