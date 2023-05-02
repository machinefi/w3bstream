package config

var (
	SupportedLanguage = []Language{"English", "中文"}
)

var (
	English Language = SupportedLanguage[0]
	Chinese Language = SupportedLanguage[1]
)

// Multi-language support
type Language string

// Config defines the config schema
type Config struct {
	Endpoint string   `json:"endpoint" yaml:"endpoint"`
	Language Language `json:"language" yaml:"language"`
}
