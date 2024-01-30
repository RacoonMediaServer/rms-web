package config

import "github.com/RacoonMediaServer/rms-packages/pkg/configuration"

// Cctv settings of VMS
type Cctv struct {
	Enabled bool
}

// Service represetns external service link
type Service struct {
	Title       string
	Description string
	Address     string
}

// Content controls fetching multimedia content
type Content struct {
	Directory string
	Backups   string
}

// Configuration represents entire service configuration
type Configuration struct {
	Http     configuration.Http
	Cctv     Cctv
	Content  Content
	Services []Service
	Bot      string
}

var config Configuration

// Load open and parses configuration file
func Load(configFilePath string) error {
	return configuration.Load(configFilePath, &config)
}

// Config returns loaded configuration
func Config() Configuration {
	return config
}
