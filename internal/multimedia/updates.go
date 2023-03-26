package multimedia

import (
	rms_library "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-library"
	"github.com/RacoonMediaServer/rms-web/internal/ui"
	"github.com/gin-gonic/gin"
	"go-micro.dev/v4/logger"
	"google.golang.org/protobuf/types/known/emptypb"
	"net/http"
)

type updatesPageContext struct {
	ui.PageContext
	Updates []*rms_library.TvSeriesUpdate
}

func (s *Service) updatesHandler(ctx *gin.Context) {
	resp, err := s.f.NewLibrary().GetTvSeriesUpdates(ctx, &emptypb.Empty{})
	if err != nil {
		logger.Errorf("Get updates failed: %s", err)
		ui.DisplayError(ctx, http.StatusInternalServerError, "Не удалось получить список новых сезонов от сервиса библиотеки")
		return
	}
	page := updatesPageContext{
		PageContext: *ui.New(),
		Updates:     resp.Updates,
	}
	ctx.HTML(http.StatusOK, "multimedia.updates.tmpl", &page)
}
