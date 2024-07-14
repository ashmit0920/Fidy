package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	Name string `json:"name"`
}

func main() {

	// info := flag.String("info", "", "Get information about fidy")
	name := flag.String("name", "", "The user's name")
	dir := flag.String("dir", "", "The directory to organize")
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

	if *dir != "" {

		files, err := os.ReadDir(*dir)
		if err != nil {
			fmt.Println("Error reading directory: ", err)
			os.Exit(1)
		}

		for _, file := range files {
			if !file.IsDir() {
				ext := filepath.Ext(file.Name())
				if ext != "" {
					ext = ext[1:] // Removing the dot
					targetDir := filepath.Join(*dir, ext)
					if _, err := os.Stat(targetDir); os.IsNotExist(err) {
						os.Mkdir(targetDir, os.ModePerm)
					}

					oldPath := filepath.Join(*dir, file.Name())
					newPath := filepath.Join(targetDir, file.Name())
					os.Rename(oldPath, newPath)
				}
			}
		}

		fmt.Println("Files organized by extension in", *dir)
	}
}
