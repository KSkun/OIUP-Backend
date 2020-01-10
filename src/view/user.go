package view

import (
    "OIUP-Backend/controller"
    "OIUP-Backend/middleware"
    "github.com/gin-gonic/gin"
)

func InitUserView(group *gin.RouterGroup) {
    group.GET("/token", controller.GetTokenHandler)
    group.GET("/info", middleware.ValidateUserToken, controller.GetInfoHandler)
}
