package v2rc1

import (
	"asyncapi-converter/pkg/asyncapi"
	"regexp"

	"github.com/pkg/errors"

	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strings"
)

const AsyncapiVersion = "2.0.0-rc1"

type converter struct {
	id             *string
	encodingFormat asyncapi.Format
	data           map[string]interface{}
	unmarshal      func(reader io.Reader) error
}

func (c *converter) Do(reader io.Reader, writer io.Writer) error {
	steps := []func() error{
		func() error {
			return c.unmarshal(reader)
		},
		c.verifyAsyncapiVersion,
		c.updateId,
		c.updateVersion,
		c.updateServers,
		c.createChannels,
		c.cleanup,
		func() error {
			encode := asyncapi.EncodeFunction(c.encodingFormat)
			return encode(c.data, writer)
		},
	}
	for _, step := range steps {
		err := step()
		if err != nil {
			return err
		}
	}
	return nil
}

type ConverterOption = func(*converter) error

func newConverter(options ...ConverterOption) (asyncapi.Converter, error) {
	converter := converter{}
	for _, f := range options {
		err := f(&converter)
		if err != nil {
			return nil, err
		}
	}
	return &converter, nil
}

func NewConverter(options ...ConverterOption) (asyncapi.Converter, error) {
	return newConverter(append([]ConverterOption{WithFallbackUnmarshal()}, options...)...)
}

func NewJsonConverter(options ...ConverterOption) (asyncapi.Converter, error) {
	return newConverter(append([]ConverterOption{WithJsonUnmarshal()}, options...)...)
}

func NewYamlConverter(options ...ConverterOption) (asyncapi.Converter, error) {
	return newConverter(append([]ConverterOption{WithYamlUnmarshal()}, options...)...)
}

func WithId(id *string) ConverterOption {
	return func(converter *converter) error {
		converter.id = id
		return nil
	}
}

func WithEncoding(encodingFormat asyncapi.Format) ConverterOption {
	return func(converter *converter) error {
		converter.encodingFormat = encodingFormat
		return nil
	}
}

func WithFallbackUnmarshal() ConverterOption {
	return func(c *converter) error {
		c.unmarshal = func(reader io.Reader) error {
			bytes, err := ioutil.ReadAll(reader)
			if err != nil {
				return err
			}
			unmarshal := asyncapi.BuildUnmarshalWithFallback(json.Unmarshal, asyncapi.UnmarshalYaml)
			var data interface{}
			err = unmarshal(bytes, &data)
			if err != nil {
				log.Fatalln(err)
			}
			var ok bool
			c.data, ok = data.(map[string]interface{})
			if !ok {
				return asyncapi.ErrInvalidDocument
			}
			return nil
		}
		return nil
	}
}

func WithJsonUnmarshal() ConverterOption {
	return func(c *converter) error {
		c.unmarshal = func(reader io.Reader) error {
			return json.NewDecoder(reader).Decode(&c.data)
		}
		return nil
	}
}

func WithYamlUnmarshal() ConverterOption {
	return func(c *converter) error {
		c.unmarshal = func(reader io.Reader) error {
			bytes, err := ioutil.ReadAll(reader)
			if err != nil {
				return err
			}
			var data interface{}
			err = asyncapi.UnmarshalYaml(bytes, &data)
			if err != nil {
				return err
			}
			var ok bool
			c.data, ok = data.(map[string]interface{})
			if !ok {
				return asyncapi.ErrInvalidDocument
			}
			return nil
		}
		return nil
	}
}

func (c *converter) updateId() error {
	if c.id != nil {
		c.data["id"] = *c.id
		return nil
	}
	info, ok := c.data["info"].(map[string]interface{})
	if !ok {
		return asyncapi.NewErrInvalidProperty("info")
	}
	title, ok := info["title"]
	if !ok {
		return asyncapi.NewErrInvalidProperty("title")
	}
	c.data["id"] = fmt.Sprintf(`urn:%s`, extractId(fmt.Sprintf("%v", title)))
	return nil
}

func (c *converter) updateVersion() error {
	c.data["asyncapi"] = AsyncapiVersion
	return nil
}

func (c *converter) updateServers() error {
	servers, ok := c.data["servers"].([]interface{})
	if !ok {
		return asyncapi.NewErrInvalidProperty("servers")
	}
	_, containsSecurity := c.data["security"]
	for _, item := range servers {
		server, ok := item.(map[string]interface{})
		if !ok {
			return asyncapi.NewErrInvalidProperty("server")
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
	return nil
}

func (c *converter) channelsFromTopics() error {
	channels := make(map[string]interface{})
	topics, ok := c.data["topics"].(map[string]interface{})
	if !ok {
		return asyncapi.NewErrInvalidProperty("topics")
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
		return asyncapi.NewErrInvalidProperty("events")
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
		return asyncapi.NewErrInvalidProperty("stream")
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
	return asyncapi.NewErrInvalidProperty("missing one of topics/stream/events")
}

func extractId(value string) string {
	title := strings.ToLower(value)
	return strings.Join(strings.Split(title, " "), ".")
}

var versionRegexp = regexp.MustCompile("^1\\.[0-2].0$")

func (c *converter) verifyAsyncapiVersion() error {
	version, ok := c.data["asyncapi"]
	if !ok {
		return asyncapi.NewErrInvalidProperty("asyncapi")
	}
	versionString := fmt.Sprintf("%v", version)
	switch versionRegexp.Match([]byte(versionString)) {
	case true:
		return nil
	default:
		return errors.Wrap(asyncapi.ErrUnsupportedAsyncapiVersion, versionString)
	}
}
