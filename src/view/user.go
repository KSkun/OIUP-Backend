/*
   OIUP-Backend Project is developed by KSkun and licensed under GPL-3.0.
   Copyright (c) KSkun, 2020
*/
package view

import (
    "OIUP-Backend/controller"
    "OIUP-Backend/middleware"
    "github.com/gin-gonic/gin"
)

func InitUserView(group *gin.RouterGroup) {
    group.GET("/token", controller.UserTokenHandler)
    group.GET("/info", middleware.ValidateUserToken, controller.UserInfoHandler)
}
