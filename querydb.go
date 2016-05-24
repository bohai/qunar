package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func NewDB() *sql.DB {
	db, err := sql.Open("sqlite3", "example.sqlite")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("create table if not exists prices(fromCity text, toCity text, fromDay datetime, toDay datetime, price integer, total integer)")
	if err != nil {
		panic(err)
	}
	return db
}

func main() {
	db = NewDB()
	var price, total int
	var fromDay, toDay time.Time
	const layout = "2006-01-02"

	rows, err := db.Query("select fromDay, toDay, price, total FROM prices")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&fromDay, &toDay, &price, &total); err != nil {
			panic(err)
		}
		fmt.Println(fromDay.Format(layout), toDay.Format(layout), price, total)
	}
	stmt, err := db.Prepare("DELETE from prices")
	if err != nil {
		fmt.Println(err)
	}
	stmt.Exec()
}
