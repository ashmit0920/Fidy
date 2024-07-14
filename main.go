package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	name := flag.String("name", "World", "Name to greet")
	age := flag.Int("age", 0, "Age of the person")
	flag.Parse()

	if *age <= 0 {
		fmt.Println("Please provide a valid age.")
		os.Exit(1)
	}

	fmt.Printf("Hello, %s! You are %d years old.\n", *name, *age)
}
