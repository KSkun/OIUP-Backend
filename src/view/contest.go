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

func InitContestView(group *gin.RouterGroup) {
    group.Use(middleware.ValidateUserToken)

    group.GET("/status", controller.ContestStatusHandler)
    group.GET("/name", controller.ContestNameHandler)
    group.GET("/problems", controller.ContestProblemsHandler)
    group.GET("/unzip", controller.ContestUnzipHandler)
}
