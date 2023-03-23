package settings

import (
	"fmt"
	"github.com/RacoonMediaServer/rms-web/internal/config"
	"github.com/RacoonMediaServer/rms-web/internal/ui"
	"github.com/gin-gonic/gin"
	"go-micro.dev/v4/logger"
	"google.golang.org/protobuf/types/known/emptypb"
	"net/http"
)

type telegramPageContext struct {
	ui.PageContext
	BotLink            string
	BotID              string
	IdentificationCode string
}

func (s *Service) getTelegramSettings(ctx *gin.Context) {
	resp, err := s.f.NewBotClient().GetIdentificationCode(ctx, &emptypb.Empty{})
	if err != nil {
		logger.Errorf("Get identification code failed: %s", err)
		ui.DisplayError(ctx, "Не удалось получить код идентификации от удаленного сервера")
		return
	}
	page := telegramPageContext{
		PageContext:        *ui.New(),
		BotLink:            fmt.Sprintf("https://t.me/%s", config.Config().Bot),
		BotID:              config.Config().Bot,
		IdentificationCode: resp.Code,
	}
	ctx.HTML(http.StatusOK, "settings.telegram.tmpl", &page)
}
