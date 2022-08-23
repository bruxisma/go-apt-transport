package transport

import (
	"bufio"
	"bytes"
	"io"
)

// MessageScanner provides a convenient interface for reading messages from any
// io.Reader. Successive calls to Scan() will step through the the reader,
// skipping the empty newline between messages. The internal function used to
// split these messages up is ScanMessages. This internal split funciton cannot
// be overriden at this time.
//
// Scanning stops unrecoverably at EOF, the first I/O error encountered, or
// when data is too large to fit into the internal buffer. Much like
// bufio.Scanner, if more control over scanning messages is required (however
// unlikely), it is recommended that users utilize a bufio.Reader in
// conjunction with ScanMessages.
//
// NOTE(bruxisma): It is unlikely that more granular control scanning is
// required by users, as the input for messages comes from os.Stdin.
// Additionally, rescanning os.Stdin is out of scope for this library, and
// indicates that "something wacky" has gone awry with the APT transport method
// protocol.
type MessageScanner struct {
	inner *bufio.Scanner
}

// NewMessageScanner will initialize a bufio.Scanner internally, call
// Scanner.Split with ScanMessages and then return. This ensures that the order
// of operations does not result in a panic when scanning.
func NewMessageScanner(reader io.Reader) *MessageScanner {
	scanner := bufio.NewScanner(reader)
	scanner.Split(ScanMessages)
	return &MessageScanner{scanner}
}

// Scan advances the MessageScanner to the next message, which will then be
// available through the Message method. It returns false when the scan has
// stopped, either by reaching the end of the input or an error. After Scan
// returns false, the Err method will return any error that occurred during
// scanning, unless it was io.EOF.
//
// NOTE(bruxisma): Scan panics if the  split function returns empty slices
// without advancing the input. This is a side effect of uing bufio.Scanner
// internally.
func (scanner *MessageScanner) Scan() bool {
	return scanner.inner.Scan()
}

// Err returns the first non-EOF error that was encountered by the
// MessageScanner.
func (scanner *MessageScanner) Err() error {
	return scanner.inner.Err()
}

// Message returns the most recent Message when scanning, or an error. The
// value returned will NOT be overwritten by subsequent calls to Scan. However,
// this is done at the cost of an allocation, as the internal bytes buffer is
// unmarshalled into the returned Message.
//
// If Scan has not been called, an error is returned.
func (scanner *MessageScanner) Message() (*Message, error) {
	message := &Message{}
	data := scanner.inner.Bytes()
	err := scanner.inner.Err()
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, ErrMessageScannerNoData
	}
	if err := message.UnmarshalBinary(data); err != nil {
		return nil, err
	}
	return message, nil
}

// ScanMessages is a SplitFunc function for bufio.Scanner that returns the
// entire set of data starting with a status code and ending in a double
// newline.
//
// Unlike bufio.ScanLines, this function will error if the last line of input
// is not a newline.
func ScanMessages(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	// messages end with two newlines. We remove the empty one so that message
	// deserialization is easier.
	if idx := bytes.Index(data, []byte("\n\n")); idx >= 0 {
		return idx + 2, data[0:idx], nil
	}
	if atEOF {
		return 0, nil, io.ErrUnexpectedEOF
	}
	// request more data
	return 0, nil, nil
}
