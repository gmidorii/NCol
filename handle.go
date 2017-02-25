package main

import (
	"errors"
	"github.com/labstack/echo"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"os"
	"time"
	"github.com/koron/go-dproxy"
	"encoding/json"
)

const layout = "2006-01-02"

type Res struct {
	Title string `json:"title"`
	Url   string `json:"url"`
}

func GetAllNews() echo.HandlerFunc {
	return func(c echo.Context) error {
		reses, err := gitHubClient(c.QueryParam("lang"))
		if err != nil {
			return err
		}
		return c.String(http.StatusOK, parseResJson(reses))
	}
}

func parseResJson(reses []Res) string {
	//var j string
	//for _, v := range reses {
	//	resJ, _ := json.Marshal(v)
	//	j += string(resJ)
	//}
	j, _ := json.Marshal(reses)
	return string(j)
}

func gitHubClient(lang string) ([]Res, error) {
	url, err := ReadUrl("GitHub")
	if err != nil {
		return nil, err
	}

	qMap := make(map[string]string, 0)
	yesterday := time.Now().AddDate(0, 0, -1)
	qMap["q"] = "language:" + lang + "+pushed:>" + yesterday.Format(layout)
	qMap["sort"] = "starts"
	qMap["order"] = "desc"
	for k, v := range qMap {
		url += k + "=" + v + "&"
	}
	println(url)

	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, errors.New("GitHub api connecting was failed")
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var v interface{}
	json.Unmarshal(body, &v)

	items, _ := dproxy.New(v).M("items").Array()
	var reses = make([]Res, 0)
	for _, item := range items {
		con := dproxy.New(item)
		var res Res
		res.Title, _ = con.M("description").String()
		res.Url, _ = con.M("html_url").String()
		reses = append(reses, res)
	}

	return reses, nil
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
