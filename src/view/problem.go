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

func InitProblemView(group *gin.RouterGroup) {
    group.Use(middleware.ValidateUserToken)

    group.GET("/list", controller.ProblemListHandler)
    group.GET("/status", controller.ProblemStatusHandler)
    group.GET("/info", controller.ProblemInfoHandler)
    group.GET("/latest", controller.ProblemLatestHandler)

    solutionGroup := group.Group("/solution")
    solutionGroup.POST("/code", controller.ProblemSubmitCodeHandler)
    solutionGroup.POST("/output", controller.ProblemSubmitOutputHandler)
    solutionGroup.POST("/confirm", controller.ProblemConfirmHandler)
}
