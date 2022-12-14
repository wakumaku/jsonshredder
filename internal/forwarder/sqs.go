package forwarder

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go/aws"
)

type isqs interface {
	GetQueueUrl(ctx context.Context, params *sqs.GetQueueUrlInput, optFns ...func(*sqs.Options)) (*sqs.GetQueueUrlOutput, error)
	SendMessage(ctx context.Context, params *sqs.SendMessageInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error)
}

type sqsForwarder struct {
	c        isqs
	queueURL string
}

// NewSQS creates a new forwarder implementing AWS SQS Queues
func NewSQS(ctx context.Context, name string, opts ...AWSOption) (Forwarder, error) {
	awsCfg := buildAWSConfigFromOptions(opts...)
	cfg, err := initAWSSession(ctx, awsCfg)
	if err != nil {
		return nil, fmt.Errorf("building aws session: %w", err)
	}

	clientOpts := make([]func(*sqs.Options), 0)
	if awsCfg.endpoint != "" {
		clientOpts = append(clientOpts, sqs.WithEndpointResolver(
			sqs.EndpointResolverFromURL(awsCfg.endpoint)))
	}

	c := sqs.NewFromConfig(cfg, clientOpts...)
	o, err := c.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
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
func (p *sqsForwarder) Publish(ctx context.Context, msg []byte) error {
	if _, err := p.c.SendMessage(ctx, &sqs.SendMessageInput{
		MessageBody: aws.String(string(msg)),
		QueueUrl:    aws.String(p.queueURL),
	}); err != nil {
		return fmt.Errorf("sqs send message: %s", err)
	}

	return nil
}
