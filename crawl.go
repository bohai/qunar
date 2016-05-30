package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/PuerkitoBio/fetchbot"
	"github.com/bohai/qunar/utils"
	_ "github.com/mattn/go-sqlite3"
)

func handler(ctx *fetchbot.Context, res *http.Response, err error) {
	fmt.Println("enter handler")
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}
	//fmt.Printf("[%d] %s %s\n", res.StatusCode, ctx.Cmd.Method(), ctx.Cmd.URL())
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	urlstr := ctx.Cmd.URL().String()
	stmt, err := db2.Prepare("INSERT INTO fetch(url, data) values(?, ?)")
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(urlstr, string(body))
	if err != nil {
		panic(err)
	}

	fmt.Println("leave handler")
}

var db1, db2 *sql.DB
var stmt sql.Stmt

func main() {
	var url string
	f := fetchbot.New(fetchbot.HandlerFunc(handler))
	f.DisablePoliteness = true
	db1 = dbutils.NewDB1()
	defer db1.Close()
	db2 = dbutils.NewDB2()
	defer db2.Close()

	rows, _ := db1.Query("select url FROM urls")
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&url)
		queue := f.Start()
		queue.SendStringGet(url)
		queue.Close()
	}
}
