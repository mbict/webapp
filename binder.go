package webappv2

import (
	"errors"
	"net/http"
	"reflect"
	"strings"
	"sync"
)

const (
	DefaultTag = "default"
	QueryTag   = "query"
	PathTag    = "path"
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
	{PathTag, paramsGetterFunc},
	{HeaderTag, headerGetterFunc},
	{CookieTag, cookieGetterFunc},
	{RequestTag, requestGetterFunc},
}

var DefaultBinder Binder

type (
	Binder interface {
		Bind(c Context, i interface{}) error
		BindBody(c Context, i interface{}) error
		BindQueryParams(c Context, i interface{}) error
	}

	BinderOption func(Binder)

	tagDecoders struct {
		tag    string
		getter func(Context) getter
	}

	binder func(Context, reflect.Value) error

	typeBinder struct {
		pre  []binder
		post []binder
	}

	ContextBinder struct {
		contentDecoders []*contentTypeDecoder
		decoders        map[reflect.Type]*typeBinder
		lock            sync.RWMutex
	}

	contentTypeDecoder struct {
		canDecode func(c Context) bool
		decode    func(c Context, i interface{}) error
	}
)

func WithContentTypeDecoder(canDecode func(Context) bool, decode func(Context, interface{}) error) BinderOption {
	return func(b Binder) {
		b.(*ContextBinder).AddDecoder(canDecode, decode)
	}
}

func NewBinder(options ...BinderOption) Binder {
	b := &ContextBinder{
		contentDecoders: nil,
		decoders:        make(map[reflect.Type]*typeBinder, 0),
		lock:            sync.RWMutex{},
	}
	for _, option := range options {
		option(b)
	}
	return b
}

func (b *ContextBinder) Bind(c Context, i interface{}) error {
	return b.bind(c, i, true)
}

func (b *ContextBinder) BindQueryParams(c Context, i interface{}) error {
	return b.bind(c, i, false)
}

func (b *ContextBinder) BindBody(c Context, i interface{}) error {
	//check if a body is provided and can be decoded
	method := c.Request().Method
	if c.Request().ContentLength > 0 && !(method == http.MethodGet || method == http.MethodDelete || method == http.MethodHead || method == http.MethodOptions) {
		for _, decoder := range b.contentDecoders {
			if decoder.canDecode(c) == true {
				return decoder.decode(c, i)
			}
		}
		return ErrUnsupportedType
	}
	return nil
}

func (b *ContextBinder) bind(c Context, i interface{}, bindBody bool) error {
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

		postDecoders, err := b.compileBinders(t, defaultPreDecoders)
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

func (b *ContextBinder) compileBinders(t reflect.Type, tagDecoders []*tagDecoders) ([]binder, error) {
	var binders []binder
	for _, d := range tagDecoders {
		if hasTag(t, d.tag) {
			dec, err := compile(t, d.tag, true)
			if err != nil {
				return nil, err
			}

			getter := d.getter
			binders = append(binders, func(c Context, v reflect.Value) error {
				return dec(v, getter(c))
			})
		}
	}
	return binders, nil
}

// AddDecoder
func (b *ContextBinder) AddDecoder(canDecode func(c Context) bool, decode func(c Context, i interface{}) error) {
	b.contentDecoders = append(b.contentDecoders, &contentTypeDecoder{
		canDecode: canDecode,
		decode:    decode,
	})
}

func init() {
	DefaultBinder = NewBinder(
		//json
		WithContentTypeDecoder(func(c Context) bool {
			return strings.HasPrefix(c.Request().Header.Get(HeaderContentType), MIMEApplicationJSON)
		}, func(c Context, i interface{}) error {
			if err := DefaultJSONEncoder.Decode(c, i); err != nil {
				switch err.(type) {
				case *HTTPError:
					return err
				default:
					return NewHTTPErrorWithInternal(http.StatusBadRequest, err)
				}
			}
			return nil
		}))

	//todo add xml

	//todo add form multipart post

	////json
	//case strings.HasPrefix(contentType, MIMEApplicationJSON):
	//	if err = c.(*context).webapp.jsonEncoder.Decode(c, i); err != nil {
	//		switch err.(type) {
	//		case *HTTPError:
	//			return err
	//		default:
	//			return NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
	//		}
	//	}
	//
	////xml
	//case strings.HasPrefix(contentType, MIMEApplicationXML), strings.HasPrefix(contentType, MIMETextXML):
	//	if err = xml.NewDecoder(req.Body).Decode(i); err != nil {
	//		if ute, ok := err.(*xml.UnsupportedTypeError); ok {
	//			return NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Unsupported type error: type=%v, error=%v", ute.Type, ute.Error())).SetInternal(err)
	//		} else if se, ok := err.(*xml.SyntaxError); ok {
	//			return NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Syntax error: line=%v, error=%v", se.Line, se.Error())).SetInternal(err)
	//		}
	//		return NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
	//	}
	//case strings.HasPrefix(contentType, MIMEApplicationForm), strings.HasPrefix(contentType, MIMEMultipartForm):
	//	params, err := c.FormParams()
	//	if err != nil {
	//		return NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
	//	}
	//	if err = b.bindData(i, params, "form"); err != nil {
	//		return NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
	//	}
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

// value getters
var defaultDefaultsGetter = defaultsGetter{}

// DefaultGetterFactory will return a value getter that always returns the key as the value
// int Property `default:"123"` will set the value 123
func defaultsGetterFunc(_ Context) getter {
	return defaultDefaultsGetter
}

type defaultsGetter struct{}

func (_ defaultsGetter) Get(key string) string {
	return key
}

func (_ defaultsGetter) Values(key string) []string {
	return []string{key}
}

func queryGetterFunc(c Context) getter {
	return mapGetter(c.Request().URL.Query())
}

func paramsGetterFunc(c Context) getter {
	return paramsGetter(c.Param)
}

type paramsGetter func(key string) string

func (ps paramsGetter) Get(key string) string {
	return ps(key)
}

func (ps paramsGetter) Values(key string) []string {
	return []string{ps(key)}
}

func cookieGetterFunc(c Context) getter {
	return cookieGetter(c.Request().Cookies())
}

type cookieGetter []*http.Cookie

func (c cookieGetter) Get(key string) string {
	for i := range c {
		if c[i].Name == key {
			return c[i].Value
		}
	}
	return ""
}

func (c cookieGetter) Values(key string) []string {
	var res []string
	for i := range c {
		if c[i].Name == key {
			res = append(res, c[i].Value)
		}
	}
	return res
}

func headerGetterFunc(c Context) getter {
	return mapGetter(c.Request().Header)
}

func requestGetterFunc(c Context) getter {
	return &requestGetter{Context: c}
}

type requestGetter struct {
	Context
}

func (r *requestGetter) Get(key string) string {
	switch key {
	case `remote-addr`:
		return (*r).RealIP()
	case `host`:
		return (*r).Request().Host
	case `method`:
		return (*r).Method()
	case `url`:
		return (*r).Request().URL.String()
	case `url:host`:
		return (*r).Request().URL.Host
	case `url:query`:
		return (*r).Request().URL.RawQuery
	case `url:path`:
		return (*r).Path()
	case `url:scheme`:
		return (*r).Request().URL.Scheme
	}

	return ""
}

func (r *requestGetter) Values(key string) []string {
	if v := r.Get(key); v != "" {
		return []string{v}
	}
	return []string{}
}
