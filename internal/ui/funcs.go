package ui

import (
	"fmt"
	"html/template"
	"path"
	"time"

	"github.com/RacoonMediaServer/rms-packages/pkg/pubsub"
	rms_backup "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-backup"
	rms_torrent "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-torrent"
	"github.com/dustin/go-humanize"
)

var Functions template.FuncMap = template.FuncMap{
	"prettyStatus":       prettyStatus,
	"prettyFloat":        prettyFloat,
	"prettyTopic":        prettyTopic,
	"fileName":           fileName,
	"prettySize":         prettySize,
	"prettyBitrate":      prettyBitrate,
	"prettyBackupType":   prettyBackupType,
	"prettyUnixTime":     prettyUnixTime,
	"prettyBackupStatus": prettyBackupStatus,
	"prettyBytes":        prettyBytes,
	"prettyTimeUnit":     prettyTimeUnit,
}

func prettyStatus(status rms_torrent.Status) string {
	switch status {
	case rms_torrent.Status_Pending:
		return "В очереди"
	case rms_torrent.Status_Downloading:
		return "Загружается"
	case rms_torrent.Status_Done:
		return "Завершено"
	case rms_torrent.Status_Failed:
		return "Ошибка"
	default:
		return ""
	}
}

func prettyFloat(f float32) string {
	return fmt.Sprintf("%.2f", f)
}

func prettyBitrate(bytesPerSecond uint64) string {
	return fmt.Sprintf("%d", uint64(float64(bytesPerSecond*8)/float64(1024*1024)))
}

func prettySize(sizeMB uint64) string {
	return humanize.Bytes(sizeMB * humanize.MByte)
}

func prettyTopic(topic string) string {
	switch topic {
	case pubsub.NotificationTopic:
		return "Уведомление"
	case pubsub.MalfunctionTopic:
		return "Сбой"
	case pubsub.AlertTopic:
		return "Тревога"
	default:
		return topic
	}
}

func fileName(p string) string {
	_, file := path.Split(p)
	return file
}

func prettyBackupType(backupType rms_backup.BackupType) string {
	switch backupType {
	case rms_backup.BackupType_Full:
		return "Полный"
	case rms_backup.BackupType_Partial:
		return "Чстичный"
	default:
		return backupType.String()
	}
}

func prettyUnixTime(unixTime uint64) string {
	t := time.Unix(int64(unixTime), 0)
	return t.Format(time.DateTime)
}

func prettyBackupStatus(status rms_backup.GetBackupStatusResponse_Status) string {
	if status == rms_backup.GetBackupStatusResponse_InProgress {
		return "Выполняется..."
	}
	return "Не запущен"
}

func prettyBytes(bytes uint64) string {
	return humanize.Bytes(bytes)
}

func prettyTimeUnit(t int) string {
	return fmt.Sprintf("%02d", t)
}
