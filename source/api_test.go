package source_test

import (
	"cuelang.org/go/cue/cuecontext"
	"errors"
	"go-valkyrie.com/cueconfig/source"
	"os"
	"testing"
)

func TestBytesSource_LoadsValidBytes(t *testing.T) {
	cuectx := cuecontext.New()
	data, err := os.ReadFile("testdata/file1.cue")
	if err != nil {
		t.Fatalf("unexpected error when reading file: %v", err)
	}

	s := source.Bytes(data)
	_, err = s.Load(cuectx)
	if err != nil {
		t.Fatalf("expected source to load, got error: %v", err)
	}
}

func TestPathSource_LoadsValidPath(t *testing.T) {
	cuectx := cuecontext.New()
	s, err := source.Path("testdata/file1.cue")

	if err != nil {
		t.Fatalf("unexpected error when constructing pathSource: %v", err)
	}

	_, err = s.Load(cuectx)
	if err != nil {
		t.Fatalf("expected source to load, got error: %v", err)
	}
}

func TestPathSource_FailsOnInvalidPath(t *testing.T) {
	cuectx := cuecontext.New()

	s, err := source.Path("testdata/thisdoesnotexist.cue")
	if err != nil {
		t.Fatalf("unexpected error when constructing pathSource: %v", err)
	}

	_, err = s.Load(cuectx)
	if err == nil {
		t.Fatalf("expected source to fail to load, got no error")
	} else if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected source to fail with os.ErrNotExist, got %v", err)
	}
}

func TestReaderSource_LoadsValidReader(t *testing.T) {
	cuectx := cuecontext.New()

	f, err := os.Open("testdata/file1.cue")
	if err != nil {
		t.Fatalf("unexpected error when opening file: %v", err)
	}
	defer f.Close()

	s, err := source.Reader(f)
	if err != nil {
		t.Fatalf("unexpected error when constructing readerSource: %v", err)
	}

	_, err = s.Load(cuectx)
	if err != nil {
		t.Fatalf("expected source to load, got error: %v", err)
	}
}

func TestReaderSource_ErrorsOnNilReader(t *testing.T) {
	_, err := source.Reader(nil)
	if err == nil {
		t.Fatalf("expected construction to fail, got no error")
	} else if !errors.Is(err, source.ErrBadReader) {
		t.Fatalf("expected construction to fail with source.ErrBadReader, got %v", err)
	}
}

func TestTypeSource_LoadsValidType(t *testing.T) {
	cuectx := cuecontext.New()
	typ := struct {
		Config struct {
			Foo string `json:"foo"`
		} `json:"config"`
	}{}

	s := source.Type(typ)
	_, err := s.Load(cuectx)
	if err != nil {
		t.Fatalf("expected source to load, got error: %v", err)
	}
}
