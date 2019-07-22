# AsyncAPI Converter

Convert AsyncAPI documents from version 1.x to 2.0.0-rc1.

## Installation

`go get github.com/asyncapi/converter-go`

## Usage

### From CLI

```bash
Usage:
asyncapi-converter <path> [--toYAML] [--id=<id>]
asyncapi-converter -h | --help | --version

Options:
--toYAML    produces results in yaml format
--id        allows to specify application id
```

Minimal example:

```bash
$ asyncapi-converter https://git.io/fjMPF --toYAML
asyncapi: 2.0.0-rc1
channels:
    /:
        publish:
            message:
                oneOf:
                  - $ref: '#/components/messages/chatMessage'
...
```

Specify the application id:

```bash
$ asyncapi-converter streetlights.yml --id=urn:com.asynapi.streetlights
{"asyncapi":"2.0.0-rc1","channels":{"/":{"publish":{"message":{"oneOf":[{"$ref":"#/components/messages/chatMessage"}...
..."id":"urn:com.asynapi.streetlights","info":{"title...
```

Save the result in a file:
```bash
$ asyncapi-converter streetlights.json --toYAML --id=urn:com.asynapi.streetlights > streetlights2.yml
```

### As a package

```go
package main

import (
	"asyncapi-converter/pkg/asyncapi"
	"asyncapi-converter/pkg/converter/v2rc1"
	"log"
	"net/http"
	"os"
)

var (
	url = "https://git.io/fjMPF"
	id  = "urn:gitter-streaming"
)

func main() {
	// get gitter-streaming.yml
	resp, err := http.Get(url)
	handleError(err)

	// create yaml converter
	converter, err := v2rc1.NewYamlConverter(
		v2rc1.WithEncoding(asyncapi.Json),
		v2rc1.WithId(&id),
	)
	handleError(err)

	// convert document
	err = converter.Do(resp.Body, os.Stdout)
	handleError(err)
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
```