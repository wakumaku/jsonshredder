package forwarder

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// httpHandlerStatusMessage helper handler for testing
func httpHandlerStatusMessage(status int, message string) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(status)
		fmt.Fprintf(rw, "%s", message)
	}
}

func TestHttpForwarderDefaultsOk(t *testing.T) {
	server := httptest.NewServer(httpHandlerStatusMessage(http.StatusOK, "ok"))
	defer server.Close()

	fwd := NewHTTP(server.URL)
	if err := fwd.Publish(context.TODO(), []byte("hello")); err != nil {
		t.Errorf("unexpected error: %s", err)
	}
}

func TestHttpForwarderRequestURLMethodError(t *testing.T) {
	fwd := NewHTTP("#")
	err := fwd.Publish(context.TODO(), []byte("hello"))
	assert.NotNil(t, err)
}

func TestHttpForwarderDefaultsFails(t *testing.T) {
	server := httptest.NewServer(httpHandlerStatusMessage(http.StatusInternalServerError, "ko"))
	defer server.Close()

	fwd := NewHTTP(server.URL)
	if err := fwd.Publish(context.TODO(), []byte("hello")); err == nil {
		t.Error("expecting an error here")
	}
}

func TestHttpForwarderErrorOnTimeOut(t *testing.T) {
	h := func() http.HandlerFunc {
		return func(rw http.ResponseWriter, r *http.Request) {
			time.Sleep(100 * time.Millisecond)
		}
	}()
	server := httptest.NewServer(h)
	defer server.Close()

	fwd := NewHTTP(server.URL, HTTPWithTimeOut(10*time.Millisecond))
	if err := fwd.Publish(context.TODO(), []byte("hello")); err == nil {
		t.Error("expecting a timeout error here")
	}
}

func TestHttpForwarderWithAuthorization(t *testing.T) {
	bearer := "Bearer 123456789"
	expectedAuthHeader := []string{bearer}

	h := func() http.HandlerFunc {
		return func(rw http.ResponseWriter, r *http.Request) {
			if _, found := r.Header["Authorization"]; found {
				assert.Equal(t, expectedAuthHeader, r.Header["Authorization"])
			}
			rw.WriteHeader(http.StatusOK)
		}
	}()

	server := httptest.NewServer(h)
	defer server.Close()

	fwd := NewHTTP(server.URL, HTTPWithHeaderAuth(bearer))
	err := fwd.Publish(context.TODO(), nil)
	assert.Nil(t, err)
}

func TestHttpForwarderWithMethod(t *testing.T) {
	h := func() http.HandlerFunc {
		return func(rw http.ResponseWriter, r *http.Request) {
			assert.Equal(t, r.Method, "POST")
		}
	}()
	server := httptest.NewServer(h)
	defer server.Close()

	fwd := NewHTTP(server.URL, HTTPWithMethod("POST"))
	err := fwd.Publish(context.TODO(), nil)
	assert.Nil(t, err)
}

func TestHttpForwarderWithExpectedStatus(t *testing.T) {
	h := func() http.HandlerFunc {
		return func(rw http.ResponseWriter, r *http.Request) {
			rw.WriteHeader(http.StatusTeapot)
		}
	}()
	server := httptest.NewServer(h)
	defer server.Close()

	fwd := NewHTTP(server.URL, HTTPWithExpectedStatus(http.StatusTeapot))
	err := fwd.Publish(context.TODO(), nil)
	assert.Nil(t, err)
}
