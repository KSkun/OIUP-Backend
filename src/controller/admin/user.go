/*
   OIUP-Backend Project is developed by KSkun and licensed under GPL-3.0.
   Copyright (c) KSkun, 2020
*/
package admin

import (
    "OIUP-Backend/model"
    "OIUP-Backend/util"
    "github.com/gin-gonic/gin"
    "net/http"
)

type RequestSearchUser struct {
    Name      string `form:"name"`
    School    string `form:"school"`
    ContestID string `form:"contest_id"`
    PersonID  string `form:"person_id"`
    Language  int    `form:"language"`
}

func AdminSearchUserHandler(context *gin.Context) {
    var request RequestSearchUser
    if err := context.Bind(&request); err != nil {
        util.ErrorResponse(context, http.StatusBadRequest, err.Error(), nil)
        return
    }

    filters := make(map[string]interface{}, 0)
    if len(request.Name) > 0 {
        filters["name"] = request.Name
    }
    if len(request.School) > 0 {
        filters["school"] = request.School
    }
    if len(request.ContestID) > 0 {
        filters["contest_id"] = request.ContestID
    }
    if len(request.PersonID) > 0 {
        filters["person_id"] = request.PersonID
    }
    if request.Language >= 1 && request.Language <= 3 {
        filters["language"] = request.Language
    }

    users, err := model.SearchUser(filters)
    if err != nil {
        util.ErrorResponse(context, http.StatusInternalServerError, err.Error(), nil)
        return
    }

    util.SuccessResponse(context, gin.H{"result": users})
}

func AdminAddUserHandler(context *gin.Context) {
    var request model.UserInfo
    if err := context.BindJSON(&request); err != nil {
        util.ErrorResponse(context, http.StatusBadRequest, err.Error(), nil)
        return
    }
    if !util.CheckContestID(request.ContestID) {
        util.ErrorResponse(context, http.StatusBadRequest, "invalid contest_id", nil)
        return
    }
    if request.Language < 1 || request.Language > 3 {
        util.ErrorResponse(context, http.StatusBadRequest, "invalid language", nil)
        return
    }

    err := model.AddUser(request)
    if err != nil {
        util.ErrorResponse(context, http.StatusInternalServerError, err.Error(), nil)
        return
    }

    util.SuccessResponse(context, nil)
}

func AdminUpdateUserHandler(context *gin.Context) {
    var request model.UserInfo
    if err := context.BindJSON(&request); err != nil {
        util.ErrorResponse(context, http.StatusBadRequest, err.Error(), nil)
        return
    }

    err := model.UpdateUser(request)
    if err != nil {
        util.ErrorResponse(context, http.StatusInternalServerError, err.Error(), nil)
        return
    }

    util.SuccessResponse(context, nil)
}

type RequestUserContestID struct {
    ContestID string `form:"contest_id" binding:"required"`
}

func AdminDeleteUserHandler(context *gin.Context) {
    var request RequestUserContestID
    if err := context.Bind(&request); err != nil {
        util.ErrorResponse(context, http.StatusBadRequest, err.Error(), nil)
        return
    }

    err := model.DeleteUser(request.ContestID)
    if err != nil {
        util.ErrorResponse(context, http.StatusInternalServerError, err.Error(), nil)
        return
    }

    util.SuccessResponse(context, nil)
}

func AdminImportUsersHandler(context *gin.Context) {
    // TODO 根据指定csv导入选手
}
