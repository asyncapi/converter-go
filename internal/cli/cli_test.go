package cli

import (
	. "github.com/onsi/gomega"

	"testing"
)

func TestCli_id_error(t *testing.T) {
	g := NewWithT(t)
	id := New(map[string]interface{}{
		"--id": 123,
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
		"--toYAML": "error",
	}).encode()
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
