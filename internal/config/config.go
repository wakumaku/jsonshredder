package config

import (
	"fmt"
	"io/ioutil"

	"github.com/rs/zerolog"
	"gopkg.in/yaml.v3"
)

// LoadFromFile loads a config file and returns an Application config
func LoadFromFile(path string) (*App, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %s", err)
	}

	cfg, err := load(b)
	if err != nil {
		return nil, fmt.Errorf("loading config: %s", err)
	}

	return cfg, nil
}

// load creates a config from a yaml file structure
func load(content []byte) (*App, error) {
	fc := &FileConfig{}
	if err := yaml.Unmarshal(content, fc); err != nil {
		return nil, fmt.Errorf("decoding config: %s", err)
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logLevel, err := zerolog.ParseLevel(fc.LogLevel)
	if err != nil {
		logLevel = zerolog.InfoLevel
	}

	ac := &App{
		LogLevel:        logLevel,
		Port:            fc.ServerPort,
		Transformations: make(map[string]Transformation, len(fc.FileTransformations)),
		Forwarders:      make(map[string]Forwarder, len(fc.ForwarderConfigs)),
	}

	for _, p := range fc.ForwarderConfigs {
		ac.Forwarders[p.Name] = Forwarder{
			Kind:     p.Kind,
			Settings: make(map[ForwarderSetting]interface{}, len(p.Params)),
		}
		for k, v := range p.Params {
			ac.Forwarders[p.Name].Settings[ForwarderSetting(k)] = v
		}
	}

	for _, t := range fc.FileTransformations {
		ac.Transformations[t.Name] = Transformation{
			Operation: Operation(t.Operation),
			Mappings:  t.Mappings,
		}
	}

	return ac, nil
}
