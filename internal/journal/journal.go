package journal

import (
	"fmt"
	"github.com/RacoonMediaServer/rms-packages/pkg/events"
	rms_notifier "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-notifier"
	"github.com/RacoonMediaServer/rms-web/internal/ui"
	"github.com/gin-gonic/gin"
	"go-micro.dev/v4/logger"
	"net/http"
	"time"
)

const eventsByPage = 20
const journalTimeout = 10 * time.Second

type journalPageContext struct {
	ui.PageContext
	Events []*event
}

type event struct {
	Image   string
	Alt     string
	Time    string
	Title   string
	Sender  string
	Details map[string]string
}

func convertEvent(e *rms_notifier.Event) *event {
	ev := event{}
	if e.Notification != nil {
		ev.Image = "notification"
		ev.Alt = "Уведомление"
		ev.Details = map[string]string{"Отправитель": e.Notification.Sender}
		switch e.Notification.Kind {
		case events.Notification_DownloadComplete:
			ev.Title = "Загрузка завершена"
		case events.Notification_DownloadFailed:
			ev.Title = "Ошибка при загрузке контента"
		case events.Notification_TranscodingDone:
			ev.Title = "Транскодирование завершено"
		case events.Notification_TranscodingFailed:
			ev.Title = "Ошибка транскодирования"
		case events.Notification_TorrentRemoved:
			ev.Title = "Торрент удален"
		}
		if e.Notification.ItemTitle != nil {
			ev.Details["Название"] = *e.Notification.ItemTitle
		}
		if e.Notification.TorrentID != nil {
			ev.Details["Торрент"] = *e.Notification.TorrentID
		}
		if e.Notification.MediaID != nil {
			ev.Details["ID"] = *e.Notification.MediaID
		}
	} else if e.Malfunction != nil {
		ev.Details = map[string]string{"Отправитель": e.Malfunction.Sender}
		ev.Image = "malfunction"
		ev.Alt = "Сбой"
		ev.Title = fmt.Sprintf("Сбой в сервисе %s", ev.Sender)
		ev.Details["Ошибка"] = e.Malfunction.Error
		ev.Details["Код"] = e.Malfunction.Code.String()
		ev.Details["Подсистема"] = e.Malfunction.System.String()
	} else if e.Alert != nil {
		ev.Details = map[string]string{"Отправитель": e.Alert.Sender}
		ev.Image = "alert"
		ev.Alt = "Тревога"
		ev.Title = "Нарушение периметра"
		ev.Details["Камера"] = e.Alert.Camera
		ev.Details["Код"] = e.Alert.Kind.String()
	}
	ev.Time = time.Unix(e.Timestamp, 0).Local().String()

	return &ev
}

func (s *Service) journalHandler(ctx *gin.Context) {
	resp, err := s.f.NewNotifier().GetJournal(ctx, &rms_notifier.GetJournalRequest{Limit: eventsByPage})
	if err != nil {
		logger.Errorf("Get journal failed: %s", err)
	}
	page := journalPageContext{
		PageContext: *ui.New(),
	}
	if resp != nil {
		for _, e := range resp.Events {
			page.Events = append(page.Events, convertEvent(e))
		}
	}
	ctx.HTML(http.StatusOK, "journal.tmpl", &page)
}
