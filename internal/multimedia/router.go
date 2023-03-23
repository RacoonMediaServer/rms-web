package multimedia

import "github.com/gin-gonic/gin"

func (s *Service) Register(router *gin.RouterGroup) {
	router.GET("/search", s.getSearchHandler)
	router.GET("/download/:id", s.getDownloadHandler)
}
