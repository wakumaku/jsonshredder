package forwarder

import "context"

type dummyForwarder struct {
}

// NewDummy for testing
func NewDummy() Forwarder {
	return &dummyForwarder{}
}

func (p *dummyForwarder) Publish(m_ context.Context, _ []byte) error {
	return nil
}
