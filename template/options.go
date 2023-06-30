package template

import (
	"html/template"
	"io/fs"
)

type Option func(t *template.Template)

func FromFile(filenames ...string) Option {
	return func(t *template.Template) {
		if _, err := t.ParseFiles(filenames...); err != nil {
			panic(err)
		}
	}
}

func FromString(templateStr ...string) Option {
	return func(t *template.Template) {
		for _, s := range templateStr {
			if _, err := t.Parse(s); err != nil {
				panic(err)
			}
		}
	}
}

func FromPattern(pattern string) Option {
	return func(t *template.Template) {
		if _, err := t.ParseGlob(pattern); err != nil {
			panic(err)
		}
	}
}

func FromFS(fs fs.FS, patterns ...string) Option {
	return func(t *template.Template) {
		_, err := t.ParseFS(fs, patterns...)
		if err != nil {
			panic(err)
		}
	}
}
