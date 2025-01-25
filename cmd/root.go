package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"runtime"
)

// Global variable to hold the current operating system
var (
	currentOS       = runtime.GOOS
	listedSaveGames map[int]string
	saveGamesFile   = "savegames.json" // File to store savegames data
)

var rootCmd = &cobra.Command{
	Use:   "wci",
	Short: "Inject Lua scripts into Factorio savegames",
	Long:  "Warp Code Injector (wci) modifies Factorio savegames by injecting custom Lua scripts.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to Warp Code Injector! Use 'wci --help' to see available commands.")
	},
}

// Execute is the entry point for the CLI application
func Execute() error {
	// Load listedSaveGames from file at startup
	if err := loadListedSaveGames(); err != nil {
		fmt.Printf("Failed to load savegames data: %v\n", err)
	}

	if err := rootCmd.Execute(); err != nil {
		// Attempt to print the error
		if _, fErr := fmt.Fprintf(os.Stderr, "Error: %v\n", err); fErr != nil {
			fmt.Printf("Unhandled error writing to Stderr: %v\n", fErr)
			return fErr
		}
		return err
	}

	// Save listedSaveGames to file on shutdown
	if err := saveListedSaveGames(); err != nil {
		fmt.Printf("Failed to save savegames data: %v\n", err)
	}

	return nil
}

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
