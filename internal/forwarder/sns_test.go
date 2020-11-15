package forwarder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSNSPublishMessage(t *testing.T) {
	sqsSVC := startAWSSNSMockServer()
	defer sqsSVC.Close()

	fwd, err := NewSNS("aws::sns::test", AWSWithEndpoint(sqsSVC.URL), AWSWithRegion("us-east-1"))
	assert.Nil(t, err, "unexpected error")

	message := []byte("hello world")
	assert.Nil(t, fwd.Publish(message), "unexpected error")
}

func TestSNSErrorPublishMessage(t *testing.T) {
	sqsSVC := startAWSSNSMockServer()
	defer sqsSVC.Close()

	fwd, err := NewSNS("aws::sns::test", AWSWithEndpoint(sqsSVC.URL+"/failPublish"), AWSWithRegion("us-east-1"))
	assert.Nil(t, err, "unexpected error")

	message := []byte("hello world")
	assert.NotNilf(t, fwd.Publish(message), "expecting an error sending message")
}
