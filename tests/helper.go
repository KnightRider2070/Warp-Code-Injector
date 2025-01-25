package tests

import (
	"archive/zip"
	"os"
)

// createTestZip creates a sample ZIP file for testing.
func createTestZip(filePath string, files map[string]string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	zipWriter := zip.NewWriter(file)
	defer zipWriter.Close()

	for name, content := range files {
		writer, err := zipWriter.Create(name)
		if err != nil {
			return err
		}
		_, err = writer.Write([]byte(content))
		if err != nil {
			return err
		}
	}

	return nil
}
