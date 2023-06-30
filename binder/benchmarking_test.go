package binder

import (
	"reflect"
	"sync"
	"webappv2"
)

type DecodeValue func(c webappv2.Context, v any) error

type Decoder struct {
	tag   string
	cache sync.Map
}

func NewDecoder(tag string) *Decoder {
	return &Decoder{tag: tag}
}

func (d *Decoder) Decode(data webappv2.getter, v any) (err error) {
	val := reflect.ValueOf(v).Elem()
	if val.Kind() != reflect.Struct {
		return webappv2.ErrUnsupportedType
	}

	t := val.Type()

	dec, ok := d.cache.Load(t)
	if !ok {
		dec, err = compile(t, d.tag, t.Kind() == reflect.Ptr)
		if err != nil {
			return err
		}

		d.cache.Store(t, dec)
	}

	return dec.(Decoder)(val, data)
}
