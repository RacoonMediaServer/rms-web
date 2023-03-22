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

// Configuration represents entire service configuration
type Configuration struct {
	Http     configuration.Http
	Cctv     Cctv
	Services []Service
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
