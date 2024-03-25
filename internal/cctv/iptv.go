package cctv

import (
	"fmt"
	"github.com/RacoonMediaServer/rms-packages/pkg/media"
	"net/http"

	rms_cctv "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-cctv"
	"github.com/RacoonMediaServer/rms-web/internal/ui"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes/empty"
	"go-micro.dev/v4/logger"
)

func (s *Service) playlistCameraHandler(ctx *gin.Context) {
	cctv := s.f.NewCctv()
	resp, err := cctv.GetCameras(ctx, &empty.Empty{})
	if err != nil {
		logger.Errorf("Get cameras list failed: %s", err)
		ui.DisplayError(ctx, http.StatusInternalServerError, "Не удалось получить список камер")
		return
	}

	useMainSource := ctx.Query("is-main") == "1"
	usePseudo := ctx.Query("pseudo") == "1"
	container := "mpegts"

	if ctx.Query("container") != "" {
		container = ctx.Query("container")
	}

	body := "#EXTM3U\n"

	for _, camera := range resp.Cameras {
		req := rms_cctv.GetLiveUriRequest{
			CameraId:    camera.Id,
			MainProfile: useMainSource,
		}
		if usePseudo {
			req.Transport = media.Transport_HTTP_HLS_ONE_CHUNK
		} else if container == "mp4" {
			req.Transport = media.Transport_HTTP_HLS_MP4
		} else {
			req.Transport = media.Transport_HTTP_HLS_MPEGTS
		}

		uri, err := cctv.GetLiveUri(ctx, &req)
		if err != nil {
			logger.Errorf("Get live URL failed: %s", err)
			continue
		}

		body += fmt.Sprintf("#EXTINF:-1 ,%s\n", camera.Name)

		body += uri.Uri + "\n\n"
	}

	ctx.Header("Content-Type", "audio/mpegurl")
	ctx.Status(http.StatusOK)
	_, _ = ctx.Writer.Write([]byte(body))
}
