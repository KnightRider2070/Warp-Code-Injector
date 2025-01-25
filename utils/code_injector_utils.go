package utils

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"io/fs"
	"path/filepath"
)

// InjectCodeIntoZip handles injecting code from an embedded file into a target file inside a savegame ZIP file.
// Parameters:
// - osName: the name of the operating system (e.g., "windows", "darwin").
// - saveGameZipName: the name of the savegame ZIP file.
// - embeddedFileName: the name of the embedded file containing the code to inject.
// - targetFileName: the name of the target file inside the ZIP to which the code should be injected/appended.
// InjectCodeIntoZip handles injecting code from an embedded file into a target file inside a savegame ZIP file.
func InjectCodeIntoZip(osName, saveGameZipName, embeddedFileName, targetFileName string, fileSystem fs.FS) error {
	// Retrieve the base savegame directory based on the OS
	baseDir, err := GetSaveGameLocation(osName)
	if err != nil {
		log.Error().
			Err(err).
			Str("osName", osName).
			Msg("Failed to retrieve savegame directory")
		return fmt.Errorf("failed to retrieve savegame directory for OS '%s': %w", osName, err)
	}

	// Construct the full path to the savegame ZIP file
	saveGameZipPath := filepath.Join(baseDir, saveGameZipName)

	log.Info().
		Str("zipPath", saveGameZipPath).
		Str("embeddedFileName", embeddedFileName).
		Str("targetFileName", targetFileName).
		Msg("Starting to inject code into ZIP")

	// Read the embedded code to inject from the provided file system
	codeToInject, err := fs.ReadFile(fileSystem, embeddedFileName)
	if err != nil {
		log.Error().
			Err(err).
			Str("file", embeddedFileName).
			Msg("Failed to read embedded file")
		return fmt.Errorf("failed to read embedded file '%s': %w", embeddedFileName, err)
	}

	// Locate the target file in the ZIP
	log.Info().
		Str("zipPath", saveGameZipPath).
		Msg("Searching for the target file inside ZIP")
	targetPathInZip, err := FindFileInZip(saveGameZipPath, targetFileName)
	if err != nil {
		log.Error().
			Err(err).
			Str("file", targetFileName).
			Msg("Failed to locate target file in ZIP")
		return fmt.Errorf("failed to locate '%s' in ZIP: %w", targetFileName, err)
	}

	// Check if the code to inject already exists in the target file
	exists, err := CheckCodeExistsInZip(saveGameZipPath, targetPathInZip, string(codeToInject))
	if err != nil {
		log.Error().
			Err(err).
			Str("file", targetPathInZip).
			Msg("Failed to check code existence in ZIP")
		return fmt.Errorf("failed to check if code exists in '%s': %w", targetPathInZip, err)
	}

	// If the code already exists, log a warning and exit
	if exists {
		log.Warn().
			Str("file", targetPathInZip).
			Msg("Code snippet already exists in the target file")
		return nil
	}

	// Append the code to the target file
	err = AppendToFileInZip(saveGameZipPath, targetPathInZip, string(codeToInject), saveGameZipPath)
	if err != nil {
		log.Error().
			Err(err).
			Str("file", targetPathInZip).
			Msg("Failed to append code to file in ZIP")
		return fmt.Errorf("failed to append code to '%s': %w", targetPathInZip, err)
	}

	log.Info().
		Str("file", targetPathInZip).
		Msg("Successfully injected code into the target file")
	return nil
}
