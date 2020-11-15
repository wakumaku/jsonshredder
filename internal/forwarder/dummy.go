package forwarder

type dummyForwarder struct {
}

// NewDummy for testing
func NewDummy() Forwarder {
	return &dummyForwarder{}
}

func (p *dummyForwarder) Publish(msg []byte) error {
	return nil
}
