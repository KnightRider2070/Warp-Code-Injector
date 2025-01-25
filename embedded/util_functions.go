package embedded

import (
	"embed"
	"fmt"
	"github.com/rs/zerolog/log"
	"path"
)

// ListEmbeddedFiles lists all files embedded from the specified directory in the provided embed.FS.
func ListEmbeddedFiles(fs embed.FS, dirName string) ([]string, error) {
	log.Info().
		Str("directory", dirName).
		Msg("Listing files in embedded directory")

	entries, err := fs.ReadDir(dirName)
	if err != nil {
		log.Error().
			Err(err).
			Str("directory", dirName).
			Msg("Failed to read embedded directory")
		return nil, fmt.Errorf("failed to read embedded directory '%s': %w", dirName, err)
	}

	var fileNames []string
	for _, entry := range entries {
		if !entry.IsDir() {
			fileNames = append(fileNames, entry.Name())
			log.Debug().
				Str("fileName", entry.Name()).
				Msg("Found embedded file")
		}
	}

	if len(fileNames) == 0 {
		log.Warn().
			Str("directory", dirName).
			Msg("No files found in embedded directory")
	}

	log.Info().
		Int("fileCount", len(fileNames)).
		Msg("Completed listing embedded files")
	return fileNames, nil
}

// ReadEmbeddedFile reads the content of an embedded file from the specified directory in the provided embed.FS.
func ReadEmbeddedFile(fs embed.FS, dirName, fileName string) ([]byte, error) {
	filePath := path.Join(dirName, fileName)
	log.Info().
		Str("filePath", filePath).
		Msg("Reading embedded file")

	content, err := fs.ReadFile(filePath)
	if err != nil {
		log.Error().
			Err(err).
			Str("filePath", filePath).
			Msg("Failed to read embedded file")
		return nil, fmt.Errorf("failed to read embedded file '%s': %w", filePath, err)
	}

	log.Debug().
		Str("filePath", filePath).
		Int("contentLength", len(content)).
		Msg("Successfully read embedded file")
	return content, nil
}
