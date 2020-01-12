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
        util.ErrorResponse(context, http.StatusBadRequest, err.Error(), nil)
        return
    }
    if !util.CheckContestID(request.ContestID) {
        util.ErrorResponse(context, http.StatusBadRequest, "invalid contest_id", nil)
        return
    }

    user, found, err := model.GetUser(request.ContestID)
    if err != nil {
        util.ErrorResponse(context, http.StatusInternalServerError, err.Error(), nil)
        return
    }
    if !found {
        util.ErrorResponse(context, http.StatusForbidden, "user not found", nil)
        return
    }
    if request.PersonID != user.PersonID {
        util.ErrorResponse(context, http.StatusForbidden, "wrong person_id", nil)
        return
    }

    token, err := util.NewUserJWTToken(user.ContestID)
    if err != nil {
        util.ErrorResponse(context, http.StatusInternalServerError, err.Error(), nil)
        return
    }

    util.SuccessResponse(context, gin.H{"token": token})
}

func UserInfoHandler(context *gin.Context) {
    contestID := util.GetIDFromContext(context)
    user, found, err := model.GetUser(contestID)
    if err != nil {
        util.ErrorResponse(context, http.StatusInternalServerError, err.Error(), nil)
        return
    }
    if !found {
        util.ErrorResponse(context, http.StatusForbidden, "user not found", nil)
        return
    }

    util.SuccessResponse(context, user)
}