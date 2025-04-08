package main

import (
	"os"

	"github.com/alexisvisco/pocketpase-gen/cmd/pb-gen/commands" // Import commands package
)

func main() {
	if err := commands.Execute(); err != nil {
		os.Exit(1)
	}
}
