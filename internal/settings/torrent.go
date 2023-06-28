package settings

import (
	rms_torrent "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-torrent"
	"github.com/RacoonMediaServer/rms-web/internal/ui"
	"github.com/gin-gonic/gin"
	"go-micro.dev/v4/logger"
	"google.golang.org/protobuf/types/known/emptypb"
	"net/http"
	"strconv"
)

type torrentPageContext struct {
	ui.PageContext
	Settings *rms_torrent.TorrentSettings
}

func (s *Service) torrentSettingsHandler(ctx *gin.Context) {
	settings, err := s.f.NewTorrent().GetSettings(ctx, &emptypb.Empty{})
	if err != nil {
		logger.Errorf("Get torrent settings failed: %s", err)
		ui.DisplayError(ctx, http.StatusInternalServerError, "Не удалось связаться с торрент-клиентом")
		return
	}
	page := torrentPageContext{
		PageContext: *ui.New(),
		Settings:    settings,
	}
	ctx.HTML(http.StatusOK, "settings.torrent.tmpl", &page)
}

func (s *Service) saveTorrentSettingsHandler(ctx *gin.Context) {
	downloadLimit, err := strconv.ParseUint(ctx.PostForm("downloadLimit"), 10, 64)
	if err != nil {
		ui.DisplayError(ctx, http.StatusBadRequest, "Указано неверное значение")
		return
	}
	uploadLimit, err := strconv.ParseUint(ctx.PostForm("uploadLimit"), 10, 64)
	if err != nil {
		ui.DisplayError(ctx, http.StatusBadRequest, "Указано неверное значение")
		return
	}

	settings := rms_torrent.TorrentSettings{
		DownloadLimit: uint64(float32(downloadLimit*1024*1024) / 8.),
		UploadLimit:   uint64(float32(uploadLimit*1024*1024) / 8.),
	}

	_, err = s.f.NewTorrent().SetSettings(ctx, &settings)
	if err != nil {
		ui.DisplayError(ctx, http.StatusInternalServerError, "Не удалось применить настройки")
		return
	}

	ui.DisplayOK(ctx, "Сохранено", "/settings/torrent")
}
