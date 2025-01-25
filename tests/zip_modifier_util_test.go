package tests

import (
	"archive/zip"
	"os"
	"testing"
	"wci/utils"

	"github.com/stretchr/testify/assert"
)

// TestModifyZipFile tests the ModifyZipFile function.
func TestModifyZipFile(t *testing.T) {
	testZipPath := "test_modify.zip"
	outputZipPath := "output_modify.zip"
	defer os.Remove(testZipPath)
	defer os.Remove(outputZipPath)

	files := map[string]string{
		"file1.txt": "This is file 1.",
		"file2.txt": "This is file 2.",
	}

	assert.NoError(t, createTestZip(testZipPath, files))

	modifiedFiles := map[string][]byte{
		"file1.txt": []byte("Modified content for file 1."),
	}

	err := utils.ModifyZipFile(testZipPath, modifiedFiles, outputZipPath)
	assert.NoError(t, err)

	outputZip, err := zip.OpenReader(outputZipPath)
	assert.NoError(t, err)
	defer outputZip.Close()

	found := false
	for _, file := range outputZip.File {
		if file.Name == "file1.txt" {
			content, err := utils.ReadZipFile(file)
			assert.NoError(t, err)
			assert.Equal(t, "Modified content for file 1.", string(content))
			found = true
		}
	}

	assert.True(t, found, "Modified file not found in ZIP")
}

// TestAppendToFileInZip tests the AppendToFileInZip function.
func TestAppendToFileInZip(t *testing.T) {
	testZipPath := "test_append.zip"
	outputZipPath := "output_append.zip"
	defer os.Remove(testZipPath)
	defer os.Remove(outputZipPath)

	files := map[string]string{
		"file1.txt": "This is file 1.",
	}

	assert.NoError(t, createTestZip(testZipPath, files))

	newContent := " Appended content."

	err := utils.AppendToFileInZip(testZipPath, "file1.txt", newContent, outputZipPath)
	assert.NoError(t, err)

	outputZip, err := zip.OpenReader(outputZipPath)
	assert.NoError(t, err)
	defer outputZip.Close()

	found := false
	for _, file := range outputZip.File {
		if file.Name == "file1.txt" {
			content, err := utils.ReadZipFile(file)
			assert.NoError(t, err)
			assert.Equal(t, "This is file 1.\n Appended content.\n", string(content))
			found = true
		}
	}

	assert.True(t, found, "Appended file not found in ZIP")
}

// TestCheckCodeExistsInZip tests the CheckCodeExistsInZip function.
func TestCheckCodeExistsInZip(t *testing.T) {
	testZipPath := "test_check.zip"
	defer os.Remove(testZipPath)

	files := map[string]string{
		"file1.txt": "This is file 1.",
		"file2.txt": "Special code snippet.",
	}

	assert.NoError(t, createTestZip(testZipPath, files))

	// Test when code exists
	exists, err := utils.CheckCodeExistsInZip(testZipPath, "file2.txt", "Special code snippet.")
	assert.NoError(t, err)
	assert.True(t, exists)

	// Test when code does not exist
	exists, err = utils.CheckCodeExistsInZip(testZipPath, "file1.txt", "Special code snippet.")
	assert.NoError(t, err)
	assert.False(t, exists)

	// Test for a non-existent file
	exists, err = utils.CheckCodeExistsInZip(testZipPath, "file3.txt", "Special code snippet.")
	assert.Error(t, err)
	assert.False(t, exists)
}

// TestFindFileInZip tests the FindFileInZip function.
func TestFindFileInZip(t *testing.T) {
	testZipPath := "test_find.zip"
	defer os.Remove(testZipPath)

	files := map[string]string{
		"dir1/file1.txt": "This is file 1.",
		"dir2/file2.txt": "This is file 2.",
	}

	assert.NoError(t, createTestZip(testZipPath, files))

	// Test finding a file
	filePath, err := utils.FindFileInZip(testZipPath, "file1.txt")
	assert.NoError(t, err)
	assert.Equal(t, "dir1/file1.txt", filePath)

	// Test for a non-existent file
	filePath, err = utils.FindFileInZip(testZipPath, "file3.txt")
	assert.Error(t, err)
	assert.Equal(t, "", filePath)
}
