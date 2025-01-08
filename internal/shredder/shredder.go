package shredder

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/wakumaku/jsonshredder/internal/config"

	"github.com/jmespath/go-jmespath"
)

// Errors when shredding
var (
	ErrInvalidJSONInput  = errors.New("invalid json input")
	ErrInvalidJSONOutput = errors.New("invalid json output")
	ErrJSONSearchValue   = errors.New("searching value")
)

// Shred processes an input based on a transformation config
func Shred(transformation config.Transformation, in []byte) ([]byte, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(in, &data); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidJSONInput, err)
	}

	mid := map[string]interface{}{}
	if transformation.Operation == config.OperationAdd {
		mid = data
	}

	for _, m := range transformation.Mappings {
		r, err := jmespath.Search(m.Path, data)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", ErrJSONSearchValue, err)
		}

		keyOut := m.Path
		if m.PathOut != "" {
			keyOut = m.PathOut
		}

		mid = deepSet(mid, setTypeOut(r, m.TypeOut, m.DefaultNull), strings.Split(keyOut, ".")...)
	}

	out, err := json.Marshal(mid)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidJSONOutput, err)
	}

	return out, nil
}

// deepSet is a recursive function that creates a map tree from the keys, putting the value in the end
func deepSet(in map[string]interface{}, value interface{}, keys ...string) map[string]interface{} {
	// We should never reach this point
	if len(keys) == 0 {
		return nil
	}

	currentKey := keys[0]
	// We reached the node where to put the value
	if len(keys) == 1 {
		in[currentKey] = value
		return in
	}

	// Check if the node already exist
	if _, alreadySet := in[currentKey]; alreadySet {
		if _, expectedType := in[currentKey].(map[string]interface{}); expectedType {
			in[currentKey] = deepSet(in[currentKey].(map[string]interface{}), value, keys[1:]...)
			return in
		}
		in[currentKey] = value
		return in
	}

	// Create new node and continue
	in[currentKey] = map[string]interface{}{}
	in[currentKey] = deepSet(in[currentKey].(map[string]interface{}), value, keys[1:]...)
	return in
}

// setTypeOut tries to convert any value to the desired one
func setTypeOut(in interface{}, typeOut, defaultNull string) interface{} {
	if typeOut == "" && defaultNull == "" {
		return in
	}

	if in == nil && defaultNull == "" {
		return in
	}

	if in == nil && defaultNull != "" {
		in = defaultNull
	}

	v := fmt.Sprint(in)
	switch typeOut {
	case "int":
		x, _ := strconv.Atoi(v)
		return x
	case "float":
		x, _ := strconv.ParseFloat(v, 64)
		return x
	}

	return v
}
