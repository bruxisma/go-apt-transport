package transport

import (
	"bytes"
	"fmt"
)

const (
	StatusCodeCapabilities             = 100
	StatusCodeLog                      = 101
	StatusCodeStatus                   = 102
	StatusCodeRedirect                 = 103
	StatusCodeWarning                  = 104
	StatusCodeURIStart                 = 200
	StatusCodeURIDone                  = 201
	StatusCodeAuxRequest               = 351
	StatusCodeURIFailure               = 400
	StatusCodeGeneralFailure           = 401
	StatusCodeAuthorizationRequired    = 402
	StatusCodeMediaFailure             = 403
	StatusCodeURIAcquire               = 600
	StatusCodeConfiguration            = 601
	StatusCodeAuthorizationCredentials = 602
	StatusCodeMediaChanged             = 603
)

// Message represents the low level data either sent or received via the APT
// transport method.
//
// Messages are very close to the same messages found in HTTP and other text
// protocols that are based off of RFC822. However, no attempt is made to
// enforce a carriage return ("\r").
//
// A raw Message makes no attempt to validate the StatusCode, Summary, or
// [Fields]. This is instead handled by higher level APIs within this library,
// and users are encouraged to use them over raw Messages.
//
// Messages can be marshaled to and from Binary data, as they have a
// well-formed "wire format".
type Message struct {
	StatusCode int    // e.g., 100
	Summary    string // e.g., Capabilities
	Fields     Fields // e.g., {"Send-Config": "true"}
}

type MessageMarshaler interface {
	// MarshalMessage encodes the receiver into a Message and returns the result.
	MarshalMessage() (*Message, error)
}

type MessageUnmarshaler interface {
	// UnmarshalMessages decodes a Message into the receiver.
	UnmarshalMessage(*Message) error
}

// StatusText returns a text for the APT status code. It returns the empty
// string if the code is unknown.
func StatusText(code int) string {
	switch code {
	case StatusCodeCapabilities:
		return "100 Capabilities"
	case StatusCodeLog:
		return "101 Log"
	case StatusCodeStatus:
		return "102 Status"
	case StatusCodeRedirect:
		return "103 Redirect"
	case StatusCodeWarning:
		return "104 Warning"
	case StatusCodeURIStart:
		return "200 URI Start"
	case StatusCodeURIDone:
		return "201 URI Done"
	case StatusCodeAuxRequest:
		return "351 Aux Request"
	case StatusCodeURIFailure:
		return "400 URI Failure"
	case StatusCodeGeneralFailure:
		return "401 General Failure"
	case StatusCodeAuthorizationRequired:
		return "402 Authorization Required"
	case StatusCodeMediaFailure:
		return "403 Media Failure"
	case StatusCodeURIAcquire:
		return "600 URI Acquire"
	case StatusCodeConfiguration:
		return "601 Configuration"
	case StatusCodeAuthorizationCredentials:
		return "602 Authorization Credentials"
	case StatusCodeMediaChanged:
		return "603 Media Changed"
	}
	return ""
}

// todo: support passing in an `error` and converting it to a [Message]
func MarshalMessage(value any) (*Message, error) {
	return nil, ErrNotImplemented
}

func UnmarshalMessage(message *Message, destination any) error {
	return ErrNotImplemented
}

// MarshalBinary serializes the receiving Message into a byte slice.
func (message *Message) MarshalBinary() ([]byte, error) {
	var buffer bytes.Buffer
	fmt.Fprintf(&buffer, "%03d %s\n", message.StatusCode, message.Summary)
	if err := message.Fields.Write(&buffer); err != nil {
		return nil, err
	}
	fmt.Fprintf(&buffer, "\n")
	return buffer.Bytes(), nil
}

// UnmarshalBinary deserializes the receiving byte slice into a [Message].
//
// This function does not perform any validation on the message's contents,
// only that it meets the correct layout.
//
// This means it is possible to receive correctly formatted but ultimately
// invalid messages.
//
// This function is dependent on the behavior of Unmarshalling a [Fields]
// object.
func (message *Message) UnmarshalBinary(data []byte) error {
	before, after, found := bytes.Cut(data, []byte("\n"))
	if !found {
		return ErrMessageHeaderNotFound
	}
	items := bytes.SplitN(before, []byte(" "), 2)
	if len(items) != 2 {
		return fmt.Errorf("%w %q", ErrMessageHeaderMalformed, string(before))
	}
	if message.Fields == nil {
		message.Fields = make(Fields)
	}
	return message.Fields.UnmarshalBinary(after)
}

func (message *Message) IsInformational() bool {
	return message.StatusCode >= 100 && message.StatusCode < 200
}

func (message *Message) IsSuccessful() bool {
	return message.StatusCode >= 200 && message.StatusCode < 300
}

func (message *Message) IsFailure() bool {
	return message.StatusCode >= 400 && message.StatusCode < 500
}

func (message *Message) IsResponse() bool {
	return message.StatusCode >= 600 && message.StatusCode < 700
}
