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
    Page      int    `form:"page"`
    ContestID string `form:"contest_id"`
    ProblemID int    `form:"problem_id"`
}

func AdminSearchSubmitHandler(context *gin.Context) {
    var request RequestSearchSubmit
    if err := context.Bind(&request); err != nil {
        util.ErrorResponse(context, http.StatusBadRequest, err.Error(), nil)
        return
    }

    filters := make(map[string]interface{}, 0)
    if len(request.ContestID) > 0 {
        filters["user"] = request.ContestID
    }
    if request.ProblemID != 0 {
        filters["problem_id"] = request.ProblemID
    }

    submits, count, err := model.SearchSubmit(filters, request.Page)
    if err != nil {
        util.ErrorResponse(context, http.StatusInternalServerError, err.Error(), nil)
        return
    }

    util.SuccessResponse(context, gin.H{"count": count, "result": submits})
}
