package config

import (
	"github.com/rs/zerolog"
)

// Operation available
type Operation string

const (
	// OperationAdd just adds new mappings to the original json
	OperationAdd Operation = "add"
	// OperationExtract creates a new json from defined mappings
	OperationExtract Operation = "extract"
)

// App holds the application config
type App struct {
	Port            string
	LogLevel        zerolog.Level
	Transformations map[string]Transformation
	Forwarders      map[string]Forwarder
}

// ForwarderKind defines the kind of proxies supported
type ForwarderKind string

// List of supported proxies
const (
	KindHTTP    ForwarderKind = "http"
	KindSNS     ForwarderKind = "sns"
	KindSQS     ForwarderKind = "sqs"
	KindKinesis ForwarderKind = "kinesis"
)

// ForwarderSetting defines proxies settings
type ForwarderSetting string

// List of available forwarder settings
const (
	SettingAWSEndpoint         ForwarderSetting = "aws_endpoint"
	SettingAWSKey              ForwarderSetting = "aws_access_key_id"
	SettingAWSSecret           ForwarderSetting = "aws_secret_access_key"
	SettingAWSProfile          ForwarderSetting = "aws_profile"
	SettingAWSRegion           ForwarderSetting = "aws_region"
	SettingAWSResourceArn      ForwarderSetting = "aws_resource_arn"
	SettingAWSResourceEndpoint ForwarderSetting = "aws_resource_endpoint"
	SettingAWSResourceName     ForwarderSetting = "aws_resource_name"
	SettingHTTPEndpoint        ForwarderSetting = "http_endpoint"
	SettingHTTPHeaderAuth      ForwarderSetting = "http_header_auth"
	SettingHTTPStatusOK        ForwarderSetting = "http_status_ok"
)

// Transformation config definition
type Transformation struct {
	Operation Operation
	Mappings  []Mapping
}

// Forwarder config definition
type Forwarder struct {
	Kind     ForwarderKind
	Settings map[ForwarderSetting]interface{}
}
