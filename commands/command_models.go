package commands

import (
	"fmt"
	"log/slog"
	"os"
	"path"
	"strings"

	"github.com/alexisvisco/pocketpase-gen/codegen"
	"github.com/spf13/cobra"
	"github.com/stoewer/go-strcase"
)

var ModelsCommand = &cobra.Command{
	Use:   "models",
	Short: "Generate models from a pocketbase instance sqlite file",
	RunE:  run,
}

var (
	packageName   string
	packageFolder string
)

func init() {
	ModelsCommand.PersistentFlags().StringVar(&packageFolder, "package-folder", "collections", "path to package folder")
	ModelsCommand.PersistentFlags().StringVar(&packageName, "package-name", "collections", "name of the package")
}

func run(cmd *cobra.Command, args []string) error {
	collections, err := GetCollections()

	if err != nil {
		return fmt.Errorf("failed to open app: %w", err)
	}

	if Verbose {
		slog.Info("found collections", slog.Int("count", len(collections)))

		for _, collection := range collections {
			slog.Info("collection", slog.String("name", collection.Name))
		}
	}

	err = os.MkdirAll(packageFolder, 0755)
	if err != nil {
		return fmt.Errorf("failed to create package folder: %w", err)
	}

	parser, err := codegen.CreateModelTemplateParser()
	if err != nil {
		return fmt.Errorf("failed to create model template parser: %w", err)
	}

	modelBuilders := make([]*codegen.ModelBuilder, 0, len(collections))
	collectionNames := make([]string, 0, len(collections))
	for _, collection := range collections {
		schema, err := codegen.ModelBuilderFromSchema(packageName, collection.Name, &collection.Fields)
		if err != nil {
			return fmt.Errorf("failed to create model builder for collection %s from schema: %w", collection.Name, err)
		}

		modelBuilders = append(modelBuilders, schema)
		collectionNames = append(collectionNames, collection.Name)
	}

	files := map[string]string{}
	for _, modelBuilder := range modelBuilders {
		file, err := modelBuilder.Gen(parser)
		if err != nil {
			return fmt.Errorf("failed to get string from model builder: %w", err)
		}

		files[fmt.Sprintf("%s.go", strcase.SnakeCase(modelBuilder.ModelName))] = file
	}

	// Generate collections.go
	collectionsFileContent := generateCollectionsFile(packageName, collectionNames)
	files["collections.go"] = collectionsFileContent

	for file, content := range files {
		if err := os.WriteFile(path.Join(packageFolder, file), []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", file, err)
		}
	}

	return nil
}

func generateCollectionsFile(packageName string, collectionNames []string) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("package %s\n\n", packageName))
	builder.WriteString("const (\n")
	for _, name := range collectionNames {
		builder.WriteString(fmt.Sprintf("\tCollection%s string = \"%s\"\n", strcase.UpperCamelCase(name), name))
	}
	builder.WriteString(")\n")
	return builder.String()
}
