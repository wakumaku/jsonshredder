package shredder

import (
	"testing"
	"wakumaku/jsonshredder/internal/config"

	"github.com/stretchr/testify/assert"
)

func TestShred(t *testing.T) {

	mappings := []config.Mapping{
		{
			Path:        "user.data.age",
			PathOut:     "aaa.bbb.ccc.ddd",
			TypeOut:     "int",
			DefaultNull: "nulli",
		},
		{
			Path:        "user.data.weight",
			PathOut:     "username.weight",
			TypeOut:     "float",
			DefaultNull: "nulli",
		},
		{
			Path:        "user.data.name",
			PathOut:     "1111.2222.3333.4444",
			TypeOut:     "string",
			DefaultNull: "nulli",
		},
		{
			Path:        "user.data.name",
			PathOut:     "1111.2222.3333.555",
			TypeOut:     "string",
			DefaultNull: "nulli",
		},
		{
			Path:        "user.data.name",
			PathOut:     "1111.2222.3333.xxx",
			TypeOut:     "string",
			DefaultNull: "nulli"},
	}

	jsonDoc := []byte(`{
    "user": {
        "data": {
            "name": "john",
            "age": "2",
            "weight": "93.504"
        }
    }
}`)

	out, err := Shred(config.Transformation{Mappings: mappings}, jsonDoc)
	assert.NoError(t, err)

	expectedOut := []byte(`{"1111":{"2222":{"3333":{"4444":"john","555":"john","xxx":"john"}}},"aaa":{"bbb":{"ccc":{"ddd":2}}},"username":{"weight":93.504}}`)

	assert.Equal(t, out, expectedOut)
}

func BenchmarkShredShort(b *testing.B) {

	mappings := []config.Mapping{
		{
			Path:        "user.data.age",
			PathOut:     "aaa.bbb.ccc.ddd",
			TypeOut:     "int",
			DefaultNull: "nulli",
		},
	}

	jsonDoc := []byte(`{
    "user": {
        "data": {
            "name": "john",
            "age": "2",
            "weight": "93.504"
        }
    }
}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Shred(config.Transformation{Mappings: mappings}, jsonDoc)
	}
}

func BenchmarkShredLong(b *testing.B) {

	mappings := []config.Mapping{
		{
			Path:        "user.data.age",
			PathOut:     "aaa.bbb.ccc.ddd",
			TypeOut:     "int",
			DefaultNull: "nulli"},
		{
			Path:        "user.data.weight",
			PathOut:     "username.weight",
			TypeOut:     "float",
			DefaultNull: "nulli",
		},
		{
			Path:        "user.data.name",
			PathOut:     "1111.2222.3333.4444",
			TypeOut:     "string",
			DefaultNull: "nulli",
		},
		{
			Path:        "user.data.name",
			PathOut:     "1111.2222.3333.555",
			TypeOut:     "string",
			DefaultNull: "nulli",
		},
		{
			Path:        "user.data.name",
			PathOut:     "1111.2222.3333.xxx1",
			TypeOut:     "string",
			DefaultNull: "nulli",
		},
		{
			Path:        "user.data.name",
			PathOut:     "1111.2222.3333.xxx2",
			TypeOut:     "string",
			DefaultNull: "nulli"},
		{
			Path:        "user.data.name",
			PathOut:     "1111.2222.3333.xxx3",
			TypeOut:     "string",
			DefaultNull: "nulli",
		},
		{
			Path:        "user.data.name",
			PathOut:     "1111.2222.3333.xxx4",
			TypeOut:     "string",
			DefaultNull: "nulli",
		},
		{
			Path:        "user.data.name",
			PathOut:     "1111.2222.3333.xxx5",
			TypeOut:     "string",
			DefaultNull: "nulli",
		},
	}

	jsonDoc := []byte(`{
    "user": {
        "data": {
            "name": "john",
            "age": "2",
            "weight": "93.504"
        }
    }
}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Shred(config.Transformation{Mappings: mappings}, jsonDoc)
	}
}

func BenchmarkDeepSetLong(b *testing.B) {

	in := map[string]interface{}{}
	keys := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = deepSet(in, "val", keys...)
	}
}

func BenchmarkDeepSetShort(b *testing.B) {

	in := map[string]interface{}{}
	keys := []string{"a", "b"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = deepSet(in, "val", keys...)
	}
}

func BenchmarkSetTypeOutToString(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = setTypeOut(1, "string", "")
	}
}

func BenchmarkSetTypeOutToInt(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = setTypeOut("1", "int", "")
	}
}

func BenchmarkSetTypeOutToFloat(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = setTypeOut("1", "float", "")
	}
}
