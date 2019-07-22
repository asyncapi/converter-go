package asyncapi

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"encoding/json"
	"io"
)

type Format = int

const (
	Json Format = iota + 1
	Yaml
)

type Encode = func(interface{}, io.Writer) error

type Converter interface {
	Do(reader io.Reader, writer io.Writer) error
}

type UnmarshalFunc func([]byte, interface{}) error

var (
	ErrInvalidProperty = errors.New("invalid property")
	ErrInvalidDocument = errors.New("unable to decode document")
	JsonEncode         = func(data interface{}, writer io.Writer) error {
		return json.NewEncoder(writer).Encode(data)
	}
	YamlEncode = func(data interface{}, writer io.Writer) error {
		return yaml.NewEncoder(writer).Encode(&data)
	}
)

func BuildUnmarshalWithFallback(primary UnmarshalFunc, fallback ...UnmarshalFunc) UnmarshalFunc {
	return func(bytes []byte, out interface{}) error {
		var err error
		for _, unmarshal := range append([]UnmarshalFunc{primary}, fallback...) {
			err = unmarshal(bytes, out)
			if err == nil {
				return nil
			}
		}
		return err
	}
}

func IsInvalidPropertyErr(e error) bool {
	return e == ErrInvalidProperty
}

func EncodeFunction(encodeFormat Format) Encode {
	switch encodeFormat {
	case Yaml:
		return YamlEncode
	default:
		return JsonEncode
	}
}
