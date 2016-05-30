package main

import (
	"net/url"
	"time"

	"github.com/bohai/qunar/utils"
)

func getLatestFriday() time.Time {
	var day0 time.Time
	for i := 0; i < 7; i++ {
		day0 = time.Now().AddDate(0, 0, i)
		if day0.Weekday() == time.Friday {
			break
		}
	}
	return day0
}

func getDayMap() map[time.Time]time.Time {
	t := getLatestFriday()
	daymap := make(map[time.Time]time.Time)
	var day0, day1 time.Time
	for i := 0; i < MAX; i++ {
		day0 = t.AddDate(0, 0, i*7)
		day1 = t.AddDate(0, 0, i*7+3)
		daymap[day0] = day1
	}
	return daymap
}

// t1: start time
// t2: end   time
// city1, code1: departure city and code
// city2, code2: arrive    city and code
func getURL(t1, t2 time.Time, city1, city2, code1, code2 string) string {
	codedURL := "http://flight.qunar.com/twell/longwell?&tags=2&http%3A%2F%2Fwww.travelco.com%2FsearchArrivalAirport=%E5%8E%A6%E9%97%A8&http%3A%2F%2Fwww.travelco.com%2FsearchDepartureAirport=%E8%A5%BF%E5%AE%89&http%3A%2F%2Fwww.travelco.com%2FsearchDepartureTime=2016-06-03&http%3A%2F%2Fwww.travelco.com%2FsearchReturnTime=2016-06-06&locale=zh&nextNDays=0&op=1&reset=true&searchLangs=zh&searchType=RoundTripFlight&version=thunder&mergeFlag=0&xd=f1464442662000&wyf=0P8FlQr%2FfPPylQy0ERPHWlr%2FfYE%2FWlP%2FbEw%2FlUd8lyeFlUd%2F%7C1441321882698&fromCity=%E8%A5%BF%E5%AE%89&toCity=%E5%8E%A6%E9%97%A8&fromDate=2016-06-03&toDate=2016-06-06&fromCode=SIA&toCode=XMN&from=flight_dom_search&lowestPrice=null&_token=30768"

	rawURL, _ := url.QueryUnescape(codedURL)
	u, _ := url.Parse(rawURL)
	q := u.Query()
	q.Set("searchDepartureAirport", city1)
	q.Set("searchArrivalAirport", city2)
	q.Set("searchDepartureTime", t1.Format(layout))
	q.Set("searchReturnTime", t2.Format(layout))
	q.Set("fromCity", city1)
	q.Set("toCity", city2)
	q.Set("fromCode", code1)
	q.Set("toCode", code2)
	q.Set("fromDate", t1.Format(layout))
	q.Set("toDate", t2.Format(layout))
	u.RawQuery = q.Encode()

	return u.String()
}

const layout = "2006-01-02"
const MAX = 1

func main() {
	FromCity := map[string]string{"西安": "XIY"}
	ToCity := map[string]string{"丽江": "LJG"}
	db := dbutils.NewDB1()
	defer db.Close()
	stmt, err := db.Prepare("INSERT INTO urls(url, date) values(?, ?)")
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	days := getDayMap()
	for d1, d2 := range days {
		for city1, code1 := range FromCity {
			for city2, code2 := range ToCity {
				stmt.Exec(getURL(d1, d2, city1, city2, code1, code2), time.Now())
				//fmt.Println(getURL(d1, d2, city1, city2, code1, code2))
			}
		}
	}
}
