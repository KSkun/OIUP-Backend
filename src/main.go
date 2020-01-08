package main

import (
	"OIUP-Backend/config"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	

	_ = router.Run(config.RouterAddress)
}
