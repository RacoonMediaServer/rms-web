package settings

import (
	rms_backup "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-backup"
	"github.com/RacoonMediaServer/rms-web/internal/ui"
	"github.com/gin-gonic/gin"
	"go-micro.dev/v4/logger"
	"google.golang.org/protobuf/types/known/emptypb"
	"net/http"
	"strconv"
)

type backupPageContext struct {
	ui.PageContext
	Settings *rms_backup.BackupSettings
	Backups  []*rms_backup.BackupInfo
	Status   rms_backup.GetBackupStatusResponse_Status
	Progress float32
}

func (s *Service) backupSettingsHandler(ctx *gin.Context) {
	backupService := s.f.NewBackup()
	settings, err := backupService.GetBackupSettings(ctx, &emptypb.Empty{})
	if err != nil {
		logger.Errorf("Get backup settings failed: %s", err)
		ui.DisplayError(ctx, http.StatusInternalServerError, "Не удалось связаться с сервисом резервного копирования")
		return
	}

	backups, err := backupService.GetBackups(ctx, &emptypb.Empty{})
	if err != nil {
		logger.Errorf("Get backups failed: %s", err)
		ui.DisplayError(ctx, http.StatusInternalServerError, "Не удалось связаться с сервисом резервного копирования")
		return
	}

	status, err := backupService.GetBackupStatus(ctx, &emptypb.Empty{})
	if err != nil {
		logger.Errorf("Get status failed: %s", err)
		ui.DisplayError(ctx, http.StatusInternalServerError, "Не удалось связаться с сервисом резервного копирования")
		return
	}

	if settings.Password == nil {
		settings.Password = new(string)
	}
	page := backupPageContext{
		PageContext: *ui.New(),
		Settings:    settings,
		Backups:     backups.Backups,
		Status:      status.Status,
		Progress:    status.Progress,
	}
	ctx.HTML(http.StatusOK, "settings.backup.tmpl", &page)
}

func parseFormInt8(ctx *gin.Context, field string) (uint64, bool) {
	parsed, err := strconv.ParseUint(ctx.PostForm(field), 10, 8)
	if err != nil {
		logger.Errorf("Parse value '%s' failed: %s", field, err)
		ui.DisplayError(ctx, http.StatusBadRequest, "Введены неверные параметры")
		return 0, false
	}
	return parsed, true
}

func (s *Service) saveBackupSettingsHandler(ctx *gin.Context) {
	period, ok := parseFormInt8(ctx, "period")
	if !ok {
		return
	}

	day, ok := parseFormInt8(ctx, "day")
	if !ok {
		return
	}

	hour, ok := parseFormInt8(ctx, "hour")
	if !ok {
		return
	}

	var password *string
	userPassword := ctx.PostForm("password")
	if userPassword != "" {
		password = &userPassword
	}

	settings := rms_backup.BackupSettings{
		Enabled:  ctx.PostForm("enabled") == "on",
		Type:     rms_backup.BackupType_Full,
		Period:   rms_backup.BackupSettings_Period(period),
		Day:      uint32(day),
		Hour:     uint32(hour),
		Password: password,
	}

	_, err := s.f.NewBackup().SetBackupSettings(ctx, &settings)
	if err != nil {
		logger.Errorf("Set backup settings failed: %s", err)
		ui.DisplayError(ctx, http.StatusInternalServerError, "Не удалось применить настройки")
		return
	}

	ui.DisplayOK(ctx, "Сохранено", "/settings/backup")
}

func (s *Service) launchBackupHandler(ctx *gin.Context) {
	backupType, ok := parseFormInt8(ctx, "type")
	if !ok {
		return
	}
	resp, err := s.f.NewBackup().LaunchBackup(ctx, &rms_backup.LaunchBackupRequest{Type: rms_backup.BackupType(backupType)})
	if err != nil {
		logger.Errorf("Launch backup failed: %s", err)
		ui.DisplayError(ctx, http.StatusInternalServerError, "Не удалось связаться с сервисом резервного копирования")
		return
	}
	if resp.AlreadyLaunch {
		ui.DisplayError(ctx, http.StatusBadRequest, "Резервное копирование уже было запущено")
		return
	}
	ui.DisplayOK(ctx, "Запущено", "/settings/backup")
}

func (s *Service) deleteBackupHandler(ctx *gin.Context) {
	_, err := s.f.NewBackup().RemoveBackup(ctx, &rms_backup.RemoveBackupRequest{FileName: ctx.Param("backup")})
	if err != nil {
		logger.Errorf("Remove backup failed: %s", err)
		ui.DisplayError(ctx, http.StatusInternalServerError, "Не удалось удалить резервную копию")
		return
	}
	ui.DisplayOK(ctx, "Удалено", "/settings/backup")
}
