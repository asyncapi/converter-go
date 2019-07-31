package error

import "fmt"

type errType = int

const (
	errInvalidProperty errType = iota + 1
	errInvalidDocument
	errUnsupportedAsyncapiVersion
)

// Error represents the conversion error.
type Error struct {
	errType
	msg string
}

func (err Error) Error() string {
	return err.msg
}

func isErrorType(errType errType, err error) bool {
	if err, ok := err.(Error); ok {
		return err.errType == errType
	}
	return false
}

// IsInvalidProperty returns true if err is the InvalidProperty error,
// otherwise it returns false.
//
// See NewInvalidProperty.
func IsInvalidProperty(err error) bool {
	return isErrorType(errInvalidProperty, err)
}

// IsInvalidDocument returns true if err is the InvalidDocument error,
// otherwise it returns false.
//
// See NewInvalidDocument.
func IsInvalidDocument(err error) bool {
	return isErrorType(errInvalidDocument, err)
}

// IsUnsupportedAsyncapiVersion returns true if err is the UnsupportedAsyncapiVersion error,
// otherwise it returns false.
//
// See UnsupportedAsyncapiVersion.
func IsUnsupportedAsyncapiVersion(err error) bool {
	return isErrorType(errUnsupportedAsyncapiVersion, err)
}

func newError(errType errType, msg string) Error {
	return Error{
		errType: errType,
		msg:     msg,
	}
}

// NewInvalidProperty creates new invalid property error.
// This error is returned by the AsynAPI Converter when one of the document
// properties is invalid or missing.
func NewInvalidProperty(context interface{}) Error {
	msg := fmt.Sprintf("asyncapi: error invalid property %v", context)
	return newError(errInvalidProperty, msg)
}

// NewInvalidDocument creates a new invalid document error.
// This error is returned by the AsyncAPI Converter when a document has an invalid structure.
func NewInvalidDocument() Error {
	return newError(errInvalidDocument, "asyncapi: unable to decode document")
}

// NewUnsupportedAsyncapiVersion creates a new unsupported AsyncAPI version error.
// This error is returned when the AsyncAPI converter does not recognize the version of the
// converted AsyncAPI document.
func NewUnsupportedAsyncapiVersion(context interface{}) Error {
	msg := fmt.Sprintf("asyncapi: unsupported asyncapi version '%v'", context)
	return newError(errUnsupportedAsyncapiVersion, msg)
}
