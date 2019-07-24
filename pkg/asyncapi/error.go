package asyncapi

import "fmt"

type errType = int

const (
	errInvalidProperty errType = iota + 1
	errInvalidDocument
	errUnsupportedAsyncapiVersion
)

var (
	ErrInvalidDocument = &asyncapiError{
		errType: errInvalidDocument,
	}
	ErrUnsupportedAsyncapiVersion = &asyncapiError{
		errType: errUnsupportedAsyncapiVersion,
	}
)

type asyncapiError struct {
	errType
	context interface{}
}

type Error interface {
	error
	InvalidProperty() bool
	InvalidDocument() bool
	UnsupportedAsyncapiVersion() bool
}

func (err asyncapiError) Error() string {
	switch err.errType {
	case errInvalidProperty:
		return fmt.Sprintf("asyncapi: error invalid property %v", err.context)
	case errInvalidDocument:
		return "unable to decode document"
	default:
		return "unsupported asyncapi version"
	}
}

func (err asyncapiError) InvalidProperty() bool {
	return err.errType == errInvalidProperty
}

func (err asyncapiError) InvalidDocument() bool {
	return err.errType == errInvalidDocument
}

func (err asyncapiError) UnsupportedAsyncapiVersion() bool {
	return err.errType == errUnsupportedAsyncapiVersion
}

func NewErrInvalidProperty(context interface{}) Error {
	return &asyncapiError{
		errType: errInvalidProperty,
		context: context,
	}
}
