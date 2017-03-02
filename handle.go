package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/koron/go-dproxy"
	"github.com/labstack/echo"
	"gopkg.in/yaml.v2"
)

const layout = "2006-01-02"
const dbLayout = "2006-01-02 15:04:05"

type Res struct {
	Items []Item `json:"items"`
}

type Item struct {
	Title string `json:"title"`
	Url   string `json:"url"`
}

func GetAllNews() echo.HandlerFunc {
	return func(c echo.Context) error {
		reses, err := gitHubClient(c.QueryParam("lang"))
		if err != nil {
			return err
		}
		println(len(reses))
		if err = insertDb(reses); err != nil {
			return err
		}
		return c.String(http.StatusOK, parseResJson(Res{Items: reses}))
	}
}

func insertDb(reses []Item) error {
	db, err := sql.Open("mysql", "root:asdfghjkl@tcp(localhost:3306)/ncol?charset=utf8")
	if err != nil {
		println("db")
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT ncol.t_news SET name=?,url=?,inserted_date=?")
	if err != nil {
		println(err.Error())
		return err
	}
	defer stmt.Close()

	t := time.Now()
	for _, v := range reses {
		println(v.Title)
		_, err := stmt.Exec(v.Title, v.Url, t.Format(dbLayout))
		if err != nil {
			println(err.Error())
			return err
		}
	}

	return nil
}

func parseResJson(res Res) string {
	j, _ := json.Marshal(res)
	return string(j)
}

func gitHubClient(lang string) ([]Item, error) {
	url, err := ReadUrl("GitHub")
	if err != nil {
		return nil, err
	}

	qMap := make(map[string]string, 0)
	yesterday := time.Now().AddDate(0, 0, -2)
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
	var reses = make([]Item, 0)
	for _, item := range items {
		con := dproxy.New(item)
		var res Item
		res.Title, _ = con.M("full_name").String()
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
