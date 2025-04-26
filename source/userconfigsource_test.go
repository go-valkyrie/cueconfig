package source

import (
	"path/filepath"
	"slices"
	"testing"
)

func TestUserConfig_SingleFile(t *testing.T) {
	cfg, err := UserConfig(&UserConfigOptions{
		Vendor:   "Valkyrie",
		App:      "cueconfig",
		Filename: "config.cue",
	})
	if err != nil {
		t.Fatal(err)
	}

	ucs := cfg.(*userConfigSource)
	is := ucs.source.(pathSource)

	cfgPath := string(is)
	parts := make([]string, 0, 6)
	for prev, curr, next := cfgPath, cfgPath, filepath.Dir(cfgPath); prev != next; prev, curr, next = curr, next, filepath.Dir(next) {
		parts = append(parts, filepath.Base(curr))
	}
	slices.Reverse(parts)

	if pl := len(parts); pl < 3 {
		t.Fatalf("expected at least 3 parts, got %d", pl)
	} else if p := parts[pl-1]; p != "config.cue" {
		t.Fatalf("expected last part to be config.cue, got %s", p)
	} else if p := parts[pl-2]; p != "cueconfig" {
		t.Fatalf("expected second to last part to be cueconfig, got %s", p)
	} else if p := parts[pl-3]; p != "Valkyrie" {
		t.Fatalf("expected third to last part to be Valkyrie, got %s", p)
	}
}

func TestUserConfig_Package(t *testing.T) {
	cfg, err := UserConfig(&UserConfigOptions{
		Vendor:   "Valkyrie",
		App:      "cueconfig",
		// No filename provided, should use packageSource
	})
	if err != nil {
		t.Fatal(err)
	}

	ucs := cfg.(*userConfigSource)
	is, ok := ucs.source.(*packageSource)
	if !ok {
		t.Fatalf("expected source to be *packageSource, got %T", ucs.source)
	}

	// Verify the path is correct
	parts := make([]string, 0, 6)
	path := is.path
	for prev, curr, next := path, path, filepath.Dir(path); prev != next; prev, curr, next = curr, next, filepath.Dir(next) {
		parts = append(parts, filepath.Base(curr))
	}
	slices.Reverse(parts)

	if pl := len(parts); pl < 2 {
		t.Fatalf("expected at least 2 parts, got %d", pl)
	} else if p := parts[pl-1]; p != "cueconfig" {
		t.Fatalf("expected last part to be cueconfig, got %s", p)
	} else if p := parts[pl-2]; p != "Valkyrie" {
		t.Fatalf("expected second to last part to be Valkyrie, got %s", p)
	}
}
