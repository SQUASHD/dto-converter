// internal/converter/csharp_converter.go

package converter

import (
	"bufio"
	"dto-converter/internal/utils"
	"os"
	"strings"
)

// CSharpConverter is a converter for C# DTOs.
type CSharpConverter struct {
	writer *TSWriter
}

// NewCSharpConverter creates a new CSharpConverter instance.
func NewCSharpConverter(writer *TSWriter) Converter {
	return &CSharpConverter{
		writer: writer,
	}
}

// Convert converts a C# DTO File to a TypeScript file.
func (c *CSharpConverter) Convert(inputFilePath string) error {
	file, err := os.Open(inputFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var tsTypes []TypeScriptType
	var currentType *TypeScriptType

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "public record") {
			typeName := extractTypeName(line)
			currentType = &TypeScriptType{Name: typeName, Fields: []string{}}
			continue
		}

		if currentType != nil {
			// If the line ends with ");", we've reached the end of the type definition.
			if strings.HasSuffix(line, ");") {
				tsTypes = append(tsTypes, *currentType)
				currentType = nil
			} else {
				// Since we check for end of type definition we can assume we are in a field definition
				// therefore we can remove the trailing comma.
				line = strings.TrimSuffix(line, ",")
				fieldName, tsType, attribute, comment := parseField(line)
				formattedField := formatField(fieldName, tsType, attribute, comment)
				currentType.Fields = append(currentType.Fields, formattedField)
			}
		}
	}

	if len(tsTypes) > 0 {
		tsFilePath := c.writer.GenerateTsFilePath(inputFilePath)
		if err := c.writer.WriteAllTypes(tsTypes, tsFilePath); err != nil {
			return err
		}
	}
	return nil
}

func extractTypeName(line string) string {
	parts := strings.Fields(line)
	return strings.TrimSuffix(parts[2], "(")
}

func (c *CSharpConverter) ConvertDirectory(inputDir string) error {
	return utils.WalkAndProcess(inputDir, ".cs", func(path string, info os.FileInfo) error {
		return c.Convert(path)
	})
}

func parseField(line string) (fieldName, tsType, attribute, comment string) {
	parts := strings.Fields(line)
	if len(parts) < 2 {
		return "", "", "", ""
	}

	fieldName = parts[len(parts)-1]
	tsType, comment = csharpTypeToTypeScript(parts[len(parts)-2])
	if len(parts) > 2 && strings.HasPrefix(parts[0], "[") {
		attribute = parts[0]
	}
	fieldName = strings.TrimRight(fieldName, ",;")
	return fieldName, tsType, attribute, comment
}

func formatField(fieldName, tsType, attribute, comment string) string {
	var formattedAttribute string
	if attribute != "" {
		formattedAttribute = " // " + attribute
	}
	return fieldName + ": " + tsType + ";" + comment + formattedAttribute
}

func csharpTypeToTypeScript(csharpType string) (string, string) {
	comment := ""
	tsType := ""

	if strings.HasSuffix(csharpType, "?") {
		nonNullableType := strings.TrimSuffix(csharpType, "?")
		tsType, comment = csharpTypeToTypeScript(nonNullableType)
		tsType += " | null"
		return tsType, comment
	}

	switch csharpType {
	case "int", "long", "short", "byte", "float", "double", "decimal":
		tsType = "number"
	case "string":
		tsType = "string"
	case "bool":
		tsType = "boolean"
	case "DateTime", "DateTimeOffset":
		tsType = "string"
		comment = " // Represents a DateTime type in string format"
	case "Guid":
		tsType = "string"
		comment = " // Represents a GUID/UUID in string format"
	case "char":
		tsType = "string"
	case "object":
		tsType = "any"
	case "dynamic":
		tsType = "any"
		comment = " // Dynamic type in C#, any type in TypeScript"
	// Example of a custom type
	case "PublicationStatus":
		tsType = "PublicationStatus"
		comment = " // Known enum type"
	default:
		tsType = csharpType // TODO: Normalize type name
		comment = " // Unknown type – could be an enum, or a custom type"
	}

	return tsType, comment
}
