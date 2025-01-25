package utils

import (
	"archive/zip"
	"bytes"
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ModifyZipFile modifies or replaces files in a ZIP archive.
func ModifyZipFile(zipPath string, modifiedFiles map[string][]byte, outputZipPath string) error {
	log.Info().
		Str("zipPath", zipPath).
		Msg("Starting ZIP modification")

	// Open the original ZIP file
	originalZip, err := zip.OpenReader(zipPath)
	if err != nil {
		log.Error().
			Err(err).
			Str("zipPath", zipPath).
			Msg("Failed to open ZIP file")
		return fmt.Errorf("failed to open ZIP file: %w", err)
	}
	defer originalZip.Close()

	var buf bytes.Buffer
	newZip := zip.NewWriter(&buf)

	for _, file := range originalZip.File {
		if newContent, exists := modifiedFiles[file.Name]; exists {
			log.Debug().
				Str("fileName", file.Name).
				Msg("Replacing file with new content")
			if err := AddFileToZip(newZip, file.Name, newContent); err != nil {
				log.Error().
					Err(err).
					Str("fileName", file.Name).
					Msg("Failed to add modified file to ZIP")
				return fmt.Errorf("failed to add modified file '%s': %w", file.Name, err)
			}
		} else {
			log.Trace().
				Str("fileName", file.Name).
				Msg("Copying file to new ZIP without modification")
			if err := CopyZipFile(file, newZip); err != nil {
				log.Error().
					Err(err).
					Str("fileName", file.Name).
					Msg("Failed to copy file to ZIP")
				return fmt.Errorf("failed to copy file '%s' to ZIP: %w", file.Name, err)
			}
		}
	}

	if err := newZip.Close(); err != nil {
		log.Error().Err(err).Msg("Failed to close ZIP writer")
		return fmt.Errorf("failed to close ZIP writer: %w", err)
	}

	if err := os.WriteFile(outputZipPath, buf.Bytes(), 0644); err != nil {
		log.Error().
			Err(err).
			Str("outputZipPath", outputZipPath).
			Msg("Failed to write new ZIP file")
		return fmt.Errorf("failed to write new ZIP file: %w", err)
	}

	log.Info().
		Str("outputZipPath", outputZipPath).
		Msg("ZIP modification completed successfully")
	return nil
}

// AppendToFileInZip modifies a specified file in the ZIP by appending new content to it.
func AppendToFileInZip(zipPath, fileName, newCode, outputZipPath string) error {
	log.Info().
		Str("zipPath", zipPath).
		Str("fileName", fileName).
		Msg("Appending content to file in ZIP")

	// Open the original ZIP file
	originalZip, err := zip.OpenReader(zipPath)
	if err != nil {
		log.Error().
			Err(err).
			Str("zipPath", zipPath).
			Msg("Failed to open ZIP file")
		return fmt.Errorf("failed to open ZIP file: %w", err)
	}
	defer originalZip.Close()

	// Create a buffer to write the new ZIP file
	var buf bytes.Buffer
	newZip := zip.NewWriter(&buf)

	// Track if the specified file was found and modified
	fileModified := false

	// Iterate through the original ZIP file's contents
	for _, file := range originalZip.File {
		if file.Name == fileName {
			log.Debug().
				Str("fileName", file.Name).
				Msg("Target file found in ZIP, appending new content")

			// Read the original file content
			originalContent, err := ReadZipFile(file)
			if err != nil {
				log.Error().
					Err(err).
					Str("fileName", fileName).
					Msg("Failed to read file from ZIP")
				return fmt.Errorf("failed to read '%s': %w", fileName, err)
			}

			// Append the new code to the content
			modifiedContent := append(originalContent, []byte("\n"+newCode+"\n")...)

			// Add the modified file to the new ZIP
			if err := AddFileToZip(newZip, fileName, modifiedContent); err != nil {
				log.Error().
					Err(err).
					Str("fileName", fileName).
					Msg("Failed to add modified file to ZIP")
				return fmt.Errorf("failed to add modified '%s' to ZIP: %w", fileName, err)
			}

			fileModified = true
		} else {
			log.Trace().
				Str("fileName", file.Name).
				Msg("Copying file to new ZIP without modification")
			// Copy other files as-is
			if err := CopyZipFile(file, newZip); err != nil {
				log.Error().
					Err(err).
					Str("fileName", file.Name).
					Msg("Failed to copy file to ZIP")
				return fmt.Errorf("failed to copy file '%s' to ZIP: %w", file.Name, err)
			}
		}
	}

	// Ensure the specified file exists in the ZIP
	if !fileModified {
		log.Warn().
			Str("fileName", fileName).
			Msg("Specified file not found in ZIP")
		return fmt.Errorf("file '%s' not found in the ZIP", fileName)
	}

	// Close the new ZIP writer
	if err := newZip.Close(); err != nil {
		log.Error().
			Err(err).
			Msg("Failed to close ZIP writer")
		return fmt.Errorf("failed to close ZIP writer: %w", err)
	}

	// Write the buffer content to the output file
	if err := os.WriteFile(outputZipPath, buf.Bytes(), 0644); err != nil {
		log.Error().
			Err(err).
			Str("outputZipPath", outputZipPath).
			Msg("Failed to write new ZIP file")
		return fmt.Errorf("failed to write new ZIP file: %w", err)
	}

	log.Info().
		Str("outputZipPath", outputZipPath).
		Msg("Successfully appended content to file in ZIP")
	return nil
}

// CopyZipFile copies a file from the original ZIP to the new ZIP.
func CopyZipFile(file *zip.File, newZip *zip.Writer) error {
	log.Trace().
		Str("fileName", file.Name).
		Msg("Copying file from original ZIP")

	r, err := file.Open()
	if err != nil {
		log.Error().
			Err(err).
			Str("fileName", file.Name).
			Msg("Failed to open file in ZIP")
		return fmt.Errorf("failed to open file '%s': %w", file.Name, err)
	}
	defer r.Close()

	w, err := newZip.Create(file.Name)
	if err != nil {
		log.Error().
			Err(err).
			Str("fileName", file.Name).
			Msg("Failed to create file in new ZIP")
		return fmt.Errorf("failed to create file '%s' in ZIP: %w", file.Name, err)
	}

	if _, err := io.Copy(w, r); err != nil {
		log.Error().
			Err(err).
			Str("fileName", file.Name).
			Msg("Failed to copy file contents to new ZIP")
		return fmt.Errorf("failed to copy contents of file '%s': %w", file.Name, err)
	}

	return nil
}

// ReadZipFile reads the content of a file inside a ZIP archive.
func ReadZipFile(file *zip.File) ([]byte, error) {
	log.Trace().
		Str("fileName", file.Name).
		Msg("Reading file content from ZIP")

	r, err := file.Open()
	if err != nil {
		log.Error().
			Err(err).
			Str("fileName", file.Name).
			Msg("Failed to open file in ZIP")
		return nil, fmt.Errorf("failed to open file '%s': %w", file.Name, err)
	}
	defer r.Close()

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		log.Error().
			Err(err).
			Str("fileName", file.Name).
			Msg("Failed to read file content from ZIP")
		return nil, fmt.Errorf("failed to read file '%s': %w", file.Name, err)
	}

	return buf.Bytes(), nil
}

// AddFileToZip adds a file with its content to the ZIP writer.
func AddFileToZip(zipWriter *zip.Writer, fileName string, content []byte) error {
	log.Trace().
		Str("fileName", fileName).
		Msg("Adding file to ZIP")

	w, err := zipWriter.Create(fileName)
	if err != nil {
		log.Error().
			Err(err).
			Str("fileName", fileName).
			Msg("Failed to create file in ZIP")
		return fmt.Errorf("failed to create file '%s' in ZIP: %w", fileName, err)
	}

	if _, err := w.Write(content); err != nil {
		log.Error().
			Err(err).
			Str("fileName", fileName).
			Msg("Failed to write content to file in ZIP")
		return fmt.Errorf("failed to write content to file '%s' in ZIP: %w", fileName, err)
	}

	return nil
}

// CheckCodeExistsInZip checks if the given code exists in a file inside the ZIP archive.
func CheckCodeExistsInZip(zipPath, fileName, codeToCheck string) (bool, error) {
	log.Info().
		Str("zipPath", zipPath).
		Str("fileName", fileName).
		Msg("Checking if code exists in file inside ZIP")

	// Open the ZIP file
	zipReader, err := zip.OpenReader(zipPath)
	if err != nil {
		log.Error().
			Err(err).
			Str("zipPath", zipPath).
			Msg("Failed to open ZIP file")
		return false, fmt.Errorf("failed to open ZIP file: %w", err)
	}
	defer zipReader.Close()

	// Normalize the target file name for comparison
	normalizedTargetFileName := filepath.Clean(fileName)

	log.Debug().
		Str("normalizedTargetFileName", normalizedTargetFileName).
		Msg("Normalized target file name for comparison")

	// Iterate through files in the ZIP
	for _, file := range zipReader.File {
		// Log all file names in the ZIP for debugging
		log.Trace().
			Str("fileName", file.Name).
			Msg("Found file in ZIP")

		// Normalize file.Name for comparison
		normalizedFileName := filepath.Clean(file.Name)

		// Compare using case-insensitive and trimmed paths
		if strings.EqualFold(normalizedFileName, normalizedTargetFileName) {
			log.Debug().
				Str("fileName", file.Name).
				Msg("Target file found in ZIP, reading content")

			// Read the file's content
			content, err := ReadZipFile(file)
			if err != nil {
				log.Error().
					Err(err).
					Str("fileName", fileName).
					Msg("Failed to read file content in ZIP")
				return false, fmt.Errorf("failed to read file '%s' in ZIP: %w", fileName, err)
			}

			// Check if the code snippet is already present
			exists := strings.Contains(string(content), codeToCheck)
			if exists {
				log.Info().
					Str("fileName", fileName).
					Msg("Code snippet already exists in the file")
			} else {
				log.Debug().
					Str("fileName", fileName).
					Msg("Code snippet not found in the file")
			}
			return exists, nil
		}
	}

	// File not found
	log.Warn().
		Str("fileName", fileName).
		Msg("File not found in ZIP archive")
	return false, fmt.Errorf("file '%s' not found in ZIP", fileName)
}

// FindFileInZip searches for a file by name inside a ZIP archive and returns its full path within the archive.
func FindFileInZip(zipPath, targetFileName string) (string, error) {
	log.Info().
		Str("zipPath", zipPath).
		Str("targetFileName", targetFileName).
		Msg("Searching for file in ZIP")

	// Open the ZIP file
	zipReader, err := zip.OpenReader(zipPath)
	if err != nil {
		log.Error().
			Err(err).
			Str("zipPath", zipPath).
			Msg("Failed to open ZIP file")
		return "", fmt.Errorf("failed to open ZIP file: %w", err)
	}
	defer zipReader.Close()

	// Iterate through files in the ZIP
	for _, file := range zipReader.File {
		log.Trace().
			Str("fileName", file.Name).
			Msg("Checking file in ZIP")

		// Check if the file name matches
		if strings.HasSuffix(file.Name, targetFileName) {
			log.Debug().
				Str("fileName", file.Name).
				Msg("Target file found in ZIP")
			return file.Name, nil
		}
	}

	// File not found
	log.Warn().
		Str("targetFileName", targetFileName).
		Msg("File not found in ZIP")
	return "", fmt.Errorf("file '%s' not found in ZIP", targetFileName)
}
