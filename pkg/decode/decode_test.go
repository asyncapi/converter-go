package decode

import (
	. "github.com/onsi/gomega"
	"strings"

	"io"
	"testing"
)

type testReader func(p []byte) (n int, err error)

func (reader testReader) Read(p []byte) (int, error) {
	return reader(p)
}

var errNoProgressReader testReader = func(_ []byte) (int, error) {
	return 0, io.ErrNoProgress
}

func TestUnmarshalYamlError(t *testing.T) {
	g := NewWithT(t)
	err := FromYaml(nil, errNoProgressReader)
	g.Expect(err).Should(HaveOccurred())
}

func TestUnmarshalYamlError2(t *testing.T) {
	g := NewWithT(t)
	reader := strings.NewReader(",")
	var out interface{}
	err := FromYaml(&out, reader)
	g.Expect(err).Should(HaveOccurred())
}

func TestUnmarshalYaml(t *testing.T) {
	g := NewWithT(t)
	reader := strings.NewReader("test: me")
	var out interface{}
	err := FromYaml(&out, reader)
	g.Expect(err).To(BeNil())
	expected := map[string]interface{}{
		"test": "me",
	}
	g.Expect(out).To(Equal(expected))
}

func TestUnmarshalYamlReaderError(t *testing.T) {
	g := NewWithT(t)
	err := FromJSONWithYamlFallback(nil, errNoProgressReader)
	g.Expect(err).Should(HaveOccurred())
}

func TestUnmarshalYamlReaderError2(t *testing.T) {
	g := NewWithT(t)
	reader := strings.NewReader(",")
	var out interface{}
	err := FromJSONWithYamlFallback(&out, reader)
	g.Expect(err).Should(HaveOccurred())
}

func TestUnmarshalYaml_err(t *testing.T) {
	g := NewWithT(t)
	var out interface{}
	err := unmarshalYaml([]byte(","), &out)
	g.Expect(err).Should(HaveOccurred())
}
