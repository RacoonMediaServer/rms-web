package multimedia

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
				Image:       "/img/library.gif",
				Title:       "Библиотека",
				Link:        "/multimedia/library",
				Description: "Просмотр и управление скачанным мультимедия-контентом",
			},
			{
				Image:       "/img/multimedia.png",
				Title:       "Поиск",
				Link:        "/multimedia/search",
				Description: "Поиск фильмов, сериалов, музыки на внешних сайтах по названию",
			},
			{
				Image:       "/img/downloads.png",
				Title:       "Загрузки",
				Link:        "/multimedia/downloads",
				Description: "Управления загрузками мультимедия-контента",
			},
		},
	}
	page.Display(ctx)
}

func (s *Service) Register(router *gin.RouterGroup) {
	router.GET("/", s.catalogHandler)

	router.GET("/search", s.searchHandler)
	router.GET("/download/:id", s.downloadHandler)

	router.GET("/downloads", s.downloadsHandler)
	router.GET("/downloads/up/:id", s.downloadsUpHandler)
	router.GET("/downloads/delete/:id", s.downloadsDeleteHandler)
}
