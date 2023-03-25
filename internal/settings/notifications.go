package settings

import (
	rms_notifier "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-notifier"
	"github.com/RacoonMediaServer/rms-web/internal/ui"
	"github.com/gin-gonic/gin"
	"go-micro.dev/v4/logger"
	"google.golang.org/protobuf/types/known/emptypb"
	"net/http"
	"strconv"
)

type notifyPageContext struct {
	ui.PageContext
	Settings *rms_notifier.Settings
}

func (s *Service) getNotificationsSettings(ctx *gin.Context) (*rms_notifier.Settings, bool) {
	settings, err := s.f.NewNotifier().GetSettings(ctx, &emptypb.Empty{})
	if err != nil {
		logger.Errorf("Get notifier settings failed: %s", err)
		ui.DisplayError(ctx, http.StatusInternalServerError, "Не удалось связаться с сервисом уведомлений")
		return nil, false
	}
	return settings, true
}

func (s *Service) setNotificationsSettings(ctx *gin.Context, settings *rms_notifier.Settings) bool {
	_, err := s.f.NewNotifier().SetSettings(ctx, settings)
	if err != nil {
		logger.Errorf("Save notifier settings failed: %s", err)
		ui.DisplayError(ctx, http.StatusInternalServerError, "Не удалось сохранить настройки")
		return false
	}
	return true
}

func (s *Service) notificationsHandler(ctx *gin.Context) {
	settings, ok := s.getNotificationsSettings(ctx)
	if !ok {
		return
	}
	page := notifyPageContext{
		PageContext: *ui.New(),
		Settings:    settings,
	}
	ctx.HTML(http.StatusOK, "settings.notifications.tmpl", &page)
}

func (s *Service) saveNotificationSettingsHandler(ctx *gin.Context) {
	settings, ok := s.getNotificationsSettings(ctx)
	if !ok {
		return
	}

	parseUint := func(s string, defaultVal uint32) uint32 {
		val, err := strconv.ParseUint(s, 10, 32)
		if err != nil {
			return defaultVal
		}
		return uint32(val)
	}

	settings.Enabled = ctx.PostForm("enabled") == "on"
	settings.FilterInterval = parseUint(ctx.PostForm("filterInterval"), settings.FilterInterval)
	settings.RotationInterval = parseUint(ctx.PostForm("rotationInterval"), settings.RotationInterval)

	if s.setNotificationsSettings(ctx, settings) {
		ctx.Redirect(http.StatusFound, "/settings/notifications")
	}
}

func (s *Service) addNotificationRuleHandler(ctx *gin.Context) {
	settings, ok := s.getNotificationsSettings(ctx)
	if !ok {
		return
	}
	method, err := strconv.ParseInt(ctx.PostForm("method"), 10, 8)
	if err != nil || method < 0 || method > int64(rms_notifier.Rule_SMS) {
		ui.DisplayError(ctx, http.StatusBadRequest, "Неверный способ доставки уведомлений")
		return
	}

	rule := rms_notifier.Rule{
		Method:      rms_notifier.Rule_Method(method),
		Destination: ctx.PostForm("address"),
	}

	topic := ctx.PostForm("topic")
	rules, ok := settings.Rules[topic]
	if !ok || rules == nil {
		settings.Rules[topic] = &rms_notifier.Settings_Rules{
			Rule: []*rms_notifier.Rule{&rule},
		}
	} else {
		rules.Rule = append(rules.Rule, &rule)
	}

	if s.setNotificationsSettings(ctx, settings) {
		ctx.Redirect(http.StatusFound, "/settings/notifications")
	}
}

func (s *Service) deleteNotificationHandler(ctx *gin.Context) {
	settings, ok := s.getNotificationsSettings(ctx)
	if !ok {
		return
	}

	topic := ctx.Param("topic")
	index, err := strconv.ParseUint(ctx.Param("index"), 10, 8)
	if err != nil {
		ui.DisplayError(ctx, http.StatusBadRequest, "Неверно указан индекс правила уведомлений")
		return
	}

	r, ok := settings.Rules[topic]
	if !ok || index >= uint64(len(r.Rule)) {
		ui.DisplayError(ctx, http.StatusNotFound, "Удаляемое правило не найдено")
		return
	}

	r.Rule = append(r.Rule[:index], r.Rule[index+1:]...)
	if s.setNotificationsSettings(ctx, settings) {
		ctx.Redirect(http.StatusFound, "/settings/notifications")
	}
}
