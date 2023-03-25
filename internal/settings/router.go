package settings

import "github.com/gin-gonic/gin"

func (s *Service) Register(router *gin.RouterGroup) {
	router.GET("/telegram", s.getTelegramSettings)

	router.GET("/notifications", s.notificationsHandler)
	router.POST("/notifications", s.saveNotificationSettingsHandler)
	router.POST("/notifications/new", s.addNotificationRuleHandler)
	router.GET("/notifications/delete/:topic/:index", s.deleteNotificationHandler)
}
