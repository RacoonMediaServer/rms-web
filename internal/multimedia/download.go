package multimedia

import (
	"net/http"
	"strconv"
	"time"

	rms_library "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-library"
	"github.com/RacoonMediaServer/rms-web/internal/ui"
	"github.com/gin-gonic/gin"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/logger"
)

const downloadTimeout = 40 * time.Second
const torrentsLimit = 15

type downloadPageContext struct {
	ui.PageContext
	Id       string
	Fast     bool
	Select   bool
	Seasons  []uint
	Torrents []*rms_library.Torrent
}

func parseSeason(season string) uint32 {
	if season == "all" {
		return 0
	}
	s, _ := strconv.ParseUint(season, 10, 8)
	return uint32(s)
}

func (s *Service) downloadHandler(ctx *gin.Context) {
	page := &downloadPageContext{
		PageContext: *ui.New(),
		Id:          ctx.Param("id"),
		Fast:        ctx.Query("fast") == "true",
		Select:      ctx.Query("select") == "true",
	}

	season := ctx.Query("season")
	if season == "" {
		data, ok := s.cache.Load(page.Id)
		if ok {
			mov, ok := data.(*rms_library.FoundMovie)
			if ok && mov.Info.Seasons != nil {
				for i := uint(1); i <= uint(*mov.Info.Seasons); i++ {
					page.Seasons = append(page.Seasons, i)
				}
				ctx.HTML(http.StatusOK, "multimedia.download.tmpl", page)
				return
			}
		}
	}

	seasonNo := parseSeason(season)

	if page.Select {
		torrent := ctx.Query("torrent")
		if torrent != "" {
			s.downloadMovieTorrent(ctx, torrent)
		} else {
			s.selectTorrent(ctx, page, seasonNo)
		}
	} else {
		s.downloadMovieAuto(ctx, page.Id, page.Fast, seasonNo)
	}
}

func (s *Service) downloadMovieAuto(ctx *gin.Context, id string, fast bool, seasonNo uint32) {
	req := rms_library.DownloadMovieAutoRequest{
		Id:     id,
		Faster: fast,
	}
	if seasonNo != 0 {
		req.Season = &seasonNo
	}
	resp, err := s.f.NewLibrary().DownloadMovieAuto(ctx, &req, client.WithRequestTimeout(downloadTimeout))
	if err != nil {
		logger.Errorf("Download auto failed: %s", err)
		ui.DisplayError(ctx, http.StatusInternalServerError, "Ошибка загрузки торрента")
		return
	}
	if !resp.Found {
		ui.DisplayError(ctx, http.StatusBadRequest, "Не удалось найти подходящую раздачу")
		return
	}
	ctx.Redirect(http.StatusFound, "/multimedia/downloads")
}

func (s *Service) selectTorrent(ctx *gin.Context, page *downloadPageContext, seasonNo uint32) {
	req := rms_library.FindMovieTorrentsRequest{
		Id:    page.Id,
		Limit: torrentsLimit,
	}
	if seasonNo != 0 {
		req.Season = &seasonNo
	}
	resp, err := s.f.NewLibrary().FindMovieTorrents(ctx, &req, client.WithRequestTimeout(downloadTimeout))
	if err != nil {
		logger.Errorf("Find torrents failed: %s", err)
		ui.DisplayError(ctx, http.StatusInternalServerError, "Ошибка при поиске раздач на торрент-трекерах")
		return
	}
	page.Torrents = resp.Results
	ctx.HTML(http.StatusOK, "multimedia.download.select.tmpl", page)
}

func (s *Service) downloadMovieTorrent(ctx *gin.Context, torrent string) {
	_, err := s.f.NewLibrary().DownloadTorrent(ctx, &rms_library.DownloadTorrentRequest{TorrentId: torrent}, client.WithRequestTimeout(downloadTimeout))
	if err != nil {
		logger.Errorf("Download torrent failed: %s", err)
		ui.DisplayError(ctx, http.StatusInternalServerError, "Не удалось скачать торрент")
		return
	}
	ctx.Redirect(http.StatusFound, "/multimedia/downloads")
}
