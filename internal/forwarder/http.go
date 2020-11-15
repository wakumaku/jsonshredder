package forwarder

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// HTTP Forwarder
type httpForwarder struct {
	c              *http.Client
	endpoint       string
	headerAuth     string
	expectedStatus int
	method         string
}

// HTTPConfig holds the options to configure the http forwarder
type HTTPConfig struct {
	headerAuth     string
	expectedStatus int
	timeOut        time.Duration
	method         string
}

// HTTPOption config option
type HTTPOption func(*HTTPConfig)

// HTTPWithHeaderAuth Adds a header auth for HTTP requests
func HTTPWithHeaderAuth(value string) HTTPOption {
	return func(c *HTTPConfig) {
		c.headerAuth = value
	}
}

// HTTPWithExpectedStatus defines the expected response http status as valid
func HTTPWithExpectedStatus(value int) HTTPOption {
	return func(c *HTTPConfig) {
		c.expectedStatus = value
	}
}

// HTTPWithTimeOut defines the timeout for HTTP connections
func HTTPWithTimeOut(value time.Duration) HTTPOption {
	return func(c *HTTPConfig) {
		c.timeOut = value
	}
}

// HTTPWithMethod defines the method to be userd GET, POST ...
func HTTPWithMethod(value string) HTTPOption {
	return func(c *HTTPConfig) {
		c.method = value
	}
}

// NewHTTP creates a new HTTP forwarder
func NewHTTP(endpoint string, opts ...HTTPOption) Forwarder {
	config := &HTTPConfig{
		expectedStatus: http.StatusOK,
		timeOut:        5 * time.Second,
		method:         http.MethodGet,
	}
	for _, o := range opts {
		if o != nil {
			o(config)
		}
	}

	return &httpForwarder{
		endpoint:       endpoint,
		headerAuth:     config.headerAuth,
		expectedStatus: config.expectedStatus,
		method:         config.method,
		c: &http.Client{
			Timeout: config.timeOut,
		},
	}
}

// Publish sends an http request
func (p *httpForwarder) Publish(msg []byte) error {
	req, err := http.NewRequest(p.method, p.endpoint, bytes.NewBuffer(msg))
	if err != nil {
		return Error("http", fmt.Errorf("create request: %s", err))
	}

	if p.headerAuth != "" {
		req.Header.Add("Authorization", p.headerAuth)
	}

	if resp, err := p.c.Do(req); err != nil {
		return Error("http", fmt.Errorf("do request: %s", err))
	} else if resp.StatusCode != p.expectedStatus {
		var r string
		if resp.Body != nil {
			defer resp.Body.Close()
			if b, err := ioutil.ReadAll(resp.Body); err == nil {
				r = fmt.Sprintf("%q", b)
			}
		}
		return Error("http", fmt.Errorf("unexpected status code: %d, message: %s", resp.StatusCode, r))
	}

	return nil
}
