package multimedia

import "github.com/gin-gonic/gin"

func (s *Service) Register(router *gin.RouterGroup) {
	router.GET("/search", s.searchHandler)
	router.GET("/download/:id", s.downloadHandler)

	router.GET("/downloads", s.downloadsHandler)
	router.GET("/downloads/up/:id", s.downloadsUpHandler)
	router.GET("/downloads/delete/:id", s.downloadsDeleteHandler)
}
