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

	userGroup := router.Group("/user")
	view.InitUserView(userGroup)

	_ = router.Run(":" + strconv.Itoa(config.Config.HTTP.Port))
}
