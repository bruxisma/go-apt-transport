package transport

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	ErrDestinationNotPointer = errors.New("destination for value is not a pointer")
	ErrDestinationNotStruct  = errors.New("destination for value is not a struct")
	ErrDestinationIsNil      = errors.New("destination for value is nil")

	ErrSourceNotPointer = errors.New("source for value is not a pointer")
	ErrSourceNotStruct  = errors.New("source for value is not a struct")
	ErrSourceIsNil      = errors.New("source for value is nil")

	ErrMessageScannerNoData = errors.New("message scanner has not data")

	ErrMessageHeaderNotFound  = errors.New("message header not found")
	ErrMessageHeaderMalformed = errors.New("message header malformed")

	ErrFieldEntryInvalid = errors.New("header field entry is invalid")
	ErrFieldsEmpty       = errors.New("header fields are empty")

	ErrEmptyInformationalMessage = errors.New("informational message is empty")

	ErrNotImplemented = errors.New("not implemented")
)

// MessageMarshalerError is used when performing automatic reflection-based
// marhsalling into a message.
type MessageMarshalerError struct {
	Type   reflect.Type
	Err    error
	source string
}

// FieldMarshalerError is used when performing automatic reflection-based
// marshalling into a field.
type FieldMarshalerError struct {
	Type   reflect.Type
	Err    error
	source string
}

func (err *FieldMarshalerError) Error() string {
	source := err.source
	if source == "" {
		source = "MarshalFields"
	}
	return fmt.Sprintf(
		"apt/transport: error calling %q for type %q: %s",
		source,
		err.Type.String(),
		err.Err.Error())
}

func (err *FieldMarshalerError) Unwrap() error {
	return err.Err
}

func (err *MessageMarshalerError) Error() string {
	source := err.source
	if source == "" {
		source = "MarshalMessage"
	}
	return fmt.Sprintf(
		"apt/transport: error calling %q for type %q: %s",
		source,
		err.Type.String(),
		err.Err.Error())
}

func (err *MessageMarshalerError) Unwrap() error {
	return err.Err
}
