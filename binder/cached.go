package binder

import (
	"reflect"
	"sync"
	"webappv2"
)

type CachedDecoder struct {
	dec webappv2.decoder
}

func NewCachedDecoder(v interface{}, tag string) (*CachedDecoder, error) {
	t, k, ptr := webappv2.typeKind(reflect.TypeOf(v))
	if k != reflect.Struct {
		return nil, webappv2.ErrUnsupportedType
	}

	dec, err := compile(t, tag, ptr)
	if err != nil {
		return nil, err
	}

	return &CachedDecoder{dec}, nil
}

func (d *CachedDecoder) Decode(data webappv2.getter, v interface{}) error {
	return d.dec(reflect.ValueOf(v).Elem(), data)
}

type CachedBinder struct {
	lock  sync.RWMutex
	types map[reflect.Type][]webappv2.decoder
}

func (d *CachedBinder) Bind(data webappv2.getter, i interface{}) error {

	v := reflect.ValueOf(i)
	t, k, ptr := webappv2.typeKind(v.Type())
	if k != reflect.Struct {
		return webappv2.ErrUnsupportedType
	}

	d.lock.RLock()
	decoders, ok := d.types[t]
	d.lock.RUnlock()

	if !ok {

		dec, err := compile(t, `schema`, ptr)
		if err != nil {
			return err
		}

		decoders = append(decoders, dec)

		d.lock.Lock()
		d.types[t] = decoders
		d.lock.Unlock()
	}

	return decoders[0](v, data)
}

type WebAppCachedBinder struct {
	lock  sync.RWMutex
	types map[reflect.Type][]webappv2.decoder
}

func (d *WebAppCachedBinder) Bind(c webappv2.Context, i interface{}) error {

	v := reflect.ValueOf(i)

	d.lock.RLock()
	decoders, ok := d.types[v.Type()]
	d.lock.RUnlock()
	if !ok {
		t, k, ptr := webappv2.typeKind(reflect.TypeOf(i))
		if k != reflect.Struct {
			return webappv2.ErrUnsupportedType
		}

		dec, err := compile(t, `schema`, ptr)
		if err != nil {
			return err
		}

		d.lock.Lock()
		d.types[v.Type()] = append(d.types[v.Type()], dec)
		d.lock.Unlock()
	}

	return decoders[0](v, webappv2.mapGetter(c.Request().URL.Query()))
}
