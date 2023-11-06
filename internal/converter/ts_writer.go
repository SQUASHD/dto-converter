package converter

import (
	"os"
	"path/filepath"
	"strings"
)

// TypeScriptType represents a TypeScript type definition.
type TypeScriptType struct {
	Name   string
	Fields []string
}

// TSWriter handles writing TypeScript type definitions to files.
type TSWriter struct {
	outPath string
}

// NewTSWriter creates a new TSWriter instance.
func NewTSWriter(outPath string) *TSWriter {
	return &TSWriter{outPath: outPath}
}

// WriteType writes a TypeScript type definition to a file.
func (w *TSWriter) WriteType(tsType TypeScriptType, filePath string) error {
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return err
	}

	var sb strings.Builder
	sb.WriteString("export type " + tsType.Name + " = {\n")
	for _, field := range tsType.Fields {
		sb.WriteString("  " + field + "\n")
	}
	sb.WriteString("}\n")

	return os.WriteFile(filePath, []byte(sb.String()), 0644)
}

// WriteAllTypes writes all TypeScript type definitions to a single file.
func (w *TSWriter) WriteAllTypes(tsTypes []TypeScriptType, filePath string) error {
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return err
	}

	var sb strings.Builder

	for _, tsType := range tsTypes {
		sb.WriteString("export type " + tsType.Name + " = {\n")
		for _, field := range tsType.Fields {
			sb.WriteString("  " + w.toCamelCase(field) + "\n")
		}
		sb.WriteString("}\n\n")
	}

	return os.WriteFile(filePath, []byte(sb.String()), 0644)
}

// GenerateTsFilePath generates a TypeScript file path based on the input file path.
func (w *TSWriter) GenerateTsFilePath(inPath string) string {
	baseFileName := strings.TrimSuffix(filepath.Base(inPath), filepath.Ext(inPath)) + ".ts"
	return filepath.Join(w.outPath, baseFileName)
}

// ToCamelCase converts a string to camel case.
func (w *TSWriter) toCamelCase(input string) string {
	if len(input) == 0 {
		return ""
	}
	return strings.ToLower(input[:1]) + input[1:]
}
