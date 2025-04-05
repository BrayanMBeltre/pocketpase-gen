package codegen

import (
	"bytes"
	_ "embed"
	"fmt"
	"text/template"

	"github.com/pocketbase/pocketbase/core"
	"github.com/stoewer/go-strcase"
)

//go:embed model.tmpl
var templateModel string

type ModelBuilder struct {
	PackageName string

	Imports Imports

	ModelName      string
	CollectionName string

	Fields []Field
}

func ModelBuilderFromSchema(pkgName string, collectionName string, s *core.FieldsList) (*ModelBuilder, error) {
	mb := &ModelBuilder{
		PackageName:    strcase.SnakeCase(pkgName),
		ModelName:      strcase.UpperCamelCase(collectionName),
		CollectionName: collectionName,
		Imports:        Imports{},
	}

	mb.Imports.addImport("github.com/pocketbase/pocketbase/core", "")

	for _, field := range s.AsMap() {

		mbField := Field{
			FieldName:    field.GetName(),
			FunctionName: strcase.UpperCamelCase(field.GetName()),
		}

		if err := resolveSchemaField(mb, field, &mbField); err != nil {
			return nil, err
		}

		mb.Fields = append(mb.Fields, mbField)
	}

	return mb, nil
}

func (mb *ModelBuilder) Gen(t *template.Template) (string, error) {
	buffer := bytes.NewBuffer(nil)
	err := t.Execute(buffer, mb)

	if err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buffer.String(), nil
}

type Imports map[string]Import

func (i Imports) addImport(importPath string, alias string) {
	i[importPath] = Import{
		Alias: alias,
		Path:  importPath,
	}
}

func (i Imports) List() []Import {
	var imports []Import
	for _, v := range i {
		imports = append(imports, v)
	}

	return imports
}

type Import struct {
	Alias string
	Path  string
}

func (i Import) ToImportPath() string {
	if i.Alias != "" {
		return fmt.Sprintf(`%s "%s"`, i.Alias, i.Path)
	} else {
		return fmt.Sprintf(`"%s"`, i.Path)
	}
}

type Field struct {
	FieldName    string
	FunctionName string

	GoType string

	GetterComment string
	GetterCall    string
	HasGetterCast bool

	SetterComment string
}

func resolveSchemaField(builder *ModelBuilder, field core.Field, f *Field) error {
	f.GetterComment = fmt.Sprintf(`// Get%s returns the value of the "%s" field`, f.FunctionName, field.GetName())
	f.SetterComment = fmt.Sprintf(`// Set%s sets the value of the "%s" field`, f.FunctionName, field.GetName())

	switch field.Type() {
	case core.FieldTypeText:
		f.GoType = "string"
		f.GetterCall = `GetString`

	case core.FieldTypeEditor:
		f.GoType = "string"
		f.GetterCall = `GetString`
		f.GetterComment = fmt.Sprintf(`// Get%s returns the value of the "%s" field as HTML`, f.FunctionName,
			field.GetName())
		f.SetterComment = fmt.Sprintf(`// Set%s sets the value of the "%s" field as HTML`, f.FunctionName,
			field.GetName())

	case core.FieldTypeURL:
		f.GoType = "string"
		f.GetterCall = `GetString`
		f.GetterComment = fmt.Sprintf(`// Get%s returns the value of the "%s" field as URL`, f.FunctionName, field.GetName())
		f.SetterComment = fmt.Sprintf(`// Set%s sets the value of the "%s" field as URL`, f.FunctionName, field.GetName())

	case core.FieldTypeEmail:
		f.GoType = "string"
		f.GetterCall = `GetString`
		f.GetterComment = fmt.Sprintf(`// Get%s returns the value of the "%s" field as email`, f.FunctionName,
			field.GetName())
		f.SetterComment = fmt.Sprintf(`// Set%s sets the value of the "%s" field as email`, f.FunctionName, field.GetName())
	case core.FieldTypeFile:
		IsMultiple := field.(*core.FileField).IsMultiple()
		f.GetterComment = fmt.Sprintf(`// Get%s returns the value of the "%s" field as file`, f.FunctionName,
			field.GetName())
		f.SetterComment = fmt.Sprintf(`// Set%s sets the value of the "%s" field as file`, f.FunctionName, field.GetName())
		if IsMultiple {
			f.GoType = "[]string"
			f.GetterCall = `GetStringSlice`
		} else {
			f.GoType = "string"
			f.GetterCall = `GetString`
		}

	case core.FieldTypeNumber:
		f.GoType = "int"
		f.GetterCall = `GetInt`

	case core.FieldTypeBool:
		f.GoType = "bool"
		f.GetterCall = `GetBool`

	case core.FieldTypePassword:
		f.GoType = "string"
		f.GetterCall = `GetString`

	case core.FieldTypeDate:
		builder.Imports.addImport("github.com/pocketbase/pocketbase/tools/types", "")
		f.GoType = "types.DateTime"
		f.GetterCall = `GetDateTime`

	case core.FieldTypeAutodate:
		builder.Imports.addImport("github.com/pocketbase/pocketbase/tools/types", "")
		f.GoType = "types.DateTime"
		f.GetterCall = `GetDateTime`

	case core.FieldTypeSelect:
		IsMultiple := field.(*core.SelectField).IsMultiple()
		values := field.(*core.SelectField).Values
		f.GetterComment += fmt.Sprintf("\n// Possible values: %s", values)
		f.SetterComment += fmt.Sprintf("\n// Possible values: %s", values)

		if IsMultiple {
			f.GoType = "[]string"
			f.GetterCall = `GetStringSlice`
		} else {
			f.GoType = "string"
			f.GetterCall = `GetString`
		}

	case core.FieldTypeJSON:
		builder.Imports.addImport("github.com/pocketbase/pocketbase/tools/types", "")
		f.GoType = "types.JSONRaw"
		f.GetterCall = `Get`
		f.HasGetterCast = true

	case core.FieldTypeRelation:
		IsMultiple := field.(*core.RelationField).IsMultiple()
		collectionId := field.(*core.RelationField).CollectionId
		f.GetterComment += fmt.Sprintf("\n// Relation collection related : %s",
			collectionId /* TODO: resolve it name */)
		f.SetterComment += fmt.Sprintf("\n// Relation collection related : %s",
			collectionId /* TODO: resolve it name */)
		if IsMultiple {
			f.GoType = "[]string"
			f.GetterCall = `GetStringSlice`
		} else {
			f.GoType = "string"
			f.GetterCall = `GetString`
		}

	default:
		return fmt.Errorf("unknown field type %s", field.Type())
	}

	return nil
}
