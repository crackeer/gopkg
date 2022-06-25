package config

const (

	// EnvDefault
	EnvDefault = "default"
)

// YamlAPIConfig
type YamlAPI struct {
	BaseURI     map[string]string       `yaml:"base_uri"`
	Timeout     int64                   `yaml:"timeout"`
	SignAlg     string                  `yaml:"sign_alg"`
	SuccessCode string                  `yaml:"success_code"`
	CodeKey     string                  `yaml:"code_key"`
	DataKey     string                  `yaml:"data_key"`
	MessageKey  string                  `yaml:"message_key"`
	Header      map[string]string       `yaml:"header"`
	API         map[string]*YamlAPIItem `yaml:"api"`
}

// YamlAPIItem
type YamlAPIItem struct {
	Path        string            `yaml:"path"`
	ContentType string            `yaml:"content_type"`
	Method      string            `yaml:"method"`
	Header      map[string]string `yaml:"header"`
	Timeout     int64             `yaml:"timeout"`
}
