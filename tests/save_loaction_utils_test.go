package tests

import (
	"os"
	"path/filepath"
	"testing"
	"wci/utils"

	"github.com/stretchr/testify/assert"
)

// TestGetSaveGameLocation tests the GetSaveGameLocation function.
func TestGetSaveGameLocation(t *testing.T) {
	// Save original environment variables to restore them later
	originalAppData := os.Getenv("APPDATA")
	originalHome := os.Getenv("HOME")
	defer func() {
		os.Setenv("APPDATA", originalAppData)
		os.Setenv("HOME", originalHome)
	}()

	t.Run("Windows - Valid APPDATA", func(t *testing.T) {
		mockAppData := t.TempDir()
		os.Setenv("APPDATA", mockAppData)

		saveGameDir := filepath.Join(mockAppData, "Factorio", "saves")
		os.MkdirAll(saveGameDir, 0755)

		result, err := utils.GetSaveGameLocation("windows")
		assert.NoError(t, err)
		assert.Equal(t, saveGameDir, result)
	})

	t.Run("Windows - Missing APPDATA", func(t *testing.T) {
		os.Unsetenv("APPDATA")

		_, err := utils.GetSaveGameLocation("windows")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "environment variable APPDATA is not set")
	})

	t.Run("MacOS - Valid HOME", func(t *testing.T) {
		mockHome := t.TempDir()
		os.Setenv("HOME", mockHome)

		saveGameDir := filepath.Join(mockHome, "Library", "Application Support", "Factorio", "saves")
		os.MkdirAll(saveGameDir, 0755)

		result, err := utils.GetSaveGameLocation("darwin")
		assert.NoError(t, err)
		assert.Equal(t, saveGameDir, result)
	})

	t.Run("MacOS - Missing HOME", func(t *testing.T) {
		os.Unsetenv("HOME")

		_, err := utils.GetSaveGameLocation("darwin")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "environment variable HOME is not set")
	})

	t.Run("Unsupported OS", func(t *testing.T) {
		_, err := utils.GetSaveGameLocation("linux")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported operating system")
	})

	t.Run("Directory Does Not Exist", func(t *testing.T) {
		mockAppData := t.TempDir()
		os.Setenv("APPDATA", mockAppData)

		_, err := utils.GetSaveGameLocation("windows")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "savegame directory does not exist")
	})
}
