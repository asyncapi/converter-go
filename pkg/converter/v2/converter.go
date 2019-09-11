package v2

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	asyncapierr "github.com/asyncapi/converter-go/pkg/error"
)

// AsyncapiVersion is the AsyncAPI version that the document will be converted to.
const AsyncapiVersion = "2.0.0-rc2"

// Decode reads an AsyncAPI document from input and stores it in the value.
type Decode = func(interface{}, io.Reader) error

// Encode writes an AsyncAPI document encoding it into a stream.
type Encode = func(interface{}, io.Writer) error

// Converter converts an AsyncAPIi document from versions 1.0.0, 1.1.1 and 1.2.0 to version 2.0.0.
type Converter interface {
	Convert(reader io.Reader, writer io.Writer) error
}

type converter struct {
	id     *string
	data   map[string]interface{}
	decode Decode
	encode Encode
}

func (c *converter) buildEncodeFunction(writer io.Writer) func() error {
	return func() error {
		return c.encode(&c.data, writer)
	}
}

func (c *converter) buildDecodeFunction(reader io.Reader) func() error {
	return func() error {
		var data interface{}
		decode := c.decode(&data, reader)
		var ok bool
		c.data, ok = data.(map[string]interface{})
		if !ok {
			return asyncapierr.NewInvalidDocument()
		}
		return decode
	}
}

func (c *converter) Convert(reader io.Reader, writer io.Writer) error {
	steps := []func() error{
		c.buildDecodeFunction(reader),
		c.verifyAsyncapiVersion,
		c.updateID,
		c.updateVersion,
		c.updateServers,
		c.createChannels,
		c.alterChannels,
		c.updateComponents,
		c.cleanup,
		c.buildEncodeFunction(writer),
	}
	for _, step := range steps {
		err := step()
		if err != nil {
			return err
		}
	}
	return nil
}

// ConverterOption is a functional option that allows you to provide
// a meaningful converter configuration that can grow over time.
type ConverterOption func(*converter) error

// New creates a new converter.
//
// See Decode, Encode and ConverterOption.
func New(decode Decode, encode Encode, options ...ConverterOption) (Converter, error) {
	converter := converter{
		encode: encode,
		decode: decode,
	}
	for _, option := range options {
		if err := option(&converter); err != nil {
			return nil, err
		}
	}
	return &converter, nil
}

// WithID is a functional option that allows you to specify the application ID.
func WithID(id *string) ConverterOption {
	return func(converter *converter) error {
		converter.id = id
		return nil
	}
}

func (c *converter) updateID() error {
	if c.id != nil {
		c.data["id"] = *c.id
		return nil
	}
	return nil
}

func (c *converter) updateVersion() error {
	c.data["asyncapi"] = AsyncapiVersion
	return nil
}

func (c *converter) updateServers() error {
	servers, ok := c.data["servers"].([]interface{})

	if !ok {
		return nil
	}

	_, containsSecurity := c.data["security"]
	for _, item := range servers {
		server, ok := item.(map[string]interface{})
		if !ok {
			return asyncapierr.NewInvalidProperty("server")
		}
		server["protocol"] = server["scheme"]
		delete(server, "scheme")
		if containsSecurity {
			server["security"] = c.data["security"]
		}
		if schemaVersion, ok := server["schemeVersion"]; ok {
			server["protocolVersion"] = schemaVersion
			delete(server, "schemeVersion")
		}
	}

	var mappedServers = make(map[string]interface{})

	for index, item := range servers {
		//done same way as in https://github.com/asyncapi/converter/blob/020946e745342a6751565406e156c499859f5763/lib/index.js#L106
		if index == 0 {
			mappedServers["default"] = item
		} else {
			mappedServers[fmt.Sprintf("server%d", index)] = item
		}
	}

	c.data["servers"] = mappedServers

	return nil
}

func (c *converter) channelsFromTopics() error {
	channels := make(map[string]interface{})
	topics, ok := c.data["topics"].(map[string]interface{})
	if !ok {
		return asyncapierr.NewInvalidProperty("topics")
	}
	for key, value := range topics {
		var topic string
		if _, ok := c.data["baseTopic"]; ok {
			topic = fmt.Sprintf("%v", c.data["baseTopic"])
		}
		if topic != "" {
			topic = fmt.Sprintf(`%s/%s`, topic, key)
		} else {
			topic = fmt.Sprintf("%v", key)
		}
		channelKey := strings.ReplaceAll(topic, ".", "/")
		if topic, ok := value.(map[string]interface{}); ok {
			switch {
			case topic["publish"] != nil:
				topic["publish"] = map[string]interface{}{
					"message": topic["publish"],
				}
			case topic["subscribe"] != nil:
				topic["subscribe"] = map[string]interface{}{
					"message": topic["subscribe"],
				}
			}
		}
		channels[channelKey] = value
	}
	c.data["channels"] = channels
	return nil
}

func (c *converter) channelsFromStream() error {
	events, ok := c.data["stream"].(map[string]interface{})
	if !ok {
		return asyncapierr.NewInvalidProperty("events")
	}
	channel := make(map[string]interface{})
	
	if _, ok := events["read"]; ok {
		channel["publish"] = map[string]map[string]interface{}{
			"message": {
				"oneOf": events["read"],
			},
		}
	}
	if _, ok := events["write"]; ok {
		channel["subscribe"] = map[string]map[string]interface{}{
			"message": {
				"oneOf": events["write"],
			},
		}
	}
	c.data["channels"] = map[string]interface{}{
		"/": channel,
	}
	return nil
}

func (c *converter) channelsFromEvents() error {
	stream, ok := c.data["events"].(map[string]interface{})
	if !ok {
		return asyncapierr.NewInvalidProperty("stream")
	}
	channel := make(map[string]interface{})
	if _, ok := stream["receive"]; ok {
		channel["publish"] = map[string]map[string]interface{}{
			"message": {
				"oneOf": stream["receive"],
			},
		}
	}
	if _, ok := stream["send"]; ok {
		channel["subscribe"] = map[string]map[string]interface{}{
			"message": {
				"oneOf": stream["send"],
			},
		}
	}
	c.data["channels"] = map[string]interface{}{
		"/": channel,
	}
	return nil
}

func (c *converter) cleanup() error {
	delete(c.data, "topics")
	delete(c.data, "baseTopic")
	delete(c.data, "stream")
	delete(c.data, "events")
	delete(c.data, "security")
	return nil
}

func (c *converter) createChannels() error {
	if _, ok := c.data["topics"]; ok {
		return c.channelsFromTopics()
	}
	if _, ok := c.data["stream"]; ok {
		return c.channelsFromStream()
	}
	if _, ok := c.data["events"]; ok {
		return c.channelsFromEvents()
	}
	return asyncapierr.NewInvalidProperty("missing one of topics/stream/events")
}

func prt(arg ...interface{}) {
	fmt.Fprintf(os.Stdout, "%v\n\n", arg)
}

func (c *converter) updateComponents() error {
	components, ok := c.data["components"].(map[string]interface{})

	if !ok {
		return nil
	}

	parameters, ok := components["parameters"].(map[string]interface{})

	if !ok {
		return nil
	}

	for _, rawParam := range parameters {
		if param, ok := rawParam.(map[string]interface{}); ok {
			delete(param, "name")
			rawParam = param
		}
	}
	return nil
}

func alterParameters(parameters []interface{}, key string) (map[string]interface{}, error) {
	paramsMap := make(map[string]interface{})
	re := regexp.MustCompile(`{([^}]+)}`)
	var paramNames []string

	for _, part := range re.FindAll([]byte(key), -1) {
		paramNames = append(paramNames, string(part))
	}

	for index, paramI := range parameters {

		param, ok := paramI.(map[string]interface{})

		if !ok {
			return nil, asyncapierr.NewInvalidProperty("malformed parameter")
		}

		name := "default"

		if paramName, ok := param["name"].(string); ok {
			name = paramName
		} else {
			if len(paramNames) > index {
				name = paramNames[index]
			}
		}
		name = strings.TrimLeft(strings.TrimRight(name, "}"), "{")

		if param["name"] != nil {
			delete(param, "name")
		}

		paramsMap[name] = param

	}
	return paramsMap, nil
}

func (c *converter) alterChannels() error {
	channels, ok := c.data["channels"].(map[string]interface{})

	if !ok {
		return asyncapierr.NewInvalidProperty("missing channels")
	}

	for key, item := range channels {
		channel, ok := item.(map[string]interface{})
		if !ok {
			return asyncapierr.NewInvalidProperty("malformed channel")
		}

		if params, ok := channel["parameters"].([]interface{}); ok {

			altered, err := alterParameters(params, key)
			if err != nil {
				return err
			}
			channel["parameters"] = altered
		}

		if publish, ok := channel["publish"].(map[string]interface{}); ok {
			alterOperation(publish)
			protocolInfoToBindings(publish)
		}
		if subscribe, ok := channel["subscribe"].(map[string]interface{}); ok {
			alterOperation(subscribe)
			protocolInfoToBindings(subscribe)
		}
		protocolInfoToBindings(channel)
	}
	return nil
}

func protocolInfoToBindings(arg map[string]interface{}) {
	if protocolInfo, ok := arg["protocolInfo"]; ok {
		arg["bindings"] = protocolInfo
		delete(arg, "protocolInfo")
	}
}

func headersToSchema(arg map[string]interface{}) {
	if headers, ok := arg["headers"]; ok {
		headers = map[string]interface{}{
			"type":       "object",
			"properties": headers,
		}
	}
}

func alterOperation(operation map[string]interface{}) {
	if message, ok := operation["message"].(map[string]interface{}); ok {
		if oneOf, ok := message["oneOf"].([]map[string]interface{}); ok {
			for _, elem := range oneOf {
				protocolInfoToBindings(elem)
				headersToSchema(elem)
			}
		} else {
			protocolInfoToBindings(message)
			headersToSchema(message)
		}
	}
}

func extractID(value string) string {
	title := strings.ToLower(value)
	return strings.Join(strings.Split(title, " "), ".")
}

var versionRegexp = regexp.MustCompile("^1\\.[0-2]\\.0$")

func (c *converter) verifyAsyncapiVersion() error {
	version, ok := c.data["asyncapi"]
	if !ok {
		return asyncapierr.NewInvalidProperty("asyncapi")
	}
	versionString := fmt.Sprintf("%v", version)
	switch {
	case versionString == AsyncapiVersion:
		return asyncapierr.NewDocumentVersionUpToDate(AsyncapiVersion)
	case versionRegexp.Match([]byte(versionString)):
		return nil
	default:
		return asyncapierr.NewUnsupportedAsyncapiVersion(versionString)
	}
}
