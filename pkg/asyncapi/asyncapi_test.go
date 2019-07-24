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
			name:     "ErrInvalidProperty",
			error:    NewErrInvalidProperty("test"),
			expected: true,
		},
		{
			name:     "ErrInvalidDocument",
			error:    ErrInvalidDocument,
			expected: false,
		},
		{
			name:     "ErrUnsupportedAsyncapiVersion",
			error:    ErrUnsupportedAsyncapiVersion,
			expected: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			g := NewWithT(t)
			if actual, ok := test.error.(Error); ok {
				g.Expect(actual.InvalidProperty()).To(Equal(test.expected))
				g.Expect(actual.Error()).ToNot(BeEmpty())
			}
		})
	}
}

func TestIsErrInvalidDocumentErr(t *testing.T) {
	tests := []struct {
		name     string
		error    error
		expected bool
	}{
		{
			name:     "ErrInvalidDocument",
			error:    ErrInvalidDocument,
			expected: true,
		},
		{
			name:     "ErrInvalidProperty",
			error:    NewErrInvalidProperty("test"),
			expected: false,
		},
		{
			name:     "ErrUnsupportedAsyncapiVersion",
			error:    ErrUnsupportedAsyncapiVersion,
			expected: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			g := NewWithT(t)
			if actual, ok := test.error.(Error); ok {
				g.Expect(actual.InvalidDocument()).To(Equal(test.expected))
				g.Expect(actual.Error()).ToNot(BeEmpty())
			}
		})
	}
}

func TestIsErrUnsupportedAsyncapiVersionErr(t *testing.T) {
	tests := []struct {
		name     string
		error    error
		expected bool
	}{
		{
			name:     "ErrUnsupportedAsyncapiVersion",
			error:    ErrUnsupportedAsyncapiVersion,
			expected: true,
		},
		{
			name:     "ErrInvalidProperty",
			error:    NewErrInvalidProperty("test"),
			expected: false,
		},
		{
			name:     "ErrInvalidDocument",
			error:    ErrInvalidDocument,
			expected: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			g := NewWithT(t)
			if actual, ok := test.error.(Error); ok {
				g.Expect(actual.UnsupportedAsyncapiVersion()).To(Equal(test.expected))
				g.Expect(actual.Error()).ToNot(BeEmpty())
			}
		})
	}
}
