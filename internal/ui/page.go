package ui

import "github.com/RacoonMediaServer/rms-web/internal/config"

// PageContext contains base fields for UI page representation
type PageContext struct {
	CctvEnabled bool
	Services    []config.Service
	Redirect    string
}

func New() *PageContext {
	cfg := config.Config()
	return &PageContext{
		CctvEnabled: cfg.Cctv.Enabled,
		Services:    cfg.Services,
	}
}
