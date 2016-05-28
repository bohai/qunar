package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/PuerkitoBio/fetchbot"
	"github.com/bitly/go-simplejson"
	"github.com/bohai/qunar/utils"
	_ "github.com/mattn/go-sqlite3"
)

func process1(db sql.DB, r io.Reader, url string) {
	body, err := ioutil.ReadAll(r)
	if err != nil {
		fmt.Println(err)
	}
	js, err := simplejson.NewJson(body[1 : len(body)-1])
	if err != nil || js == nil {
		fmt.Println(err)
		return
	}
	for _, v := range js.Get("rt_data").Get("flightGroupInfo").MustMap() {
		a := v.(map[string]interface{})
		rank, _ := a["rank"].(json.Number).Int64()
		if rank == 1 {
			stmt, err := db.Prepare("INSERT INTO prices(fromCity, toCity, fromDay, toDay, price, total) values(?, ?, ?, ?, ?, ?)")
			if err != nil {
				panic(err)
			}
			defer stmt.Close()
			lowpr, _ := a["lowpr"].(json.Number).Int64()
			op, _ := a["op"].(json.Number).Int64()
			fmt.Println("xian", "lijiang", k, l, lowpr, op)
			_, err = stmt.Exec("xian", "lijiang", k, l, lowpr, op)
			if err != nil {
				panic(err)
			}
		}
	}
}

func handler(ctx *fetchbot.Context, res *http.Response, err error) {
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}
	//fmt.Printf("[%d] %s %s\n", res.StatusCode, ctx.Cmd.Method(), ctx.Cmd.URL())
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	js, err := simplejson.NewJson(body[1 : len(body)-1])
	if err != nil || js == nil {
		fmt.Println(err)
		return
	}
	for _, v := range js.Get("rt_data").Get("flightGroupInfo").MustMap() {
		a := v.(map[string]interface{})
		rank, _ := a["rank"].(json.Number).Int64()
		if rank == 1 {
			stmt, err := db.Prepare("INSERT INTO prices(fromCity, toCity, fromDay, toDay, price, total) values(?, ?, ?, ?, ?, ?)")
			if err != nil {
				panic(err)
			}
			defer stmt.Close()
			lowpr, _ := a["lowpr"].(json.Number).Int64()
			op, _ := a["op"].(json.Number).Int64()
			fmt.Println("xian", "lijiang", k, l, lowpr, op)
			_, err = stmt.Exec("xian", "lijiang", k, l, lowpr, op)
			if err != nil {
				panic(err)
			}
		}
	}
}

var db *sql.DB

func main() {
	var url string
	f := fetchbot.New(fetchbot.HandlerFunc(handler))
	f.DisablePoliteness = true
	db = dbutils.NewDB()
	rows, _ := db.Query("select url FROM urls")
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&url)
		queue := f.Start()
		queue.SendStringGet(url)
		queue.Close()
	}
}
