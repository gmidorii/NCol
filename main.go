package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

const urls = "urls.yaml"
const configfile = "config.toml"

var fEnv string
var env Config

type Configs struct {
	Settings []Config
}
type Config struct {
	Tag    string
	Newsdb string
	User   string
	Pass   string
}

type Url struct {
	Name string
	Urls []string
}

func main() {
	flag.StringVar(&fEnv, "e", "local", "production environment")
	flag.Parse()
	var configs Configs
	_, err := toml.DecodeFile(configfile, &configs)
	if err != nil {
		println(err)
		return
	}
	for _, v := range configs.Settings {
		if v.Tag == fEnv {
			env = v
			break
		}
	}
	if env.Tag == "" {
		println("not environment")
		return
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/news", GetAllNews())

	e.Logger.Fatal(e.Start(":1323"))
}
