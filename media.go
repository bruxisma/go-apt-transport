package transport

// MediaFailure (status code 403) indicates new media must be inserted.
//
// This method is executed primarily when the transport method deals with
// multiple media to install packages. This can include resources like disks
// that mounted, FUSE mounts, or any other transparent "media".
type MediaFailure struct {
	Media string
	Drive string
}

// MediaChanged (status code 603) is sent in response to a 403 Media Failure
// message.
//
// This message is sent in response to a 403 Media Failure message. It
// indicates the user has changed media and it is safe to proceed.
type MediaChanged struct {
	Media string
	Fail  string
}
