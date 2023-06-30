package webapp

import (
	"encoding"
	"reflect"
	"strconv"
	"strings"
	"unsafe"
)

type (
	decoder func(reflect.Value, getter) error

	getter interface {
		Get(string) string
		Values(string) []string
	}

	mapGetter map[string][]string
)

func (m mapGetter) Get(key string) string {
	if m == nil {
		return ""
	}

	if vs, ok := m[key]; ok && len(vs) >= 1 {
		return vs[0]
	}

	return ""
}

func (m mapGetter) Values(key string) []string {
	if m == nil {
		return nil
	}

	return m[key]
}

var unmarshalerType = reflect.TypeOf(new(encoding.TextUnmarshaler)).Elem()

func compile(typ reflect.Type, tagKey string, isPtr bool) (decoder, error) {
	decoders := []decoder{}

	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		if f.PkgPath != "" {
			continue // skip unexported fields
		}

		t, k, ptr := typeKind(f.Type)

		tag, ok := f.Tag.Lookup(tagKey)
		if !ok && k != reflect.Struct {
			continue
		}

		tag, options := parseTag(tag)

		if reflect.PointerTo(t).Implements(unmarshalerType) {
			decoders = append(decoders, decodeTextUnmarshaler(get(ptr, i, t), tag))
			continue
		}

		switch k {
		case reflect.Struct:
			dec, err := compile(t, tagKey, ptr)
			if err != nil {
				return nil, err
			}

			index := i

			decoders = append(decoders, func(v reflect.Value, m getter) error {
				return dec(v.Field(index), m)
			})
		case reflect.String:
			decoders = append(decoders, decodeString(set[string](ptr, i, t), tag))
		case reflect.Int:
			decoders = append(decoders, decodeInt(set[int](ptr, i, t), tag))
		case reflect.Int8:
			decoders = append(decoders, decodeInt8(set[int8](ptr, i, t), tag))
		case reflect.Int16:
			decoders = append(decoders, decodeInt16(set[int16](ptr, i, t), tag))
		case reflect.Int32:
			decoders = append(decoders, decodeInt32(set[int32](ptr, i, t), tag))
		case reflect.Int64:
			decoders = append(decoders, decodeInt64(set[int64](ptr, i, t), tag))
		case reflect.Uint:
			decoders = append(decoders, decodeUint(set[uint](ptr, i, t), tag))
		case reflect.Uint8:
			decoders = append(decoders, decodeUint8(set[uint8](ptr, i, t), tag))
		case reflect.Uint16:
			decoders = append(decoders, decodeUint16(set[uint16](ptr, i, t), tag))
		case reflect.Uint32:
			decoders = append(decoders, decodeUint32(set[uint32](ptr, i, t), tag))
		case reflect.Uint64:
			decoders = append(decoders, decodeUint64(set[uint64](ptr, i, t), tag))
		case reflect.Float32:
			decoders = append(decoders, decodeFloat32(set[float32](ptr, i, t), tag))
		case reflect.Float64:
			decoders = append(decoders, decodeFloat64(set[float64](ptr, i, t), tag))
		case reflect.Bool:
			decoders = append(decoders, decodeBool(set[bool](ptr, i, t), tag))
		case reflect.Slice:

			//slice with a text unmarshaller, time and uuid for example
			if reflect.PointerTo(t.Elem()).Implements(unmarshalerType) {
				decoders = append(decoders, decodeTextUnmarshalerSlice(i, get(ptr, i, t), tag, getDelimiterFromOptions(options)))
				continue
			}

			_, sk, _ := typeKind(t.Elem())
			switch sk {
			case reflect.String:
				decoders = append(decoders, decodeStrings(set[[]string](ptr, i, t), tag, getDelimiterFromOptions(options)))
			case reflect.Uint8:
				decoders = append(decoders, decodeBytes(set[[]byte](ptr, i, t), tag))
			}
		default:
			return nil, ErrUnsupportedType
		}
	}

	if len(decoders) == 0 {
		return func(reflect.Value, getter) error { return nil }, nil
	}

	return func(v reflect.Value, d getter) error {
		if isPtr {
			if v.IsNil() {
				v.Set(reflect.New(typ))
			}

			v = v.Elem()
		}

		for _, dec := range decoders {
			if err := dec(v, d); err != nil {
				return err
			}
		}

		return nil
	}, nil
}

func typeKind(t reflect.Type) (reflect.Type, reflect.Kind, bool) {
	var isPtr bool

	k := t.Kind()
	if k == reflect.Pointer {
		t = t.Elem()
		k = t.Kind()
		isPtr = true
	}

	return t, k, isPtr
}

func typeHasTag(t reflect.Type, tag string) bool {
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
			if typeHasTag(ft, tag) {
				return true
			}
		}
	}

	return false
}

func set[T any](ptr bool, i int, t reflect.Type) func(reflect.Value, T) {
	if ptr {
		return func(v reflect.Value, d T) {
			f := v.Field(i)
			if f.IsNil() {
				f.Set(reflect.New(t))
			}

			*(*T)(unsafe.Pointer(f.Elem().UnsafeAddr())) = d
		}
	}

	return func(v reflect.Value, d T) {
		*(*T)(unsafe.Pointer(v.Field(i).UnsafeAddr())) = d
	}
}

func get(ptr bool, i int, t reflect.Type) func(v reflect.Value) reflect.Value {
	if ptr {
		return func(v reflect.Value) reflect.Value {
			f := v.Field(i)
			if f.IsNil() {
				f.Set(reflect.New(t))
			}

			return f
		}
	}

	return func(v reflect.Value) reflect.Value {
		return v.Field(i).Addr()
	}
}

func decodeTextUnmarshaler(get func(reflect.Value) reflect.Value, k string) decoder {
	return func(v reflect.Value, g getter) error {
		if s := g.Get(k); s != "" {
			return get(v).Interface().(encoding.TextUnmarshaler).UnmarshalText(atob(s))
		}

		return nil
	}
}

func decodeString(set func(reflect.Value, string), k string) decoder {
	return func(v reflect.Value, g getter) error {
		if s := g.Get(k); s != "" {
			set(v, s)
		}

		return nil
	}
}

func decodeInt(set func(reflect.Value, int), k string) decoder {
	return func(v reflect.Value, g getter) error {
		if s := g.Get(k); s != "" {
			n, err := strconv.Atoi(s)
			if err != nil {
				return err
			}

			set(v, n)
		}

		return nil
	}
}

func decodeInt8(set func(reflect.Value, int8), k string) decoder {
	return func(v reflect.Value, g getter) error {
		if s := g.Get(k); s != "" {
			n, err := strconv.ParseInt(s, 10, 8)
			if err != nil {
				return err
			}

			set(v, int8(n))
		}

		return nil
	}
}

func decodeInt16(set func(reflect.Value, int16), k string) decoder {
	return func(v reflect.Value, g getter) error {
		if s := g.Get(k); s != "" {
			n, err := strconv.ParseInt(s, 10, 16)
			if err != nil {
				return err
			}

			set(v, int16(n))
		}

		return nil
	}
}

func decodeInt32(set func(reflect.Value, int32), k string) decoder {
	return func(v reflect.Value, g getter) error {
		if s := g.Get(k); s != "" {
			n, err := strconv.ParseInt(s, 10, 32)
			if err != nil {
				return err
			}

			set(v, int32(n))
		}

		return nil
	}
}

func decodeInt64(set func(reflect.Value, int64), k string) decoder {
	return func(v reflect.Value, g getter) error {
		if s := g.Get(k); s != "" {
			n, err := strconv.Atoi(s)
			if err != nil {
				return err
			}

			set(v, int64(n))
		}

		return nil
	}
}

func decodeFloat32(set func(reflect.Value, float32), k string) decoder {
	return func(v reflect.Value, g getter) error {
		if s := g.Get(k); s != "" {
			f, err := strconv.ParseFloat(s, 32)
			if err != nil {
				return err
			}

			set(v, float32(f))
		}

		return nil
	}
}

func decodeFloat64(set func(reflect.Value, float64), k string) decoder {
	return func(v reflect.Value, g getter) error {
		if s := g.Get(k); s != "" {
			f, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return err
			}

			set(v, f)
		}

		return nil
	}
}

func decodeUint(set func(reflect.Value, uint), k string) decoder {
	return func(v reflect.Value, g getter) error {
		if s := g.Get(k); s != "" {
			n, err := strconv.ParseUint(s, 10, strconv.IntSize)
			if err != nil {
				return err
			}

			set(v, uint(n))
		}

		return nil
	}
}

func decodeUint8(set func(reflect.Value, uint8), k string) decoder {
	return func(v reflect.Value, g getter) error {
		if s := g.Get(k); s != "" {
			n, err := strconv.ParseUint(s, 10, 8)
			if err != nil {
				return err
			}

			set(v, uint8(n))
		}

		return nil
	}
}

func decodeUint16(set func(reflect.Value, uint16), k string) decoder {
	return func(v reflect.Value, g getter) error {
		if s := g.Get(k); s != "" {
			n, err := strconv.ParseUint(s, 10, 16)
			if err != nil {
				return err
			}

			set(v, uint16(n))
		}

		return nil
	}
}

func decodeUint32(set func(reflect.Value, uint32), k string) decoder {
	return func(v reflect.Value, g getter) error {
		if s := g.Get(k); s != "" {
			n, err := strconv.ParseUint(s, 10, 32)
			if err != nil {
				return err
			}

			set(v, uint32(n))
		}

		return nil
	}
}

func decodeUint64(set func(reflect.Value, uint64), k string) decoder {
	return func(v reflect.Value, g getter) error {
		if s := g.Get(k); s != "" {
			n, err := strconv.ParseUint(s, 10, 64)
			if err != nil {
				return err
			}

			set(v, n)
		}

		return nil
	}
}

func decodeBool(set func(reflect.Value, bool), k string) decoder {
	return func(v reflect.Value, g getter) error {
		if s := g.Get(k); s != "" {
			b, err := strconv.ParseBool(s)
			if err != nil {
				return err
			}

			set(v, b)
		}

		return nil
	}
}

func decodeBytes(set func(reflect.Value, []byte), k string) decoder {
	return func(v reflect.Value, g getter) error {
		if s := g.Get(k); s != "" {
			set(v, atob(s))
		}
		return nil
	}
}

func decodeStrings(set func(reflect.Value, []string), k string, delimiter string) decoder {
	return func(v reflect.Value, g getter) error {
		if s := g.Values(k); s != nil {

			if delimiter != "" {
				var res []string
				for _, value := range s {
					res = append(res, strings.Split(value, delimiter)...)
				}
				s = res
			}
			set(v, s)
		}

		return nil
	}
}

func decodeTextUnmarshalerSlice(i int, get func(reflect.Value) reflect.Value, k string, delimiter string) decoder {
	return func(v reflect.Value, g getter) error {
		if s := g.Values(k); s != nil {

			if delimiter != "" {
				var res []string
				for _, value := range s {
					res = append(res, strings.Split(value, delimiter)...)
				}
				s = res
			}

			var values []reflect.Value
			for _, val := range s {
				marshVal := reflect.New(get(v).Type().Elem().Elem())
				if err := marshVal.Interface().(encoding.TextUnmarshaler).UnmarshalText(atob(val)); err != nil {
					return err
				}

				values = append(values, reflect.Indirect(marshVal))
			}
			v.Field(i).Set(reflect.Append(v.Field(i), values...))
		}
		return nil
	}
}

func parseTag(tag string) (string, string) {
	tag, opt, _ := strings.Cut(tag, ",")
	return tag, opt
}

func getDelimiterFromOptions(o string) string {

	if len(o) == 0 {
		return ""
	}
	s := string(o)
	for s != "" {
		var option string
		option, s, _ = strings.Cut(s, ",")
		switch {
		case option == "comma-delimited":
			return ","
		case option == "semicolon-delimited":
			return ";"
		case option == "pipe-delimited":
			return "|"
		case option == "space-delimited":
			return " "
		case option == "tab-delimited":
			return "\t"
		case strings.HasPrefix(option, "delimiter=") == true:
			_, delimiter, _ := strings.Cut(option, "=")
			switch delimiter {
			case "space":
				return " "
			case "comma":
				return ","
			case "semicolon":
				return ";"
			case "pipe":
				return "|"
			case "": //naive way for comma delimiters, as we dont need to escape `field,delimiter=,,nextoption`
				return ","
			default:
				return delimiter
			}
		}
	}
	return ""
}

func btoa(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func atob(s string) []byte {
	sp := unsafe.Pointer(&s)
	b := *(*[]byte)(sp)
	(*reflect.SliceHeader)(unsafe.Pointer(&b)).Cap = (*reflect.StringHeader)(sp).Len
	return b
}
