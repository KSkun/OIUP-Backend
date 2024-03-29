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

type RequestSearchSubmit struct {
    ContestID string `form:"contest_id"`
    ProblemID int    `form:"problem_id"`
}

func AdminSearchSubmitHandler(context *gin.Context) {
    var request RequestSearchSubmit
    if err := context.Bind(&request); err != nil {
        util.ErrorResponse(context, http.StatusBadRequest, "解析请求错误：" + err.Error(), nil)
        return
    }

    filters := make(map[string]interface{}, 0)
    if len(request.ContestID) > 0 {
        filters["user"] = request.ContestID
    }
    if request.ProblemID != 0 {
        filters["problem_id"] = request.ProblemID
    }

    submits, err := model.SearchSubmit(filters)
    if err != nil {
        util.ErrorResponse(context, http.StatusInternalServerError, "数据库错误：" + err.Error(), nil)
        return
    }

    util.SuccessResponse(context, gin.H{"result": submits})
}
