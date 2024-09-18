package decoder

import (
	"encoding"
	"errors"
	"golang.org/x/exp/constraints"
	"reflect"
	"strconv"
	"strings"
	"unsafe"
)

type BindingError struct {
	Field        string `json:"field"`
	Value        string `json:"value"`
	ErrorMessage string `json:"error"`
}

func (b *BindingError) Error() string {
	return "binding error for field " + b.Field + " with error :" + b.ErrorMessage
}

var ErrUnsupportedType = errors.New("decoder: unsupported type")

type Decoder func(reflect.Value, Getter) error

//nolint:cyclop
func Compile(typ reflect.Type, tagKey string, isPtr bool) (Decoder, error) {
	decoders := []Decoder{}

	unmarshalerType := reflect.TypeOf(new(encoding.TextUnmarshaler)).Elem()

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
			dec, err := Compile(t, tagKey, ptr)
			if err != nil {
				return nil, err
			}

			index := i

			decoders = append(decoders, func(v reflect.Value, m Getter) error {
				return dec(v.Field(index), m)
			})
		case reflect.String:
			decoders = append(decoders, decodeString(set[string](ptr, i, t), tag))
		case reflect.Int:
			decoders = append(decoders, decodeInt(set[int](ptr, i, t), tag, strconv.IntSize))
		case reflect.Int8:
			decoders = append(decoders, decodeInt(set[int8](ptr, i, t), tag, 8))
		case reflect.Int16:
			decoders = append(decoders, decodeInt(set[int16](ptr, i, t), tag, 16))
		case reflect.Int32:
			decoders = append(decoders, decodeInt(set[int32](ptr, i, t), tag, 32))
		case reflect.Int64:
			decoders = append(decoders, decodeInt(set[int64](ptr, i, t), tag, 64))
		case reflect.Uint:
			decoders = append(decoders, decodeUint(set[uint](ptr, i, t), tag, strconv.IntSize))
		case reflect.Uint8:
			decoders = append(decoders, decodeUint(set[uint8](ptr, i, t), tag, 8))
		case reflect.Uint16:
			decoders = append(decoders, decodeUint(set[uint16](ptr, i, t), tag, 16))
		case reflect.Uint32:
			decoders = append(decoders, decodeUint(set[uint32](ptr, i, t), tag, 32))
		case reflect.Uint64:
			decoders = append(decoders, decodeUint(set[uint64](ptr, i, t), tag, 64))
		case reflect.Float32:
			decoders = append(decoders, decodeFloat(set[float32](ptr, i, t), tag, 32))
		case reflect.Float64:
			decoders = append(decoders, decodeFloat(set[float64](ptr, i, t), tag, 64))
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
			case reflect.Int:
				decoders = append(decoders, decodeSignedIntSlice(set[[]int](ptr, i, t), tag, getDelimiterFromOptions(options), strconv.IntSize))
			case reflect.Int8:
				decoders = append(decoders, decodeSignedIntSlice(set[[]int8](ptr, i, t), tag, getDelimiterFromOptions(options), 8))
			case reflect.Int16:
				decoders = append(decoders, decodeSignedIntSlice(set[[]int16](ptr, i, t), tag, getDelimiterFromOptions(options), 16))
			case reflect.Int32:
				decoders = append(decoders, decodeSignedIntSlice(set[[]int32](ptr, i, t), tag, getDelimiterFromOptions(options), 32))
			case reflect.Int64:
				decoders = append(decoders, decodeSignedIntSlice(set[[]int64](ptr, i, t), tag, getDelimiterFromOptions(options), 64))
			case reflect.Uint:
				decoders = append(decoders, decodeUnsignedIntSlice(set[[]uint](ptr, i, t), tag, getDelimiterFromOptions(options), strconv.IntSize))
			//case reflect.Uint8: //clashed with []byte
			//	decoders = append(decoders, decodeUnsignedIntSlice(set[[]uint8](ptr, i, t), tag, getDelimiterFromOptions(options), 8))
			case reflect.Uint16:
				decoders = append(decoders, decodeUnsignedIntSlice(set[[]uint16](ptr, i, t), tag, getDelimiterFromOptions(options), 16))
			case reflect.Uint32:
				decoders = append(decoders, decodeUnsignedIntSlice(set[[]uint32](ptr, i, t), tag, getDelimiterFromOptions(options), 32))
			case reflect.Uint64:
				decoders = append(decoders, decodeUnsignedIntSlice(set[[]uint64](ptr, i, t), tag, getDelimiterFromOptions(options), 64))

			default:
				return nil, ErrUnsupportedType
			}
		default:
			return nil, ErrUnsupportedType
		}
	}

	if len(decoders) == 0 {
		return func(reflect.Value, Getter) error { return nil }, nil
	}

	return func(v reflect.Value, d Getter) error {
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

func decodeTextUnmarshaler(get func(reflect.Value) reflect.Value, k string) Decoder {
	return func(v reflect.Value, g Getter) error {
		if s := g.Get(k); s != "" {
			if err := get(v).Interface().(encoding.TextUnmarshaler).UnmarshalText(atob(s)); err != nil {
				return &BindingError{
					Field:        k,
					Value:        s,
					ErrorMessage: err.Error(),
				}
			}
		}

		return nil
	}
}

func decodeString(set func(reflect.Value, string), k string) Decoder {
	return func(v reflect.Value, g Getter) error {
		if s := g.Get(k); s != "" {
			set(v, s)
		}

		return nil
	}
}

func decodeInt[T constraints.Signed](set func(reflect.Value, T), k string, bitSize int) Decoder {
	return func(v reflect.Value, g Getter) error {
		if s := g.Get(k); s != "" {
			n, err := strconv.ParseInt(s, 10, bitSize)
			if err != nil {
				return &BindingError{
					Field:        k,
					Value:        s,
					ErrorMessage: "cannot convert to integer",
				}
			}

			set(v, T(n))
		}

		return nil
	}
}

func decodeUint[T constraints.Unsigned](set func(reflect.Value, T), k string, bitSize int) Decoder {
	return func(v reflect.Value, g Getter) error {
		if s := g.Get(k); s != "" {
			n, err := strconv.ParseUint(s, 10, bitSize)
			if err != nil {
				return &BindingError{
					Field:        k,
					Value:        s,
					ErrorMessage: "cannot convert to unsigned integer",
				}
			}

			set(v, T(n))
		}

		return nil
	}
}

func decodeFloat[T constraints.Float](set func(reflect.Value, T), k string, bitSize int) Decoder {
	return func(v reflect.Value, g Getter) error {
		if s := g.Get(k); s != "" {
			f, err := strconv.ParseFloat(s, bitSize)
			if err != nil {
				return &BindingError{
					Field:        k,
					Value:        s,
					ErrorMessage: "cannot convert to float",
				}
			}
			set(v, T(f))
		}

		return nil
	}
}

func decodeSignedIntSlice[T constraints.Signed](set func(reflect.Value, []T), k string, delimiter string, bitSize int) Decoder {
	return func(v reflect.Value, g Getter) error {
		if s := g.Values(k); s != nil {
			var res []T

			if delimiter != "" {
				for _, values := range s {
					for _, rawString := range strings.Split(values, delimiter) {
						n, err := strconv.ParseInt(rawString, 10, bitSize)
						if err != nil {
							return &BindingError{
								Field:        k,
								Value:        rawString,
								ErrorMessage: "cannot convert to integer",
							}
						}
						res = append(res, T(n))
					}
				}
				set(v, res)

			} else {
				for _, rawString := range s {
					n, err := strconv.ParseInt(rawString, 10, bitSize)
					if err != nil {
						return &BindingError{
							Field:        k,
							Value:        rawString,
							ErrorMessage: "cannot convert to integer",
						}
					}
					res = append(res, T(n))
				}
			}
			set(v, res)
		}
		return nil
	}
}

func decodeUnsignedIntSlice[T constraints.Unsigned](set func(reflect.Value, []T), k string, delimiter string, bitSize int) Decoder {
	return func(v reflect.Value, g Getter) error {
		if s := g.Values(k); s != nil {
			var res []T

			if delimiter != "" {
				for _, values := range s {
					for _, rawString := range strings.Split(values, delimiter) {
						n, err := strconv.ParseUint(rawString, 10, bitSize)
						if err != nil {
							return &BindingError{
								Field:        k,
								Value:        rawString,
								ErrorMessage: "cannot convert to unsigned integer",
							}
						}
						res = append(res, T(n))
					}
				}
				set(v, res)

			} else {
				for _, rawString := range s {
					n, err := strconv.ParseUint(rawString, 10, bitSize)
					if err != nil {
						return &BindingError{
							Field:        k,
							Value:        rawString,
							ErrorMessage: "cannot convert to unsigned integer",
						}
					}
					res = append(res, T(n))
				}
			}
			set(v, res)
		}
		return nil
	}
}

func decodeBool(set func(reflect.Value, bool), k string) Decoder {
	return func(v reflect.Value, g Getter) error {
		if s := g.Get(k); s != "" {
			b, err := strconv.ParseBool(s)
			if err != nil {
				return &BindingError{
					Field:        k,
					Value:        s,
					ErrorMessage: "cannot convert to boolean",
				}
			}

			set(v, b)
		}

		return nil
	}
}

func decodeBytes(set func(reflect.Value, []byte), k string) Decoder {
	return func(v reflect.Value, g Getter) error {
		if s := g.Get(k); s != "" {
			set(v, atob(s))
		}
		return nil
	}
}

func decodeStrings(set func(reflect.Value, []string), k string, delimiter string) Decoder {
	return func(v reflect.Value, g Getter) error {
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

func decodeTextUnmarshalerSlice(i int, get func(reflect.Value) reflect.Value, k string, delimiter string) Decoder {
	return func(v reflect.Value, g Getter) error {
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
					return &BindingError{
						Field:        k,
						Value:        val,
						ErrorMessage: err.Error(),
					}
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
