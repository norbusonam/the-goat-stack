package handlers

import (
	"goat-module/pkg/templates"

	"github.com/labstack/echo/v4"
)

func Root(c echo.Context) error {
	return templates.Hello("world 🐐").Render(c.Request().Context(), c.Response().Writer)
}
