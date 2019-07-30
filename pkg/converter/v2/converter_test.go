package v2

import (
	"github.com/asyncapi/converter-go/pkg/decode"
	"github.com/asyncapi/converter-go/pkg/encode"
	. "github.com/onsi/gomega"

	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func TestNewJsonConverter(t *testing.T) {
	testID := "test"
	tests := []struct {
		inputFilePath    string
		expectedFilePath string
		options          []ConverterOption
	}{
		{
			inputFilePath:    "./testdata/input/streetlights1.0.0.json",
			expectedFilePath: "./testdata/output/streetlights.json",
		},
		{
			inputFilePath:    "./testdata/input/streetlights1.1.0.json",
			expectedFilePath: "./testdata/output/streetlights.json",
		},
		{
			inputFilePath:    "./testdata/input/streetlights1.2.0.json",
			expectedFilePath: "./testdata/output/streetlights.json",
		},
		{
			inputFilePath:    "./testdata/input/gitter-streaming1.2.0.json",
			expectedFilePath: "./testdata/output/gitter-streaming.json",
		},
		{
			inputFilePath:    "./testdata/input/gitter-streaming1.2.0_modified_write.json",
			expectedFilePath: "./testdata/output/gitter-streaming_modified_write.json",
		},
		{
			inputFilePath:    "./testdata/input/gitter-streaming1.2.0_with_id_option.json",
			expectedFilePath: "./testdata/output/gitter-streaming_with_id_option.json",
			options: []ConverterOption{
				WithID(&testID),
			},
		},
		{
			inputFilePath:    "./testdata/input/slack-rtm1.2.0.json",
			expectedFilePath: "./testdata/output/slack-rtm.json",
		},
		{
			inputFilePath:    "./testdata/input/streetlights1.0.0_no_base_topic.json",
			expectedFilePath: "./testdata/output/streetlights_no_base_topic.json",
		},
	}
	for _, test := range tests {
		t.Run(test.inputFilePath, func(t *testing.T) {
			g := NewWithT(t)
			converter, err := New(decode.FromJSON, encode.ToJSON, test.options...)
			g.Expect(err).To(BeNil(), "error while creating converter")
			result := convertFile(converter, test.inputFilePath, g)
			expected, err := ioutil.ReadFile(test.expectedFilePath)
			g.Expect(err).To(BeNil(), "error while reading file containing expected results")
			g.Expect(result).To(MatchJSON(string(expected)))
		})
	}
}

func TestNewConverter(t *testing.T) {
	testID := "test"
	tests := []struct {
		inputFilePath    string
		expectedFilePath string
		options          []ConverterOption
	}{
		{
			inputFilePath:    "./testdata/input/streetlights1.0.0.json",
			expectedFilePath: "./testdata/output/streetlights.yaml",
		},
		{
			inputFilePath:    "./testdata/input/streetlights1.0.0.yaml",
			expectedFilePath: "./testdata/output/streetlights.yaml",
		},
		{
			inputFilePath:    "./testdata/input/streetlights1.1.0.json",
			expectedFilePath: "./testdata/output/streetlights.yaml",
		},
		{
			inputFilePath:    "./testdata/input/streetlights1.1.0.yaml",
			expectedFilePath: "./testdata/output/streetlights.yaml",
		},
		{
			inputFilePath:    "./testdata/input/streetlights1.2.0.yaml",
			expectedFilePath: "./testdata/output/streetlights.yaml",
		},
		{
			inputFilePath:    "./testdata/input/streetlights1.2.0.json",
			expectedFilePath: "./testdata/output/streetlights.yaml",
		},
		{
			inputFilePath:    "./testdata/input/gitter-streaming1.2.0.json",
			expectedFilePath: "./testdata/output/gitter-streaming.yml",
		},
		{
			inputFilePath:    "./testdata/input/gitter-streaming1.2.0_modified_write.json",
			expectedFilePath: "./testdata/output/gitter-streaming_modified_write.yml",
		},
		{
			inputFilePath:    "./testdata/input/gitter-streaming1.2.0_with_id_option.json",
			expectedFilePath: "./testdata/output/gitter-streaming_with_id_option.yml",
			options: []ConverterOption{
				WithID(&testID),
			},
		},
		{
			inputFilePath:    "./testdata/input/slack-rtm1.2.0.json",
			expectedFilePath: "./testdata/output/slack-rtm.yml",
		},
		{
			inputFilePath:    "./testdata/input/streetlights1.0.0_no_base_topic.json",
			expectedFilePath: "./testdata/output/streetlights_no_base_topic.yaml",
		},
	}
	for _, test := range tests {
		t.Run(test.inputFilePath, func(t *testing.T) {
			g := NewWithT(t)
			converter, err := New(decode.FromJSONWithYamlFallback, encode.ToYaml, test.options...)
			g.Expect(err).To(BeNil(), "error while creating converter")
			result := convertFile(converter, test.inputFilePath, g)
			expected, err := ioutil.ReadFile(test.expectedFilePath)
			g.Expect(err).To(BeNil(), "error while reading file containing expected results")
			g.Expect(result).To(MatchYAML(string(expected)))
		})
	}
}

func TestNewYamlConverter(t *testing.T) {
	tests := []struct {
		inputFilePath    string
		expectedFilePath string
		options          []ConverterOption
	}{
		{
			inputFilePath:    "./testdata/input/streetlights1.0.0.yaml",
			expectedFilePath: "./testdata/output/streetlights.yaml",
		},
		{
			inputFilePath:    "./testdata/input/streetlights1.1.0.yaml",
			expectedFilePath: "./testdata/output/streetlights.yaml",
		},
		{
			inputFilePath:    "./testdata/input/streetlights1.2.0.yaml",
			expectedFilePath: "./testdata/output/streetlights.yaml",
		},
		{
			inputFilePath:    "./testdata/input/slack-rtm1.2.0.yaml",
			expectedFilePath: "./testdata/output/slack-rtm.yml",
		},
	}
	for _, test := range tests {
		t.Run(test.inputFilePath, func(t *testing.T) {
			g := NewWithT(t)
			converter, err := New(decode.FromJSONWithYamlFallback, encode.ToYaml)
			g.Expect(err).To(BeNil(), "error while creating converter")
			result := convertFile(converter, test.inputFilePath, g)
			expected, err := ioutil.ReadFile(test.expectedFilePath)
			g.Expect(err).To(BeNil(), "error while reading file containing expected results")
			g.Expect(result).To(MatchYAML(string(expected)))
		})
	}
}

func TestConverter_Do_Invalid(t *testing.T) {
	tests := []struct {
		inputFilePath string
	}{
		{
			inputFilePath: "./testdata/input/invalid/gitter-streaming1.2.0_invalid1.json",
		},
		{
			inputFilePath: "./testdata/input/invalid/gitter-streaming1.2.0_invalid2.json",
		},
		{
			inputFilePath: "./testdata/input/invalid/gitter-streaming1.2.0_invalid_version.json",
		},
		{
			inputFilePath: "./testdata/input/invalid/streetlights1.0.0_invalid2.json",
		},
		{
			inputFilePath: "./testdata/input/invalid/streetlights1.0.0_invalid3.json",
		},
		{
			inputFilePath: "./testdata/input/invalid/streetlights1.0.0_invalid4.json",
		},
		{
			inputFilePath: "./testdata/input/invalid/streetlights1.0.0_invalid5.json",
		},
		{
			inputFilePath: "./testdata/input/invalid/slack-rtm1.2.0_invalid1.json",
		},
	}
	for _, test := range tests {
		t.Run(test.inputFilePath, func(t *testing.T) {
			g := NewWithT(t)
			converter, err := New(decode.FromJSON, encode.ToJSON)
			g.Expect(err).To(BeNil(), "error while creating converter")
			_, err = readDataFromFile(converter, test.inputFilePath, g)
			g.Expect(err).Should(HaveOccurred())
		})
	}
}

func getFileReader(filePath string) (io.Reader, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func convertFile(converter Converter, filePath string, g *WithT) string {
	resultWriter, err := readDataFromFile(converter, filePath, g)
	g.Expect(err).To(BeNil(), "error while converting input data")
	return resultWriter.String()
}

func readDataFromFile(converter Converter, filePath string, g *WithT) (*bytes.Buffer, error) {
	resultWriter := bytes.NewBufferString("")
	resultReader, err := getFileReader(filePath)
	g.Expect(err).To(BeNil(), fmt.Sprintf("error while reading file: %s", filePath))
	err = converter.Convert(resultReader, resultWriter)
	return resultWriter, err
}
