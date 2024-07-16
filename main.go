package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Name string `json:"name"`
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configDir := filepath.Join(homeDir, ".fidy")
	if err := os.MkdirAll(configDir, os.ModePerm); err != nil {
		return "", err
	}
	return filepath.Join(configDir, "config.json"), nil
}

func loadConfig() (*Config, error) {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return nil, err
	}
	file, err := os.ReadFile(configFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{}, nil
		}
		return nil, err
	}
	var config Config
	err = json.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func saveConfig(config *Config) error {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return err
	}
	file, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configFilePath, file, os.ModePerm)
}

func main() {

	name := flag.String("name", "", "The user's name")
	dir := flag.String("dir", "", "The directory to organize")
	help := flag.Bool("help", false, "Show information about fidy")
	include := flag.String("include", "", "Comma-separated list of extensions to include")
	exclude := flag.String("exclude", "", "Comma-separated list of extensions to exclude")
	verbose := flag.Bool("verbose", false, "Enable verbose output")
	dryrun := flag.Bool("dryrun", false, "Simulate the file organization without doing any actual changes")
	cleanAll := flag.Bool("cleanAll", false, "Delete all the empty folders and sub-folders in the specified directory after organizing files.")

	flag.Parse()

	config, err := loadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	// Update the config file with name (if provided)
	if *name != "" {
		config.Name = *name
		if err := saveConfig(config); err != nil {
			fmt.Println("Error saving config:", err)
			return
		}
		fmt.Println("Name updated! Nice to meet you", *name)
	}

	// Default message without flags
	if len(os.Args) == 1 {
		if config.Name != "" {
			fmt.Printf("Hey, I am Fidy. Nice to see you, %s!\n", config.Name)
		} else {
			fmt.Println("Hey, I am Fidy. You can let me know your name by using 'fidy -name YOUR_NAME' for our future conversations!")
		}
	}

	if *help {
		fmt.Println("")
		fmt.Println("     ________   _______     _______      ___     ___")
		fmt.Println("    / ______/  /__  __/    / _____ \\    /  /    /  /")
		fmt.Println("   / /____       / /      / /    / /   /  /____/  /")
		fmt.Println("  / _____/      / /      / /    / /   /___   ____/")
		fmt.Println(" / /         __/ /__    / /____/ /       /  /")
		fmt.Println("/_/         /______/   /________/       /__/")

		fmt.Println("\n---------- The File Organizer CLI Tool ----------")
		fmt.Println("\nFidy helps you organize your files by sorting them into directories based on their extensions.")
		fmt.Println("\nUsage:")
		fmt.Println("")
		fmt.Println("  -help           : Show information about Fidy.")
		fmt.Println("  -name <name>    : Set your name to personalize Fidy's greetings.")
		fmt.Println("  -dir <path>     : Specify the directory to organize. Use 'fidy -dir .' for current directory.")
		fmt.Println("  -include <exts> : Comma-separated list of extensions to include.")
		fmt.Println("  -exclude <exts> : Comma-separated list of extensions to exclude.")
		fmt.Println("  -verbose        : Enable verbose output.")
		fmt.Println("  -dryrun         : Simulate the file organization without doing any actual changes.")
		fmt.Println("  -cleanAll       : Delete all the empty folders and sub-folders in the specified directory after organizing files.")
		fmt.Println("")
		return
	}

	if *dir != "" {
		excludeExtensions := strings.Split(*exclude, ",") // list of excluded extensions
		includeExtensions := strings.Split(*include, ",") // list of included extensions

		createdDirs := make(map[string]bool) // Tracking created directories for verbose/dryrun mode

		files, err := os.ReadDir(*dir)
		if err != nil {
			fmt.Println("Error reading directory: ", err)
			os.Exit(1)
		}

		for _, file := range files {
			if !file.IsDir() { // only run the loop for a file, not a folder
				ext := filepath.Ext(file.Name())
				if ext != "" {
					ext = ext[1:] // Removing the dot

					// Checking if extension is in exclude list
					exc := false
					for _, excludeExt := range excludeExtensions {
						if ext == excludeExt {
							exc = true
							break
						}
					}
					if exc {
						continue
					}

					if *include != "" {
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
					}

					targetDir := filepath.Join(*dir, ext)

					if _, err := os.Stat(targetDir); os.IsNotExist(err) {
						// checking if directory is already created
						if _, exists := createdDirs[targetDir]; !exists {

							if *verbose || *dryrun {
								fmt.Printf("\nCreating directory %s\n", targetDir)
							}
							if !*dryrun {
								os.Mkdir(targetDir, os.ModePerm) // make the new dir if dryrun is false
							}
							createdDirs[targetDir] = true
						}
					}

					oldPath := filepath.Join(*dir, file.Name())
					newPath := filepath.Join(targetDir, file.Name())
					if *verbose || *dryrun {
						fmt.Printf("Moving file: %s -> %s \n", oldPath, newPath)
					}
					if !*dryrun {
						os.Rename(oldPath, newPath)
					}
				}
			}
		}

		fmt.Printf("\nFiles organized by extension in %s \n", *dir)
	}

	if *cleanAll {
		cleanEmptyDirs(*dir, *dryrun)
	}
}

func cleanEmptyDirs(dir string, dryrun bool) {
	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	for _, file := range files {
		if file.IsDir() {
			subDir := filepath.Join(dir, file.Name())
			cleanEmptyDirs(subDir, dryrun)
			subFiles, err := os.ReadDir(subDir)
			if err != nil {
				fmt.Println("Error reading subdirectory:", err)
				continue
			}
			if len(subFiles) == 0 {
				fmt.Printf("Deleting empty directory: %s\n", subDir)
				if !dryrun {
					if err := os.Remove(subDir); err != nil {
						fmt.Println("Error deleting directory:", err)
					}
				}
			}
		}
	}
}
