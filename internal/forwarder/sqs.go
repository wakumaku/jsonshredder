package forwarder

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
)

type sqsForwarder struct {
	c        sqsiface.SQSAPI
	queueURL string
}

// NewSQS creates a new forwarder implementing AWS SQS Queues
func NewSQS(name string, opts ...AWSOption) (Forwarder, error) {
	s, err := initAWSSession(buildAWSConfigFromOptions(opts...))
	if err != nil {
		return nil, fmt.Errorf("building aws session: %w", err)
	}

	c := sqs.New(s)
	o, err := c.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(name),
	})
	if err != nil {
		return nil, fmt.Errorf("getting queue URL: %s", err)
	}

	return &sqsForwarder{
		c:        c,
		queueURL: *o.QueueUrl,
	}, nil
}

// Publish sends the message to an SQS
func (p *sqsForwarder) Publish(msg []byte) error {
	if _, err := p.c.SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String(string(msg)),
		QueueUrl:    aws.String(p.queueURL),
	}); err != nil {
		return fmt.Errorf("sqs send message: %s", err)
	}

	return nil
}
