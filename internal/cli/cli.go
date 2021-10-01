package cli

import (
	"github.com/docopt/docopt-go"
	"github.com/pkg/errors"

	v2 "github.com/asyncapi/converter-go/pkg/converter/v2"
	"github.com/asyncapi/converter-go/pkg/decode"
	asyncapiEncode "github.com/asyncapi/converter-go/pkg/encode"

	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

var (
	errFileDoesNotExist = errors.New("file does not exist")
	errInvalidArgument  = errors.New("invalid argument")
)

const (
	encodeOptionYAML  = "--toYAML"
	fileOptionPath = "<PATH>"
	idOption = "--id"
)

type encode = func(interface{}, io.Writer) error

// Converter converts an AsyncAPI document.
type Converter interface {
	Convert(reader io.Reader, writer io.Writer) error
}

var _ v2.Converter = Converter(nil)

// Cli is a helper type that allows you to instantiate the AsyncAPI Converter and io.Reader of
// the converted document with arguments passed from the terminal.
type Cli struct {
	docopt.Opts
}

// New returns a new Cli instance.
func New(opts docopt.Opts) Cli {
	return Cli{
		Opts: opts,
	}
}

func (h Cli) id() *string {
	idOption, ok := h.Opts[idOption]
	if !ok || idOption == nil {
		return nil
	}
	id := fmt.Sprintf("%v", idOption)
	return &id
}

func (h Cli) encode() (encode, error) {
	if _, ok := h.Opts[encodeOptionYAML]; !ok {
		return asyncapiEncode.ToJSON, nil
	}
	toYaml, ok := h.Opts[encodeOptionYAML].(bool)
	if !ok {
		return nil, errors.Wrap(errInvalidArgument, encodeOptionYAML)
	}
	if toYaml {
		return asyncapiEncode.ToYaml, nil
	}
	return asyncapiEncode.ToJSON, nil
}

func isURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func (h Cli) reader() (io.Reader, error) {
	fileOption := h.Opts[fileOptionPath]
	path := fmt.Sprintf("%v", fileOption)
	if isURL(path) {
		resp, err := http.Get(path)
		if err != nil {
			return nil, err
		}
		return resp.Body, nil
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, errors.Wrap(errFileDoesNotExist, path)
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// NewConverterAndReader creates both a converter and a reader of the converted AsyncAPI document.
func (h Cli) NewConverterAndReader() (Converter, io.Reader, error) {
	reader, err := h.reader()
	if err != nil {
		return nil, nil, err
	}
	encode, err := h.encode()
	if err != nil {
		return nil, nil, err
	}
	converter, err := v2.New(decode.FromJSONWithYamlFallback, encode, v2.WithID(h.id()))
	return converter, reader, err
}
