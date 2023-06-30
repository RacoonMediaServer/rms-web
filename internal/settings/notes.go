package settings

import (
	rms_notes "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-notes"
	"github.com/RacoonMediaServer/rms-web/internal/ui"
	"github.com/gin-gonic/gin"
	"go-micro.dev/v4/logger"
	"google.golang.org/protobuf/types/known/emptypb"
	"net/http"
	"strconv"
)

type notesPageContext struct {
	ui.PageContext
	Settings *rms_notes.NotesSettings
}

func (s *Service) notesSettingsHandler(ctx *gin.Context) {
	settings, err := s.f.NewNotes().GetSettings(ctx, &emptypb.Empty{})
	if err != nil {
		logger.Errorf("Get notes settings failed: %s", err)
		ui.DisplayError(ctx, http.StatusInternalServerError, "Не удалось связаться с сервисом управления заметками")
		return
	}

	page := notesPageContext{
		PageContext: *ui.New(),
		Settings:    settings,
	}
	ctx.HTML(http.StatusOK, "settings.notes.tmpl", &page)
}

func (s *Service) saveNotesSettingsHandler(ctx *gin.Context) {
	notificationTime, err := strconv.ParseUint(ctx.PostForm("notificationTime"), 10, 8)
	if err != nil || notificationTime > 23 {
		ui.DisplayError(ctx, http.StatusBadRequest, "Неверно указано время оповещения")
		return
	}

	settings := rms_notes.NotesSettings{
		Directory:        ctx.PostForm("directory"),
		NotesDirectory:   ctx.PostForm("notesDirectory"),
		TasksFile:        ctx.PostForm("tasksFile"),
		NotificationTime: uint32(notificationTime),
	}

	_, err = s.f.NewNotes().SetSettings(ctx, &settings)
	if err != nil {
		logger.Errorf("Apply notes settings failed: %s", err)
		ui.DisplayError(ctx, http.StatusInternalServerError, "Не удалось применить настройки")
		return
	}

	ui.DisplayOK(ctx, "Сохранено", "/settings/notes")
}
