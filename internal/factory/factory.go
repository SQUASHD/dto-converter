package factory

import (
	"fmt"
	"github.com/SQUASHD/dto-converter/internal/converter"
	"github.com/SQUASHD/dto-converter/internal/structs"
)

var SupportedLanguages = []string{"csharp"}

func ConverterFactory(project *structs.Project) (converter.Converter, error) {
	tsWriter := converter.NewTSWriter(project.OutPath)
	switch project.Language {
	case "csharp":
		return converter.NewCSharpConverter(tsWriter), nil
	// TODO: Add more languages here
	default:
		return nil, fmt.Errorf("unsupported language: %s", project.Language)
	}
}
