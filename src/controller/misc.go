package controller

import (
    "OIUP-Backend/util"
    "github.com/gin-gonic/gin"
    "time"
)

func TimeHandler(context *gin.Context) {
    util.SuccessResponse(context, gin.H{"time": time.Now().Unix()})
}
