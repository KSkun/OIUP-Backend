/*
   OIUP-Backend Project is developed by KSkun and licensed under GPL-3.0.
   Copyright (c) KSkun, 2020
*/
package main

import (
	"OIUP-Backend/config"
	"OIUP-Backend/middleware"
	"OIUP-Backend/view"
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"strconv"
	"time"
)

func main() {
	err := os.MkdirAll("logs/", os.ModePerm)
	if err != nil {
		panic(err)
	}
	logFile, err := os.Create("logs/access-" + time.Now().Format("200601021504") + ".log")
	if err != nil {
		panic(err)
	}
	gin.DefaultWriter = io.MultiWriter(logFile)

	router := gin.Default()
	router.Use(middleware.CrossOriginMiddleware)
	apiGroup := router.Group("/api/v1")

	userGroup := apiGroup.Group("/user")
	view.InitUserView(userGroup)

	contestGroup := apiGroup.Group("/contest")
	view.InitContestView(contestGroup)

	problemGroup := apiGroup.Group("/problem")
	view.InitProblemView(problemGroup)

	miscGroup := apiGroup.Group("")
	view.InitMiscView(miscGroup)

	adminGroup := apiGroup.Group("/admin")
	view.InitAdminView(adminGroup)

	err = router.Run(":" + strconv.Itoa(config.Config.HTTP.Port))
	if err != nil {
		panic(err)
	}
}
