package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	// name := flag.String("name", "World", "Name to greet")
	// age := flag.Int("age", 0, "Age of the person")
	// flag.Parse()

	// if *age <= 0 {
	// 	fmt.Println("Please provide a valid age.")
	// 	os.Exit(1)
	// }

	// fmt.Printf("Hello, %s! You are %d years old.\n", *name, *age)

	dir := flag.String("dir", ".", "The directory to organize")
	flag.Parse()

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
