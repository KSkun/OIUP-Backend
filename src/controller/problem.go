/*
   OIUP-Backend Project is developed by KSkun and licensed under GPL-3.0.
   Copyright (c) KSkun, 2020
*/
package controller

import (
    "OIUP-Backend/config"
    "OIUP-Backend/model"
    "OIUP-Backend/util"
    "crypto/sha256"
    "encoding/hex"
    "github.com/gin-gonic/gin"
    uuid "github.com/satori/go.uuid"
    "io/ioutil"
    "net/http"
    "os"
    "strconv"
    "time"
)

type RequestProblemID struct {
    ID int `form:"id" binding:"required"`
}

func ProblemStatusHandler(context *gin.Context) {
    var request RequestProblemID
    if err := context.Bind(&request); err != nil {
        util.ErrorResponse(context, http.StatusBadRequest, "解析请求错误：" + err.Error(), nil)
        return
    }

    _, found := config.GetProblemConfig(request.ID)
    if !found {
        util.ErrorResponse(context, http.StatusForbidden,
            "找不到编号为 " + strconv.Itoa(request.ID) + " 的题目！", nil)
        return
    }

    user := util.GetIDFromContext(context)
    _, found, err := model.GetLatestSubmit(user, request.ID)
    if err != nil {
        util.ErrorResponse(context, http.StatusInternalServerError, "数据库错误：" + err.Error(), nil)
        return
    }
    util.SuccessResponse(context, gin.H{"is_submit": found})
}

func ProblemListHandler(context *gin.Context) {
    util.SuccessResponse(context, gin.H{"problems": config.Config.Contest.ProblemSet})
}

func ProblemLatestHandler(context *gin.Context) {
    var request RequestProblemID
    if err := context.Bind(&request); err != nil {
        util.ErrorResponse(context, http.StatusBadRequest, "解析请求错误：" + err.Error(), nil)
        return
    }

    problem, found := config.GetProblemConfig(request.ID)
    if !found {
        util.ErrorResponse(context, http.StatusForbidden,
            "找不到编号为 " + strconv.Itoa(request.ID) + " 的题目！", nil)
        return
    }

    user := util.GetIDFromContext(context)
    submit, found, err := model.GetLatestSubmit(user, request.ID)
    if err != nil {
        util.ErrorResponse(context, http.StatusInternalServerError, "数据库错误：" + err.Error(), nil)
        return
    }
    if !found {
        util.ErrorResponse(context, http.StatusForbidden,
            "您还未提交解答至编号为 " + strconv.Itoa(request.ID) + " 的题目！", nil)
        return
    }

    response := gin.H{"id": submit.ID, "hash": submit.HashSet, "data": ""}
    if problem.Type != config.ProblemAnswer {
        response["hash"] = submit.Hash
        file, err := ioutil.ReadFile(util.GetUploadPath(submit.ID) + "code")
        if err != nil {
            util.ErrorResponse(context, http.StatusInternalServerError, "文件错误：" + err.Error(), nil)
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
        util.ErrorResponse(context, http.StatusBadRequest, "解析请求错误：" + err.Error(), nil)
        return
    }
    if getContestStatus() != ContestStatusRunning {
        util.ErrorResponse(context, http.StatusForbidden, "现在不允许提交解答！", nil)
        return
    }

    problem, found := config.GetProblemConfig(request.ID)
    if !found {
        util.ErrorResponse(context, http.StatusForbidden,
            "找不到编号为 " + strconv.Itoa(request.ID) + " 的题目！", nil)
        return
    }
    if problem.Type == config.ProblemAnswer {
        util.ErrorResponse(context, http.StatusForbidden, "题目要求提交输出文件！", nil)
        return
    }

    submitID := uuid.NewV4()
    err := os.MkdirAll(util.GetUploadPath(submitID.String()), os.ModePerm)
    if err != nil {
        util.ErrorResponse(context, http.StatusInternalServerError, "文件错误：" + err.Error(), nil)
        return
    }
    err = ioutil.WriteFile(util.GetUploadPath(submitID.String()) + "code",
        []byte(request.Code), os.ModePerm)
    if err != nil {
        util.ErrorResponse(context, http.StatusInternalServerError, "文件错误：" + err.Error(), nil)
        return
    }

    user := util.GetIDFromContext(context)
    hashRes := sha256.Sum256([]byte(request.Code))
    hashStr := hex.EncodeToString(hashRes[:])
    err = model.AddCodeSubmit(submitID.String(), user, hashStr, time.Now(), request.ID)
    if err != nil {
        util.ErrorResponse(context, http.StatusInternalServerError, "数据库错误：" + err.Error(), nil)
        return
    }

    util.SuccessResponse(context, gin.H{
        "id":   submitID.String(),
        "hash": hashStr,
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
        util.ErrorResponse(context, http.StatusBadRequest, "解析请求错误：" + err.Error(), nil)
        return
    }
    if getContestStatus() != ContestStatusRunning {
        util.ErrorResponse(context, http.StatusForbidden, "现在不允许提交解答！", nil)
        return
    }

    problem, found := config.GetProblemConfig(request.ID)
    if !found {
        util.ErrorResponse(context, http.StatusForbidden,
            "找不到编号为 " + strconv.Itoa(request.ID) + " 的题目！", nil)
        return
    }
    if problem.Type != config.ProblemAnswer {
        util.ErrorResponse(context, http.StatusForbidden, "题目要求提交源代码文件！", nil)
        return
    }
    testID := make(map[int]int, 0)
    for _, output := range request.Outputs {
        if output.TestID < 1 || output.TestID > problem.Testcase {
            util.ErrorResponse(context, http.StatusBadRequest,
                "测试点编号 " + strconv.Itoa(output.TestID) + " 超出范围！", nil)
            return
        }
        if testID[output.TestID] == 1 {
            util.ErrorResponse(context, http.StatusBadRequest,
                "测试点编号 " + strconv.Itoa(output.TestID) + " 重复上传！", nil)
            return
        }
        testID[output.TestID] = 1
    }

    submitID := uuid.NewV4()
    err := os.MkdirAll(util.GetUploadPath(submitID.String()), os.ModePerm)
    if err != nil {
        util.ErrorResponse(context, http.StatusInternalServerError, "文件错误：" + err.Error(), nil)
        return
    }
    for _, output := range request.Outputs {
        err := ioutil.WriteFile(util.GetUploadPath(submitID.String()) + strconv.Itoa(output.TestID),
            []byte(output.Output), os.ModePerm)
        if err != nil {
            util.ErrorResponse(context, http.StatusInternalServerError, "文件错误：" + err.Error(), nil)
            return
        }
    }

    user := util.GetIDFromContext(context)
    hashSet := make([]model.HashInfo, 0)
    for _, output := range request.Outputs {
        hashRes := sha256.Sum256([]byte(output.Output))
        hashStr := hex.EncodeToString(hashRes[:])
        hashSet = append(hashSet, model.HashInfo{TestID: output.TestID, Hash: hashStr})
    }
    err = model.AddOutputSubmit(submitID.String(), user, hashSet, time.Now(), request.ID)
    if err != nil {
        util.ErrorResponse(context, http.StatusInternalServerError, "数据库错误：" + err.Error(), nil)
        return
    }

    util.SuccessResponse(context, gin.H{
        "id":   submitID.String(),
        "hash": hashSet,
        "data": "",
    })
}

type RequestConfirm struct {
    ID string `json:"id" binding:"required"`
}

func ProblemConfirmHandler(context *gin.Context) {
    var request RequestConfirm
    if err := context.BindJSON(&request); err != nil {
        util.ErrorResponse(context, http.StatusBadRequest, "解析请求错误：" + err.Error(), nil)
        return
    }

    contestID := util.GetIDFromContext(context)
    user, found, err := model.GetUser(contestID)
    if err != nil {
        util.ErrorResponse(context, http.StatusInternalServerError, "数据库错误：" + err.Error(), nil)
        return
    }
    if !found {
        util.ErrorResponse(context, http.StatusForbidden, "找不到考号为 " + request.ID + " 的考生信息！", nil)
        return
    }

    submit, found, err := model.GetSubmit(request.ID)
    if err != nil {
        util.ErrorResponse(context, http.StatusInternalServerError, "数据库错误：" + err.Error(), nil)
        return
    }
    if !found {
        util.ErrorResponse(context, http.StatusForbidden, "找不到编号为 " + request.ID + " 的提交信息！", nil)
        return
    }
    if submit.Confirm == model.SubmitConfirmed {
        util.ErrorResponse(context, http.StatusForbidden, "提交已被确认！", nil)
        return
    }
    if submit.User != contestID {
        util.ErrorResponse(context, http.StatusForbidden, "不能确认其他考生的提交！", nil)
        return
    }

    problem, _ := config.GetProblemConfig(submit.ProblemID)
    err = os.MkdirAll(util.GetSourcePath(contestID, problem.Filename), os.ModePerm)
    if err != nil {
        util.ErrorResponse(context, http.StatusInternalServerError, "文件错误：" + err.Error(), nil)
        return
    }
    if problem.Type != config.ProblemAnswer {
        data, err := ioutil.ReadFile(util.GetUploadPath(submit.ID) + "code")
        if err != nil {
            util.ErrorResponse(context, http.StatusInternalServerError, "文件错误：" + err.Error(), nil)
            return
        }

        suffix, err := util.GetCodeSuffix(user.Language)
        if err != nil {
            util.ErrorResponse(context, http.StatusInternalServerError, "错误：" + err.Error(), nil)
            return
        }
        err = ioutil.WriteFile(util.GetSourcePath(contestID, problem.Filename) + problem.Filename + suffix, data, os.ModePerm)
        if err != nil {
            util.ErrorResponse(context, http.StatusInternalServerError, "文件错误：" + err.Error(), nil)
            return
        }
    }
    if problem.Type == config.ProblemAnswer {
        for _, hashInfo := range submit.HashSet {
            data, err := ioutil.ReadFile(util.GetUploadPath(submit.ID) + strconv.Itoa(hashInfo.TestID))
            if err != nil {
                util.ErrorResponse(context, http.StatusInternalServerError, "文件错误：" + err.Error(), nil)
                return
            }

            err = ioutil.WriteFile(util.GetSourcePath(contestID, problem.Filename) + problem.Filename +
                strconv.Itoa(hashInfo.TestID) + ".out", data, os.ModePerm)
            if err != nil {
                util.ErrorResponse(context, http.StatusInternalServerError, "文件错误：" + err.Error(), nil)
                return
            }
        }
    }

    err = model.ConfirmSubmit(request.ID)
    if err != nil {
        util.ErrorResponse(context, http.StatusInternalServerError, "数据库错误：" + err.Error(), nil)
        return
    }

    util.SuccessResponse(context, nil)
}
