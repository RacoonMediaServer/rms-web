package ui

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type OkPageContext struct {
	PageContext
	Text string
}

func DisplayOK(ctx *gin.Context, text, redirect string) {
	page := OkPageContext{
		PageContext: *New(),
		Text:        text,
	}
	page.Redirect = redirect
	ctx.HTML(http.StatusOK, "success.tmpl", &page)
}
