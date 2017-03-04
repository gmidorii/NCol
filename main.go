package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

const config = "config.yaml"

type Url struct {
	Name string
	Urls  []string
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/news", GetAllNews())

	e.Logger.Fatal(e.Start(":1323"))
}
