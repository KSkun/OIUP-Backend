package view

import (
    "OIUP-Backend/controller"
    "github.com/gin-gonic/gin"
)

func InitMiscView(group *gin.RouterGroup) {
    group.GET("/time", controller.TimeHandler)
}
