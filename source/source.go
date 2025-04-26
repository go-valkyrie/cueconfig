package source

import (
	"cuelang.org/go/cue"
	"cuelang.org/go/cue/load"
	"errors"
	"fmt"
	"go-valkyrie.com/cueconfig/pkgerrs"
	"io"
	"os"
	"path/filepath"
	"slices"
)

var (
	ErrBadReader             = errors.New("cueconfig: cannot create source from nil reader")
	ErrInvalidUserConfigPath = errors.New("cueconfig: cannot create user config source with empty vendor or app")
)

type Source interface {
	Load(ctx *cue.Context) (cue.Value, error)
}

var (
	_ Source = (*bytesSource)(nil)
	_ Source = bytesSource{}
	_ Source = pathSource("")
	_ Source = (*typeSource)(nil)
	_ Source = (*packageSource)(nil)
)

func Bytes(value []byte) Source {
	return bytesSource(slices.Clone(value))
}

type bytesSource []byte

func (s bytesSource) Load(ctx *cue.Context) (cue.Value, error) {
	return ctx.CompileBytes(s), nil
}

func Path(path string) (Source, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	return pathSource(path), nil
}

type pathSource string

func (s pathSource) Load(ctx *cue.Context) (cue.Value, error) {
	if data, err := os.ReadFile(string(s)); err != nil {
		return cue.Value{}, err
	} else {
		return ctx.CompileBytes(data), nil
	}
}

func Reader(reader io.Reader) (Source, error) {
	if reader == nil {
		return nil, ErrBadReader
	}
	if bytes, err := io.ReadAll(reader); err != nil {
		return nil, err
	} else {
		return bytesSource(bytes), nil
	}
}

func Type(value any) Source {
	return &typeSource{value}
}

type typeSource struct {
	value any
}

func (s *typeSource) Load(ctx *cue.Context) (cue.Value, error) {
	return ctx.EncodeType(s.value), nil
}

type userConfigSource struct {
	source    Source
	directory string
	template  map[string][]byte
}

func (s *userConfigSource) Load(ctx *cue.Context) (cue.Value, error) {
	if err := s.ensureExists(); err != nil {
		return cue.Value{}, err
	}

	return s.source.Load(ctx)
}

func (s *userConfigSource) ensureExists() error {
	if _, err := os.Stat(s.directory); err == nil {
		return nil
	}

	if err := os.MkdirAll(s.directory, 0755); err != nil {
		return err
	}

	for k, v := range s.template {
		if err := os.WriteFile(filepath.Join(s.directory, k), v, 0644); err != nil {
			return fmt.Errorf("failed to write template file %s: %w", k, err)
		}
	}

	return nil
}

type UserConfigOptions struct {
	Vendor   string
	App      string
	Filename string
	Template map[string][]byte
}

func UserConfig(opts *UserConfigOptions) (Source, error) {
	if opts == nil {
		return nil, &pkgerrs.ErrInvalidArgument{Argument: "opts", Reason: "cannot be nil"}
	}

	configPath, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	if opts.Vendor == "" || opts.App == "" {
		return nil, ErrInvalidUserConfigPath
	}

	baseDirectory := filepath.Join(configPath, opts.Vendor, opts.App)
	var innerSource Source

	if opts.Filename != "" {
		is, err := Path(filepath.Join(baseDirectory, opts.Filename))
		if err != nil {
			return nil, err
		}

		innerSource = is
	} else {
		is, err := Package(baseDirectory)
		if err != nil {
			return nil, err
		}

		innerSource = is
	}

	template := make(map[string][]byte)
	if opts.Template != nil {
		for k, v := range opts.Template {
			template[k] = slices.Clone(v)
		}
	}

	s := &userConfigSource{
		source:    innerSource,
		directory: baseDirectory,
		template:  template,
	}

	return s, nil
}

func Package(path string) (Source, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	return &packageSource{path: path}, nil
}

type packageSource struct {
	path string
}

func (s *packageSource) Load(ctx *cue.Context) (cue.Value, error) {
	// Use the load package to load the CUE package
	config := &load.Config{
		Dir: s.path,
	}

	// Load the package
	instances := load.Instances([]string{}, config)
	if len(instances) == 0 {
		return cue.Value{}, fmt.Errorf("no CUE files found in directory: %s", s.path)
	}

	// Build the first instance
	value := ctx.BuildInstance(instances[0])
	if value.Err() != nil {
		return cue.Value{}, fmt.Errorf("failed to build instance: %w", value.Err())
	}

	return value, nil
}
