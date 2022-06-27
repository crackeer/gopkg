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

// YamlGroupPage
type YamlGroupPage struct {
	Tag  string               `yaml:"tag"`
	List map[string]*YamlPage `yaml:"list"`
}

// YamlAPIConfig
type YamlPage struct {
	DataAPI       string                 `yaml:"data_api"`
	Type          string                 `yaml:"type"`
	Extension     map[string]interface{} `yaml:"extension"`
	DefaultParams map[string]interface{} `yaml:"default_params"`
	Title         string                 `yaml:"title"`
	HrefTitle     string                 `yaml:"href_title"`

	FrameFile   string `yaml:"frame_file"`
	ContentFile string `yaml:"content_file"`
}

// DBConfig
type DBConfig struct {
	Driver        string `yaml:"driver"`
	File          string `yaml:"file"`
	ReadHost      string `yaml:"read_host"`
	WriteHost     string `yaml:"write_host"`
	WriteUser     string `yaml:"write_user"`
	WritePassword string `yaml:"writer_password"`
	ReadUser      string `yaml:"read_user"`
	ReadPassword  string `yaml:"read_password"`
	Database      string `yaml:"database"`
	Charset       string `yaml:"charset"`
}
