package main

import (
	"asyncapi-converter/pkg/asyncapi"
	"asyncapi-converter/pkg/converter/v2rc1"
	"net/http"
	"net/url"

	"github.com/docopt/docopt-go"
	"github.com/pkg/errors"

	"fmt"
	"io"
	"log"
	"os"
)

const version = "asyncapi-converter 0.1.0rc1"

var (
	errFileDoesNotExist = errors.New("file does not exist")
	errInvalidArgument  = errors.New("invalid argument")
)

type helper struct {
	docopt.Opts
	data interface{}
}

func newHelper(opts docopt.Opts) helper {
	return helper{
		Opts: opts,
	}
}

func (h helper) id() *string {
	idOption, ok := h.Opts["--id"]
	if !ok || idOption == nil {
		return nil
	}
	id := fmt.Sprintf("%v", idOption)
	return &id
}

func (h helper) encodeFormat() (asyncapi.Format, error) {
	if _, ok := h.Opts["--toYAML"]; !ok {
		return asyncapi.Json, nil
	}
	toYaml, ok := h.Opts["--toYAML"].(bool)
	if !ok {
		return 0, errors.Wrap(errInvalidArgument, "--toYAML")
	}
	if toYaml {
		return asyncapi.Yaml, nil
	}
	return asyncapi.Json, nil
}

func (h helper) reader() (io.Reader, error) {
	fileOption := h.Opts["<PATH>"]
	path := fmt.Sprintf("%v", fileOption)
	if _, err := url.ParseRequestURI(path); err == nil {
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

func (h helper) newConverterAndReader() (asyncapi.Converter, io.Reader, error) {
	reader, err := h.reader()
	if err != nil {
		return nil, nil, err
	}
	format, err := h.encodeFormat()
	if err != nil {
		return nil, nil, err
	}
	converter, err := v2rc1.NewConverter(v2rc1.WithEncoding(format), v2rc1.WithId(h.id()))
	return converter, reader, err
}

func main() {
	usage := fmt.Sprintf(`
  Convert AsyncAPI documents from version 1.x to %s. 

  Usage:
    asyncapi-converter <PATH> [--toYAML] [--id=<id>]
    asyncapi-converter -h | --help | --version

  Arguments:
    PATH        a path to asyncapi document (either url or local file, supports json and yaml format)  	

  Options:
    --toYAML    produces results in yaml format instead json
    --id=<id>   allows to specify application id`, v2rc1.AsyncapiVersion)

	opts, err := docopt.ParseArgs(usage, nil, version)
	if err != nil {
		log.Fatal(err)
	}
	helper := newHelper(opts)
	converter, reader, err := helper.newConverterAndReader()
	if err != nil {
		log.Fatal(err)
	}
	err = converter.Do(reader, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}
