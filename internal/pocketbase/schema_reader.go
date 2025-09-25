package pocketbase

import (
	"database/sql"
	"fmt"
	"log/slog"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

type Field struct {
	Name   string
	Values []string
}

type CollectionSchema struct {
	Name   string
	Fields []Field
}

func GetCollections(dbPath string, verbose bool) ([]CollectionSchema, error) {
	db, err := sql.Open("sqlite3", dbPath)

	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite3 database: %w", err)
	}

	open := dbx.NewFromDB(db, "sqlite3")

	var collections []core.Collection

	err = open.NewQuery("SELECT * FROM _collections").All(&collections)

	if err != nil {
		return nil, fmt.Errorf("failed to get collections: %w", err)
	}

	if verbose {
		slog.Info("found collections", slog.Int("count", len(collections)))
	}

	result := make([]CollectionSchema, 0, len(collections))
	for _, collection := range collections {
		if verbose {
			slog.Info("processing collection", slog.String("name", collection.Name), slog.String("id", collection.Id))
		}

		schema := CollectionSchema{
			Name:   collection.Name,
			Fields: make([]Field, 0, len(collection.Fields)),
		}

		for _, field := range collection.Fields {
			fieldName := field.GetName()
			var fieldValues []string

			if selectField, ok := field.(*core.SelectField); ok {
				fieldValues = selectField.Values
			}

			schema.Fields = append(schema.Fields, Field{
				Name:   fieldName,
				Values: slugifyFieldValues(fieldValues),
			})
		}

		result = append(result, schema)
	}

	return result, nil
}

func slugifyFieldValues(values []string) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		result = append(result, strings.Replace(value, " ", "_", -1))
	}
	return result
}
