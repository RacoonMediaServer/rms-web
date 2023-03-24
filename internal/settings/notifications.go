package settings

import (
	"github.com/RacoonMediaServer/rms-web/internal/ui"
	"github.com/gin-gonic/gin"
	"net/http"
)

type notifyPageContext struct {
	ui.PageContext
}

func (s *Service) getNotificationsHandler(ctx *gin.Context) {
	page := notifyPageContext{
		PageContext: *ui.New(),
	}
	ctx.HTML(http.StatusOK, "settings.notifications.tmpl", &page)
}
