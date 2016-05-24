package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/PuerkitoBio/fetchbot"
	"github.com/bitly/go-simplejson"
	_ "github.com/mattn/go-sqlite3"
)

func getLatestFriday() time.Time {
	day0 := time.Now()
	i := time.Friday - day0.Weekday()
	day1 := day0.AddDate(0, 0, int(i))
	if day1.Before(day0) {
		day1.AddDate(0, 0, 7)
	}
	return day1
}

func getDayMap() map[time.Time]time.Time {
	t := getLatestFriday()
	daymap := make(map[time.Time]time.Time)
	var day0, day1 time.Time
	for i := 0; i < MAX; i++ {
		day0 = t.AddDate(0, 0, i*7)
		day1 = t.AddDate(0, 0, i*7+4)
		daymap[day0] = day1
	}
	return daymap
}

func getURL(t1 time.Time, t2 time.Time) string {
	str0 := "http://flight.qunar.com/twell/longwell?&tags=2&http%3A%2F%2Fwww.travelco.com%2FsearchArrivalAirport=%E4%B8%BD%E6%B1%9F&http%3A%2F%2Fwww.travelco.com%2FsearchDepartureAirport=%E8%A5%BF%E5%AE%89&http%3A%2F%2Fwww.travelco.com%2F"
	str1 := "searchDepartureTime=%s"
	str2 := "&http%3A%2F%2Fwww.travelco.com%2F"
	str3 := "searchReturnTime=%s"
	str4 := "&locale=zh&nextNDays=0&op=1&reset=true&searchLangs=zh&searchType=RoundTripFlight&version=thunder&mergeFlag=0&xd=f1452344371000&wyf=0P8%2Fflr0fRPFflU%2FERPHWlt%2F0YA%2FWlr%2F%2FQPFuUd8lyeFlUd%2F%7C1441321882698&fromCity=%E8%A5%BF%E5%AE%89&toCity=%E4%B8%BD%E6%B1%9F&"
	str5 := "fromDate=%s&toDate=%s"
	str6 := "&fromCode=SIA&toCode=LJG&from=fi_re_search&lowestPrice=null&_token=46688"

	str1 = fmt.Sprintf(str1, t1.Format(layout))
	str3 = fmt.Sprintf(str3, t2.Format(layout))
	str5 = fmt.Sprintf(str5, t1.Format(layout), t2.Format(layout))
	return str0 + str1 + str2 + str3 + str4 + str5 + str6
}

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

const layout = "2006-01-02"
const MAX = 8

var k, l time.Time
var db *sql.DB

func main() {
	db = NewDB()

	for k, l = range getDayMap() {
		f := fetchbot.New(fetchbot.HandlerFunc(handler))
		f.DisablePoliteness = true
		queue := f.Start()
		queue.SendStringGet(getURL(k, l))
		queue.Close()
	}
}
