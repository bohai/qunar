package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/bohai/qunar/utils"
	_ "github.com/mattn/go-sqlite3"
)

var db2, db3 *sql.DB

const layout = "2006-01-02"

func main() {
	var urlstr, data string
	db2 = dbutils.NewDB2()
	defer db2.Close()
	db3 = dbutils.NewDB3()
	defer db3.Close()
	rows, _ := db2.Query("select url,data FROM fetch")
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&urlstr, &data)
		process(urlstr, data)
	}
}

func process(urlstr, data string) {
	data_byte := []byte(data)
	js, err := simplejson.NewJson(data_byte[1 : len(data_byte)-1])
	if err != nil || js == nil {
		fmt.Println(err)
		return
	}
	for _, v := range js.Get("rt_data").Get("flightGroupInfo").MustMap() {
		a := v.(map[string]interface{})
		rank, _ := a["rank"].(json.Number).Int64()
		if rank == 1 {
			lowpr, _ := a["lowpr"].(json.Number).Int64()
			op, _ := a["op"].(json.Number).Int64()
			url_decode, _ := url.QueryUnescape(urlstr)
			u, _ := url.Parse(url_decode)
			q := u.Query()
			fromCity := q.Get("fromCity")
			toCity := q.Get("toCity")
			fromDay := q.Get("fromDate")
			toDay := q.Get("toDate")
			fromDate, _ := time.Parse(layout, fromDay)
			toDate, _ := time.Parse(layout, toDay)
			insertDB(fromCity, toCity, fromDate, toDate, lowpr, op)
		}
	}

}

func insertDB(fromCity, toCity string, fromDay, toDay time.Time, price, total int64) {
	stmt, err := db3.Prepare("INSERT INTO prices(fromCity, toCity, fromDay, toDay, price, total) values(?, ?, ?, ?, ?, ?)")
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(fromCity, toCity, fromDay, toDay, price, total)
	if err != nil {
		panic(err)
	}
}
