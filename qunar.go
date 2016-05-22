package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/PuerkitoBio/fetchbot"
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

func getDayMap(t time.Time) map[string]string {
	daymap := make(map[string]string)
	const layout = "2006-01-02"
	const MAX = 1
	var day0, day1 time.Time
	for i := 0; i < MAX; i++ {
		day0 = t.AddDate(0, 0, i*7)
		day1 = t.AddDate(0, 0, i*7+4)
		daymap[day0.Format(layout)] = day1.Format(layout)
	}

	return daymap
}

func getURL(t1 string, t2 string) string {
	str0 := "http://flight.qunar.com/twell/longwell?&tags=2&http%3A%2F%2Fwww.travelco.com%2FsearchArrivalAirport=%E4%B8%BD%E6%B1%9F&http%3A%2F%2Fwww.travelco.com%2FsearchDepartureAirport=%E8%A5%BF%E5%AE%89&http%3A%2F%2Fwww.travelco.com%2F"
	str1 := "searchDepartureTime=%s"
	str2 := "&http%3A%2F%2Fwww.travelco.com%2F"
	str3 := "searchReturnTime=%s"
	str4 := "&locale=zh&nextNDays=0&op=1&reset=true&searchLangs=zh&searchType=RoundTripFlight&version=thunder&mergeFlag=0&xd=f1452344371000&wyf=0P8%2Fflr0fRPFflU%2FERPHWlt%2F0YA%2FWlr%2F%2FQPFuUd8lyeFlUd%2F%7C1441321882698&fromCity=%E8%A5%BF%E5%AE%89&toCity=%E4%B8%BD%E6%B1%9F&"
	str5 := "fromDate=%s&toDate=%s"
	str6 := "&fromCode=SIA&toCode=LJG&from=fi_re_search&lowestPrice=null&_token=46688"

	str1 = fmt.Sprintf(str1, t1)
	str3 = fmt.Sprintf(str3, t2)
	str5 = fmt.Sprintf(str5, t1, t2)
	return str0 + str1 + str2 + str3 + str4 + str5 + str6
}

func main() {
	f := fetchbot.New(fetchbot.HandlerFunc(handler))
	f.DisablePoliteness = true
	queue := f.Start()

	j := getDayMap(getLatestFriday())
	for k, l := range j {
		queue.SendStringGet(getURL(k, l))
	}
	queue.Close()
}

func handler(ctx *fetchbot.Context, res *http.Response, err error) {
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}
	//fmt.Printf("[%d] %s %s\n", res.StatusCode, ctx.Cmd.Method(), ctx.Cmd.URL())
	body, err := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))
}
