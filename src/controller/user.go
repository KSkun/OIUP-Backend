/*
   OIUP-Backend Project is developed by KSkun and licensed under GPL-3.0.
   Copyright (c) KSkun, 2020
*/
package controller

import (
    "OIUP-Backend/model"
    "OIUP-Backend/util"
    "github.com/gin-gonic/gin"
    "net/http"
)

type GetTokenRequest struct {
    ContestID string `form:"contest_id" binding:"required"`
    PersonID  string `form:"person_id" binding:"required"`
}

func UserTokenHandler(context *gin.Context) {
    var request GetTokenRequest
    if err := context.Bind(&request); err != nil {
        util.ErrorResponse(context, http.StatusBadRequest, "解析请求错误：" + err.Error(), nil)
        return
    }
    if !util.CheckContestID(request.ContestID) {
        util.ErrorResponse(context, http.StatusBadRequest, "考号不合法！", nil)
        return
    }

    user, found, err := model.GetUser(request.ContestID)
    if err != nil {
        util.ErrorResponse(context, http.StatusInternalServerError, "数据库错误：" + err.Error(), nil)
        return
    }
    if !found {
        util.ErrorResponse(context, http.StatusForbidden, "不存在该考生！", nil)
        return
    }
    if request.PersonID != user.PersonID {
        util.ErrorResponse(context, http.StatusForbidden, "证件号错误！", nil)
        return
    }

    token, err := util.NewUserJWTToken(user.ContestID)
    if err != nil {
        util.ErrorResponse(context, http.StatusInternalServerError, "JWT 错误：" + err.Error(), nil)
        return
    }

    util.SuccessResponse(context, gin.H{"token": token})
}

func UserInfoHandler(context *gin.Context) {
    contestID := util.GetIDFromContext(context)
    user, found, err := model.GetUser(contestID)
    if err != nil {
        util.ErrorResponse(context, http.StatusInternalServerError, "数据库错误：" + err.Error(), nil)
        return
    }
    if !found {
        util.ErrorResponse(context, http.StatusForbidden, "不存在该考生！", nil)
        return
    }

    util.SuccessResponse(context, user)
}