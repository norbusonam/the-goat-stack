package main

import (
	"my-module/pkg/templates"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.Static("/", "public")

	e.GET("/", func(c echo.Context) error {
		return templates.Hello("world 🐐").Render(c.Request().Context(), c.Response().Writer)
	})

	e.Logger.Fatal(e.Start(":8080"))
}
