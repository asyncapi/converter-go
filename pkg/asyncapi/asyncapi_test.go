package asyncapi

import (
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"

	"testing"
)

var errExpected = errors.New("test error")

func TestBuildUnmarshalWithFallback(t *testing.T) {
	g := NewWithT(t)
	unmarshalFunc := BuildUnmarshalWithFallback(func(_ []byte, _ interface{}) error {
		return errExpected
	})
	var i interface{}
	err := unmarshalFunc([]byte("test"), &i)
	g.Expect(err).Should(HaveOccurred())
}

func TestUnmarshalYaml_err(t *testing.T) {
	g := NewWithT(t)
	var out interface{}
	err := UnmarshalYaml([]byte(","), &out)
	g.Expect(err).Should(HaveOccurred())
}

func TestIsInvalidPropertyErr(t *testing.T) {
	tests := []struct {
		name     string
		error    error
		expected bool
	}{
		{
			name:     "is ErrInvalidProperty",
			error:    ErrInvalidProperty,
			expected: true,
		},
		{
			name:     "is not ErrInvalidProperty",
			error:    errors.New("error"),
			expected: false,
		},
		{
			name:     "is not ErrInvalidProperty",
			expected: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := IsInvalidPropertyErr(test.error)
			NewWithT(t).Expect(actual).To(Equal(test.expected))
		})
	}
}
