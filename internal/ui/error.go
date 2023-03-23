package ui

import (
	"github.com/gin-gonic/gin"
)

type ErrorContext struct {
	PageContext
	Error string
}

func DisplayError(ctx *gin.Context, status int, err string) {
	ctx.HTML(status, "error.tmpl", &ErrorContext{
		PageContext: *New(),
		Error:       err,
	})
}
