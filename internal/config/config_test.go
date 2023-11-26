package config

import (
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestLoadYML(t *testing.T) {
	config := []byte(`
port: 8080 
loglevel: debug
transformations:
    - name: testingtransformation1
      operation: extract
      mappings:
        - path: user.data.name
          path_out: name
          type_out: string
          default_null: there is a null name
forwarders:
  - name: sendtosns
    kind: sns
    params:
      aws_access_key_id: foo
`)

	f, err := os.CreateTemp("", "config.yml")
	assert.Nil(t, err, "cannot create config temp file")
	_, err = f.Write(config)
	assert.Nil(t, err, "cannot write config temp file")
	f.Close()

	ac, err := LoadFromFile(f.Name())
	if err != nil {
		t.Error(err)
	}
	expected := &App{
		Port:     "8080",
		LogLevel: zerolog.DebugLevel,
		Transformations: map[string]Transformation{
			"testingtransformation1": {
				Operation: OperationExtract,
				Mappings: []Mapping{
					{
						Path:        "user.data.name",
						PathOut:     "name",
						TypeOut:     "string",
						DefaultNull: "there is a null name"}},
			}},
		Forwarders: map[string]Forwarder{
			"sendtosns": {
				Kind: KindSNS,
				Settings: map[ForwarderSetting]interface{}{
					SettingAWSKey: "foo",
				},
			},
		},
	}

	assert.Equal(t, *ac, *expected)
}
