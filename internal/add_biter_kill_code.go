package internal

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"wci/embedded"
	"wci/utils"
)

// AddBiterKillCode appends the biter-killer code to control.lua in the savegame ZIP file if not already present.
func AddBiterKillCode(osName, saveGameZipName string) error {
	log.Info().
		Str("saveGameZipName", saveGameZipName).
		Msg("Starting to inject biter-killer code")

	// Call the generalized code injection function
	err := utils.InjectCodeIntoZip(osName, saveGameZipName, "biter_killer.lua", "control.lua", embedded.LuaInjections)
	if err != nil {
		log.Error().
			Err(err).
			Str("saveGameZipName", saveGameZipName).
			Msg("Failed to inject biter-killer code")
		return fmt.Errorf("failed to inject biter-killer code into '%s': %w", saveGameZipName, err)
	}

	log.Info().
		Str("saveGameZipName", saveGameZipName).
		Msg("Successfully injected biter-killer code")
	return nil
}
