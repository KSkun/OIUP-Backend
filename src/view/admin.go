/*
   OIUP-Backend Project is developed by KSkun and licensed under GPL-3.0.
   Copyright (c) KSkun, 2020
*/
package view

import (
    "OIUP-Backend/controller/admin"
    "OIUP-Backend/middleware"
    "github.com/gin-gonic/gin"
)

func InitAdminView(group *gin.RouterGroup) {
    group.Use(middleware.ValidateBackstageKey)

    configGroup := group.Group("/config")
    configGroup.GET("", admin.AdminGetConfigHandler)
    configGroup.PUT("", admin.AdminModifyConfigHandler)
    configGroup.POST("/reload", admin.AdminReloadConfigHandler)

    userGroup := group.Group("/user")
    userGroup.GET("", admin.AdminSearchUserHandler)
    userGroup.POST("", admin.AdminAddUserHandler)
    userGroup.PUT("", admin.AdminUpdateUserHandler)
    userGroup.DELETE("", admin.AdminDeleteUserHandler)
    userGroup.POST("/csv", admin.AdminImportUsersHandler)

    // TODO 提交记录
}
