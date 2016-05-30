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
	db = dbutils.NewDB3()
	var fromCity, toCity string
	var price, total int
	var fromDay, toDay time.Time
	const layout = "2006-01-02"

	rows, _ := db.Query("select fromCity, toCity, fromDay, toDay, price, total FROM prices")
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&fromCity, &toCity, &fromDay, &toDay, &price, &total)
		fmt.Println(fromCity, toCity, fromDay.Format(layout), toDay.Format(layout), price, total)
	}

	//stmt, _ := db.Prepare("DELETE from prices")
	//stmt.Exec()
}
