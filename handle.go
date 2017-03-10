package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"database/sql"

	"fmt"
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

func getAllNews() echo.HandlerFunc {
	return func(c echo.Context) error {
		items, err := selectItems()
		if err != nil {
			return err
		}
		println(len(items))
		if len(items) == 0 {
			urls, err := readUrls("GitHubTrend")
			if err != nil {
				return err
			}
			for _, url := range urls {
				itemList, err := gitHubClient(url)
				if err != nil {
					return err
				}
				items = append(items, itemList...)
			}
			if err = insertDb(items); err != nil {
				return err
			}
		}
		println(len(items))
		return c.String(http.StatusOK, parseResJson(Res{Items: items}))
	}
}

func selectItems() ([]Item, error) {
	db, err := sql.Open("mysql", env.User + ":" + env.Pass + env.Newsdb);
	if err != nil {
		return nil, err
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT name, url FROM t_news WHERE inserted_date LIKE ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	res, err := stmt.Query(time.Now().Format(layout) + "%")
	if err != nil {
		return nil, err
	}

	items := make([]Item, 0)
	for res.Next() {
		var item Item
		if err := res.Scan(&item.Title, &item.Url); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func insertDb(reses []Item) error {
	db, err := sql.Open("mysql", env.User + ":" + env.Pass + env.Newsdb);
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

func gitHubClient(url string) ([]Item, error) {
	week := time.Now().AddDate(0, 0, -7)
	client := http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf(url, week.Format(layout)), nil)
	fmt.Printf(url, week.Format(layout))
	println()

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

func readUrls(name string) ([]string, error) {
	c, err := os.Open(urls)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	file, err := ioutil.ReadAll(c)
	if err != nil {
		return nil, err
	}
	var urls []Url
	err = yaml.Unmarshal(file, &urls)
	if err != nil {
		return nil, err
	}

	for _, v := range urls {
		if name == v.Name {
			return v.Urls, nil
		}
	}
	return nil, errors.New("Param name name is not found")
}
