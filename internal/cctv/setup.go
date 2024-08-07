package cctv

import (
	"net/http"
	"strconv"

	"github.com/RacoonMediaServer/rms-packages/pkg/schedule"
	rms_cctv "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-cctv"
	"github.com/RacoonMediaServer/rms-web/internal/ui"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes/empty"
	"go-micro.dev/v4/logger"
)

type camerasPageContext struct {
	ui.PageContext
	Cameras []*rms_cctv.Camera
}

type scheduleItem struct {
	ID   string
	Name string
}

type periodItem struct {
	Enabled  bool
	Interval schedule.Interval
}

type cameraPageContext struct {
	ui.PageContext
	Camera    *rms_cctv.Camera
	Schedules []scheduleItem
}

type cameraForm struct {
	Name     string                 `form:"name"`
	Url      string                 `form:"url"`
	User     string                 `form:"user"`
	Password string                 `form:"password"`
	Mode     rms_cctv.RecordingMode `form:"mode"`
	KeepDays uint                   `form:"keep_days"`
	Schedule string                 `form:"schedule"`
}

func (s *Service) getCamerasHandler(ctx *gin.Context) {
	resp, err := s.f.NewCctvCameras().GetCameras(ctx, &empty.Empty{})
	if err != nil {
		logger.Errorf("Get cameras list failed: %s", err)
		ui.DisplayError(ctx, http.StatusInternalServerError, "Не удалось получить список камер")
		return
	}

	page := camerasPageContext{
		PageContext: *ui.New(),
		Cameras:     resp.Cameras,
	}
	ctx.HTML(http.StatusOK, "cctv.cameras.setup.tmpl", &page)
}

func (s *Service) getNewCameraHandler(ctx *gin.Context) {
	page := struct {
		ui.PageContext
		Schedules []scheduleItem
	}{PageContext: *ui.New()}

	resp, err := s.f.NewCctvSchedules().GetSchedulesList(ctx, &empty.Empty{})
	if err != nil {
		logger.Errorf("Get schedules failed: %s", err)
		ui.DisplayError(ctx, http.StatusInternalServerError, "Не удалось связаться с системой видеонаблюдения")
		return
	}

	page.Schedules = make([]scheduleItem, len(resp.Result))
	for i, sched := range resp.Result {
		page.Schedules[i] = scheduleItem{Name: sched.Name, ID: sched.Id}
	}

	ctx.HTML(http.StatusOK, "cctv.cameras.setup.new.tmpl", page)
}

func (s *Service) postNewCameraHandler(ctx *gin.Context) {
	form := cameraForm{}
	if err := ctx.ShouldBind(&form); err != nil {
		ui.DisplayError(ctx, http.StatusBadRequest, "Ошибка в полях формы")
		return
	}

	cam := rms_cctv.Camera{
		Name:     form.Name,
		Url:      form.Url,
		User:     form.User,
		Password: form.Password,
		Mode:     form.Mode,
		KeepDays: uint32(form.KeepDays),
		Schedule: form.Schedule,
	}

	if _, err := s.f.NewCctvCameras().AddCamera(ctx, &cam); err != nil {
		logger.Errorf("Add camera failed: %s", err)
		ui.DisplayError(ctx, http.StatusInternalServerError, "Не удалось добавить камеру")
		return
	}
	ui.DisplayOK(ctx, "Камера добавлена", "/cctv/cameras/setup")
}

func (s *Service) getCameraHandler(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("camera"), 10, 32)
	if err != nil {
		ui.DisplayError(ctx, http.StatusNotFound, "Камера не найдена")
		return
	}

	cctv := s.f.NewCctvCameras()
	resp, err := cctv.GetCameras(ctx, &empty.Empty{})
	if err != nil {
		logger.Errorf("Get cameras list failed: %s", err)
		ui.DisplayError(ctx, http.StatusInternalServerError, "Не удалось получить список камер")
		return
	}

	schedResp, err := s.f.NewCctvSchedules().GetSchedulesList(ctx, &empty.Empty{})
	if err != nil {
		logger.Errorf("Get schedules failed: %s", err)
		ui.DisplayError(ctx, http.StatusInternalServerError, "Не удалось связаться с системой видеонаблюдения")
		return
	}

	var camera *rms_cctv.Camera
	for _, cam := range resp.Cameras {
		if cam.Id == uint32(id) {
			camera = cam
			break
		}
	}
	if camera == nil {
		ui.DisplayError(ctx, http.StatusNotFound, "Камера не найдена")
		return
	}

	page := cameraPageContext{
		PageContext: *ui.New(),
		Camera:      camera,
		Schedules:   make([]scheduleItem, len(schedResp.Result)),
	}
	for i, sched := range schedResp.Result {
		page.Schedules[i] = scheduleItem{Name: sched.Name, ID: sched.Id}
	}
	ctx.HTML(http.StatusOK, "cctv.cameras.setup.edit.tmpl", &page)
}

func (s *Service) postCameraHandler(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("camera"), 10, 32)
	if err != nil {
		ui.DisplayError(ctx, http.StatusNotFound, "Камера не найдена")
		return
	}

	form := cameraForm{}
	if err := ctx.ShouldBind(&form); err != nil {
		ui.DisplayError(ctx, http.StatusBadRequest, "Ошибка в полях формы")
		return
	}

	cam := rms_cctv.ModifyCameraRequest{
		Id:       uint32(id),
		Name:     form.Name,
		Mode:     form.Mode,
		KeepDays: uint32(form.KeepDays),
		Schedule: form.Schedule,
	}

	if _, err := s.f.NewCctvCameras().ModifyCamera(ctx, &cam); err != nil {
		logger.Errorf("Modify camera failed: %s", err)
		ui.DisplayError(ctx, http.StatusInternalServerError, "Не удалось изменить настройки камеры")
		return
	}
	ui.DisplayOK(ctx, "Камера изменена", "/cctv/cameras/setup")
}

func (s *Service) deleteCameraHandler(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("camera"), 10, 32)
	if err != nil {
		ui.DisplayError(ctx, http.StatusNotFound, "Камера не найдена")
		return
	}
	_, err = s.f.NewCctvCameras().DeleteCamera(ctx, &rms_cctv.DeleteCameraRequest{CameraId: uint32(id)})
	if err != nil {
		ui.DisplayError(ctx, http.StatusInternalServerError, "Не удалось удалить камеру")
		return
	}
	ui.DisplayOK(ctx, "Камера удалена", "/cctv/cameras/setup")
}
