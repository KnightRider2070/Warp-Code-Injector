package internal

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"path/filepath"
	"wci/embedded"
	"wci/utils"
)

// AddBiterKillCode appends the biter-killer code to control.lua in the savegame ZIP file if not already present.
func AddBiterKillCode(osName, saveGameZipName string) error {
	embeddedFileName := "biter_killer.lua"
	targetFileName := "control.lua" // Target file to modify inside the ZIP

	// Retrieve the base savegame directory based on the OS
	baseDir, err := utils.GetSaveGameLocation(osName)
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
		Msg("Starting to append biter-killer code")

	// Read the embedded biter-killer code
	biterCodeToAdd, err := embedded.ReadEmbeddedFile(embedded.LuaInjections, "lua_injections", embeddedFileName)
	if err != nil {
		log.Error().
			Err(err).
			Str("file", embeddedFileName).
			Msg("Failed to read embedded file")
		return fmt.Errorf("failed to read embedded file '%s': %w", embeddedFileName, err)
	}

	// Locate the target file in the ZIP (e.g., "Test - NoBiter/control.lua")
	log.Info().
		Str("zipPath", saveGameZipPath).
		Msg("Searching for the target file inside ZIP")
	targetPathInZip, err := utils.FindFileInZip(saveGameZipPath, targetFileName)
	if err != nil {
		log.Error().
			Err(err).
			Str("file", targetFileName).
			Msg("Failed to locate target file in ZIP")
		return fmt.Errorf("failed to locate '%s' in ZIP: %w", targetFileName, err)
	}

	// Check if the biter-killer code already exists in control.lua
	exists, err := utils.CheckCodeExistsInZip(saveGameZipPath, targetPathInZip, string(biterCodeToAdd))
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
			Msg("Biter-killer code already exists in the file")
		return nil
	}

	// Append the biter-killer code to control.lua
	err = utils.AppendToFileInZip(saveGameZipPath, targetPathInZip, string(biterCodeToAdd), saveGameZipPath)
	if err != nil {
		log.Error().
			Err(err).
			Str("file", targetPathInZip).
			Msg("Failed to append code to file in ZIP")
		return fmt.Errorf("failed to append code to '%s': %w", targetPathInZip, err)
	}

	log.Info().
		Str("file", targetPathInZip).
		Msg("Successfully appended biter-killer code")
	return nil
}
