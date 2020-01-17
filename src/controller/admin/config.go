package admin

import (
    "OIUP-Backend/config"
    "OIUP-Backend/util"
    "github.com/gin-gonic/gin"
    "net/http"
)

func AdminGetConfigHandler(context *gin.Context) {
    util.SuccessResponse(context, config.Config)
}

func AdminModifyConfigHandler(context *gin.Context) {
    var configObj config.ConfigObject
    if err := context.BindJSON(&configObj); err != nil {
        util.ErrorResponse(context, http.StatusBadRequest, err.Error(), nil)
        return
    }

    err := config.ApplyConfig(configObj)
    if err != nil {
        util.ErrorResponse(context, http.StatusInternalServerError, err.Error(), nil)
        return
    }

    err = config.SaveConfig()
    if err != nil {
        util.ErrorResponse(context, http.StatusInternalServerError, err.Error(), nil)
        return
    }

    util.SuccessResponse(context, nil)
}

func AdminReloadConfigHandler(context *gin.Context) {
    err := config.LoadConfig()
    if err != nil {
        util.ErrorResponse(context, http.StatusInternalServerError, err.Error(), nil)
        return
    }

    util.SuccessResponse(context, nil)
}
