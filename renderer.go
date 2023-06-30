package webapp

import "io"

// Renderer is the interface that wraps the Render function. Usally used for templates
type Renderer interface {
	Render(Context, io.Writer, string, interface{}) error
}
