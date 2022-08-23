package transport

import (
	"fmt"
	"io"
)

// MessageWriter is used to send additional messages back to the consumer.
//
// These messages are sent immediately once called, and can result in a handler
// being cancelled if an error is sent.
type MessageWriter struct {
	inner io.Writer
}

func NewMessageWriter(writer io.Writer) *MessageWriter {
	return &MessageWriter{writer}
}

// Configuration returns a copy of configuration sent to the Method from APT.
func (writer *MessageWriter) Configuration() Configuration {
	return Configuration{}
}

// Write attempts to marshal the provided message into a binary wire format,
// and then write it all at once to the underlying writer.
func (writer *MessageWriter) Write(message *Message) error {
	data, err := message.MarshalBinary()
	if err != nil {
		return err
	}
	writer.inner.Write(data)
	return nil
}

// Writes a [transport.Warning] message to the communication stream.
func (writer *MessageWriter) Warning(message string) {
	// we know this won't actually error.
	value, _ := Warning(message).MarshalMessage()
	writer.Write(value)
}

// Writes a [transport.Status] message to the communication stream.
func (writer *MessageWriter) Status(message string) {
	// we know this won't actually error.
	value, _ := Status(message).MarshalMessage()
	writer.Write(value)
}

// Print is an alias for [MessageWriter.Log]
func (writer *MessageWriter) Print(message string) {
	writer.Log(message)
}

// Debug is an alias for [MessageWriter.Log]
func (writer *MessageWriter) Debug(message string) {
	writer.Log(message)
}

// Writes a [transport.Log] message to the communication stream.
func (writer *MessageWriter) Log(message string) {
	// we know this won't actually error.
	value, _ := Log(message).MarshalMessage()
	writer.Write(value)
}

// Warningf writes a [transport.Warning] message to the communication stream,
// using the provided format specifier.
func (writer *MessageWriter) Warningf(format string, args ...any) {
	writer.Warning(fmt.Sprintf(format, args...))
}

// Statusf writes a [transport.Status] message to the communication stream,
// using the provided format specifier.
func (writer *MessageWriter) Statusf(format string, args ...any) {
	writer.Status(fmt.Sprintf(format, args...))
}

// Printf is an alias for [MessageWriter.Logf]
func (writer *MessageWriter) Printf(format string, args ...any) {
	writer.Logf(format, args...)
}

// Debugf is an alias for [MessageWriter.Logf]
func (writer *MessageWriter) Debugf(format string, args ...any) {
	writer.Logf(format, args...)
}

// Logf writes a [transport.Log] message to the communication stream, using the
// provided format specifier.
func (writer *MessageWriter) Logf(format string, args ...any) {
	writer.Log(fmt.Sprintf(format, args...))
}
