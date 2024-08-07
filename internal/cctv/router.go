package cctv

import (
	"github.com/RacoonMediaServer/rms-web/internal/ui"
	"github.com/gin-gonic/gin"
)

func (s *Service) catalogHandler(ctx *gin.Context) {
	page := ui.CatalogPageContext{
		PageContext: *ui.New(),
		Title:       "Мультимедиа",
		Parts: []ui.CatalogPart{
			{
				Image:       "/img/play.png",
				Title:       "Просмотр",
				Link:        "/cctv/cameras/view",
				Description: "Просмотр видео с IP-камер",
			},
			{
				Image:       "/img/cctv.png",
				Title:       "Настройка",
				Link:        "/cctv/cameras/setup",
				Description: "Добавление и удаления IP-камер в систему видеонаблюдения",
			},
			{
				Image:       "/img/schedule.png",
				Title:       "Расписания",
				Link:        "/cctv/schedules",
				Description: "Управление расписаниями получения уведомлений с камер",
			},
		},
	}
	page.Display(ctx)
}

func (s *Service) Register(router *gin.RouterGroup) {
	router.GET("/", s.catalogHandler)

	router.GET("/cameras/setup", s.getCamerasHandler)

	router.GET("/cameras/setup/new", s.getNewCameraHandler)
	router.POST("/cameras/setup/new", s.postNewCameraHandler)

	router.GET("/cameras/setup/edit/:camera", s.getCameraHandler)
	router.POST("/cameras/setup/edit/:camera", s.postCameraHandler)
	router.GET("/cameras/setup/delete/:camera", s.deleteCameraHandler)

	router.GET("/cameras/view", s.viewCamerasHandler)

	router.GET("/iptv.m3u8", s.playlistCameraHandler)

	router.GET("/schedules", s.getSchedulesHandler)
	router.GET("/schedules/new", s.getNewScheduleHandler)
	router.POST("/schedules/new", s.postNewScheduleHandler)
	router.GET("/schedules/edit/:id", s.getScheduleHandler)
	router.POST("/schedules/edit/:id", s.postScheduleHandler)
	router.GET("/schedules/delete/:id", s.deleteScheduleHandler)
}
