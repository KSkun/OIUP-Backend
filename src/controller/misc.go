/*
   OIUP-Backend Project is developed by KSkun and licensed under GPL-3.0.
   Copyright (c) KSkun, 2020
*/
package controller

import (
    "OIUP-Backend/util"
    "github.com/gin-gonic/gin"
    "time"
)

func TimeHandler(context *gin.Context) {
    util.SuccessResponse(context, gin.H{"time": time.Now().Unix()})
}
