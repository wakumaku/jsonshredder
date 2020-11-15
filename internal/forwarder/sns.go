package forwarder

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
)

type snsForwarder struct {
	c        snsiface.SNSAPI
	topicARN string
}

// NewSNS creates a new SNS forwarder
func NewSNS(topicARN string, opts ...AWSOption) (Forwarder, error) {
	s, err := initAWSSession(buildAWSConfigFromOptions(opts...))
	if err != nil {
		return nil, fmt.Errorf("initializing AWS session: %s", err)
	}

	return &snsForwarder{
		c:        sns.New(s),
		topicARN: topicARN,
	}, nil
}

func (p *snsForwarder) Publish(msg []byte) error {
	if _, err := p.c.Publish(&sns.PublishInput{
		TopicArn: aws.String(p.topicARN),
		Message:  aws.String(string(msg)),
	}); err != nil {
		return fmt.Errorf("sns publish message: %s", err)
	}

	return nil
}
