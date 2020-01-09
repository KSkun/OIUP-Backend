/*
   OIUP-Backend Project is developed by KSkun and licensed under GPL-3.0.
   Copyright (c) KSkun, 2020
*/
package main

import (
	"OIUP-Backend/config"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()



	_ = router.Run(":" + string(config.Config.HTTP.Port))
}
