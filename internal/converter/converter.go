package converter

// Converter is an interface for converting DTOs to a specific language.
type Converter interface {
	Convert(inputFilePath string) error
	ConvertDirectory(inputDir string) error
}
