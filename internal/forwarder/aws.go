package forwarder

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

type AWSConfig struct {
	endpoint         string
	key              string
	secret           string
	profile          string
	region           string
	resourceArn      string
	resourceEndpoint string
	resourceName     string
}

type AWSOption func(*AWSConfig)

func AWSWithEndpoint(value string) AWSOption {
	return func(c *AWSConfig) {
		c.endpoint = value
	}
}

func AWSWithKeyID(value string) AWSOption {
	return func(c *AWSConfig) {
		c.key = value
	}
}

func AWSWithSecret(value string) AWSOption {
	return func(c *AWSConfig) {
		c.secret = value
	}
}

func AWSWithProfile(value string) AWSOption {
	return func(c *AWSConfig) {
		c.profile = value
	}
}

func AWSWithRegion(value string) AWSOption {
	return func(c *AWSConfig) {
		c.region = value
	}
}

func AWSWithResourceARN(value string) AWSOption {
	return func(c *AWSConfig) {
		c.resourceArn = value
	}
}

func AWSWithResourceEndpoint(value string) AWSOption {
	return func(c *AWSConfig) {
		c.resourceEndpoint = value
	}
}

func AWSWithResourceName(value string) AWSOption {
	return func(c *AWSConfig) {
		c.resourceName = value
	}
}

func buildAWSConfigFromOptions(opts ...AWSOption) *AWSConfig {
	config := AWSConfig{}
	for _, o := range opts {
		if o != nil {
			o(&config)
		}
	}
	return &config
}

func initAWSSession(ctx context.Context, awsCfg *AWSConfig) (aws.Config, error) {
	cfgOptions := make([]func(*config.LoadOptions) error, 0)

	if awsCfg.key != "" && awsCfg.secret != "" {
		cfgOptions = append(cfgOptions,
			config.WithCredentialsProvider(
				credentials.NewStaticCredentialsProvider(awsCfg.key, awsCfg.secret, "")))
	}

	if awsCfg.profile != "" {
		cfgOptions = append(cfgOptions, config.WithSharedConfigProfile(awsCfg.profile))
	}

	if awsCfg.region != "" {
		cfgOptions = append(cfgOptions, config.WithRegion(awsCfg.region))
	}

	return config.LoadDefaultConfig(ctx, cfgOptions...)
}
