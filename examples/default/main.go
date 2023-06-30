package main

import (
	"fmt"
	"github.com/mbict/webapp"
	"github.com/mbict/webapp/router"
	"log"
	"net/http"
)

type JSON map[string]interface{}

func main() {
	notFoundRouter := router.New()
	notFoundRouter.GET("/{user}/{grid}", func(c webapp.Context) error {
		return nil
	})

	r := router.New(router.WithNotFoundHandler(notFoundRouter))

	app := webapp.New(
		webapp.WithRouter(r), // custom router
		webapp.WithErrorHandlerFallback(customErrorHandler),
	)

	app.GET("/", func(c webapp.Context) error {
		return c.JSON(http.StatusOK, JSON{"path": "/"})
	})

	app.GET("/test", func(c webapp.Context) error {
		return c.JSON(http.StatusOK, JSON{"path": "/test"})
	})

	app.GET("/test/*wildcard", func(c webapp.Context) error {
		return c.JSON(http.StatusOK, JSON{"path": "/test/*wildcard", "paramNames": c.ParamNames(), "paramValues": c.ParamValues()})
	})

	//intention driven urls with  {id}:{intention}
	app.GET("/intention/{id}:rename", func(c webapp.Context) error {
		return c.JSON(http.StatusOK, JSON{"path": "/intention/{id}:rename", "paramNames": c.ParamNames(), "paramValues": c.ParamValues()})
	})

	app.GET("/intention/{id}:activate", func(c webapp.Context) error {
		return c.JSON(http.StatusOK, JSON{"path": "/intention/{id}:activate", "paramNames": c.ParamNames(), "paramValues": c.ParamValues()})
	})

	app.GET("/intention/{id}@activate", func(c webapp.Context) error {
		return c.JSON(http.StatusOK, JSON{"path": "/intention/{id}@activate", "paramNames": c.ParamNames(), "paramValues": c.ParamValues()})
	})

	app.GET("/intention/{id}/activate", func(c webapp.Context) error {
		return c.JSON(http.StatusOK, JSON{"path": "/intention/{id}/activate", "paramNames": c.ParamNames(), "paramValues": c.ParamValues()})
	})

	app.GET("/intention/{id}", func(c webapp.Context) error {
		return c.JSON(http.StatusOK, JSON{"path": "/intention/{id}", "paramNames": c.ParamNames(), "paramValues": c.ParamValues()})
	})

	app.GET("/intention/{id}/test", func(c webapp.Context) error {
		return c.JSON(http.StatusOK, JSON{"path": "/intention/{id}/test", "paramNames": c.ParamNames(), "paramValues": c.ParamValues()})
	})

	//binding
	type bind struct {
		Foo   string `path:"foo"`
		Bar   string `query:"bar"`
		Baz   string `query:"baz" default:"dflt value"`
		Lorem string `json:"lorem"`
	}

	app.GET("/xyz/{foo}", func(c webapp.Context) error {
		request := &bind{}

		if err := c.Bind(request); err != nil {
			return err
		}

		return c.JSON(http.StatusOK, request)
	})

	app.POST("/test/{foo}", func(c webapp.Context) error {
		request := &bind{}

		if err := c.Bind(request); err != nil {
			return err
		}

		return c.JSON(http.StatusOK, request)
	})

	log.Print(app.Start(":8088"))
}

func customErrorHandler(c webapp.Context, err error) error {

	fmt.Println("custom error handler handles", err)

	//we did not handle the error, we let the default error handlder handle it
	return err
}
