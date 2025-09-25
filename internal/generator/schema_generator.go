package generator

import (
	"bytes"
	_ "embed"
	"fmt"
	"strings"
	"text/template"

	"github.com/brayanmbeltre/pocketpase-gen/internal/pocketbase"
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
	Options     []string
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

		for _, field := range collection.Fields {
			collectionData.Fields = append(collectionData.Fields, FieldTemplateData{
				GoFieldName: strcase.UpperCamelCase(field.Name),
				FieldName:   field.Name,
				Options:     field.Values,
			})
		}

		data.Collections = append(data.Collections, collectionData)
	}

	funcMap := template.FuncMap{
		"ToUpper":   strings.ToUpper,
		"SnakeCase": strcase.SnakeCase,
	}

	tmpl, err := template.New("schema_template.tmpl").Funcs(funcMap).Parse(schemaTemplate)
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
