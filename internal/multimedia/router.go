package multimedia

import (
	"github.com/RacoonMediaServer/rms-web/internal/config"
	"github.com/RacoonMediaServer/rms-web/internal/ui"
	"github.com/gin-gonic/gin"
	"path"
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
				Image:       "/img/updates.png",
				Title:       "Обновления",
				Link:        "/multimedia/updates",
				Description: "Посмотреть наличие новых сезонов для уже имеющихся в Библиотеке сериалов",
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

	router.GET("/library", s.libraryHandler)
	router.GET("/library/movie/:id", s.playHandler)
	router.GET("/library/delete/:id", s.deleteMovieHandler)
	router.Static("/library/content", path.Join(config.Config().Content.Directory, "movies"))

	router.GET("/updates", s.updatesHandler)

	router.GET("/search", s.searchHandler)
	router.GET("/download/:id", s.downloadHandler)

	router.GET("/downloads", s.downloadsHandler)
	router.GET("/downloads/up/:id", s.downloadsUpHandler)
	router.GET("/downloads/delete/:id", s.downloadsDeleteHandler)

	router.GET("/upload", s.getUploadFileHandler)
	router.POST("/upload", s.postUploadFileHandler)
}
