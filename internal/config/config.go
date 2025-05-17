package config

import (
	"os"
	"path/filepath"
	"runtime"
)

// Config representa a configuração do agente
type Config struct {
	LogPath    string
	DataPath   string
	ConfigPath string
	Platform   string
}

// New retorna uma nova configuração baseada no sistema operacional
func New() *Config {
	config := &Config{
		Platform: runtime.GOOS,
	}

	// Obtém o diretório atual
	currentDir, err := os.Getwd()
	if err != nil {
		currentDir = "."
	}

	// Define os caminhos base dependendo do sistema operacional
	if runtime.GOOS == "windows" {
		config.LogPath = filepath.Join(currentDir, "logs")
		config.DataPath = filepath.Join(currentDir, "data")
		config.ConfigPath = filepath.Join(currentDir, "config")
	} else {
		config.LogPath = filepath.Join(currentDir, "logs")
		config.DataPath = filepath.Join(currentDir, "data")
		config.ConfigPath = filepath.Join(currentDir, "config")
	}

	return config
}
