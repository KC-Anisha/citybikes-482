package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	loggly "github.com/JamesPEarly/loggly"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
)

type Item struct {
	Time     string
	Id       string
	Stations []Station
}

type Network struct {
	Stations Stations `json:"network"`
}

type Stations struct {
	Stations []Station `json:"stations"`
}

type Station struct {
	EmptySlots int    `json:"empty_slots"`
	FreeBikes  int    `json:"free_bikes"`
	Name       string `json:"name"`
	Extra      Extra  `json:"extra"`
	Id         string `json:"id"`
}

type Extra struct {
	Renting   int `json:"renting"`
	Returning int `json:"returning"`
}

func pollData() {
	// Tag + client init for Loggly
	var tag string = "citybikes-482"
	client := loggly.New(tag)

	// Call Citybikes API
	resp, err := http.Get("https://api.citybik.es/v2/networks/austin")

	if err != nil {
		client.EchoSend("error", "Failed with error: "+err.Error())
	}

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		client.EchoSend("error", "Failed with error: "+err.Error())
	}

	// Parse the JSON and display some info to the terminal
	var network Network
	json.Unmarshal(body, &network)
	formattedData, _ := json.MarshalIndent(network, "", "    ")
	fmt.Println(string(formattedData))

	// Send success message to loggly with response size
	var respSize string = strconv.Itoa(len(body))
	logErr := client.EchoSend("info", "Successful data collection of size: "+respSize)
	if logErr != nil {
		fmt.Println("err: ", logErr)
	}

	// Initialize a session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)
	if err != nil {
		log.Fatalf("Got error initializing AWS: %s", err)
	}

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	// Create item to be added to DynamoDB
	var item Item
	id := uuid.New().String()
	item.Stations = network.Stations.Stations
	item.Id = id
	item.Time = time.Now().Format(time.RFC3339)

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		log.Fatalf("Got error marshalling new network item: %s", err)
	}

	// Create item in table citybikes
	tableName := "akc-citybikes"

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		log.Fatalf("Got error calling PutItem: %s", err)
	}

	fmt.Println("Successfully added to table " + tableName)
}

func main() {
	for range time.Tick(time.Hour * 3) {
		pollData()
	}
}
