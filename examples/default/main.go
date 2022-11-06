package main

import (
	"fmt"
	"log"
	"net/http"
	"webappv2"
	"webappv2/router"
)

func main() {

	app := webappv2.New(
		webappv2.WithRouter(router.New()), // custom router
		webappv2.WithErrorHandlerFallback(customErrorHandler),
	)

	app.GET("/", func(c webappv2.Context) error {
		return c.JSON(http.StatusOK, "hello")
	})

	//app.GET("/test/@test", func(c webappv2.Context) error {
	//	return c.JSON(http.StatusOK, "hello test")
	//})

	//app.GET("/test/@test", func(c webappv2.Context) error {
	//	return c.JSON(http.StatusOK, "hello test")
	//})

	app.GET("/test/blup:test", func(c webappv2.Context) error {
		return c.JSON(http.StatusOK, "hello blup:test")
	})

	app.GET("/test/test:test", func(c webappv2.Context) error {
		return c.JSON(http.StatusOK, "hello test:test")
	})

	app.GET("/foo/@test:test", func(c webappv2.Context) error {
		return c.JSON(http.StatusOK, "hello @test:test")
	})

	app.GET("/foo/@test:foo", func(c webappv2.Context) error {
		return c.JSON(http.StatusOK, "hello @test:foo")
	})

	log.Print(app.Start(":8088"))
}

func customErrorHandler(c webappv2.Context, err error) error {

	fmt.Println("custom error handler handles", err)

	//we did not handle the error, we let the default error handlder handle it
	return err
}
