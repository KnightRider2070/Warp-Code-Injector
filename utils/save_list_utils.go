package utils

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
	"sort"
	"strings"
)

// saveGameInfo represents metadata for a savegame file.
type saveGameInfo struct {
	Name    string
	ModTime int64
}

// ListSaveGamesWithNumbers lists all savegames in the savegame directory for the given OS.
// The returned map contains the numbered savegames, with their full names (including .zip) for internal use,
// but the numbers printed exclude the .zip extension.
func ListSaveGamesWithNumbers(osName string) (map[int]string, error) {
	log.Debug().Str("osName", osName).Msg("Getting savegame location")
	saveGameDir, err := GetSaveGameLocation(osName)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get savegame directory")
		return nil, err
	}

	log.Info().Str("directory", saveGameDir).Msg("Reading savegame directory")
	files, err := os.ReadDir(saveGameDir)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read savegame directory")
		return nil, fmt.Errorf("failed to read savegame directory: %v", err)
	}

	var saveGameInfos []saveGameInfo
	for _, file := range files {
		if file.Type().IsRegular() && strings.HasSuffix(file.Name(), ".zip") {
			info, err := file.Info()
			if err != nil {
				log.Warn().Str("file", file.Name()).Msg("Skipping file due to metadata read failure")
				continue
			}
			saveGameInfos = append(saveGameInfos, saveGameInfo{
				Name:    file.Name(),
				ModTime: info.ModTime().Unix(),
			})
		}
	}

	if len(saveGameInfos) == 0 {
		log.Warn().Msg("No savegames found in the directory")
		return nil, fmt.Errorf("no savegames found in the directory: %s", saveGameDir)
	}

	log.Debug().Msg("Sorting savegames by modification time")
	sort.Slice(saveGameInfos, func(i, j int) bool {
		return saveGameInfos[i].ModTime < saveGameInfos[j].ModTime
	})

	saveGames := make(map[int]string)
	for i, saveGame := range saveGameInfos {
		saveGames[i+1] = saveGame.Name
		log.Trace().
			Int("index", i+1).
			Str("savegame", saveGame.Name).
			Msg("Added savegame to the list")
	}

	return saveGames, nil
}

// PrintSaveGames prints the list of savegames with numbers and their modification dates for user selection.
// The `.zip` extension is omitted in the printed output, but the numbers correspond to the internal saveGames map.
func PrintSaveGames(saveGames map[int]string) {
	// Extract keys from the map and sort them
	keys := make([]int, 0, len(saveGames))
	for k := range saveGames {
		keys = append(keys, k)
	}
	sort.Ints(keys) // Sort keys numerically

	// Print savegames in sorted order of keys
	fmt.Println("Available savegames (sorted by date):")
	for _, key := range keys {
		// Remove the .zip extension for display purposes
		fmt.Printf("%d. %s\n", key, strings.TrimSuffix(saveGames[key], ".zip"))
	}
}
