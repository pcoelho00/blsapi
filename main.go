package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type BLSResponse struct {
	Status       string   `json:"status"`
	ResponseTime int      `json:"responseTime"`
	Message      []string `json:"message"`
	Results      Result   `json:"Results"`
}

type Result struct {
	Series []Series `json:"series"`
}

type Series struct {
	Data     []Data `json:"data"`
	SeriesID string `json:"seriesID"`
}

type Data struct {
	Year       string     `json:"year"`
	Period     string     `json:"period"`
	PeriodName string     `json:"periodName"`
	Value      string     `json:"value"`
	Footnotes  []Footnote `json:"footnotes"`
}

type Footnote struct {
	Code string `json:"code"`
	Text string `json:"text"`
}

func BLSPrint(data BLSResponse) {
	fmt.Println("Request Status: ", data.Status)
	for _, series := range data.Results.Series {
		fmt.Println("Series ID", series.SeriesID)
		for _, data := range series.Data {
			fmt.Println(data.Year, data.Period, data.Value)
		}
	}
}

func PostRequest(data map[string]interface{}) {
	url := "https://api.bls.gov/publicAPI/v2/timeseries/data/"

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error making Get request", err)
		return
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// var data map[string]interface{}
	var responseData BLSResponse
	if err := json.Unmarshal(body, &responseData); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	BLSPrint(responseData)

}

func GetRequest(seriesID, registrationKey string) {
	// Make the GET request
	url := "https://api.bls.gov/publicAPI/v2/timeseries/data/"

	full_url := url + seriesID + "?registrationkey=" + registrationKey

	response, err := http.Get(full_url)
	if err != nil {
		fmt.Println("Error making GET request:", err)
		return
	}
	defer response.Body.Close()

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// Unmarshal JSON response
	var responseData BLSResponse
	if err := json.Unmarshal(body, &responseData); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	BLSPrint(responseData)

}

func main() {
	args := os.Args[1:] // Exclude the first argument, which is the program name

	// Check if any arguments are provided
	if len(args) > 0 {
		registry := args[0]

		data := map[string]interface{}{
			"seriesid":        []string{"CUUR0000SA0"},
			"startyear":       "2022",
			"endyear":         "2024",
			"catalog":         false,
			"calculations":    false,
			"annualaverage":   false,
			"aspects":         false,
			"registrationkey": registry,
		}

		GetRequest("CUUR0000SA0", registry)
		PostRequest(data)
	} else {
		fmt.Println("No arguments provided")
	}
}
