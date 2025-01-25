package utils

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"
)

// GetSaveGameLocation returns the default savegame location for Factorio based on the OS.
// Only Windows and macOS are supported.
func GetSaveGameLocation(osName string) (string, error) {
	log.Debug().Str("osName", osName).Msg("Determining savegame location")
	var saveGameDir string

	switch osName {
	case "windows":
		appData := os.Getenv("APPDATA")
		if appData == "" {
			log.Error().Msg("Environment variable APPDATA is not set")
			return "", fmt.Errorf("environment variable APPDATA is not set")
		}
		saveGameDir = filepath.Join(appData, "Factorio", "saves")
	case "darwin":
		homeDir := os.Getenv("HOME")
		if homeDir == "" {
			log.Error().Msg("Environment variable HOME is not set")
			return "", fmt.Errorf("environment variable HOME is not set")
		}
		saveGameDir = filepath.Join(homeDir, "Library", "Application Support", "Factorio", "saves")
	default:
		log.Warn().Str("osName", osName).Msg("Unsupported operating system")
		return "", fmt.Errorf("unsupported operating system: %s", osName)
	}

	if _, err := os.Stat(saveGameDir); os.IsNotExist(err) {
		log.Warn().Str("directory", saveGameDir).Msg("Savegame directory does not exist")
		return "", fmt.Errorf("savegame directory does not exist: %s", saveGameDir)
	}

	log.Info().Str("directory", saveGameDir).Msg("Savegame directory found")
	return saveGameDir, nil
}
