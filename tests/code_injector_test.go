package tests

import (
	"archive/zip"
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"wci/utils"

	"github.com/stretchr/testify/assert"
)

func TestInjectCodeIntoZipAtFactorioLocation(t *testing.T) {

	// Retrive os of the test system
	osName := runtime.GOOS

	// Retrieve the savegame location for the current OS
	saveGameDir, err := utils.GetSaveGameLocation(osName)
	assert.NoError(t, err)

	// Ensure the savegame directory exists
	err = os.MkdirAll(saveGameDir, 0755)
	assert.NoError(t, err)

	// Create a test savegame ZIP file in the Factorio savegame directory
	saveGameZipPath := filepath.Join(saveGameDir, "TestSave.zip")
	testFiles := map[string]string{
		"control.lua": "original content",
	}
	err = createTestZip(saveGameZipPath, testFiles)
	assert.NoError(t, err)

	// Create a mock embedded file in a temporary directory
	tempDir := t.TempDir()
	embeddedFilePath := filepath.Join(tempDir, "lua_injections", "biter_killer.lua")
	err = os.MkdirAll(filepath.Dir(embeddedFilePath), 0755)
	assert.NoError(t, err)

	err = os.WriteFile(embeddedFilePath, []byte("-- biter killer code"), 0644)
	assert.NoError(t, err)

	// Use os.DirFS to emulate the embedded FS
	embeddedFS := os.DirFS(tempDir)

	// Run the InjectCodeIntoZip function
	err = utils.InjectCodeIntoZip("windows", "TestSave.zip", "lua_injections/biter_killer.lua", "control.lua", embeddedFS)
	assert.NoError(t, err)

	// Verify the ZIP file was modified correctly
	zipReader, err := zip.OpenReader(saveGameZipPath)
	assert.NoError(t, err)
	defer zipReader.Close()

	var controlLuaContent string
	for _, file := range zipReader.File {
		if file.Name == "control.lua" {
			reader, err := file.Open()
			assert.NoError(t, err)
			defer reader.Close()

			buf := new(bytes.Buffer)
			_, err = buf.ReadFrom(reader)
			assert.NoError(t, err)

			controlLuaContent = buf.String()
			break
		}
	}

	// Check that the new content was appended to the file
	assert.Contains(t, controlLuaContent, "original content")
	assert.Contains(t, controlLuaContent, "-- biter killer code")
}
