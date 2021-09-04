package controller

import (
	"github.com/gin-gonic/gin"
)

type IController interface {
	Handle(ctx *gin.Context)
}
