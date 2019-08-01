package decode

import (
	"gopkg.in/yaml.v3"

	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
)

type unmarshalFunc func([]byte, interface{}) error

// FromJSON reads an AsyncAPI document from input in a JSON format
// and stores it in the value. If operation fails function returns an error.
//
// See InvalidProperty, InvalidDocument, UnsupportedAsyncapiVersion in pkg/error.
func FromJSON(v interface{}, reader io.Reader) error {
	return json.NewDecoder(reader).Decode(&v)
}

// FromYaml reads an AsyncAPI document from input in the YAML format
// and stores it in the value. If the operation fails, the function returns an error.
//
// See InvalidProperty, InvalidDocument, UnsupportedAsyncapiVersion in pkg/error.
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

// FromJSONWithYamlFallback reads an AsyncAPI document from input in the JSON format.
// If the operation fails, the function tries to read the AsyncAPI document in the YAML format.
// If any of the decoding attempts succeeds, the result is stored in the value.
// If both decoding attempts fail, the function returns an error.
//
// See InvalidProperty, InvalidDocument, UnsupportedAsyncapiVersion in pkg/error.
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
