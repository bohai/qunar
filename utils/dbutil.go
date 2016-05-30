package dbutils

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func NewDB3() *sql.DB {
	db, err := sql.Open("sqlite3", "prices.sqlite")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("create table if not exists prices(fromCity text, toCity text, fromDay datetime, toDay datetime, price integer, total integer)")
	if err != nil {
		panic(err)
	}

	return db
}

func NewDB1() *sql.DB {
	db, err := sql.Open("sqlite3", "urls.sqlite")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("create table if not exists urls(url text, date datetime)")
	if err != nil {
		panic(err)
	}

	return db
}

func NewDB2() *sql.DB {
	db, err := sql.Open("sqlite3", "fetch.sqlite")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("create table if not exists fetch(url text, data text)")
	if err != nil {
		panic(err)
	}
	return db
}
