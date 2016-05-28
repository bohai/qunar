package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/bohai/qunar/utils"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func main() {
	db = dbutils.NewDB()
	var price, total int
	var fromDay, toDay, day time.Time
	var url string
	const layout = "2006-01-02"

	rows, _ := db.Query("select fromDay, toDay, price, total FROM prices")
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&fromDay, &toDay, &price, &total)
		fmt.Println(fromDay.Format(layout), toDay.Format(layout), price, total)
	}

	rows2, _ := db.Query("select url, date FROM urls")
	defer rows2.Close()
	for rows2.Next() {
		rows2.Scan(&url, &day)
		fmt.Println(url, day.Format(layout))
	}

	stmt, _ := db.Prepare("DELETE from prices")
	stmt.Exec()
	stmt, _ = db.Prepare("DELETE from urls")
	defer stmt.Close()
	stmt.Exec()
}
