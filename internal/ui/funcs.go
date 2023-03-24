package ui

import (
	"fmt"
	rms_torrent "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-torrent"
	"html/template"
)

var Functions template.FuncMap = template.FuncMap{
	"prettyStatus": prettyStatus,
	"prettyFloat":  prettyFloat,
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
