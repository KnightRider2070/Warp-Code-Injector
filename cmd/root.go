package cmd

import (
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

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = false
}

// Execute is the entry point for the CLI application
func Execute() error {

	rootCmd.SetUsageTemplate(`
Usage:
  wci [command]

Available Commands:
  list       List all savegames
  add-biter-killer   Injects the biter killer script
  clean      Cleans up temporary files

Examples:
  # List all savegames
  wci list

  # Inject the biter killer script into a savegame
  wci add-biter-killer 2
`)

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
