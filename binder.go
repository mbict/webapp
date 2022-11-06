package webappv2

type Binder interface {
	Bind(c Context, i interface{}) error
}
