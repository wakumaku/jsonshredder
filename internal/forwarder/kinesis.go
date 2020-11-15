package forwarder

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/aws/aws-sdk-go/service/kinesis/kinesisiface"
)

type kinesisForwarder struct {
	c          kinesisiface.KinesisAPI
	streamName string
}

// NewKinesis creates a new Kinesis Datastream forwarder
func NewKinesis(streamName string, opts ...AWSOption) (Forwarder, error) {
	s, err := initAWSSession(buildAWSConfigFromOptions(opts...))
	if err != nil {
		return nil, fmt.Errorf("initializing AWS session: %s", err)
	}

	return &kinesisForwarder{
		c:          kinesis.New(s),
		streamName: streamName,
	}, nil
}

// Publish puts a record to the kinesis stream
func (p *kinesisForwarder) Publish(msg []byte) error {
	if _, err := p.c.PutRecord(&kinesis.PutRecordInput{
		StreamName:   aws.String(p.streamName),
		PartitionKey: aws.String(fmt.Sprint(time.Now().UnixNano() / int64(time.Millisecond))),
		Data:         msg,
	}); err != nil {
		return fmt.Errorf("kinesis put record: %s", err)
	}

	return nil
}
