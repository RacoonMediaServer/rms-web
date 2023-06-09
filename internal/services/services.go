package services

import (
	"fmt"
	"github.com/RacoonMediaServer/rms-web/internal/config"
	"github.com/RacoonMediaServer/rms-web/internal/ui"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func Register(router *gin.RouterGroup) {
	router.GET("/", catalogHandler)
	router.GET("/:id", goToService)
}

func catalogHandler(ctx *gin.Context) {
	page := ui.CatalogPageContext{
		PageContext: *ui.New(),
		Title:       "Сервисы",
	}

	services := config.Config().Services
	for i, service := range services {
		page.Parts = append(page.Parts, ui.CatalogPart{
			Image:       fmt.Sprintf("/img/%s.png", strings.ToLower(service.Title)),
			Title:       service.Title,
			Link:        fmt.Sprintf("/services/%d", i),
			Description: service.Description,
		})
	}

	page.Display(ctx)
}

func goToService(ctx *gin.Context) {
	services := config.Config().Services
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 32)
	if err != nil || id < 0 || id >= int64(len(services)) {
		ui.DisplayError(ctx, http.StatusNotFound, "Попытка обратиться к несуществующему сервису")
		return
	}
	u, err := url.Parse(services[id].Address)
	if err != nil {
		ctx.Redirect(http.StatusFound, services[id].Address)
		return
	}
	// TODO: rewrite IP address
	ctx.Redirect(http.StatusFound, u.String())
}
