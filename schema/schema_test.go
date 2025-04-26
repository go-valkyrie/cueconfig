package schema_test

import (
	"go-valkyrie.com/cueconfig/schema"
	"go-valkyrie.com/cueconfig/source"
	"testing"
)

func TestSchema_Simple(t *testing.T) {
	if src, err := source.Path("testdata/config1_schema.cue"); err != nil {
		t.Fatal(err)
	} else if schema, err := schema.New(src); err != nil {
		t.Fatal(err)
	} else if !schema.Valid() {
		t.Fatal("schema is invalid")
	}
}

//
//func TestManagerOne(t *testing.T) {
//	manager, err := New(&Options{
//		Schema: FileSchema("testdata/schema2.cue"),
//		Schema: FileSource("testdata/schema2_config.cue"),
//	})
//
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	value := &struct {
//		Schema struct {
//			Path string `json:"path"`
//		} `json:"value"`
//	}{}
//
//	if err := manager.Unmarshal(value); err != nil {
//		t.Fatal(err)
//	}
//
//	if value.Schema.Path != "/test" {
//		t.Fatalf("expected /test, got %v", value.Schema.Path)
//	}
//}
