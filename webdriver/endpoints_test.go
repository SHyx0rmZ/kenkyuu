package webdriver

import (
	"context"
	"reflect"
	"strings"
	"testing"
)

var contextType = reflect.TypeOf((*context.Context)(nil)).Elem()

func TestImplemented(t *testing.T) {
	for _, e := range endpoints {
		t.Run(e.Description, func(t *testing.T) {
			if e.Implementation == nil {
				t.Fatal("not implemented")
			}
			rt := reflect.TypeOf(e.Implementation)
			if rt.Kind() != reflect.Func {
				t.Fatal("not a func")
			}
			no := strings.Count(e.Path, "{")
			nc := strings.Count(e.Path, "}")
			if no != nc {
				t.Skip("could not make sense of path")
			}
			if rt.NumIn() > no+2 || rt.NumIn() < 2 {
				t.Fatalf("wrong number of arguments, got %d, want %d (or less, but at least one)", rt.NumIn()-1, no+1)
			}
			if !rt.In(1).ConvertibleTo(contextType) {
				t.Fatal("first argument must be context.Context")
			}
		})
	}
}
