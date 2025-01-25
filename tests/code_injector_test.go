package tests

import (
	"archive/zip"
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"wci/utils"

	"github.com/stretchr/testify/assert"
)

func TestInjectCodeIntoZipAtFactorioLocation(t *testing.T) {
	// Setup: Create a temporary directory for the test
	tempDir := t.TempDir()

	// Simulate environment variables for Windows and macOS
	os.Setenv("APPDATA", filepath.Join(tempDir, "windows_appdata"))
	os.Setenv("HOME", filepath.Join(tempDir, "macos_home"))

	// Create the savegame directory for Windows
	saveGameDir := filepath.Join(tempDir, "windows_appdata", "Factorio", "saves")
	err := os.MkdirAll(saveGameDir, 0755)
	assert.NoError(t, err, "Failed to create savegame directory")

	// Create a test savegame ZIP file in the savegame directory
	saveGameZipPath := filepath.Join(saveGameDir, "TestSave.zip")
	testFiles := map[string]string{
		"control.lua": "original content",
	}
	err = createTestZip(saveGameZipPath, testFiles)
	assert.NoError(t, err, "Failed to create test ZIP file")

	// Create a mock embedded file for "biter_killer.lua" in the temp directory
	injectionsDir := filepath.Join(tempDir, "lua_injections")
	err = os.MkdirAll(injectionsDir, 0755)
	assert.NoError(t, err, "Failed to create lua_injections directory")

	embeddedFilePath := filepath.Join(injectionsDir, "biter_killer.lua")
	err = os.WriteFile(embeddedFilePath, []byte("-- biter killer code"), 0644)
	assert.NoError(t, err, "Failed to create biter_killer.lua file")

	// Use os.DirFS to emulate the embedded file system
	embeddedFS := os.DirFS(tempDir)

	// Inject code into the ZIP file
	err = utils.InjectCodeIntoZip("windows", "TestSave.zip", "lua_injections/biter_killer.lua", "control.lua", embeddedFS)
	assert.NoError(t, err, "InjectCodeIntoZip failed")

	// Verify the ZIP file was modified correctly
	zipReader, err := zip.OpenReader(saveGameZipPath)
	assert.NoError(t, err, "Failed to open modified ZIP file")
	defer zipReader.Close()

	var controlLuaContent string
	for _, file := range zipReader.File {
		if file.Name == "control.lua" {
			reader, err := file.Open()
			assert.NoError(t, err, "Failed to open control.lua")
			defer reader.Close()

			buf := new(bytes.Buffer)
			_, err = buf.ReadFrom(reader)
			assert.NoError(t, err, "Failed to read control.lua")

			controlLuaContent = buf.String()
			break
		}
	}

	// Validate that the original content and the injected content are present
	assert.Contains(t, controlLuaContent, "original content", "Original content missing in control.lua")
	assert.Contains(t, controlLuaContent, "-- biter killer code", "Injected code missing in control.lua")
}
