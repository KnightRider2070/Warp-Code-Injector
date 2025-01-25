package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"wci/utils"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all savegames in the default directory, sorted by creation date",
	Run: func(cmd *cobra.Command, args []string) {
		// List savegames with numbers
		saveGames, err := utils.ListSaveGamesWithNumbers(currentOS)
		if err != nil {
			// Print error and exit if no savegames are found or an error occurs
			fmt.Fprintf(os.Stderr, "Error listing savegames: %v\n", err)
			os.Exit(1)
		}

		listedSaveGames = saveGames

		// Print savegames to the user
		utils.PrintSaveGames(saveGames)
	},
}

func init() {
	// Add the "list" subcommand to the root command
	rootCmd.AddCommand(listCmd)
}
