package main

import (
	v2 "github.com/asyncapi/converter-go/pkg/converter/v2"
	"github.com/asyncapi/converter-go/pkg/decode"
	"github.com/asyncapi/converter-go/pkg/encode"

	"log"
	"os"
	"strings"
)

func main() {
	reader := strings.NewReader(schema)

	// create json to yaml converter
	converter, err := v2.New(decode.FromJSON, encode.ToYaml)
	if err != nil {
		log.Fatal(err)
	}

	// convert document
	err = converter.Convert(reader, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}

var schema = `{
    "info": {
        "version": "1.0.0",
        "title": "Not example"
    },
    "topics": {
        "test": {
            "publish": {
                "$ref": "#/components/messages/testMessages"
            }
        }
    },
    "asyncapi": "1.1.0",
    "components": {
        "messages": {
            "testMessages": {
                "payload": {
                    "$ref": "#/components/schemas/testSchema"
                }
            }
        },
        "schemas": {
            "testSchema": {
                "type": "object",
                "properties": {
                    "key": {
                        "not": {
                            "type": "integer"
                        }
                    }
                }
            }
        }
    }
}`
