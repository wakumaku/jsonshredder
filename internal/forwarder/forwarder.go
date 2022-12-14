package forwarder

import (
	"context"
	"fmt"
)

// Forwarder defines a forwarder
type Forwarder interface {
	Publish(context.Context, []byte) error
}

// ErrForwarder is the error type to be returned by forwarders
type ErrForwarder struct {
	name string
	err  error
}

func (e *ErrForwarder) Unwrap(err error) error {
	return e.err
}

// Error returns the stringified error
func (e *ErrForwarder) Error() string {
	return fmt.Sprintf("forwarder '%s': %s", e.name, e.err)
}

// Error creates a new forwarder error
func Error(name string, err error) error {
	return &ErrForwarder{
		name: name,
		err:  err,
	}
}
