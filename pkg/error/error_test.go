package error

import (
	. "github.com/onsi/gomega"

	"errors"
	"testing"
)

func TestIsInvalidPropertyErr(t *testing.T) {
	err := errors.New("test error")
	tests := []struct {
		name     string
		error    error
		expected bool
	}{
		{
			name:     "ErrInvalidProperty",
			error:    NewInvalidProperty("test"),
			expected: true,
		},
		{
			name:     "ErrInvalidDocument",
			error:    NewInvalidDocument(),
			expected: false,
		},
		{
			name:     "ErrUnsupportedAsyncapiVersion",
			error:    NewUnsupportedAsyncapiVersion("test"),
			expected: false,
		},
		{
			name:     "random err",
			error:    err,
			expected: false,
		},
		{
			name:     "DocumentVersionUpToDate",
			error:    NewDocumentVersionUpToDate("test"),
			expected: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			g := NewWithT(t)
			g.Expect(IsInvalidProperty(test.error)).To(Equal(test.expected))
			if actual, ok := test.error.(Error); ok {
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
			name:     "ErrInvalidProperty",
			error:    NewInvalidProperty("test"),
			expected: false,
		},
		{
			name:     "ErrInvalidDocument",
			error:    NewInvalidDocument(),
			expected: true,
		},
		{
			name:     "ErrUnsupportedAsyncapiVersion",
			error:    NewUnsupportedAsyncapiVersion("test"),
			expected: false,
		},
		{
			name:     "DocumentVersionUpToDate",
			error:    NewDocumentVersionUpToDate("test"),
			expected: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			g := NewWithT(t)
			g.Expect(IsInvalidDocument(test.error)).To(Equal(test.expected))
			if actual, ok := test.error.(Error); ok {
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
			name:     "ErrInvalidProperty",
			error:    NewInvalidProperty("test"),
			expected: false,
		},
		{
			name:     "ErrInvalidDocument",
			error:    NewInvalidDocument(),
			expected: false,
		},
		{
			name:     "ErrUnsupportedAsyncapiVersion",
			error:    NewUnsupportedAsyncapiVersion("test"),
			expected: true,
		},
		{
			name:     "DocumentVersionUpToDate",
			error:    NewDocumentVersionUpToDate("test"),
			expected: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			g := NewWithT(t)
			g.Expect(IsUnsupportedAsyncapiVersion(test.error)).To(Equal(test.expected))
			if actual, ok := test.error.(Error); ok {
				g.Expect(actual.Error()).ToNot(BeEmpty())
			}
		})
	}
}

func TestDocumentVersionUpToDateErr(t *testing.T) {
	tests := []struct {
		name     string
		error    error
		expected bool
	}{
		{
			name:     "ErrInvalidProperty",
			error:    NewInvalidProperty("test"),
			expected: false,
		},
		{
			name:     "ErrInvalidDocument",
			error:    NewInvalidDocument(),
			expected: false,
		},
		{
			name:     "ErrUnsupportedAsyncapiVersion",
			error:    NewUnsupportedAsyncapiVersion("test"),
			expected: false,
		},
		{
			name:     "DocumentVersionUpToDate",
			error:    NewDocumentVersionUpToDate("test"),
			expected: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			g := NewWithT(t)
			g.Expect(IsDocumentVersionUpToDate(test.error)).To(Equal(test.expected))
			if actual, ok := test.error.(Error); ok {
				g.Expect(actual.Error()).ToNot(BeEmpty())
			}
		})
	}
}
