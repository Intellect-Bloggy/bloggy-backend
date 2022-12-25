package http

import (
    "github.com/gin-gonic/gin"
    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
)

type docsHandler struct {
}

func NewDocsHandler() *docsHandler {
    return &docsHandler{}
}

func (h *docsHandler) Register(router *gin.Engine) {
    router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
