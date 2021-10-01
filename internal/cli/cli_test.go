package cli

import (
	. "github.com/onsi/gomega"

	"testing"
)

func TestCli_id_error(t *testing.T) {
	g := NewWithT(t)
	id := New(map[string]interface{}{
		idOption: 123,
	}).id()
	g.Expect(*id).To(Equal("123"))
}

func TestCli_id_ok(t *testing.T) {
	g := NewWithT(t)
	id := New(map[string]interface{}{}).id()
	g.Expect(id).To(BeNil())
}

func TestCli_encode_err(t *testing.T) {
	g := NewWithT(t)
	_, err := New(map[string]interface{}{
		encodeOptionYAML: "error",
	}).encode()
	g.Expect(err).Should(HaveOccurred())
}

func TestCli_encode_ToJSON(t *testing.T) {
	g := NewWithT(t)
	_, err := New(map[string]interface{}{}).encode()
	g.Expect(err).ShouldNot(HaveOccurred())
}

func TestCli_encode_ToYaml_true(t *testing.T) {
	g := NewWithT(t)
	_, err := New(map[string]interface{}{
		encodeOptionYAML: true,
	}).encode()
	g.Expect(err).ShouldNot(HaveOccurred())
}

func TestCli_encode_ToYaml_false(t *testing.T) {
	g := NewWithT(t)
	_, err := New(map[string]interface{}{
		encodeOptionYAML: false,
	}).encode()
	g.Expect(err).ShouldNot(HaveOccurred())
}

func TestCli_reader_error_no_path(t *testing.T) {
	g := NewWithT(t)
	_, err := New(map[string]interface{}{}).reader()
	g.Expect(err).Should(HaveOccurred())
}

func TestCli_reader_http_error(t *testing.T) {
	g := NewWithT(t)
	_, err := New(map[string]interface{}{
		fileOptionPath: "http://atest",
	}).reader()
	g.Expect(err).Should(HaveOccurred())
}

func TestCli_reader_file_error(t *testing.T) {
	g := NewWithT(t)
	_, err := New(map[string]interface{}{
		fileOptionPath: "/invalid/path/to/a/file",
	}).reader()
	g.Expect(err).Should(HaveOccurred())
}

func TestIsUrl(t *testing.T) {
	tests := []struct {
		url   string
		isURL bool
	}{
		{"/test/me", false},
		{"http://test.it", true},
		{"aaa", false},
	}
	for _, test := range tests {
		t.Run(test.url, func(t *testing.T) {
			g := NewWithT(t)
			g.Expect(isURL(test.url)).To(Equal(test.isURL))
		})
	}
}
