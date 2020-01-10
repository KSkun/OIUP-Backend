package controller

import (
    "OIUP-Backend/config"
    "OIUP-Backend/util"
    "github.com/gin-gonic/gin"
    "net/http"
    "time"
)

const (
    ContestStatusPreparing = 0
    ContestStatusRunning   = 1
    ContestStatusEnd       = 2
    ContestStatusError     = -1
)

func StatusHandler(context *gin.Context) {
    if config.Config.Contest.Status == config.ContestStatusError {
        util.SuccessResponse(context, gin.H{
            "status":  ContestStatusError,
            "message": config.Config.Contest.Message,
        })
        return
    }

    startTime, err := time.Parse("2006-01-02 15:04", config.Config.Contest.StartTime)
    if err != nil {
        util.ErrorResponse(context, http.StatusInternalServerError,
            "config error: invalid start_time", nil)
        return
    }
    if time.Now().Unix() < startTime.Unix() {
        util.SuccessResponse(context, gin.H{
            "status":  ContestStatusPreparing,
            "message": "",
        })
        return
    }
    if time.Now().Unix() > startTime.Add(time.Minute * time.Duration(config.Config.Contest.Duration)).Unix() {
        util.SuccessResponse(context, gin.H{
            "status":  ContestStatusEnd,
            "message": "",
        })
        return
    }
    util.SuccessResponse(context, gin.H{
        "status":  ContestStatusRunning,
        "message": "",
    })
}

func SubmitStatusHandler(context *gin.Context) {

}