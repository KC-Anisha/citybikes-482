package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"fmt"
	"net/http"
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
}

type Extra struct {
    Renting int `json:"renting"`
    Returning  int `json:"returning"`
}

func main() {
	resp, err := http.Get("http://api.citybik.es/v2/networks/citi-bike-nyc")
	if err != nil {
		log.Fatalln(err)
	}

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}


	var network Network

	json.Unmarshal(body, &network)

	for i := 0; i < len(network.Stations.Stations); i++ {
		fmt.Println("Station Name: ", network.Stations.Stations[i].Name)
		fmt.Println("Empty Slots: ", network.Stations.Stations[i].EmptySlots)
		fmt.Println("Free Bikes: ", network.Stations.Stations[i].FreeBikes)
		fmt.Println("-------------------")
	}
}
