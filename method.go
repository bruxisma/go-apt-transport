package transport

import (
	"context"
	"net/url"
	"os"
	"time"
)

// A Handler respons to a URI Acquire message.
type Handler interface {
	AcquireResource(*MessageWriter, *Request) error
}

type Request struct {
	Modified time.Time `transport:"Last-Modified"`
	Source   *url.URL  `tranport:"URI"`
	Target   string    `transport:"Filename"`
}

type HandlerFunc func(*MessageWriter, *Request) error
type MethodOption func(*Method) error

type Method struct {
	stream       *Stream
	capabilities Capabilities
	requests     chan *Request
	ctx          context.Context
	Handler      Handler
}

// WithCapabilities sets the capabilities of the Method.
//
// NOTE: The SendConfig and Version fields are always overwritten to be true,
// and the version that was passed in to the call to NewMethod.
func WithCapabilities(capabilities Capabilities) MethodOption {
	return func(method *Method) error {
		method.capabilities = capabilities
		return nil
	}
}

// WithHandlerFunction sets the Handler for the Method to the provided function
func WithHandlerFunction(function func(*MessageWriter, *Request) error) MethodOption {
	return func(method *Method) error {
		method.Handler = HandlerFunc(function)
		return nil
	}
}

// WithStream allows overwriting the default IO stream for the Method.
func WithStream(stream *Stream) MethodOption {
	return func(method *Method) error {
		method.stream = stream
		return nil
	}
}

// WithHandler sets the Handler for the Method
func WithHandler(handler Handler) MethodOption {
	return func(method *Method) error {
		method.Handler = handler
		return nil
	}
}

func NewMethod(ctx context.Context, version string, options ...MethodOption) (*Method, error) {
	method := &Method{
		stream:   NewStream(),
		requests: make(chan *Request),
		ctx:      ctx,
	}
	for _, option := range options {
		if err := option(method); err != nil {
			return nil, err
		}
	}
	// These are ALWAYS set.
	method.capabilities.SendConfig = true
	method.capabilities.Version = version
	// TODO(bruxisma): Ensure that method.Handler is not nil
	return method, nil
}

// SendAndReceive is the Method's main loop, and can be considered equivalent
// to http.Server's ListenAndServe.
//
// This function will perform the initial handshake, launch an event queue, and
// then block on stdin, until it is closed, cannot be read from any longer, or
// the context is cancelled.
//
// NOTE(bruxisma): The use of context.Context is currently barebones, and will
// most likely improve over time.
func (method *Method) SendAndReceive() error {
	// TODO(bruxisma): Should we create a "root" span here?
	scanner := NewMessageScanner(os.Stdin)
	if err := method.handshake(); err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(method.ctx)
	defer cancel()
	go method.handleRequests(ctx)
	for scanner.Scan() {
		message, err := scanner.Message()
		if err != nil {
			return err
		}
		switch message.StatusCode {
		case 600:
			request := &Request{}
			if err := UnmarshalMessage(message, request); err != nil {
				return err
			}
			method.requests <- request
		}
	}
	return scanner.Err()
}

func (method *Method) handleRequests(ctx context.Context) {
	// TODO(bruxisma): Add a span here, to at least act as a "root" span.
	for {
		select {
		case request := <-method.requests:
			writer := NewMessageWriter(method.stream)
			go func() {
				// TODO(bruxisma): pass a span or span context into the Handler
				// AcquireResource call (maybe adjust the signature?).
				// TODO(bruxisma): media failure means we need to pause all other acquire
				// resource calls until we are unblocked. We will need to do some work with a
				// sync.WaitGroup, but have it so that anything that returns a media failure
				// dynamically becomes the controller, and all other handlers are paused.
				// TODO(bruxisma): When authorization credentials are needed, a
				// condition should be used *somehow* to allow us to then signal a
				// goroutine to resume (and read from) data to allow for the
				// authorization process to continue.
				//				_, span := tracer.Start(ctx, "AcquireResource")
				//				defer span.End()
				//				setSpanRequest(span, request)
				err := method.Handler.AcquireResource(writer, request)
				if err == nil {
					return
				}
				//				span.SetStatus(codes.Error, err.Error())
				//				span.RecordError(err)
				message, err := MarshalMessage(err)
				if err != nil {
					return
				}
				//				_, span = tracer.Start(ctx, "MessageWriter.Write")
				//				defer span.End()
				err = writer.Write(message)
				if err != nil {
					//					span.RecordError(err)
					//					span.SetStatus(codes.Error, err.Error())
				}
			}()
		case <-ctx.Done():
			// TODO(bruxisma): Handle ctx.Err() here
			return
		}
	}
}

func (method *Method) handshake() error {
	// TODO(bruxisma): Add a span here for the handshake.
	writer := NewMessageWriter(method.stream)
	// We don't bother using a MessageWriter here.
	message, err := MarshalMessage(&method.capabilities)
	if err != nil {
		return err
	}
	return writer.Write(message)
	// TODO(bruxisma): We need to receive the configuration from APT.
}

func (handler HandlerFunc) AcquireResource(writer *MessageWriter, request *Request) error {
	// TODO(bruxisma): Fire off a URI Start here.
	err := handler(writer, request)
	if err != nil {
		return err
	}
	// if writer.IsDone() { send URI Done here }
	// return URI failure here if we reach this point.
	return nil
}
