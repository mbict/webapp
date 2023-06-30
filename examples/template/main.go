package main

import (
	"fmt"
	"github.com/mbict/webapp"
	"github.com/mbict/webapp/router"
	"github.com/mbict/webapp/template"
	"log"
)

type JSON map[string]interface{}

func main() {
	router := router.New()
	renderer := template.HTMLTemplate(
		template.FromString(`{{define "index"}}Hello {{ .name }}{{end}}`),
	)

	app := webapp.New(
		webapp.WithRouter(router), // custom router
		webapp.WithRenderer(renderer),
		webapp.WithErrorHandlerFallback(customErrorHandler),
	)

	app.GET("/", func(c webapp.Context) error {
		return c.Render(200, "index", JSON{"name": "world"})
	})

	log.Print(app.Start(":8088"))
}

func customErrorHandler(c webapp.Context, err error) error {

	fmt.Println("custom error handler handles", err)

	//we did not handle the error, we let the default error handlder handle it
	return err
}
