package main

import (
	"time"
	"fmt"
	"flag"
	"net/http"
	"encoding/json"
	"strings"
)

type Item struct {
	Title string `json:"title"`
	Description string `json:"description"`
	Link string `json:"link"`
}

type Result struct{
	Item []*Item `json:"item"`
}

type Query struct {
	Count int `json:"count"`
	Created string `json:"created"`
	Lang string `json:"lang"`
	Results Result `json:"results"`
}

type Response struct {
	Query Query `json:"query"`
}



func Checking(numberToCheck string, day int, month int, quitChan *chan bool){
	rawResponse, err := http.Get("https://query.yahooapis.com/v1/public/yql?q=select%20*%20from%20rss%20where%20url%20IN%20(%22http%3A%2F%2Fxskt.com.vn%2Frss-feed%2Fmien-nam-xsmn.rss%22%2C%20%22http%3A%2F%2Fxskt.com.vn%2Frss-feed%2Fmien-trung-xsmt.rss%22%2C%20%22http%3A%2F%2Fxskt.com.vn%2Frss-feed%2Fmien-bac-xsmb.rss%22%20)%20and%20title%20like%20%22%25%20" + fmt.Sprintf("%02d", day) + "%2F" + fmt.Sprintf("%02d", month) + "%20%25%22&format=json&diagnostics=true&callback=")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rawResponse.Body.Close()
	responseDecoder := json.NewDecoder(rawResponse.Body)
	var response Response
	if err = responseDecoder.Decode(&response); err != nil {
		fmt.Println(err)
		return
	}
	if response.Query.Count == 0 {
		fmt.Println("No result found.")
		return
	}
	for _, item := range response.Query.Results.Item{
		item.Description = strings.Replace(item.Description, "\n", " ", -1)
		if(strings.Contains(item.Description, numberToCheck)){
			fmt.Println(item.Title)
			fmt.Println(fmt.Sprintf("Found a number ended with %s, Checkit out at %s ", numberToCheck, item.Link))

			close(*quitChan)
			return
		}
	}
}

func Watch(numberToCheck string, day int, month int){
	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan bool)
	for {
		select {
			case  <- ticker.C:
				Checking(numberToCheck, day, month, &quit)
			case <- quit:
				ticker.Stop()
				return
		}
	}
}

func main(){
	var day, month int
	var number string
	flag.StringVar(&number, "number", "", "Number to check")
	flag.IntVar(&day, "day", -1, "Day to check")
	flag.IntVar(&month, "month", -1, "Month to check")
	flag.Parse()
	if(number == "" || day == -1 || month == -1){
		fmt.Println("You must provide number/day/month to check.")
	}else{
		Watch(number, day, month)
	}
}
