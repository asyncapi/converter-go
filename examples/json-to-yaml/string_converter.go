package main

import (
	v2 "asyncapi-converter/pkg/converter/v2"
	"asyncapi-converter/pkg/decode"
	"asyncapi-converter/pkg/encode"

	"log"
	"os"
	"strings"
)

func main() {
	reader := strings.NewReader(schema)

	// create json to yaml converter
	converter, err := v2.NewConverter(decode.FromJson, encode.ToYaml)
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
