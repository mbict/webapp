package webappv2

import "io"

// Renderer is the interface that wraps the Render function.
type Renderer interface {
	Render(Context, io.Writer, string, interface{}) error
}
