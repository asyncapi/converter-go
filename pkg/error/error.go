package error

import "fmt"

type errType = int

const (
	errInvalidProperty errType = iota + 1
	errInvalidDocument
	errUnsupportedAsyncapiVersion
)

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

func IsInvalidProperty(err error) bool {
	return isErrorType(errInvalidProperty, err)
}

func IsInvalidDocument(err error) bool {
	return isErrorType(errInvalidDocument, err)
}

func IsUnsupportedAsyncapiVersion(err error) bool {
	return isErrorType(errUnsupportedAsyncapiVersion, err)
}

func newError(errType errType, msg string) Error {
	return Error{
		errType: errType,
		msg:     msg,
	}
}

func NewInvalidProperty(context interface{}) Error {
	msg := fmt.Sprintf("asyncapi: error invalid property %v", context)
	return newError(errInvalidProperty, msg)
}

func NewInvalidDocument() Error {
	return newError(errInvalidDocument, "asyncapi: unable to decode document")
}

func NewUnsupportedAsyncapiVersion(context interface{}) Error {
	msg := fmt.Sprintf("asyncapi: unsupported asyncapi version '%v'", context)
	return newError(errUnsupportedAsyncapiVersion, msg)
}
