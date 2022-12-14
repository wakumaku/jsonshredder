package forwarder

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	kindSNS     = "sns"
	kindSQS     = "sqs"
	kindKinesis = "kinesis"
)

var (
	awsForwardersTest = map[string]string{
		kindSNS:     "Publish",
		kindSQS:     "SendMessage",
		kindKinesis: "PutRecord",
	}
)

func TestAWSForwarders(t *testing.T) {
	for kind := range awsForwardersTest {
		sqsSVC := startAWSMockServer(kind)
		defer sqsSVC.Close()

		fwd, err := buildForwarder(kind, sqsSVC.URL)
		assert.Nil(t, err, "unexpected error creating forwarder")

		message := []byte("hello world")
		assert.Nil(t, fwd.Publish(context.TODO(), message), "unexpected error publishing message")
	}
}

func TestSNSErrorPublishMessage(t *testing.T) {
	for kind, action := range awsForwardersTest {
		sqsSVC := startAWSMockServer(kind)
		defer sqsSVC.Close()

		fwd, err := buildForwarder(kind, sqsSVC.URL+"/fail"+action)
		assert.Nil(t, err, "unexpected error")

		message := []byte("hello world")
		assert.NotNilf(t, fwd.Publish(context.TODO(), message), "expecting an error sending message")
	}
}

func buildForwarder(kind, endpoint string) (Forwarder, error) {
	switch kind {
	case kindSQS:
		return NewSQS(context.TODO(), "queuename", AWSWithEndpoint(endpoint), AWSWithRegion("us-east-1"))
	case kindSNS:
		return NewSNS(context.TODO(), "topicname", AWSWithEndpoint(endpoint), AWSWithRegion("us-east-1"))
	case kindKinesis:
		return NewKinesis(context.TODO(), "streamname", AWSWithEndpoint(endpoint), AWSWithRegion("us-east-1"))
	}
	return nil, errors.New("unknown kind of forwarder")
}

func TestConfigOptions(t *testing.T) {
	const (
		endpoint         = "Endpoint"
		keyID            = "KeyID"
		secret           = "Secret"
		profile          = "Profile"
		region           = "Region"
		resourceARN      = "ResourceARN"
		resourceEndpoint = "ResourceEndpoint"
		resourceName     = "ResourceName"
	)

	cfg := buildAWSConfigFromOptions(
		AWSWithEndpoint(endpoint),
		AWSWithKeyID(keyID),
		AWSWithSecret(secret),
		AWSWithProfile(profile),
		AWSWithRegion(region),
		AWSWithResourceARN(resourceARN),
		AWSWithResourceEndpoint(resourceEndpoint),
		AWSWithResourceName(resourceName),
	)

	assert.Equal(t, endpoint, cfg.endpoint)
	assert.Equal(t, keyID, cfg.key)
	assert.Equal(t, secret, cfg.secret)
	assert.Equal(t, profile, cfg.profile)
	assert.Equal(t, region, cfg.region)
	assert.Equal(t, resourceARN, cfg.resourceArn)
	assert.Equal(t, resourceEndpoint, cfg.resourceEndpoint)
	assert.Equal(t, resourceName, cfg.resourceName)
}

func TestInitAWSSession(t *testing.T) {
	cfg := &AWSConfig{
		endpoint: "endpoint",
		profile:  "profile",
		region:   "region",
	}
	s, _ := initAWSSession(context.TODO(), cfg)

	assert.Equal(t, "region", s.Region)

	cfg = &AWSConfig{
		endpoint: "endpoint",
		key:      "key",
		secret:   "secret",
		region:   "region",
	}
	s, _ = initAWSSession(context.TODO(), cfg)

	p, _ := s.Credentials.Retrieve(context.TODO())

	assert.Equal(t, "region", s.Region)
	assert.Equal(t, "key", p.AccessKeyID)
	assert.Equal(t, "secret", p.SecretAccessKey)
}
