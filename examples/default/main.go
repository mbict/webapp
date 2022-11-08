package main

import (
	"fmt"
	"log"
	"net/http"
	"webappv2"
	"webappv2/router"
)

type JSON map[string]interface{}

func main() {

	app := webappv2.New(
		webappv2.WithRouter(router.New()), // custom router
		webappv2.WithErrorHandlerFallback(customErrorHandler),
	)

	app.GET("/", func(c webappv2.Context) error {
		return c.JSON(http.StatusOK, JSON{"path": "/"})
	})

	app.GET("/test", func(c webappv2.Context) error {
		return c.JSON(http.StatusOK, JSON{"path": "/test"})
	})

	app.GET("/test/*wildcard", func(c webappv2.Context) error {
		return c.JSON(http.StatusOK, JSON{"path": "/test/*wildcard", "paramNames": c.ParamNames(), "paramValues": c.ParamValues()})
	})

	//intention driven urls with  {id}:{intention}
	app.GET("/intention/{id}:rename", func(c webappv2.Context) error {
		return c.JSON(http.StatusOK, JSON{"path": "/intention/{id}:rename", "paramNames": c.ParamNames(), "paramValues": c.ParamValues()})
	})

	app.GET("/intention/{id}:activate", func(c webappv2.Context) error {
		return c.JSON(http.StatusOK, JSON{"path": "/intention/{id}:activate", "paramNames": c.ParamNames(), "paramValues": c.ParamValues()})
	})

	app.GET("/intention/{id}@activate", func(c webappv2.Context) error {
		return c.JSON(http.StatusOK, JSON{"path": "/intention/{id}@activate", "paramNames": c.ParamNames(), "paramValues": c.ParamValues()})
	})

	app.GET("/intention/{id}/activate", func(c webappv2.Context) error {
		return c.JSON(http.StatusOK, JSON{"path": "/intention/{id}/activate", "paramNames": c.ParamNames(), "paramValues": c.ParamValues()})
	})

	app.GET("/intention/{id}", func(c webappv2.Context) error {
		return c.JSON(http.StatusOK, JSON{"path": "/intention/{id}", "paramNames": c.ParamNames(), "paramValues": c.ParamValues()})
	})

	app.GET("/intention/{id}/test", func(c webappv2.Context) error {
		return c.JSON(http.StatusOK, JSON{"path": "/intention/{id}/test", "paramNames": c.ParamNames(), "paramValues": c.ParamValues()})
	})

	//binding
	type bind struct {
		Foo   string `path:"foo"`
		Bar   string `query:"bar"`
		Baz   string `query:"baz" default:"dflt value"`
		Lorem string `json:"lorem"`
	}

	app.GET("/xyz/{foo}", func(c webappv2.Context) error {
		request := &bind{}

		if err := c.Bind(request); err != nil {
			return err
		}

		return c.JSON(http.StatusOK, request)
	})

	app.POST("/test/{foo}", func(c webappv2.Context) error {
		request := &bind{}

		if err := c.Bind(request); err != nil {
			return err
		}

		return c.JSON(http.StatusOK, request)
	})

	log.Print(app.Start(":8088"))
}

func customErrorHandler(c webappv2.Context, err error) error {

	fmt.Println("custom error handler handles", err)

	//we did not handle the error, we let the default error handlder handle it
	return err
}
