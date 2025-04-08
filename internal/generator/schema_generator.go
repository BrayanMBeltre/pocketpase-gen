package generator

import (
	"bytes"
	_ "embed"
	"fmt"
	"text/template"

	"github.com/alexisvisco/pocketpase-gen/internal/pocketbase"
	"github.com/stoewer/go-strcase"
)

//go:embed schema_template.tmpl
var schemaTemplate string

type CollectionData struct {
	PackageName string
	Collections []CollectionTemplateData
}

type CollectionTemplateData struct {
	StructName     string
	CollectionName string
	Fields         []FieldTemplateData
}

type FieldTemplateData struct {
	GoFieldName string
	FieldName   string
}

func GenerateCollectionSchemaFileContent(packageName string, collections []pocketbase.CollectionSchema) (string, error) {
	data := CollectionData{
		PackageName: strcase.SnakeCase(packageName),
	}

	for _, collection := range collections {
		collectionData := CollectionTemplateData{
			StructName:     strcase.UpperCamelCase(collection.Name),
			CollectionName: collection.Name,
		}

		for _, fieldName := range collection.Fields {
			collectionData.Fields = append(collectionData.Fields, FieldTemplateData{
				GoFieldName: strcase.UpperCamelCase(fieldName),
				FieldName:   fieldName,
			})
		}

		data.Collections = append(data.Collections, collectionData)
	}

	tmpl, err := template.New("schema_template.tmpl").Parse(schemaTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to load template: %w", err)
	}

	// Execute the template with the data
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}
