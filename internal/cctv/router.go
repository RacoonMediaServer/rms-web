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
				Link:        "/cameras/view",
				Description: "Просмотр видео с IP-камер",
			},
			{
				Image:       "/img/cctv.png",
				Title:       "Настройка",
				Link:        "/cameras/setup",
				Description: "Добавление и удаления IP-камер в систему видеонаблюдения",
			},
		},
	}
	page.Display(ctx)
}

func (s *Service) Register(router *gin.RouterGroup) {
	router.GET("/", s.catalogHandler)
}
