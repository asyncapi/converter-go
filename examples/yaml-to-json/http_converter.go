package main

import (
	v2 "asyncapi-converter/pkg/converter/v2"
	"asyncapi-converter/pkg/decode"
	"asyncapi-converter/pkg/encode"

	"log"
	"net/http"
	"os"
)

func main() {

	url := "https://raw.githubusercontent.com/asyncapi/converter/master/test/input/1.2.0/gitter-streaming.yml"

	// get gitter-streaming.yml from url
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	// create yaml to json converter
	converter, err := v2.NewConverter(decode.YamlDecoder, encode.JsonEncoder)
	if err != nil {
		log.Fatal(err)
	}

	// convert document
	err = converter.Convert(resp.Body, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}
