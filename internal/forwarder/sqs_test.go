package forwarder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSQSPublish(t *testing.T) {
	sqsSVC := startAWSSQSMockServer()
	defer sqsSVC.Close()

	fwd, err := NewSQS("queuename", AWSWithEndpoint(sqsSVC.URL), AWSWithRegion("us-east-1"))

	assert.Nilf(t, err, "unexpected error")
	message := []byte("hello world")

	assert.Nilf(t, fwd.Publish(message), "unexpected error")
}

func TestSQSErrorGettingQueue(t *testing.T) {
	sqsSVC := startAWSSQSMockServer()
	defer sqsSVC.Close()

	_, err := NewSQS("queuename", AWSWithEndpoint(sqsSVC.URL+"/failGetQueueUrl"), AWSWithRegion("us-east-1"))

	assert.NotNilf(t, err, "expecting an error initializing and getting the queue url")
}

func TestSQSErrorSendMessage(t *testing.T) {
	sqsSVC := startAWSSQSMockServer()
	defer sqsSVC.Close()

	fwd, err := NewSQS("queuename", AWSWithEndpoint(sqsSVC.URL+"/failSendMessage"), AWSWithRegion("us-east-1"))

	assert.Nilf(t, err, "unexpected error")

	message := []byte("hello world")
	assert.NotNilf(t, fwd.Publish(message), "expecting an error sending message")
}
