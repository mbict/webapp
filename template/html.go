package template

import (
	"github.com/Masterminds/sprig/v3"
	"github.com/mbict/webapp"
	"html/template"
	"io"
)

// HTMLTemplate will return a new html template renderer
func HTMLTemplate(options ...Option) webapp.Renderer {
	funcMap := sprig.FuncMap()

	funcMap["toHTML"] = func(s string) template.HTML {
		return template.HTML(s)
	}

	funcMap["toJS"] = func(s string) template.JS {
		return template.JS(s)
	}

	t := template.New("").Funcs(funcMap)
	for _, option := range options {
		option(t)
	}

	return &htmlTemplate{template: t}
}

type htmlTemplate struct {
	template *template.Template
}

func (t *htmlTemplate) Render(_ webapp.Context, w io.Writer, name string, data interface{}) error {
	return t.template.ExecuteTemplate(w, name, data)
}
