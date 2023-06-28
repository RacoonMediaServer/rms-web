package settings

import (
	"github.com/RacoonMediaServer/rms-web/internal/config"
	"github.com/RacoonMediaServer/rms-web/internal/ui"
	"github.com/gin-gonic/gin"
)

func (s *Service) catalogHandler(ctx *gin.Context) {
	page := ui.CatalogPageContext{
		PageContext: *ui.New(),
		Title:       "Настройки",
		Parts: []ui.CatalogPart{
			{
				Image:       "/img/telegram.png",
				Title:       "Telegram-бот",
				Link:        "/settings/telegram",
				Description: "Управление привязкой устройства к Telegram-боту",
			},
			{
				Image:       "/img/torrent.webp",
				Title:       "Торрент-клиент",
				Link:        "/settings/torrent",
				Description: "Настройки загрузки контекст с торрент-трекеров",
			},
			{
				Image:       "/img/notification.png",
				Title:       "Уведомления",
				Link:        "/settings/notifications",
				Description: "Управления отправкой уведомлений о событиях",
			},
			{
				Image:       "/img/video.png",
				Title:       "Транскодирование",
				Link:        "/settings/transcoding",
				Description: "Настройки транскодирования видео",
			},
		},
	}
	if config.Config().Cctv.Enabled {
		page.Parts = append(page.Parts, ui.CatalogPart{
			Image:       "/img/alert.png",
			Title:       "Тревога",
			Link:        "/settings/alert",
			Description: "Настройки тревоги системы видеонаблюдения",
		})
	}
	page.Display(ctx)
}

func (s *Service) Register(router *gin.RouterGroup) {
	router.GET("/", s.catalogHandler)

	router.GET("/telegram", s.getTelegramSettings)

	router.GET("/notifications", s.notificationsHandler)
	router.POST("/notifications", s.saveNotificationSettingsHandler)
	router.POST("/notifications/new", s.addNotificationRuleHandler)
	router.GET("/notifications/delete/:topic/:index", s.deleteNotificationHandler)
	router.GET("/notifications/test/:topic", s.testNotificationHandler)
}
