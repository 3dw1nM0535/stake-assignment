package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
)

// Response type we get from smartbit.com.au API
type Response struct {
	Success bool  `json:"success"`
	Chart   Chart `json:"chart"`
}

// Chart type returns stats from response
type Chart struct {
	Type       string `json:"type"`
	Unit       string `json:"unit"`
	DayAverage int64  `json:"day_average"`
	Data       []Data `json:"data"`
}

// Data type returns data points(block intervals) from response
type Data struct {
	When  string  `json:"x"`
	Value float64 `json:"y"`
}

func main() {
	// Get bitcoin block intervals from smartbit.com.au API since 2009
	response, err := http.Get("https://api.smartbit.com.au/v1/blockchain/chart/block-interval?&unit=minute")
	if err != nil {
		log.Fatal(err)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var cleanResponse Response
	json.Unmarshal(responseData, &cleanResponse)

	// Available events from smartbit.com.au
	var all float64

	// Holds occurrences of block intervals > 2hrs(120min) from smartbit.com.au
	var gtTwoHrs = make([]float64, 0)

	fmt.Println("-----------------------")
	fmt.Println("How often does the Bitcoin network see two consecutive blocks mined more than 2 hours apart from each other?")

	// HOW OFTEN CALCULATION

	// Looping to get all possible outcomes(total block interval processed daily since 2009)
	for i := 0; i < len(cleanResponse.Chart.Data); i++ {
		all += cleanResponse.Chart.Data[i].Value
	}
	events := all // According to smartbit.com.au source, we have all the likely events to have consecutive block mined more than 2 hours apart from each other(averaged daily)

	// Looping all events to find occurrence of the above event in question
	for i := 0; i < len(cleanResponse.Chart.Data); i++ {
		tmp := cleanResponse.Chart.Data[i]
		if math.Round(tmp.Value) > 120 {
			gtTwoHrs = append(gtTwoHrs, tmp.Value/events) // divide event occurrence by possible outcomes
		}
	}

	// Sum all occurrences
	var likelyHood float64
	for i := 0; i < len(gtTwoHrs); i++ {
		likelyHood += gtTwoHrs[i]
	}
	fmt.Println(likelyHood)
	fmt.Println("-----------------------")
	fmt.Println("How many times the above had happened so far in the history of Bitcoin?")
	fmt.Println(len(gtTwoHrs))
}
