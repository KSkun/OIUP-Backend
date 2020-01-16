/*
   OIUP-Backend Project is developed by KSkun and licensed under GPL-3.0.
   Copyright (c) KSkun, 2020
*/
package view

import (
    "OIUP-Backend/controller"
    "github.com/gin-gonic/gin"
)

func InitMiscView(group *gin.RouterGroup) {
    group.GET("/time", controller.TimeHandler)
}
