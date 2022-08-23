package transport

import (
	"io"
	"os"
)

// Stream acts as an io.ReadWriteCloser, that operates on stdin (for reading),
// and stdout (for writing). There are several ways to open a stream, the
// simplest of which is to use the default os.Stdin and os.Stdout pointers.
//
// This type is only thread-safe if the underlying io.Writer and io.Reader are
// thread-safe.
//
// Users must call close on the underlying io.Writer or io.Reader.
type Stream struct {
	output io.Writer
	input  io.Reader
}

// NewStream initializes a new Stream with os.Stdin and os.Stdout.
//
// This function does *not* check if the files are actually open and valid.
func NewStream() *Stream {
	return &Stream{
		output: os.Stdout,
		input:  os.Stdin,
	}
}

// NewStreamWith initializes a new Stream with the given input and output.
//
// This function does *not* check if the provided interfaces are not nil.
func NewStreamWith(input io.Reader, output io.Writer) *Stream {
	return &Stream{
		output: output,
		input:  input,
	}
}

func (stream *Stream) Read(data []byte) (int, error) {
	return stream.input.Read(data)
}

func (stream *Stream) Write(data []byte) (int, error) {
	return stream.output.Write(data)
}

// Close is a no-op, as we let the operating system take care of closing
// os.Stdout and os.Stdin. We also make no effort to know whether we can close
// the underlying io.Writer and io.Reader
func (stream *Stream) Close() error {
	return nil
}
