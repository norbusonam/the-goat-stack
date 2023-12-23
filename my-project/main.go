package main

import (
	"github.com/labstack/echo/v4"
	"my-module/pkg/templates"
)

func main() {
	e := echo.New()

	e.Static("/", "public")

	e.GET("/", func(c echo.Context) error {
		return templates.Hello("world").Render(c.Request().Context(), c.Response().Writer)
	})

	e.Logger.Fatal(e.Start(":8080"))
}
