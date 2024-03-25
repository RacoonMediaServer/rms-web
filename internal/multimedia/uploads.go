package multimedia

import (
	rms_library "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-library"
	"github.com/RacoonMediaServer/rms-web/internal/ui"
	"github.com/gin-gonic/gin"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/logger"
	"io"
	"net/http"
	"time"
)

const uploadTimeout = 30 * time.Second

type uploadPageContext struct {
	ui.PageContext
	ID string
}

func (s *Service) getUploadFileHandler(ctx *gin.Context) {
	page := uploadPageContext{
		PageContext: *ui.New(),
		ID:          ctx.Query("id"),
	}
	ctx.HTML(http.StatusOK, "multimedia.upload.tmpl", &page)
}

func (s *Service) postUploadFileHandler(ctx *gin.Context) {
	id := ctx.PostForm("id")
	mov, ok := s.movieFromCache(id)
	if !ok {
		logger.Errorf("Movie not found in cache: %s", id)
		ui.DisplayError(ctx, http.StatusBadRequest, "Неверно указан фильм")
		return
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		logger.Errorf("Upload torrent file failed: %s", err)
		ui.DisplayError(ctx, http.StatusBadRequest, "Не удается загрузить файл")
		return
	}
	f, err := file.Open()
	if err != nil {
		logger.Errorf("Upload torrent file failed: %s", err)
		ui.DisplayError(ctx, http.StatusBadRequest, "Не удается загрузить файл")
		return
	}

	defer f.Close()
	buf, err := io.ReadAll(f)
	if err != nil {
		logger.Errorf("Upload torrent file failed: %s", err)
		ui.DisplayError(ctx, http.StatusBadRequest, "Не удается загрузить файл")
		return
	}

	req := rms_library.UploadMovieRequest{
		Id:          id,
		Info:        mov.Info,
		TorrentFile: buf,
	}
	_, err = s.f.NewLibrary().UploadMovie(ctx, &req, client.WithRequestTimeout(uploadTimeout))
	if err != nil {
		logger.Errorf("Upload file to library failed: %s", err)
		ui.DisplayError(ctx, http.StatusInternalServerError, "Не удалось связаться с сервисом библиотеки медиа")
		return
	}
	ui.DisplayOK(ctx, "Файл загружен", "/multimedia/downloads")
}
