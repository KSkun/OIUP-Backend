/*
   OIUP-Backend Project is developed by KSkun and licensed under GPL-3.0.
   Copyright (c) KSkun, 2020
*/
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
    "strconv"
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
    ID int `form:"id" binding:"required"`
}

func ProblemStatusHandler(context *gin.Context) {
    var request RequestProblemID
    if err := context.Bind(&request); err != nil {
        util.ErrorResponse(context, http.StatusBadRequest, err.Error(), nil)
        return
    }

    _, found := config.GetProblemConfig(request.ID)
    if !found {
        util.ErrorResponse(context, http.StatusForbidden,
            "problem with id " + strconv.Itoa(request.ID) + " not found", nil)
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
    problem, found := config.GetProblemConfig(request.ID)
    if !found {
        util.ErrorResponse(context, http.StatusForbidden,
            "problem with id " + strconv.Itoa(request.ID) + " not found", nil)
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

    problem, found := config.GetProblemConfig(request.ID)
    if !found {
        util.ErrorResponse(context, http.StatusForbidden,
            "problem with id " + strconv.Itoa(request.ID) + " not found", nil)
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
    if problem.Type != config.ProblemAnswer {
        response["md5"] = submit.MD5
        file, err := ioutil.ReadFile(util.GetUploadPath(submit.ID) + "code")
        if err != nil {
            util.ErrorResponse(context, http.StatusInternalServerError, err.Error(), nil)
            return
        }
        response["data"] = string(file)
    }
    util.SuccessResponse(context, response)
}

type RequestCode struct {
    ID   int    `json:"id" binding:"required"`
    Code string `json:"code" binding:"required"`
}

func ProblemSubmitCodeHandler(context *gin.Context) {
    var request RequestCode
    if err := context.BindJSON(&request); err != nil {
        util.ErrorResponse(context, http.StatusBadRequest, err.Error(), nil)
        return
    }

    problem, found := config.GetProblemConfig(request.ID)
    if !found {
        util.ErrorResponse(context, http.StatusForbidden,
            "problem with id " + strconv.Itoa(request.ID) + " not found", nil)
        return
    }
    if problem.Type == config.ProblemAnswer {
        util.ErrorResponse(context, http.StatusForbidden,
            "problem requires submit output files", nil)
        return
    }

    submitID := uuid.NewV4()
    err := os.MkdirAll(util.GetUploadPath(submitID.String()), os.ModePerm)
    if err != nil {
        util.ErrorResponse(context, http.StatusInternalServerError, err.Error(), nil)
        return
    }
    err = ioutil.WriteFile(util.GetUploadPath(submitID.String()) + "code",
        []byte(request.Code), os.ModePerm)
    if err != nil {
        util.ErrorResponse(context, http.StatusInternalServerError, err.Error(), nil)
        return
    }

    user := util.GetIDFromContext(context)
    md5Res := md5.Sum([]byte(request.Code))
    md5Str := hex.EncodeToString(md5Res[:])
    err = model.AddCodeSubmit(submitID.String(), user, md5Str, time.Now(), request.ID)
    if err != nil {
        util.ErrorResponse(context, http.StatusInternalServerError, err.Error(), nil)
        return
    }

    util.SuccessResponse(context, gin.H{
        "id":   submitID.String(),
        "md5":  md5Str,
        "data": request.Code,
    })
}

type OutputInfo struct {
    TestID int    `json:"test_id" binding:"required"`
    Output string `json:"output" binding:"required"`
}

type RequestOutput struct {
    ID      int          `json:"id" binding:"required"`
    Outputs []OutputInfo `json:"outputs" binding:"required"`
}

func ProblemSubmitOutputHandler(context *gin.Context) {
    var request RequestOutput
    if err := context.BindJSON(&request); err != nil {
        util.ErrorResponse(context, http.StatusBadRequest, err.Error(), nil)
        return
    }

    problem, found := config.GetProblemConfig(request.ID)
    if !found {
        util.ErrorResponse(context, http.StatusForbidden,
            "problem with id " + strconv.Itoa(request.ID) + " not found", nil)
        return
    }
    if problem.Type != config.ProblemAnswer {
        util.ErrorResponse(context, http.StatusForbidden,
            "problem requires submit code file", nil)
        return
    }

    submitID := uuid.NewV4()
    err := os.MkdirAll(util.GetUploadPath(submitID.String()), os.ModePerm)
    if err != nil {
        util.ErrorResponse(context, http.StatusInternalServerError, err.Error(), nil)
        return
    }
    for _, output := range request.Outputs {
        err := ioutil.WriteFile(util.GetUploadPath(submitID.String()) + strconv.Itoa(output.TestID),
            []byte(output.Output), os.ModePerm)
        if err != nil {
            util.ErrorResponse(context, http.StatusInternalServerError, err.Error(), nil)
            return
        }
    }

    user := util.GetIDFromContext(context)
    md5Set := make([]model.MD5Info, 0)
    for _, output := range request.Outputs {
        md5Res := md5.Sum([]byte(output.Output))
        md5Str := hex.EncodeToString(md5Res[:])
        md5Set = append(md5Set, model.MD5Info{TestID: output.TestID, MD5: md5Str})
    }
    err = model.AddOutputSubmit(submitID.String(), user, md5Set, time.Now(), request.ID)
    if err != nil {
        util.ErrorResponse(context, http.StatusInternalServerError, err.Error(), nil)
        return
    }

    util.SuccessResponse(context, gin.H{
        "id":   submitID.String(),
        "md5":  md5Set,
        "data": "",
    })
}

type RequestConfirm struct {
    ID string `json:"id" binding:"required"`
}

func ProblemConfirmHandler(context *gin.Context) {
    var request RequestConfirm
    if err := context.BindJSON(&request); err != nil {
        util.ErrorResponse(context, http.StatusBadRequest, err.Error(), nil)
        return
    }

    contestID := util.GetIDFromContext(context)
    user, found, err := model.GetUser(contestID)
    if err != nil {
        util.ErrorResponse(context, http.StatusInternalServerError, err.Error(), nil)
        return
    }
    if !found {
        util.ErrorResponse(context, http.StatusForbidden, "user with contest_id " + request.ID + " not found", nil)
        return
    }

    submit, found, err := model.GetSubmit(request.ID)
    if err != nil {
        util.ErrorResponse(context, http.StatusInternalServerError, err.Error(), nil)
        return
    }
    if !found {
        util.ErrorResponse(context, http.StatusForbidden, "submit with id " + request.ID + " not found", nil)
        return
    }
    if submit.Confirm == model.SubmitConfirmed {
        util.ErrorResponse(context, http.StatusForbidden, "submit has been confirmed", nil)
        return
    }
    if submit.User != contestID {
        util.ErrorResponse(context, http.StatusForbidden, "can not confirm other's submission", nil)
        return
    }

    problem, _ := config.GetProblemConfig(submit.ProblemID)
    err = os.MkdirAll(util.GetSourcePath(contestID, problem.Filename), os.ModePerm)
    if err != nil {
        util.ErrorResponse(context, http.StatusInternalServerError, err.Error(), nil)
        return
    }
    if problem.Type != config.ProblemAnswer {
        data, err := ioutil.ReadFile(util.GetUploadPath(submit.ID) + "code")
        if err != nil {
            util.ErrorResponse(context, http.StatusInternalServerError, err.Error(), nil)
            return
        }

        suffix, err := util.GetCodeSuffix(user.Language)
        if err != nil {
            util.ErrorResponse(context, http.StatusInternalServerError, err.Error(), nil)
            return
        }
        err = ioutil.WriteFile(util.GetSourcePath(contestID, problem.Filename) + problem.Filename + suffix, data, os.ModePerm)
        if err != nil {
            util.ErrorResponse(context, http.StatusInternalServerError, err.Error(), nil)
            return
        }
    }
    if problem.Type == config.ProblemAnswer {
        for _, md5Info := range submit.MD5Set {
            data, err := ioutil.ReadFile(util.GetUploadPath(submit.ID) + strconv.Itoa(md5Info.TestID))
            if err != nil {
                util.ErrorResponse(context, http.StatusInternalServerError, err.Error(), nil)
                return
            }

            err = ioutil.WriteFile(util.GetSourcePath(contestID, problem.Filename) + problem.Filename + strconv.Itoa(md5Info.TestID) + ".out",
                data, os.ModePerm)
            if err != nil {
                util.ErrorResponse(context, http.StatusInternalServerError, err.Error(), nil)
                return
            }
        }
    }

    err = model.ConfirmSubmit(request.ID)
    if err != nil {
        util.ErrorResponse(context, http.StatusInternalServerError, err.Error(), nil)
        return
    }

    util.SuccessResponse(context, nil)
}
