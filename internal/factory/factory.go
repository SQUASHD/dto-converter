package factory

import (
	"dto-converter/internal/converter"
	"dto-converter/internal/structs"
	"fmt"
)

var SupportedLanguages = []string{"csharp"}

func ConverterFactory(project structs.Project) (converter.Converter, error) {
	tsWriter := converter.NewTSWriter(project.OutPath)
	switch project.Language {
	case "csharp":
		return converter.NewCSharpConverter(tsWriter), nil
	// TODO: Add more languages here
	default:
		return nil, fmt.Errorf("unsupported language: %s", project.Language)
	}
}
