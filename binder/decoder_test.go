package binder

import (
	"reflect"
	"testing"
	"webappv2"
)

func BenchmarkBinder(b *testing.B) {
	type child struct {
		String string `schema:"string"`
	}

	type test struct {
		String    string   `schema:"string"`
		StringPtr *string  `schema:"string"`
		Int       int      `schema:"int"`
		Int8      int8     `schema:"int8"`
		Int16     int16    `schema:"int16"`
		Int32     int32    `schema:"int32"`
		Int64     int64    `schema:"int64"`
		Uint      uint     `schema:"uint"`
		Uint8     uint8    `schema:"uint8"`
		Uint16    uint16   `schema:"uint16"`
		Uint32    uint32   `schema:"uint32"`
		Uint64    uint64   `schema:"uint64"`
		Float32   float32  `schema:"float32"`
		Float64   float64  `schema:"float64"`
		Bool      bool     `schema:"bool"`
		Strings   []string `schema:"strings"`
		Nested    child
		NestedPtr *child
	}

	in := webappv2.mapGetter{
		"string":  {"string"},
		"strings": {"string", "string"},
		"int":     {"1"},
		"int8":    {"1"},
		"int16":   {"1"},
		"int32":   {"1"},
		"int64":   {"1"},
		"uint":    {"1"},
		"uint8":   {"1"},
		"uint16":  {"1"},
		"uint32":  {"1"},
		"uint64":  {"1"},
		"float32": {"1"},
		"float64": {"1"},
		"bool":    {"true"},
	}

	b.Run("BinderInstance", func(b *testing.B) {

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			out := &test{}

			decoder, err := Bind(out, `schema`)
			if err != nil {
				b.Fatal(err)
			}

			v := reflect.ValueOf(out)
			if err := decoder(v, in); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("BinderReused", func(b *testing.B) {

		decoder, err := Bind(&test{}, `schema`)
		if err != nil {
			b.Fatal(err)
		}

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			out := &test{}
			v := reflect.ValueOf(out)

			if err := decoder(v, in); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("BinderCached", func(b *testing.B) {

		cachedBinder := &CachedBinder{
			types: make(map[reflect.Type][]webappv2.decoder, 0),
		}

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			out := &test{}
			if err := cachedBinder.Bind(in, out); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("WebappContextBinder", func(b *testing.B) {

		cachedBinder := &CachedBinder{
			types: make(map[reflect.Type][]webappv2.decoder, 0),
		}

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			out := &test{}
			if err := cachedBinder.Bind(in, out); err != nil {
				b.Fatal(err)
			}
		}
	})
}
