package main

import (
	"goat-module/pkg/handlers"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.Static("/", "public")

	e.GET("/", handlers.Root)

	e.Logger.Fatal(e.Start(":8080"))
}
