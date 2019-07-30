package decode

import (
	"gopkg.in/yaml.v3"

	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
)

type unmarshalFunc func([]byte, interface{}) error

// FromJSON reads an asyncapi document from an input in a JSON format
// and stores it in the value. If operation fails function returns an error.
//
// See InvalidProperty, InvalidDocument, UnsupportedAsyncapiVersion in pkg error.
func FromJSON(v interface{}, reader io.Reader) error {
	return json.NewDecoder(reader).Decode(&v)
}

// FromYaml reads an asyncapi document from an input in yaml format
// and stores it in the value. If operation fails it returns error.
//
// See InvalidProperty, InvalidDocument, UnsupportedAsyncapiVersion in pkg error.
func FromYaml(v interface{}, reader io.Reader) error {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	err = unmarshalYaml(data, v)
	if err != nil {
		return err
	}
	return nil
}

// FromJSONWithYamlFallback reads an asyncapi document from an input in JSON format,
// if this operation fails, the function tries to read an asyncapi document in yaml format.
// If any decoding attempt succeeds, the result is stored in the value.
// If both decoding attempts fail, function returns an error.
//
// See InvalidProperty, InvalidDocument, UnsupportedAsyncapiVersion in pkg error.
func FromJSONWithYamlFallback(out interface{}, reader io.Reader) error {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	for _, unmarshal := range []unmarshalFunc{json.Unmarshal, unmarshalYaml} {
		err = unmarshal(data, out)
		if err == nil {
			return nil
		}
	}
	return err
}

func unmarshalYaml(in []byte, out interface{}) error {
	var result interface{}
	if err := yaml.Unmarshal(in, &result); err != nil {
		return err
	}
	*out.(*interface{}) = convertMap(result)
	return nil
}

func convertInterfaceArray(in []interface{}) []interface{} {
	result := make([]interface{}, len(in))
	for i, v := range in {
		result[i] = convertMap(v)
	}
	return result
}

func convertInterfaceMap(in map[interface{}]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range in {
		result[fmt.Sprintf("%v", k)] = convertMap(v)
	}
	return result
}

func convertMap(v interface{}) interface{} {
	switch v := v.(type) {
	case []interface{}:
		return convertInterfaceArray(v)
	case map[interface{}]interface{}:
		return convertInterfaceMap(v)
	default:
		return v
	}
}
