package encode

import (
	"gopkg.in/yaml.v3"

	"encoding/json"
	"io"
)

func ToJSON(i interface{}, writer io.Writer) error {
	return json.NewEncoder(writer).Encode(i)
}

func ToYaml(i interface{}, writer io.Writer) error {
	return yaml.NewEncoder(writer).Encode(i)
}
