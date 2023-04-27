package coinbase

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var errInvalidRoundTripArgs = fmt.Errorf("invalid auth arguments")

// roundTripper is an HTTP round tripper that acts as a middleware to add
// auth requirements to HTTP requests.
type roundTripper struct {
	roundTrip func(*http.Request) (*http.Response, error)
}

// RoundTrip implements the "http.RoundTripper" interface.
func (rtripper *roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return rtripper.roundTrip(req)
}

// newRoundTrip signs the given HTTP request with the provided Coinbase API
// key and secret, and sends the request using the default HTTP transport. The
// signed request includes the current timestamp, HTTP method, request path, and
// request body (if present). The function returns the HTTP response and any
// error that occurred during the request. If an error occurs during the
// request, it is wrapped with additional context information.
func newRoundTrip(req *http.Request, key, secret string) (*http.Response, error) {
	var body []byte
	if req.Body != nil {
		body, _ = io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewBuffer(body))
	}

	signature := hmac.New(sha256.New, []byte(secret))

	rpath := req.URL.Path
	if req.URL.RawQuery != "" {
		rpath = fmt.Sprintf("%s?%s", req.URL.Path, req.URL.RawQuery)
	}

	formatBase := 10
	unix := strconv.FormatInt(time.Now().Unix(), formatBase)

	msg := strings.Join([]string{unix, req.Method, rpath, string(body)}, "")

	// Don't handle error because hash.Write method never returns an
	// error.
	signature.Write([]byte(msg))
	sig := hex.EncodeToString(signature.Sum(nil))

	req.Header.Add("cb-access-key", key)
	req.Header.Add("cb-access-sign", sig)
	req.Header.Add("cb-access-timestamp", unix)

	rsp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}

	return rsp, nil
}

// newRoundTripper will return a "RoundTrip" function that can be used
// as a "RoundTrip" function in an "http.RoundTripper" interface to authenticate
// requests to the Coinbase Cloud API.
func newRoundTripper(key, secret string) (*roundTripper, error) {
	if key == "" || secret == "" {
		return nil, errInvalidRoundTripArgs
	}

	rtripper := &roundTripper{
		roundTrip: func(req *http.Request) (*http.Response, error) {
			return newRoundTrip(req, key, secret)
		},
	}

	return rtripper, nil
}
