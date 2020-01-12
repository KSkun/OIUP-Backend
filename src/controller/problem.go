package controller

import (
    "OIUP-Backend/config"
    "OIUP-Backend/model"
    "OIUP-Backend/util"
    "crypto/md5"
    "encoding/hex"
    "github.com/gin-gonic/gin"
    uuid "github.com/satori/go.uuid"
    "io/ioutil"
    "net/http"
    "os"
    "time"
)

func ProblemListHandler(context *gin.Context) {
    response := make([]gin.H, 0)
    for _, problem := range config.Config.Contest.Problems {
        response = append(response, gin.H{"id": problem.ID, "filename": problem.Filename})
    }
    util.SuccessResponse(context, response)
}

type RequestProblemID struct {
    ID int
}

func ProblemStatusHandler(context *gin.Context) {
    var request RequestProblemID
    if err := context.Bind(&request); err != nil {
        util.ErrorResponse(context, http.StatusBadRequest, err.Error(), nil)
        return
    }

    user := util.GetIDFromContext(context)
    _, found, err := model.GetLatestSubmit(user, request.ID)
    if err != nil {
        util.ErrorResponse(context, http.StatusInternalServerError, err.Error(), nil)
        return
    }
    util.SuccessResponse(context, gin.H{"is_submit": found})
}

func ProblemInfoHandler(context *gin.Context) {
    var request RequestProblemID
    if err := context.Bind(&request); err != nil {
        util.ErrorResponse(context, http.StatusBadRequest, err.Error(), nil)
        return
    }
    problem, err := config.GetProblemConfig(request.ID)
    if err != nil {
        util.ErrorResponse(context, http.StatusInternalServerError, err.Error(), nil)
        return
    }
    util.SuccessResponse(context, problem)
}

func ProblemLatestHandler(context *gin.Context) {
    var request RequestProblemID
    if err := context.Bind(&request); err != nil {
        util.ErrorResponse(context, http.StatusBadRequest, err.Error(), nil)
        return
    }

    user := util.GetIDFromContext(context)
    submit, found, err := model.GetLatestSubmit(user, request.ID)
    if err != nil {
        util.ErrorResponse(context, http.StatusInternalServerError, err.Error(), nil)
        return
    }
    if !found {
        util.ErrorResponse(context, http.StatusForbidden, "user has not submitted to this problem", nil)
        return
    }

    response := gin.H{"id": submit.ID, "md5": submit.MD5Set, "data": ""}
    problem, err := config.GetProblemConfig(request.ID)
    if err != nil {
        util.ErrorResponse(context, http.StatusInternalServerError, err.Error(), nil)
        return
    }
    if problem.Type != config.ProblemAnswer {
        response["md5"] = submit.MD5
        file, err := ioutil.ReadFile(config.Config.File.DirectoryUpload + "/" + submit.ID + "/code")
        if err != nil {
            util.ErrorResponse(context, http.StatusInternalServerError, err.Error(), nil)
            return
        }
        response["data"] = string(file)
    }
    util.SuccessResponse(context, response)
}

type RequestCode struct {
    ID   int    `json:"id"`
    Code string `json:"code"`
}

func ProblemSubmitCodeHandler(context *gin.Context) {
    var request RequestCode
    if err := context.BindJSON(&request); err != nil {
        util.ErrorResponse(context, http.StatusBadRequest, err.Error(), nil)
        return
    }

    _, err := config.GetProblemConfig(request.ID)
    if err != nil {
        util.ErrorResponse(context, http.StatusInternalServerError, err.Error(), nil)
        return
    }

    submitID := uuid.NewV4()
    err = ioutil.WriteFile(config.Config.File.DirectoryUpload + "/" + submitID.String() + "/code",
        []byte(request.Code), os.ModePerm)
    if err != nil {
        util.ErrorResponse(context, http.StatusInternalServerError, err.Error(), nil)
        return
    }

    user := util.GetIDFromContext(context)
    md5Str := md5.Sum([]byte(request.Code))
    err = model.AddCodeSubmit(submitID.String(), user, hex.EncodeToString(md5Str[:]), time.Now(), request.ID)
    if err != nil {
        util.ErrorResponse(context, http.StatusInternalServerError, err.Error(), nil)
        return
    }
}

type OutputInfo struct {
    TestID int    `json:"test_id"`
    Output string `json:"output"`
}

type RequestOutput struct {
    ID      int          `json:"id"`
    Outputs []OutputInfo `json:"outputs"`
}

func ProblemSubmitOutputHandler(context *gin.Context) {
    var request RequestCode
    if err := context.BindJSON(&request); err != nil {
        util.ErrorResponse(context, http.StatusBadRequest, err.Error(), nil)
        return
    }

    _, err := config.GetProblemConfig(request.ID)
    if err != nil {
        util.ErrorResponse(context, http.StatusInternalServerError, err.Error(), nil)
        return
    }

    // TODO Submit Output
    /*submitID := uuid.NewV4()
    err = ioutil.WriteFile(config.Config.File.DirectoryUpload + "/" + submitID.String() + "/code",
        []byte(request.Code), os.ModePerm)
    if err != nil {
        util.ErrorResponse(context, http.StatusInternalServerError, err.Error(), nil)
        return
    }

    user := util.GetIDFromContext(context)
    md5Str := md5.Sum([]byte(request.Code))
    err = model.AddCodeSubmit(submitID.String(), user, hex.EncodeToString(md5Str[:]), time.Now(), request.ID)
    if err != nil {
        util.ErrorResponse(context, http.StatusInternalServerError, err.Error(), nil)
        return
    }*/
}
