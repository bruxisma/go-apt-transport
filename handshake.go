package transport

import (
	"strings"
)

// Capabilities (status code 100) messages are used to inform APT of a transport
// method's feature set.
//
// Displays the capabilities of the transport method. Methods should set the
// pipeline bit if their underlying protocol supports pipelining. The only
// known built-in method that does support pipelining is http(s).
type Capabilities struct {
	SingleInstance bool `transport:"Single-Instance"`
	NeedsCleanup   bool `transport:"Needs-Cleanup"`
	Pipeline       bool
	SendURIEncoded bool `transport:"Send-URI-Encoded"`
	SendConfig     bool `transport:"Send-Config"`
	Removable      bool
	AuxRequests    bool
	PreScan        string `transport:"Pre-Scan"`
	Version        string
}

// Configuration (status code 601) indicates the configuration was sent to the
// method.
//
// APT will send the configuration 'space' to the transport method. A series of
// Config-Items fields are sent, each containing an entry from the APT
// configuration. Each item is placed verbatim into the map directly. No
// canonicalization occurs, as these values come in directly from the APT
// configuration space, and not as field headers.
//
// While files found in `/etc/apt.conf.d/*.conf` will sometimes have several
// settings:
//
//  APT {
//    Install-Recommends "false";
//    Get {
//      Assume-Yes "true";
//    };
//  }
//
// Apt will then convert these into their full namespaced form before the
// transport receives them:
//
//   APT::Install-Recommends "false";
//   APT::Get::Assume-Yes "true";
//
// This is the format that the transport method will receive, and is what users
// should expect to look for. No parsing is done for the user as the content of
// a setting can be, quite literally, anything.
type Configuration map[string]string

// Section returns a subsection of the configuration, based on the given
// prefix. All keys in the returned subsection will have the prefix trimmed.
// Multiple section lookups can be performed by combining sections together
// with "::" (e.g., "APT::Get") This allows taking all parts of a configuration
// and putting it into a smaller lookup:
//
//   confg := cfg.Section("APT::Get")
//   fmt.Println(confg["Assume-Yes"])
//
// This reduces the amount of work and lookup required when a method has its
// own configuration section
func (cfg Configuration) Section(prefix string) Configuration {
	if !strings.HasSuffix(prefix, "::") {
		prefix += "::"
	}
	section := make(Configuration)
	for key, value := range cfg {
		if strings.HasPrefix(key, prefix) {
			section[strings.TrimPrefix(key, prefix)] = value
		}
	}
	return section
}

func (cfg Configuration) UnmarshalFields(fields Fields) error {
	values := fields.Values("Config-Item")
	for _, value := range values {
		parts := strings.SplitN(value, "=", 2)
		if len(parts) != 2 {
			// TODO(bruxisma): log this, so we can figure out if this ever happens, and then return an error.
		}
		key := strings.TrimSpace(parts[0])
		cfg[key] = strings.TrimSpace(parts[1])
	}
	return nil
}
