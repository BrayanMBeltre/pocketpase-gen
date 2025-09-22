package commands

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/brayanmbeltre/pocketpase-gen/internal/generator"
	"github.com/brayanmbeltre/pocketpase-gen/internal/pocketbase"
	"github.com/spf13/cobra"
)

var (
	dbPath        string
	packageName   string
	packageFolder string
	verbose       bool
)

// GenerateCmd represents the generate command
var GenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a Go collections schema file from a PocketBase SQLite DB",
	RunE:  runGenerate,
}

func init() {
	GenerateCmd.Flags().StringVarP(&dbPath, "db", "d", "./pb_data/data.db", "Path to the PocketBase SQLite database file")
	GenerateCmd.Flags().StringVarP(&packageFolder, "output-dir", "o", "collections", "Directory path for the generated Go file")
	GenerateCmd.Flags().StringVarP(&packageName, "package-name", "p", "collections", "Go package name for the generated file")
	GenerateCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")
}

func runGenerate(cmd *cobra.Command, args []string) error {
	if verbose {
		// Configure slog for verbose output if needed
		// opts := &slog.HandlerOptions{Level: slog.LevelDebug}
		// handler := slog.NewTextHandler(os.Stderr, opts)
		// slog.SetDefault(slog.New(handler))
		slog.Info("Starting schema generation", "db", dbPath, "outputDir", packageFolder, "packageName", packageName)
	}

	// --- Get Collections ---
	collections, err := pocketbase.GetCollections(dbPath, verbose)
	if err != nil {
		return fmt.Errorf("failed to get collections from PocketBase: %w", err)
	}

	if len(collections) == 0 {
		slog.Warn("No collections found in the database.")
		return nil // Nothing to generate
	}

	// --- Generate Code String ---
	slog.Info("Generating Go code...")
	goCode, err := generator.GenerateCollectionSchemaFileContent(packageName, collections)
	if err != nil {
		return fmt.Errorf("failed to generate Go code: %w", err)
	}

	// --- Write File ---
	// Ensure the output directory exists
	err = os.MkdirAll(packageFolder, 0755)
	if err != nil {
		return fmt.Errorf("failed to create output directory %s: %w", packageFolder, err)
	}

	// Define the output file path
	// Using package name for the file is common, e.g., collections/collections.go
	outputFilePath := filepath.Join(packageFolder, fmt.Sprintf("%s.go", packageName))

	slog.Info("Writing generated code", "file", outputFilePath)
	err = os.WriteFile(outputFilePath, []byte(goCode), 0644)
	if err != nil {
		return fmt.Errorf("failed to write generated file %s: %w", outputFilePath, err)
	}

	slog.Info("Collections schema generation complete.")
	return nil
}
