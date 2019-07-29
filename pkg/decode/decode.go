package decode

import (
	"gopkg.in/yaml.v3"

	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
)

type unmarshalFunc func([]byte, interface{}) error

func FromJSON(v interface{}, reader io.Reader) error {
	return json.NewDecoder(reader).Decode(&v)
}

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
