package ui

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ErrorContext struct {
	PageContext
	Error string
}

func DisplayError(ctx *gin.Context, err string) {
	ctx.HTML(http.StatusOK, "error.tmpl", &ErrorContext{
		PageContext: *New(),
		Error:       err,
	})
}
