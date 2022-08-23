package transport

import "net/url"

// Redirect (status code 103) is currently undocumented by the APT transport
// method protocol.
type Redirect struct {
	URI        *url.URL
	NewURI     *url.URL `transport:"New-URI"`
	AltURIs    *url.URL `transport:"Alt-URIs"`
	UsedMirror bool     `transport:"Used-Mirror"`
}

// AuxRequest (status code 351) indicates a request for an auxiliary file to be
// downloaded by the acquire system (via another method) and made available for
// the requesting method.
//
// The requestor will get a 600 URI Acquire with the URI it requested and the
// filename will either be an existing file if the request was a success or if
// the acquire failed for the some reason the file will not exist.
type AuxRequest struct {
	MaximumSize int64  `transport:"MaximumSize"`
	ShortDesc   string `transport:"Aux-ShortDesc"`
	Description string `transport:"Aux-Description"`
	URI         string `transport:"Aux-URI"`
}
