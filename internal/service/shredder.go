package service

import (
	"errors"
	"fmt"

	"github.com/wakumaku/jsonshredder/internal/shredder"

	"github.com/wakumaku/jsonshredder/internal/config"

	"github.com/rs/zerolog"
)

// Shredder service
type Shredder struct {
	transformations map[string]config.Transformation
	logger          *zerolog.Logger
}

// Errors
var (
	ErrShredProcessing        = errors.New("processing input")
	ErrShredTransformNotFound = errors.New("unknown transformation")
)

// NewShredder creates a new shredder service
func NewShredder(transformations map[string]config.Transformation, logger *zerolog.Logger) *Shredder {
	lgr := logger.With().Str("section", "service.shredder").Logger()

	return &Shredder{
		transformations: transformations,
		logger:          &lgr,
	}
}

// Shred applies a transformation to an input, returns the shredded result
func (c *Shredder) Shred(transformName string, input []byte) ([]byte, error) {
	if transformConfig, err := c.getTransformConfig(transformName); err == nil {
		out, err := shredder.Shred(transformConfig, input)
		if err != nil {
			err = fmt.Errorf("%w: %s", ErrShredProcessing, err)
			c.logger.Debug().Err(err).Send()
			return nil, err
		}
		c.logger.Debug().Str("transformation", transformName).Msg("ok")
		return out, nil
	}

	err := fmt.Errorf("%w: '%s'", ErrShredTransformNotFound, transformName)
	c.logger.Debug().Err(err).Send()
	return nil, err
}

// getTransformConfig searches a transformation config by name in the list
func (c *Shredder) getTransformConfig(name string) (config.Transformation, error) {
	if t, found := c.transformations[name]; found {
		return t, nil
	}

	return config.Transformation{}, errors.New("not found")
}
