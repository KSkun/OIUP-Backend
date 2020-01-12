package view

import (
    "OIUP-Backend/controller"
    "OIUP-Backend/middleware"
    "github.com/gin-gonic/gin"
)

func InitContestView(group *gin.RouterGroup) {
    group.Use(middleware.ValidateUserToken)

    group.GET("/status", controller.StatusHandler)
    group.GET("/name", controller.NameHandler)
    group.GET("/problems", controller.ProblemsHandler)
    group.GET("/unzip", controller.UnzipHandler)
}
