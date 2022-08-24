package transport

import (
	"fmt"
	"time"
)

// URIStart (status code 200) indicates that the resource located at the given
// URI is to begin transferring.
//
// The URI is specified along with stats regarding the file itself.
type URIStart struct {
	LastModified string `transport:"Last-Modified"`
	ResumePoint  string `transport:"Resume-Point"`
	URI          string
	Size         int64
}

// URIDone (status code 201) indicates that a URI has completed transferrence.
//
// It is possible to specify a 201 URI Done without a URI Start which would
// mean no data was transferred, but the file is now available. A
// [transport.URIDone.Filename] field is specified when the URI is directly
// available in the local pathname space. APT will either directly use that
// file or copy it into another location. It is possible to return fields
// prefix with Alt- to indicate that another possible for the URI has been
// found in the local pathname space. This is done if a decompressed version of
// a gunzip file is found.
//
// BUG(bruxisma): We do not currently support the Alt- prefixed fields.
type URIDone struct {
	URI          string
	LastModified string `transport:"Last-Modified"`
	IMSHit       string `transport:"IMS-Hit"`
	Filename     string
	MD5Hash      string `transport:"MD5-Hash"`
	Size         int64
}

// URIFailure (status code 400) indicates the URI is not retrievable from this
// source.
//
// Indicates a fatal URI failure. As with 201 URI Done, 200 URI start is not
// required to precede this message.
type URIFailure struct {
	URI     string
	Message string
}

// URIAcquire (status code 600) indicates that APT is requesting a new URI be
// added to the acquire list.
//
// The deserialized [transport.URIAcquire.LastModified] field has the time
// stamp of the current cache file if applicable.
// [transport.URIAcquire.Filename] is the name of the file that the acquired
// URI should be written to. It is safe for the method to assume it has correct
// write permissions.
//
// NOTE(bruxisma): This message is effectively "repeated" by the
// [transport.Request] type passed to Method's Handler.
type URIAcquire struct {
	LastModified *time.Time `transport:"Last-Modified"`
	URI          string
	Filename     string
}

func (failure *URIFailure) Error() string {
	return fmt.Sprintf("failure acquiring uri %q: %s", failure.URI, failure.Message)
}

func (failure *URIFailure) Is(target error) bool {
	_, ok := target.(*URIFailure)
	return ok
}
