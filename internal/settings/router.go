package settings

import "github.com/gin-gonic/gin"

func (s *Service) Register(router *gin.RouterGroup) {
	router.GET("/telegram", s.getTelegramSettings)
}
