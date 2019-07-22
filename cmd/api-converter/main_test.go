package main

import (
	"asyncapi-converter/pkg/asyncapi"

	. "github.com/onsi/gomega"

	"testing"
)

func TestOptionHelper_id_error(t *testing.T) {
	g := NewWithT(t)
	id := newHelper(map[string]interface{}{
		"--id": 123,
	}).id()
	g.Expect(*id).To(Equal("123"))
}

func TestOptionHelper_id_ok(t *testing.T) {
	g := NewWithT(t)
	id := newHelper(map[string]interface{}{}).id()
	g.Expect(id).To(BeNil())
}

func TestOptionHelper_ecodeFormat_err(t *testing.T) {
	g := NewWithT(t)
	_, err := newHelper(map[string]interface{}{
		"--toYAML": "error",
	}).encodeFormat()
	g.Expect(err).Should(HaveOccurred())
}

func TestOptionHelper_ecodeFormat_ok(t *testing.T) {
	tests := []struct {
		name     string
		opts     map[string]interface{}
		expected asyncapi.Format
	}{
		{
			name:     "json",
			opts:     map[string]interface{}{},
			expected: asyncapi.Json,
		},
		{
			name: "json",
			opts: map[string]interface{}{
				"--toYAML": false,
			},
			expected: asyncapi.Json,
		},
		{
			name: "yaml",
			opts: map[string]interface{}{
				"--toYAML": true,
			},
			expected: asyncapi.Yaml,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			g := NewWithT(t)
			format, err := newHelper(test.opts).encodeFormat()
			g.Expect(err).To(BeNil())
			g.Expect(format).To(Equal(test.expected))
		})
	}
}
