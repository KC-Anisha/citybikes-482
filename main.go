package main

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
	"fmt"
	"net/http"
	loggly "github.com/JamesPEarly/loggly"
)

type Network struct {
	Stations Stations `json:"network"`
}

type Stations struct {
    Stations []Station `json:"stations"`
}

type Station struct {
    EmptySlots   int `json:"empty_slots"`
    FreeBikes   int `json:"free_bikes"`
    Name    string    `json:"name"`
    Extra Extra `json:"extra"`
	Id string `json:"id"`
}

type Extra struct {
    Renting int `json:"renting"`
    Returning  int `json:"returning"`
}

func main() {

	// Tag + client init for Loggly
	var tag string = "citybikes-482"
	client := loggly.New(tag)

	// Call Citybikes API
	resp, err := http.Get("http://api.citybik.es/v2/networks/citi-bike-nyc")

	if err != nil {
		client.EchoSend("error", "Failed with error: " + err.Error())
	}

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		client.EchoSend("error", "Failed with error: " + err.Error())
	}

	// Parse the JSON and display some info to the terminal
	var network Network
	json.Unmarshal(body, &network)
	var stations = network.Stations

	for i := 0; i < len(stations.Stations); i++ {
		fmt.Println("Station Name: ", stations.Stations[i].Name)
		fmt.Println("Empty Slots: ", stations.Stations[i].EmptySlots)
		fmt.Println("Free Bikes: ", stations.Stations[i].FreeBikes)
		fmt.Println("-------------------")
	}


	// Send success message to loggly with response size
	var respSize string = strconv.Itoa(len(body))
	logErr := client.EchoSend("info", "Successful data collection of size: " + respSize)
	if (logErr != nil) {
		fmt.Println("err: ", logErr)
	}
}
