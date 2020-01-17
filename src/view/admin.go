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

    // TODO 选手和提交记录
}
