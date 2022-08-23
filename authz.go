package transport

// AuthorizationRequired (status code 402) is sent to APT to request credentials.
//
// The transport method requires a User and Password pair to continue. After
// sending this message, a Method will expect APT to send a 602 Authorization
// Credentials message with the required information. It is possible for a
// transport method to send this message to APT multiple times (both for
// multiple credential steps as well as retries and timeouts)
type AuthorizationRequired struct {
	Site string
}

// AuthorizationCredentials (status code 602) is sent in response to a 402
// Authorization Required.
//
// When received, it will contain the entered username and password.
type AuthorizationCredentials struct {
	Password string
	User     string
	Site     string
}
