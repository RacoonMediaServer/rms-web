package cctv

import (
	"github.com/RacoonMediaServer/rms-packages/pkg/media"
	"net/http"

	rms_cctv "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-cctv"
	"github.com/RacoonMediaServer/rms-web/internal/ui"
	"github.com/gin-gonic/gin"
	"go-micro.dev/v4/logger"
	"google.golang.org/protobuf/types/known/emptypb"
)

type cameraView struct {
	Name string
	URL  string
}

func (s *Service) viewCamerasHandler(ctx *gin.Context) {
	cctv := s.f.NewCctvCameras()
	resp, err := cctv.GetCameras(ctx, &emptypb.Empty{})
	if err != nil {
		logger.Errorf("Get cameras list failed: %s", err)
		ui.DisplayError(ctx, http.StatusInternalServerError, "Не удалось получить список камер")
		return
	}

	pageCtx := struct {
		ui.PageContext
		Cameras  []cameraView
		BeginURL string
	}{}

	req := rms_cctv.GetLiveUriRequest{
		Transport:   media.Transport_HTTP_MP4,
		MainProfile: true,
	}
	for _, camera := range resp.Cameras {
		req.CameraId = camera.Id
		resp, err := cctv.GetLiveUri(ctx, &req)
		if err != nil {
			logger.Errorf("Get live url failed: %s", err)
			continue
		}
		pageCtx.Cameras = append(pageCtx.Cameras, cameraView{Name: camera.Name, URL: resp.Uri})
	}

	if len(resp.Cameras) > 0 {
		pageCtx.BeginURL = pageCtx.Cameras[0].URL
	}

	pageCtx.PageContext = *ui.New()

	ctx.HTML(http.StatusOK, "cctv.cameras.view.tmpl", &pageCtx)
}
