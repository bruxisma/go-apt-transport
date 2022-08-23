package transport

import (
	"reflect"
)

// Log (status code 101) messages are used when debugging is enabled.
//
// These messages are ONLY used for debugging.
type Log string

// Status (status code 102) messages are used regarding transportation
// progress.
//
// Status gives a progress indication for the transportation method. It can be
// used to show pre-transfer status for internet enabled transport methods.
type Status string

// Warning (status code 104) messages are used to indicate a warning.
//
// Warnings can be used to alert users to possibly problematic conditions.
type Warning string

// GeneralFailure (status code 401) indicates that some unspecific failure has
// occurred.
//
// This is used when the transport method is no longer able to continue.
// Shortly after sending this, the transport method SHOULD terminate. It is
// intended to for invalid configuration options or other severe conditions.
//
// When using the [transport.Method.SendAndReceive] method, this is
// automatically sent if the [transport.Method.Handler] returns an error.
type GeneralFailure string

func (log Log) MarshalMessage() (*Message, error) {
	fields, err := log.MarshalFields()
	if err != nil {
		return nil, err
	}
	message := &Message{
		StatusCode: StatusCodeLog,
		Summary:    "Log",
		Fields:     fields,
	}
	return message, nil
}

func (log Log) MarshalFields() (Fields, error) {
	return marshalTextField(log)
}

func (status Status) MarshalMessage() (*Message, error) {
	fields, err := status.MarshalFields()
	if err != nil {
		return nil, err
	}
	message := &Message{
		StatusCode: StatusCodeStatus,
		Summary:    "Status",
		Fields:     fields,
	}
	return message, nil
}

func (status Status) MarshalFields() (Fields, error) {
	return marshalTextField(status)
}

func (warning Warning) MarshalMessage() (*Message, error) {
	fields, err := warning.MarshalFields()
	if err != nil {
		return nil, err
	}
	message := &Message{
		StatusCode: StatusCodeWarning,
		Summary:    "Warning",
		Fields:     fields,
	}
	return message, nil
}

func (warning Warning) MarshalFields() (Fields, error) {
	return marshalTextField(warning)
}

func (failure GeneralFailure) MarshalMessage() (*Message, error) {
	fields, err := failure.MarshalFields()
	if err != nil {
		return nil, err
	}
	message := &Message{
		StatusCode: StatusCodeGeneralFailure,
		Summary:    "General Failure",
		Fields:     fields,
	}
	return message, nil
}

func (failure GeneralFailure) MarshalFields() (Fields, error) {
	return marshalTextField(failure)
}

func (failure GeneralFailure) Error() string {
	return string(failure)
}

func (failure GeneralFailure) Is(target error) bool {
	switch target.(type) {
	case *GeneralFailure, GeneralFailure:
		return true
	default:
		return false
	}
}

func marshalTextField[T ~string](message T) (Fields, error) {
	if message == "" {
		return nil, &FieldMarshalerError{
			Type: reflect.TypeOf(message),
			Err:  ErrEmptyInformationalMessage,
		}
	}
	fields := make(Fields)
	fields.Add("message", string(message))
	return fields, nil
}
