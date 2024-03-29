/*
   OIUP-Backend Project is developed by KSkun and licensed under GPL-3.0.
   Copyright (c) KSkun, 2020
*/
package admin

import (
    "OIUP-Backend/model"
    "OIUP-Backend/util"
    "encoding/csv"
    "github.com/gin-gonic/gin"
    uuid "github.com/satori/go.uuid"
    "io"
    "io/ioutil"
    "net/http"
    "os"
    "strings"
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
        util.ErrorResponse(context, http.StatusBadRequest, "解析请求错误：" + err.Error(), nil)
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
        util.ErrorResponse(context, http.StatusInternalServerError, "数据库错误：" + err.Error(), nil)
        return
    }

    util.SuccessResponse(context, gin.H{"result": users})
}

func AdminAddUserHandler(context *gin.Context) {
    var request model.UserInfo
    if err := context.BindJSON(&request); err != nil {
        util.ErrorResponse(context, http.StatusBadRequest, "解析请求错误：" + err.Error(), nil)
        return
    }
    if !util.CheckContestID(request.ContestID) {
        util.ErrorResponse(context, http.StatusBadRequest, "考号不合法！", nil)
        return
    }
    if request.Language < 1 || request.Language > 3 {
        util.ErrorResponse(context, http.StatusBadRequest, "语言类型不合法！", nil)
        return
    }

    err := model.AddUser(request)
    if err != nil {
        util.ErrorResponse(context, http.StatusInternalServerError, "数据库错误：" + err.Error(), nil)
        return
    }

    util.SuccessResponse(context, nil)
}

func AdminUpdateUserHandler(context *gin.Context) {
    var request model.UserInfo
    if err := context.BindJSON(&request); err != nil {
        util.ErrorResponse(context, http.StatusBadRequest, "解析请求错误：" + err.Error(), nil)
        return
    }

    err := model.UpdateUser(request)
    if err != nil {
        util.ErrorResponse(context, http.StatusInternalServerError, "数据库错误：" + err.Error(), nil)
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
        util.ErrorResponse(context, http.StatusBadRequest, "解析请求错误：" + err.Error(), nil)
        return
    }

    err := model.DeleteUser(request.ContestID)
    if err != nil {
        util.ErrorResponse(context, http.StatusInternalServerError, "数据库错误：" + err.Error(), nil)
        return
    }

    util.SuccessResponse(context, nil)
}

func AdminImportUsersHandler(context *gin.Context) {
    file, err := context.FormFile("data")
    if err != nil {
        util.ErrorResponse(context, http.StatusInternalServerError, "文件错误：" + err.Error(), nil)
        return
    }
    err = os.MkdirAll(util.GetTempPath(""), os.ModePerm)
    if err != nil {
        util.ErrorResponse(context, http.StatusInternalServerError, "文件错误：" + err.Error(), nil)
        return
    }
    filename := uuid.NewV4().String() + ".csv"
    err = context.SaveUploadedFile(file, util.GetTempPath(filename))
    if err != nil {
        util.ErrorResponse(context, http.StatusInternalServerError, "文件错误：" + err.Error(), nil)
        return
    }

    data, err := ioutil.ReadFile(util.GetTempPath(filename))
    if err != nil {
        util.ErrorResponse(context, http.StatusInternalServerError, "文件错误：" + err.Error(), nil)
        return
    }
    reader := csv.NewReader(strings.NewReader(string(data)))
    for {
        line, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            util.ErrorResponse(context, http.StatusInternalServerError, "文件错误：" + err.Error(), nil)
            return
        }

        err = model.AddUser(model.UserInfo{
            Name:      line[2],
            School:    line[4],
            ContestID: line[0],
            PersonID:  line[3],
            Language:  model.LanguageCPlusPlus, // manually change it if necessary
        })
        if err != nil {
            util.ErrorResponse(context, http.StatusInternalServerError, "数据库错误：" + err.Error(), nil)
            return
        }
    }
    err = os.Remove(util.GetTempPath(filename))
    if err != nil {
        util.ErrorResponse(context, http.StatusInternalServerError, "文件错误：" + err.Error(), nil)
        return
    }

    util.SuccessResponse(context, nil)
}
