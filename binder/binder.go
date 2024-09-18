package binder

import (
	"errors"
	"github.com/mbict/webapp"
	"github.com/mbict/webapp/binder/decoder"
	"net/http"
	"reflect"
	"strings"
	"sync"
)

const (
	DefaultTag = "default"
	QueryTag   = "query"
	ParamTag   = "param"
	HeaderTag  = "header"
	CookieTag  = "cookie"
	RequestTag = "request"
)

var ErrUnsupportedType = errors.New("decoder: unsupported type")

var defaultPreDecoders = []*tagDecoders{
	{DefaultTag, defaultsGetterFunc},
}

var defaultPostDecoders = []*tagDecoders{
	{QueryTag, queryGetterFunc},
	{ParamTag, paramsGetterFunc},
	{HeaderTag, headerGetterFunc},
	{CookieTag, cookieGetterFunc},
	{RequestTag, requestGetterFunc},
}

type (
	BinderOption func(*contextBinder)

	tagDecoders struct {
		tag    string
		getter func(webapp.Context) getter
	}

	binder func(webapp.Context, reflect.Value) error

	typeBinder struct {
		pre  []binder
		post []binder
	}

	contextBinder struct {
		contentDecoders        []*contentTypeDecoder
		fallbackContentDecoder func(webapp.Context, interface{}) error
		decoders               map[reflect.Type]*typeBinder
		lock                   sync.RWMutex
	}

	contentTypeDecoder struct {
		canDecode func(c webapp.Context) bool
		decode    func(c webapp.Context, i interface{}) error
	}
)

func WithContentTypeDecoder(canDecode func(webapp.Context) bool, decode func(webapp.Context, interface{}) error) BinderOption {
	return func(b *contextBinder) {
		b.AddDecoder(canDecode, decode)
	}
}

func WithFallbackContentTypeDecoder(canDecode func(webapp.Context) bool, decode func(webapp.Context, interface{}) error) BinderOption {
	return func(b *contextBinder) {
		b.AddDecoder(canDecode, decode)
		b.fallbackContentDecoder = decode
	}
}

func New(options ...BinderOption) webapp.Binder {
	b := &contextBinder{
		contentDecoders: nil,
		decoders:        make(map[reflect.Type]*typeBinder, 0),
		lock:            sync.RWMutex{},
	}

	//add default json decoder
	b.AddDecoder(
		func(c webapp.Context) bool {
			return strings.HasPrefix(c.Request().Header.Get(webapp.HeaderContentType), webapp.MIMEApplicationJSON)
		}, func(c webapp.Context, i interface{}) error {
			if err := webapp.DefaultJSONEncoder.Decode(c, i); err != nil {
				return errors.New("BindError")
			}
			return nil
		})

	////add default xml decoder
	//b.AddDecoder(
	//	func(c webapp.Context) bool {
	//		return strings.HasPrefix(c.Request().Header.Get(webapp.HeaderContentType), webapp.MIMETextXML)
	//	}, func(c webapp.Context, i interface{}) error {
	//		if err := webapp.DefaultXMLEncoder.Decode(c, i); err != nil {
	//			return errors.New("BindError")
	//		}
	//		return nil
	//	})
	//
	////add form multipart post decoder
	//b.AddDecoder(
	//	func(c webapp.Context) bool {
	//		return strings.HasPrefix(c.Request().Header.Get(webapp.HeaderContentType), webapp.MIMEMultipartForm)
	//	}, func(c webapp.Context, i interface{}) error {
	//		if err := webapp.DefaultFormMultipartEncoder.Decode(c, i); err != nil {
	//			return errors.New("BindError")
	//		}
	//		return nil
	//	})

	//run options
	for _, option := range options {
		option(b)
	}
	return b
}

func (b *contextBinder) Bind(c webapp.Context, i interface{}) error {
	return b.bind(c, i, true)
}

func (b *contextBinder) BindQueryParams(c webapp.Context, i interface{}) error {
	return b.bind(c, i, false)
}

func (b *contextBinder) BindBody(c webapp.Context, i interface{}) error {
	//check if a body is provided and can be decoded
	method := c.Request().Method
	if c.Request().ContentLength > 0 && !(method == http.MethodGet || method == http.MethodDelete || method == http.MethodHead || method == http.MethodOptions) {
		for _, decoder := range b.contentDecoders {
			if decoder.canDecode(c) == true {
				return decoder.decode(c, i)
			}
		}

		if b.fallbackContentDecoder != nil {
			return b.fallbackContentDecoder(c, i)
		}
		return ErrUnsupportedType
	}
	return nil
}

func (b *contextBinder) bind(c webapp.Context, i interface{}, bindBody bool) error {
	v := reflect.ValueOf(i)
	t := reflect.Indirect(v).Type()
	if t.Kind() != reflect.Struct {
		return ErrUnsupportedType
	}

	b.lock.RLock()
	decoders, ok := b.decoders[t]
	b.lock.RUnlock()

	//add decoders for this struct
	if !ok {
		preDecoders, err := b.compileBinders(t, defaultPreDecoders)
		if err != nil {
			return err
		}

		postDecoders, err := b.compileBinders(t, defaultPostDecoders)
		if err != nil {
			return err
		}

		decoders = &typeBinder{
			pre:  preDecoders,
			post: postDecoders,
		}

		b.lock.Lock()
		b.decoders[t] = decoders
		b.lock.Unlock()
	}

	//run the decoders that should run before serialize content
	for _, dec := range decoders.pre {
		if err := dec(c, v); err != nil {
			return err
		}
	}

	if bindBody == true {
		if err := b.BindBody(c, i); err != nil {
			return err
		}
	}

	//run the decoder that should run after a body decode, like query, params cookies and headers take precedence over values deserialized from a json
	for _, dec := range decoders.post {
		if err := dec(c, v); err != nil {
			return err
		}
	}

	return nil
}

func (b *contextBinder) compileBinders(t reflect.Type, tagDecoders []*tagDecoders) ([]binder, error) {
	var binders []binder
	for _, d := range tagDecoders {
		if hasTag(t, d.tag) {
			dec, err := decoder.Compile(t, d.tag, true)
			if err != nil {
				return nil, err
			}

			getter := d.getter
			binders = append(binders, func(c webapp.Context, v reflect.Value) error {
				return dec(v, getter(c))
			})
		}
	}
	return binders, nil
}

// AddDecoder
func (b *contextBinder) AddDecoder(canDecode func(c webapp.Context) bool, decode func(c webapp.Context, i interface{}) error) {
	b.contentDecoders = append(b.contentDecoders, &contentTypeDecoder{
		canDecode: canDecode,
		decode:    decode,
	})
}

func hasTag(t reflect.Type, tag string) bool {
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		if f.PkgPath != "" {
			continue // skip unexported fields
		}

		if _, ok := f.Tag.Lookup(tag); ok {
			return true
		}

		ft := f.Type
		if ft.Kind() == reflect.Ptr {
			ft = ft.Elem()
		}

		if ft.Kind() == reflect.Struct {
			if hasTag(ft, tag) {
				return true
			}
		}
	}

	return false
}
