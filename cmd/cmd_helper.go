package cmd

import (
	"encoding/json"
	"fmt"
	"os"
)

// saveListedSaveGames saves the listedSaveGames map to a file
func saveListedSaveGames() error {
	file, err := os.Create(saveGamesFile)
	if err != nil {
		return fmt.Errorf("failed to create savegames file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(listedSaveGames); err != nil {
		return fmt.Errorf("failed to encode savegames data: %w", err)
	}

	return nil
}

// loadListedSaveGames loads the listedSaveGames map from a file
func loadListedSaveGames() error {
	if _, err := os.Stat(saveGamesFile); os.IsNotExist(err) {
		listedSaveGames = make(map[int]string) // Initialize with an empty map if file doesn't exist
		return nil
	}

	file, err := os.Open(saveGamesFile)
	if err != nil {
		return fmt.Errorf("failed to open savegames file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&listedSaveGames); err != nil {
		return fmt.Errorf("failed to decode savegames data: %w", err)
	}

	return nil
}
