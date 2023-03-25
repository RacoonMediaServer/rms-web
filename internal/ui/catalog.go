package ui

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type CatalogPart struct {
	Image       string
	Title       string
	Link        string
	Description string
}

type CatalogPageContext struct {
	PageContext
	Title string
	Parts []CatalogPart
}

func (p CatalogPageContext) Display(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "catalog.tmpl", &p)
}
