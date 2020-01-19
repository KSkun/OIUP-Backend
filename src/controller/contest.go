/*
   OIUP-Backend Project is developed by KSkun and licensed under GPL-3.0.
   Copyright (c) KSkun, 2020
*/
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

func getContestStatus() int {
    if config.Config.Contest.Status == config.ContestStatusError {
        return ContestStatusError
    }

    startTime := config.Config.Contest.StartTime
    nowTime := time.Now()
    if nowTime.Before(startTime) {
        return ContestStatusPreparing
    }
    if nowTime.After(startTime.Add(time.Minute * time.Duration(config.Config.Contest.Duration * 60))) {
        return ContestStatusEnd
    }
    return ContestStatusRunning
}

func ContestStatusHandler(context *gin.Context) {
    status := getContestStatus()

    startTime := config.Config.Contest.StartTime
    endTime := startTime.Add(time.Minute * time.Duration(config.Config.Contest.Duration * 60))
    response := gin.H{
        "status":     status,
        "start_time": startTime.Unix(),
        "end_time":   endTime.Unix(),
        "message":    "",
    }
    if status == ContestStatusError {
        response["message"] = config.Config.Contest.Message
    }
    util.SuccessResponse(context, response)
}

func ContestNameHandler(context *gin.Context) {
    util.SuccessResponse(context, gin.H{"name": config.Config.Contest.Name})
}

func ContestProblemsHandler(context *gin.Context) {
    util.SuccessResponse(context, gin.H{"url": config.Config.Contest.Download})
}

func ContestUnzipHandler(context *gin.Context) {
    nowTime := time.Now()
    validTime := config.Config.Contest.StartTime.Add(-time.Minute * time.Duration(config.Config.Contest.UnzipShift))
    if nowTime.Before(validTime) {
        util.ErrorResponse(context, http.StatusForbidden, "contest has not started", nil)
        return
    }

    util.SuccessResponse(context, gin.H{"unzip_token": config.Config.Contest.UnzipToken})
}
