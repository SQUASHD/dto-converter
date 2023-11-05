package structs

// Config holds the entire configuration with multiple conversion aliases
type Config struct {
	Projects []Project `json:"conversionAliases"`
}

type Project struct {
	Name     string `json:"alias"`
	Language string `json:"language"`
	InPath   string `json:"inPath"`
	OutPath  string `json:"outPath"`
}
