package forwarder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKinesisPutRecord(t *testing.T) {
	sqsSVC := startAWSKinesisMockServer()
	defer sqsSVC.Close()

	fwd, err := NewKinesis("stream", AWSWithEndpoint(sqsSVC.URL), AWSWithRegion("us-east-1"))
	assert.Nil(t, err, "unexpected error")

	message := []byte("hello world")
	assert.Nil(t, fwd.Publish(message), "unexpected error")
}

func TestKinesisErrorPutRecord(t *testing.T) {
	sqsSVC := startAWSKinesisMockServer()
	defer sqsSVC.Close()

	fwd, err := NewKinesis("stream", AWSWithEndpoint(sqsSVC.URL+"/failPutRecord"), AWSWithRegion("us-east-1"))
	assert.Nil(t, err, "unexpected error")

	message := []byte("hello world")
	assert.NotNilf(t, fwd.Publish(message), "expecting an error sending message")
}
