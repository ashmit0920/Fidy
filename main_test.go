package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func test(t *testing.T) {
	// Create a temporary directory for testing
	testDir, err := os.MkdirTemp("", "fidy_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(testDir) // Clean up after the test

	// test files with different extensions
	testFiles := []string{
		"file1.txt",
		"file2.txt",
		"file3.jpg",
		"file4.pdf",
	}
	for _, fileName := range testFiles {
		_, err := os.Create(filepath.Join(testDir, fileName))
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", fileName, err)
		}
	}

	// Call the function to organize files
	organizeFiles(testDir, []string{}, []string{}, false, false)

	// Check if the files were moved to the correct directories
	expectedDirs := []string{"txt", "jpg", "pdf"}
	for _, dir := range expectedDirs {
		dirPath := filepath.Join(testDir, dir)
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			t.Errorf("Expected directory %s does not exist", dirPath)
		}
	}

	// Check if the files are in the correct directories
	expectedFiles := map[string]string{
		"file1.txt": "txt",
		"file2.txt": "txt",
		"file3.jpg": "jpg",
		"file4.pdf": "pdf",
	}
	for fileName, dir := range expectedFiles {
		filePath := filepath.Join(testDir, dir, fileName)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Errorf("Expected file %s does not exist in directory %s", fileName, dir)
		}
	}
}

// function to replicate fidy's logic
func organizeFiles(dir string, excludeExtensions []string, includeExtensions []string, verbose bool, dryRun bool) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	createdDirs := make(map[string]bool)

	for _, file := range files {
		if !file.IsDir() {
			ext := filepath.Ext(file.Name())
			if ext != "" {
				ext = ext[1:] // Remove the leading dot

				// Check if the extension is in the exclude list
				exclude := false
				for _, excludeExt := range excludeExtensions {
					if ext == excludeExt {
						exclude = true
						break
					}
				}
				if exclude {
					continue
				}

				inc := false
				for _, includeExt := range includeExtensions {
					if ext == includeExt {
						inc = true
						break
					}
				}
				if !inc {
					continue
				}

				targetDir := filepath.Join(dir, ext)
				if _, exists := createdDirs[targetDir]; !exists {
					if verbose || dryRun {
						fmt.Printf("Creating directory: %s\n", targetDir)
					}
					if !dryRun {
						os.Mkdir(targetDir, os.ModePerm)
					}
					createdDirs[targetDir] = true
				}

				oldPath := filepath.Join(dir, file.Name())
				newPath := filepath.Join(targetDir, file.Name())
				if verbose || dryRun {
					fmt.Printf("Moving file: %s -> %s\n", oldPath, newPath)
				}
				if !dryRun {
					os.Rename(oldPath, newPath)
				}
			}
		}
	}
}
