/*
   OIUP-Backend Project is developed by KSkun and licensed under GPL-3.0.
   Copyright (c) KSkun, 2020
*/
package util

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func SuccessResponse(context *gin.Context, data interface{}) {
	context.JSON(http.StatusOK, gin.H{
		"message": 	"",
		"data": 	data,
	})
}

func ErrorResponse(context *gin.Context, code int, message string, data interface{}) {
	context.JSON(code, gin.H{
		"message": 	message,
		"data": 	data,
	})
}