package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var cmdCleanupSaveFile = &cobra.Command{
	Use:   "clean",
	Short: "Removes the savegames.json file from the executable's running location",
	Long: `This command removes the savegames.json file from the directory where the executable is running.
It is useful for cleaning up outdated or unnecessary savegame data files.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get the current directory of the executable
		exePath, err := os.Executable()
		if err != nil {
			log.Error().Err(err).Msg("Failed to get executable path")
			fmt.Println("Error: Could not determine executable path.")
			return
		}

		exeDir := filepath.Dir(exePath)
		saveFilePath := filepath.Join(exeDir, "savegames.json")

		// Log the file path
		absPath, _ := filepath.Abs(saveFilePath)
		log.Debug().Str("absPath", absPath).Msg("Absolute path to savegames.json")

		// Ensure the file exists
		fileInfo, err := os.Stat(saveFilePath)
		if os.IsNotExist(err) {
			log.Warn().Str("file", saveFilePath).Msg("File not found")
			fmt.Println("The savegames.json file does not exist. Nothing to clean up.")
			return
		} else if err != nil {
			log.Error().Err(err).Str("file", saveFilePath).Msg("Error checking file existence")
			fmt.Printf("Error: Could not access savegames.json file: %v\n", err)
			return
		}

		log.Debug().
			Str("file", saveFilePath).
			Int64("size", fileInfo.Size()).
			Bool("isDir", fileInfo.IsDir()).
			Msg("File details before removal")

		// Attempt to remove the file
		err = os.Remove(saveFilePath)
		if err != nil {
			log.Error().Err(err).Str("file", saveFilePath).Msg("Failed to remove file")
			fmt.Printf("Error: Could not remove savegames.json file: %v\n", err)
			return
		}

		// Double-check if the file still exists
		time.Sleep(1 * time.Second) // Allow time for any external interference
		if _, err := os.Stat(saveFilePath); err == nil {
			log.Warn().Str("file", saveFilePath).Msg("File still exists after removal")
			fmt.Println("Warning: savegames.json file still exists after removal.")
			return
		}

		log.Info().Str("file", saveFilePath).Msg("File removed successfully")
		fmt.Println("The savegames.json file has been successfully removed.")
	},
}

func init() {
	rootCmd.AddCommand(cmdCleanupSaveFile)
}
