package services

import (
	"github.com/RacoonMediaServer/rms-web/internal/config"
	"github.com/RacoonMediaServer/rms-web/internal/ui"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"strconv"
)

func Register(router *gin.RouterGroup) {
	router.GET("/:id", goToService)
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
