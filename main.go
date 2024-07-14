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

func main() {

	name := flag.String("name", "", "The user's name")
	dir := flag.String("dir", "", "The directory to organize")
	help := flag.Bool("help", false, "Show information about fidy")
	include := flag.String("include", "", "Comma-separated list of extensions to include")
	exclude := flag.String("exclude", "", "Comma-separated list of extensions to exclude")
	verbose := flag.Bool("verbose", false, "Enable verbose output")
	dryrun := flag.Bool("dryrun", false, "Simulate the file organization without doing any actual changes")

	flag.Parse()

	configFile := "config.json"
	var config Config

	// Read existing configuration if available
	if _, err := os.Stat(configFile); err == nil {
		data, err := os.ReadFile(configFile)
		if err == nil {
			json.Unmarshal(data, &config)
		}
	}

	// Update the config file with name (if provided)
	if *name != "" {
		config.Name = *name
		data, _ := json.Marshal(config)
		os.WriteFile(configFile, data, os.ModePerm)
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
		fmt.Println("\n---------- Fidy - The File Organizer CLI Tool ----------")
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

		fmt.Println("\nFiles organized by extension in", *dir)
	}
}
