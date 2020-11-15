package service

import (
	"errors"
	"fmt"
	"wakumaku/jsonshredder/internal/config"
	"wakumaku/jsonshredder/internal/forwarder"

	"github.com/rs/zerolog"
)

// Forwarder service
type Forwarder struct {
	forwarders map[string]forwarder.Forwarder
	logger     *zerolog.Logger
}

// Errors
var (
	ErrForwardPublishing = errors.New("forwarder publishing")
	ErrForwardNotFound   = errors.New("forwarder not found")
)

// NewForwarder creates a new forwarder service
func NewForwarder(fwdConfig map[string]config.Forwarder, logger *zerolog.Logger) *Forwarder {
	lgr := logger.With().Str("section", "service.forwarder").Logger()

	forwarders := buildForwarders(fwdConfig, &lgr)

	return &Forwarder{
		forwarders: forwarders,
		logger:     &lgr,
	}
}

// Forward resolves the forwarder to be used and sends the data
func (c *Forwarder) Forward(forwarderName string, input []byte) error {
	fwd, err := c.getForwarder(forwarderName)
	if err == nil {
		if err := fwd.Publish(input); err != nil {
			err = fmt.Errorf("%w: %s", ErrForwardPublishing, err)
			c.logger.Debug().Err(err).Send()
			return err
		}
		c.logger.Debug().Str("forwarder", forwarderName).Msg("ok")
		return nil
	}

	err = fmt.Errorf("%w: '%s'", ErrForwardNotFound, forwarderName)
	c.logger.Debug().Err(err).Send()
	return err
}

// getTransformConfig searches a transformation config by name in the list
func (c *Forwarder) getForwarder(name string) (forwarder.Forwarder, error) {
	if t, found := c.forwarders[name]; found && t != nil {
		return t, nil
	}

	return nil, errors.New("not found")
}

func buildForwarders(fwdConfig map[string]config.Forwarder, logger *zerolog.Logger) map[string]forwarder.Forwarder {
	forwarders := make(map[string]forwarder.Forwarder, len(fwdConfig))
	for name, cfg := range fwdConfig {
		var err error
		var f forwarder.Forwarder
		switch cfg.Kind {
		case config.KindHTTP:
			endpoint, _ := cfg.Settings[config.SettingHTTPEndpoint].(string)
			params := make([]forwarder.HTTPOption, 0)
			if v, ok := cfg.Settings[config.SettingHTTPHeaderAuth].(string); ok {
				params = append(params, forwarder.HTTPWithHeaderAuth(v))
			}
			if v, ok := cfg.Settings[config.SettingHTTPStatusOK].(int); ok {
				params = append(params, forwarder.HTTPWithExpectedStatus(v))
			}
			f = forwarder.NewHTTP(endpoint, params...)
		case config.KindSNS:
			topicARN, _ := cfg.Settings[config.SettingAWSResourceArn].(string)
			params := getCommonAWSParams(cfg.Settings)
			f, err = forwarder.NewSNS(topicARN, params...)
		case config.KindSQS:
			queueName, _ := cfg.Settings[config.SettingAWSResourceName].(string)
			params := getCommonAWSParams(cfg.Settings)
			f, err = forwarder.NewSQS(queueName, params...)
		case config.KindKinesis:
			streamName, _ := cfg.Settings[config.SettingAWSResourceName].(string)
			params := getCommonAWSParams(cfg.Settings)
			f, err = forwarder.NewKinesis(streamName, params...)
		default:
			logger.Error().Msgf("forwarder: %s, unknown kind: '%s'!", name, cfg.Kind)
		}
		if err != nil {
			logger.Warn().Err(err).Msgf("forwarder: %s IS DISABLED! kind: '%s'", name, cfg.Kind)
			continue
		}
		forwarders[name] = f
	}
	return forwarders
}

func getCommonAWSParams(c map[config.ForwarderSetting]interface{}) []forwarder.AWSOption {
	params := make([]forwarder.AWSOption, 0)

	if v, ok := c[config.SettingAWSEndpoint].(string); ok {
		params = append(params, forwarder.AWSWithEndpoint(v))
	}

	if v, ok := c[config.SettingAWSKey].(string); ok {
		params = append(params, forwarder.AWSWithKeyID(v))
	}

	if v, ok := c[config.SettingAWSSecret].(string); ok {
		params = append(params, forwarder.AWSWithSecret(v))
	}

	if v, ok := c[config.SettingAWSProfile].(string); ok {
		params = append(params, forwarder.AWSWithProfile(v))
	}

	if v, ok := c[config.SettingAWSRegion].(string); ok {
		params = append(params, forwarder.AWSWithRegion(v))
	}
	return params
}
