package pkg

import (
	"github.com/gin-gonic/gin"
)

type ResponseError struct {
	StatusCode int
	Message    interface{}
}

func SendErrorResponse(ctx *gin.Context, err ResponseError) {
	ctx.AbortWithStatusJSON(err.StatusCode, gin.H{
		"message": err.Message,
	})
}
