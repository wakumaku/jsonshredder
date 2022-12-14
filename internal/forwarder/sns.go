package forwarder

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

type snsi interface {
	Publish(ctx context.Context, params *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error)
}

type snsForwarder struct {
	c        snsi
	topicARN string
}

// NewSNS creates a new SNS forwarder
func NewSNS(ctx context.Context, topicARN string, opts ...AWSOption) (Forwarder, error) {
	awsCfg := buildAWSConfigFromOptions(opts...)
	cfg, err := initAWSSession(ctx, awsCfg)
	if err != nil {
		return nil, fmt.Errorf("initializing AWS session: %s", err)
	}

	clientOpts := make([]func(*sns.Options), 0)
	if awsCfg.endpoint != "" {
		clientOpts = append(clientOpts, sns.WithEndpointResolver(
			sns.EndpointResolverFromURL(awsCfg.endpoint)))
	}

	return &snsForwarder{
		c:        sns.NewFromConfig(cfg, clientOpts...),
		topicARN: topicARN,
	}, nil
}

func (p *snsForwarder) Publish(ctx context.Context, msg []byte) error {
	if _, err := p.c.Publish(ctx, &sns.PublishInput{
		TopicArn: aws.String(p.topicARN),
		Message:  aws.String(string(msg)),
	}); err != nil {
		return fmt.Errorf("sns publish message: %s", err)
	}

	return nil
}
