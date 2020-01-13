/*
   OIUP-Backend Project is developed by KSkun and licensed under GPL-3.0.
   Copyright (c) KSkun, 2020
*/
package main

import (
	"OIUP-Backend/config"
	"OIUP-Backend/view"
	"github.com/gin-gonic/gin"
	"strconv"
)

func main() {
	router := gin.Default()
	apiGroup := router.Group("/api/v1")

	userGroup := apiGroup.Group("/user")
	view.InitUserView(userGroup)

	contestGroup := apiGroup.Group("/contest")
	view.InitContestView(contestGroup)

	problemGroup := apiGroup.Group("/problem")
	view.InitProblemView(problemGroup)

	miscGroup := apiGroup.Group("")
	view.InitMiscView(miscGroup)

	_ = router.Run(":" + strconv.Itoa(config.Config.HTTP.Port))
}
