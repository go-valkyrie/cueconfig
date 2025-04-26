package cueconfig

import (
	"cuelang.org/go/cue"
	"errors"
	"go-valkyrie.com/cueconfig/source"
	"slices"
	"sync"
)

type Config struct {
	schema  *Schema
	config  *cue.Value
	loadMu  sync.Mutex
	sources []source.Source
}

type Option func(c *Config) error

func WithSources(sources ...source.Source) Option {
	return func(c *Config) error {
		c.sources = append(c.sources, sources...)
		return nil
	}
}

func New(schema *Schema, options ...Option) (*Config, error) {
	config := &Config{
		schema:  schema,
		sources: make([]source.Source, 0),
	}

	for _, option := range options {
		if err := option(config); err != nil {
			return nil, err
		}
	}

	if err := config.Reload(); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) AddSource(source source.Source) error {
	c.sources = append(c.sources, source)
	return nil
}

func (c *Config) ClearSources() {
	c.sources = make([]source.Source, 0)
}

func (c *Config) Raw() *cue.Value {
	return c.config
}

func (c *Config) Reload() error {
	c.loadMu.Lock()
	defer c.loadMu.Unlock()

	schema := c.schema.Value()
	sources := make([]cue.Value, len(c.sources))

	for i, src := range c.sources {
		if v, err := src.Load(schema.Context()); err != nil {
			return err
		} else {
			sources[i] = v
		}
	}

	value := schema
	for _, src := range sources {
		value = value.Unify(src)
	}

	c.config = &value
	return nil
}

func (c *Config) Sources() []source.Source {
	return slices.Clone(c.sources)
}

func (c *Config) UnmarshalKey(key string, v any) error {
	if c.config == nil {
		return errors.New("config not loaded")
	}

	if val := c.config.LookupPath(cue.ParsePath(key)); !val.Exists() {
		return errors.New("key not found")
	} else if val.Err() != nil {
		return val.Err()
	} else {
		return val.Decode(v)
	}
}

func (c *Config) Unmarshal(v any) error {
	if c.config == nil {
		return errors.New("config not loaded")
	}

	if err := c.config.Decode(v); err != nil {
		return err
	}

	return nil
}

func (c *Config) ValueAt(path string) cue.Value {
	return c.config.LookupPath(cue.ParsePath(path))
}
