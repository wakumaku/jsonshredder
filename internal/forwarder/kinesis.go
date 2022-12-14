package forwarder

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
)

type kinesisi interface {
	PutRecord(ctx context.Context, params *kinesis.PutRecordInput, optFns ...func(*kinesis.Options)) (*kinesis.PutRecordOutput, error)
}

type kinesisForwarder struct {
	c          kinesisi
	streamName string
}

// NewKinesis creates a new Kinesis Datastream forwarder
func NewKinesis(ctx context.Context, streamName string, opts ...AWSOption) (Forwarder, error) {
	awsCfg := buildAWSConfigFromOptions(opts...)
	s, err := initAWSSession(ctx, awsCfg)
	if err != nil {
		return nil, fmt.Errorf("initializing AWS session: %s", err)
	}

	clientOpts := make([]func(*kinesis.Options), 0)
	if awsCfg.endpoint != "" {
		clientOpts = append(clientOpts, kinesis.WithEndpointResolver(
			kinesis.EndpointResolverFromURL(awsCfg.endpoint)))
	}

	return &kinesisForwarder{
		c:          kinesis.NewFromConfig(s, clientOpts...),
		streamName: streamName,
	}, nil
}

// Publish puts a record to the kinesis stream
func (p *kinesisForwarder) Publish(ctx context.Context, msg []byte) error {
	if _, err := p.c.PutRecord(ctx, &kinesis.PutRecordInput{
		StreamName:   aws.String(p.streamName),
		PartitionKey: aws.String(fmt.Sprint(time.Now().UnixNano() / int64(time.Millisecond))),
		Data:         msg,
	}); err != nil {
		return fmt.Errorf("kinesis put record: %s", err)
	}

	return nil
}
