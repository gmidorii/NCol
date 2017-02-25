package main

import (
	"errors"
	"github.com/labstack/echo"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const layout = "2006-01-02"

func GetAllNews() echo.HandlerFunc {
	return func(c echo.Context) error {
		url, err := ReadUrl("GitHub")
		if err != nil {
			return err
		}

		qMap := make(map[string]string, 0)
		yesterday := time.Now().AddDate(0, 0, -1)
		qMap["q"] = "language:" + c.QueryParam("lang") + "+pushed:>" + yesterday.Format(layout)
		qMap["sort"] = "starts"
		qMap["order"] = "desc"
		for k, v := range qMap {
			url += k + "=" + v + "&"
		}

		client := http.Client{}
		req, err := http.NewRequest("GET", url, nil)
		res, err := client.Do(req)
		if err != nil {
			return err
		}
		if res.StatusCode != http.StatusOK {
			return c.String(res.StatusCode, "Api connecting was failed")
		}
		body, err := ioutil.ReadAll(res.Body)
		if  err != nil {
			return err
		}

		return c.String(http.StatusOK, string(body))
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
