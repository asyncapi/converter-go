# AsyncAPI Converter

[![Build Status](https://godoc.org/github.com/asyncapi/converter-go?status.svg)](https://godoc.org/github.com/asyncapi/converter-go) [![Build Status](https://travis-ci.org/asyncapi/converter-go.svg?branch=master)](https://travis-ci.org/asyncapi/converter-go) [![Go Report Card](https://goreportcard.com/badge/github.com/asyncapi/converter-go)](https://goreportcard.com/report/github.com/asyncapi/converter-go) 

## Overview

The AsyncAPI Converter converts AsyncAPI documents from versions 1.0.0, 1.1.0 and 1.2.0 to version 2.0.0-rc1. It supports both `json` and `yaml` formats on input and output. By default, the AsyncAPI Converter converts a document into the `json` format.

## Prerequisites

- [Golang](https://golang.org/dl/) version 1.11+

## Installation

 To install the AsyncAPI Converter package, run:

```bash
go get github.com/asyncapi/converter-go
```

## Usage

You can use the AsyncAPI Converter in the terminal or as a package.

### In CLI

Before you use the AsyncAPI Converter in the terminal, build the application. Run:

```bash
git clone https://github.com/asyncapi/converter-go.git
cd ./converter-go
go build -o=asyncapi-converter ./cmd/api-converter/main.go
```

To convert a document use the following command:

```text
asyncapi-converter <document_path> [--toYAML] [--id=<id>]
```

where:

- `document_path` is a mandatory argument that is either a URL or  a file path to an AsyncAPI document
- `--toYAML` is an optional argument that allows producing results in the `yaml` format instead of `json`
- `--id` is an optional argument that allows specifying the application `id`

**Examples**

See the following minimal examples of the AsyncAPI Converter usage in the terminal:


- `gitter-streaming` conversion from version 1.2.0 to 2.0.0-rc1 in the `json` format

  ```text
  asyncapi-converter https://git.io/fjMPF
  ```

- `gitter-streaming` conversion from version 1.2.0 to 2.0.0-rc1 in the `yaml` format

  ```bash
  asyncapi-converter https://git.io/fjMPF --toYAML
  ```

- `gitter-streaming` conversion from version 1.2.0 to 2.0.0-rc1 in the `json` format specifying the application `id`

  ```bash
  asyncapi-converter https://git.io/fjMXl --id=urn:com.asynapi.streetlights
  ```

### As a package

To see examples of how to use the AsyncAPI Converter as a package, go to the [README.md](./examples/README.md).

## Contribution

If you have a feature request, add it as an issue or propose changes in a pull request (PR).
If you create a feature request, use the dedicated **Feature request** issue template. When you create a PR, follow the contributing rules described in the [`CONTRIBUTING.md`](CONTRIBUTING.md) document.

## Credits

<p align="center">
 <a href="https://kyma-project.io/" target="_blank">
  <img src="https://raw.githubusercontent.com/kyma-project/kyma/master/logo.png" width="235">
 </a>
</p>
