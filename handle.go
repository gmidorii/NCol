package main

import (
	"errors"
	"github.com/labstack/echo"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"os"
)

func GetAllNews() echo.HandlerFunc {
	return func(c echo.Context) error {
		url, err := ReadUrl("GitHub")
		if err != nil {
			return err
		}
		return c.String(http.StatusOK, url)
	}
}

func ReadUrl(url string) (string, error) {
	c, err := os.Open(config)
	if err != nil {
		return "", err
	}
	defer c.Close()

	file, err := ioutil.ReadAll(c)
	if err != nil {
		return "", err
	}
	var urls []Url
	err = yaml.Unmarshal(file, &urls)
	if err != nil {
		return "", err
	}

	for _, v := range urls {
		if url == v.Name {
			return v.Url, nil
		}
	}
	return "", errors.New("Param name url is not found")
}
