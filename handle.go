package main

import (
	"github.com/labstack/echo"
	"net/http"
)

func GetAllNews() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, "All News")
	}
}
