package main

import (
	"github.com/mbict/webapp"
	"github.com/mbict/webapp/router"
	"github.com/mbict/webapp/validator"
	"log"
)

type JSON map[string]interface{}

func main() {
	router := router.New()

	app := webapp.New(
		webapp.WithRouter(router),                             // custom router
		webapp.WithValidator(validator.PlaygroundValidator()), //use playground as validator middleware
	)

	app.POST("/{pathParam}", func(c webapp.Context) error {
		payload := struct {
			PathParam string `path:"pathParam" validate:"required"`
			ID        int    `json:"id" validate:"required,gte=1,lte=99999" default:"123"`
			Text      string `json:"text" validate:"required"`
		}{}

		if err := c.Bind(&payload); err != nil {
			return c.JSON(400, JSON{
				"type":  "binding",
				"error": err,
			})
		}

		if err := c.Validate(payload); err != nil {
			return c.JSON(400, JSON{
				"type":  "validation",
				"error": err.Error(),
			})
		}

		return c.JSON(200, JSON{
			"payload": payload,
		})

	})

	log.Print(app.Start(":8088"))
}
