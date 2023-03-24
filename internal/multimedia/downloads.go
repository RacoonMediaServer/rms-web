package multimedia

import (
	rms_torrent "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-torrent"
	"github.com/RacoonMediaServer/rms-web/internal/ui"
	"github.com/gin-gonic/gin"
	"go-micro.dev/v4/logger"
	"net/http"
	"sort"
)

type downloadsPageContext struct {
	ui.PageContext
	Torrents []*rms_torrent.TorrentInfo
}

func reorderTorrentStatus(status rms_torrent.Status) int {
	switch status {
	case rms_torrent.Status_Downloading:
		return 0
	case rms_torrent.Status_Pending:
		return 1
	case rms_torrent.Status_Failed:
		return 2
	case rms_torrent.Status_Done:
		return 3
	default:
		return 4
	}
}

func (s *Service) downloadsHandler(ctx *gin.Context) {
	req := rms_torrent.GetTorrentsRequest{IncludeDoneTorrents: true}
	resp, err := s.f.NewTorrent().GetTorrents(ctx, &req)
	if err != nil {
		logger.Errorf("Get torrents failed: %s", err)
		ui.DisplayError(ctx, http.StatusInternalServerError, "Не удалось обратиться к сервису загрузок")
		return
	}
	sort.SliceStable(resp.Torrents, func(i, j int) bool {
		return reorderTorrentStatus(resp.Torrents[i].Status) < reorderTorrentStatus(resp.Torrents[j].Status)
	})
	page := downloadsPageContext{
		PageContext: *ui.New(),
		Torrents:    resp.Torrents,
	}
	ctx.HTML(http.StatusOK, "multimedia.downloads.tmpl", &page)
}

func (s *Service) downloadsDeleteHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	_, err := s.f.NewTorrent().RemoveTorrent(ctx, &rms_torrent.RemoveTorrentRequest{Id: id})
	if err != nil {
		logger.Errorf("Remove torrent failed: %s", err)
		ui.DisplayError(ctx, http.StatusInternalServerError, "Не удалось удалить загрузку")
		return
	}
	ctx.Redirect(http.StatusFound, "/multimedia/downloads")
}

func (s *Service) downloadsUpHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	_, err := s.f.NewTorrent().UpPriority(ctx, &rms_torrent.UpPriorityRequest{Id: id})
	if err != nil {
		logger.Errorf("Up torrent failed: %s", err)
		ui.DisplayError(ctx, http.StatusInternalServerError, "Ошибка при обращении к сервису загрузок")
		return
	}
	ctx.Redirect(http.StatusFound, "/multimedia/downloads")
}
