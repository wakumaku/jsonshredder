package forwarder

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
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

func AWSWithKResourceEndpoint(value string) AWSOption {
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
	config := &AWSConfig{}
	for _, o := range opts {
		if o != nil {
			o(config)
		}
	}
	return config
}

func initAWSSession(cfg *AWSConfig) (*session.Session, error) {
	var cred *credentials.Credentials
	if cfg.key != "" && cfg.secret != "" {
		cred = credentials.NewStaticCredentials(cfg.key, cfg.secret, "")
	}

	if cfg.profile != "" {
		os.Setenv("AWS_SDK_LOAD_CONFIG", "true")
		return session.NewSessionWithOptions(session.Options{
			Profile: cfg.profile,
			Config: aws.Config{
				Endpoint:    aws.String(cfg.endpoint),
				Credentials: cred,
				Region:      aws.String(cfg.region),
			},
		})
	}

	return session.NewSession(&aws.Config{
		Endpoint:    aws.String(cfg.endpoint),
		Credentials: cred,
		Region:      aws.String(cfg.region),
	})
}
