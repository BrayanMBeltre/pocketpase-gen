package commands

import (
	"fmt"
	"os"
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	err := GenerateCmd.Execute()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing command: %v\n", err)
	}

	return err // Return the error for main.go to handle exit codes
}
