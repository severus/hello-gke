package kubernetes

import "net/http"

// AuthTransport is the implementation of http.RoundTripper
// to provide HTTP Basic Authorization.
type AuthTransport struct {
	username string
	password string
}

// NewAuthTransport returns a new HTTP Transport.
func NewAuthTransport(cfg *Config) *AuthTransport {
	return &AuthTransport{
		username: cfg.Username,
		password: cfg.Password,
	}
}

// RoundTrip implements the http.RoundTripper interface.
func (t *AuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req2 := cloneRequest(req)
	req2.SetBasicAuth(t.username, t.password)

	return http.DefaultTransport.RoundTrip(req2)
}

// cloneRequest returns a clone of the provided *http.Request.
// The clone is a shallow copy of the struct and its Header map.
func cloneRequest(req *http.Request) *http.Request {
	// shallow copy of the struct
	req2 := new(http.Request)
	*req2 = *req
	// deep copy of the Header
	req2.Header = make(http.Header, len(req.Header))
	for k, s := range req.Header {
		req2.Header[k] = append([]string(nil), s...)
	}
	return req2
}
