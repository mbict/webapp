package binder

import (
	"reflect"
	"sync"
	"webappv2"
)

const DefaultTag = "default"
const PathTag = "path"
const QueryTag = "query"
const HeaderTag = "header"
const CookieTag = "cookie"
const RequestTag = "request"

type Decode func(c webappv2.Context, v reflect.Value) error

type GetterFactory func(c webappv2.Context) webappv2.getter

type bindConfig struct {
	tag           string
	getterFactory GetterFactory
}

type binder struct {
	tag  string
	bind func(c webappv2.Context, v reflect.Value) error
}

type contextBinder struct {
	binders map[reflect.Type][]*binder
	lock    sync.RWMutex
	configs []*bindConfig
}

func (b *contextBinder) Bind(c webappv2.Context, i interface{}, binderTags ...string) error {

	v := reflect.ValueOf(i)
	t, k, ptr := webappv2.typeKind(v.Type())
	if k != reflect.Struct {
		return webappv2.ErrUnsupportedType
	}

	b.lock.RLock()
	decoders, ok := b.binders[t]
	b.lock.RUnlock()

	if !ok {
		for _, bc := range b.configs {
			if webappv2.typeHasTag(t, bc.tag) {
				dec, err := compile(t, bc.tag, ptr)
				if err != nil {
					return err
				}

				decoders = append(decoders, &binder{
					tag: bc.tag,
					bind: func(c webappv2.Context, v reflect.Value) error {
						return dec(v, bc.getterFactory(c))
					},
				})
			}
		}

		b.lock.Lock()
		b.binders[t] = decoders
		b.lock.Unlock()
	}

	for _, d := range decoders {
		if binderTags != nil && !sliceContains(binderTags, d.tag) {
			continue
		}

		if err := d.bind(c, v); err != nil {
			return err
		}
	}
	return nil
}

func sliceContains[T comparable](tags []T, tag T) bool {
	for _, v := range tags {
		if v == tag {
			return true
		}
	}
	return false
}

func (b *contextBinder) RegisterDecoder(tag string, factory GetterFactory) {
	b.lock.Lock()
	defer b.lock.Unlock()

	b.configs = append(b.configs, &bindConfig{
		tag:           tag,
		getterFactory: factory,
	})
}

var DefaultBinder Binder

func init() {
	b := &contextBinder{}
	b.RegisterDecoder(DefaultTag, DefaultGetterFactory)
	b.RegisterDecoder(QueryTag, QueryGetterFactory)
	b.RegisterDecoder(PathTag, ParamsGetterFactory)
	b.RegisterDecoder(HeaderTag, HeaderGetterFactory)
	b.RegisterDecoder(CookieTag, CookieGetterFactory)
	b.RegisterDecoder(RequestTag, RequestGetterFactory)

	DefaultBinder = b
}
