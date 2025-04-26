package cueconfig_test

import (
	_ "embed"
	"go-valkyrie.com/cueconfig"
	"go-valkyrie.com/cueconfig/schema"
	"go-valkyrie.com/cueconfig/source"
	"testing"
)

//go:embed testdata/config1_schema.cue
var config1Schema []byte

//go:embed testdata/config1_config.cue
var config1Config []byte

func TestConfig_Unmarshal(t *testing.T) {
	cfgSchema, err := schema.New(source.Bytes(config1Schema))
	if err != nil {
		t.Fatal(err)
	}

	cfg, err := cueconfig.New(cfgSchema, cueconfig.WithSources(
		source.Bytes(config1Config)))
	if err != nil {
		t.Fatal(err)
	}

	config := struct {
		Foo struct {
			Bar string `json:"bar"`
		} `json:"foo"`
	}{}

	if err := cfg.Unmarshal(&config); err != nil {
		t.Fatal(err)
	}

	if config.Foo.Bar != "baz" {
		t.Fatalf("expected baz, got %s", config.Foo.Bar)
	}
}

func TestConfig_AddSource(t *testing.T) {
	cfgSchema, err := schema.New(source.Bytes(config1Schema))
	if err != nil {
		t.Fatal(err)
	}

	cfg, err := cueconfig.New(cfgSchema)
	if err != nil {
		t.Fatal(err)
	}

	config := struct {
		Foo struct {
			Bar string `json:"bar"`
		} `json:"foo"`
	}{}

	if err := cfg.Unmarshal(&config); err == nil {
		t.Fatal("expected error when unmarshalling config with no sources")
	}

	if err := cfg.AddSource(source.Bytes(config1Config)); err != nil {
		t.Fatal(err)
	}

	if err := cfg.Reload(); err != nil {
		t.Fatal(err)
	}

	if err := cfg.Unmarshal(&config); err != nil {
		t.Fatal(err)
	}

	if config.Foo.Bar != "baz" {
		t.Fatalf("expected baz, got %s", config.Foo.Bar)
	}
}
