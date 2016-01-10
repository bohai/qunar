package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/PuerkitoBio/fetchbot"
)

func main() {
	fetch_url := "http://flight.qunar.com/twell/longwell?&tags=2&http%3A%2F%2Fwww.travelco.com%2FsearchArrivalAirport=%E4%B8%BD%E6%B1%9F&http%3A%2F%2Fwww.travelco.com%2FsearchDepartureAirport=%E8%A5%BF%E5%AE%89&http%3A%2F%2Fwww.travelco.com%2FsearchDepartureTime=2016-01-21&http%3A%2F%2Fwww.travelco.com%2FsearchReturnTime=2016-01-24&locale=zh&nextNDays=0&op=1&reset=true&searchLangs=zh&searchType=RoundTripFlight&version=thunder&mergeFlag=0&xd=f1452344371000&wyf=0P8%2Fflr0fRPFflU%2FERPHWlt%2F0YA%2FWlr%2F%2FQPFuUd8lyeFlUd%2F%7C1441321882698&fromCity=%E8%A5%BF%E5%AE%89&toCity=%E4%B8%BD%E6%B1%9F&fromDate=2016-01-21&toDate=2016-01-24&fromCode=SIA&toCode=LJG&from=fi_re_search&lowestPrice=null&_token=46688"
	f := fetchbot.New(fetchbot.HandlerFunc(handler))
	f.DisablePoliteness = true
	queue := f.Start()
	queue.SendStringGet(fetch_url)
	queue.Close()
}

func handler(ctx *fetchbot.Context, res *http.Response, err error) {
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}
	fmt.Printf("[%d] %s %s\n", res.StatusCode, ctx.Cmd.Method(), ctx.Cmd.URL())
	body, err := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))
}
