package asyncapi

import (
	"fmt"
	"gopkg.in/yaml.v3"
)

func UnmarshalYaml(in []byte, out interface{}) error {
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
