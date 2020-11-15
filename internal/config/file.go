package config

// FileConfig represents the YAML config file
type FileConfig struct {
	ServerPort          string               `yaml:"port"`
	LogLevel            string               `yaml:"loglevel"`
	ForwarderConfigs    []ForwarderConfig    `yaml:"forwarders"`
	FileTransformations []FileTransformation `yaml:"transformations"`
}

// ForwarderConfig holds the forwarder configuration
type ForwarderConfig struct {
	Name   string                 `yaml:"name"`
	Kind   ForwarderKind          `yaml:"kind"`
	Params map[string]interface{} `yaml:"params"`
}

// FileTransformation defines a transformation specified in the yaml config file
type FileTransformation struct {
	Name      string    `yaml:"name"`
	Operation string    `yaml:"operation"`
	Mappings  []Mapping `yaml:"mappings"`
}

// Mapping defines a mapping
type Mapping struct {
	Path        string `yaml:"path"`
	PathOut     string `yaml:"path_out"`
	TypeOut     string `yaml:"type_out"`
	DefaultNull string `yaml:"default_null"`
}
