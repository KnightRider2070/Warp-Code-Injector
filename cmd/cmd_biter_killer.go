package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"wci/internal"
)

var addBiterKillerCmd = &cobra.Command{
	Use:   "add-biter-killer [number]",
	Short: "Add biter-killer Lua script to the selected savegame",
	Long: `Appends the biter-killer Lua script to the 'control.lua' file of the selected savegame ZIP file
based on the savegame number obtained from the 'list' command.`,
	Args: cobra.ExactArgs(1), // Requires exactly one argument (the savegame number)
	Run: func(cmd *cobra.Command, args []string) {
		// Ensure savegames were listed before this command
		if len(listedSaveGames) == 0 {
			fmt.Fprintln(os.Stderr, "Error: No savegames listed. Run 'wci list' first.")
			os.Exit(1)
		}

		// Parse the input number
		var saveGameNumber int
		_, err := fmt.Sscanf(args[0], "%d", &saveGameNumber)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Invalid savegame number '%s'. Please provide a valid number.\n", args[0])
			os.Exit(1)
		}

		// Validate the savegame number
		saveGamePath, exists := listedSaveGames[saveGameNumber]
		if !exists {
			fmt.Fprintf(os.Stderr, "Error: Savegame number '%d' not found. Run 'wci list' to see available savegames.\n", saveGameNumber)
			os.Exit(1)
		}

		// Full path to the selected savegame ZIP file
		saveGameZipPath := saveGamePath

		// Inject the biter-killer code
		err = internal.AddBiterKillCode(currentOS, saveGameZipPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error adding biter-killer code to '%s': %v\n", saveGameZipPath, err)
			os.Exit(1)
		}

		fmt.Printf("Successfully added biter-killer code to '%s'.\n", saveGameZipPath)
	},
}

func init() {
	rootCmd.AddCommand(addBiterKillerCmd)
}
