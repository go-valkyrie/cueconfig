package schema

import (
	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"go-valkyrie.com/cueconfig/pkgerrs"
	"go-valkyrie.com/cueconfig/source"
)

type Schema struct {
	ctx    *cue.Context
	schema cue.Value
}

type Options struct {
	Schema source.Source
}

func New(source source.Source) (*Schema, error) {
	if source == nil {
		return nil, &pkgerrs.ErrInvalidArgument{Argument: "source", Reason: "cannot be nil"}
	}

	ctx := cuecontext.New()
	if schema, err := source.Load(ctx); err != nil {
		return nil, err
	} else {
		return &Schema{
			ctx:    ctx,
			schema: schema,
		}, nil
	}
}

func (c *Schema) Context() *cue.Context {
	return c.ctx
}

func (c *Schema) Value() cue.Value {
	return c.schema
}

func (c *Schema) Valid() bool {
	return c.schema.Err() == nil
}
