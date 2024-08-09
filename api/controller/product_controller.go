package controller

import "github.com/gin-gonic/gin"

type ProductController interface {
	GetAll(ctx *gin.Context)
	GetByID(ctx *gin.Context)
	UpdateData(ctx *gin.Context)
}
