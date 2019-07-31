package encode

import (
	"gopkg.in/yaml.v3"

	"encoding/json"
	"io"
)

// ToJSON writes an asyncapi document in JSON format into a stream
func ToJSON(i interface{}, writer io.Writer) error {
	return json.NewEncoder(writer).Encode(i)
}

// ToYaml writes an AsyncAPI document in the YAML format encoding it into a stream.
func ToYaml(i interface{}, writer io.Writer) error {
	return yaml.NewEncoder(writer).Encode(i)
}
