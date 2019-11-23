package main

import (
	"fmt"
	"log"
	"noonhack/client"
	"time"
)

func main() {
	client := client.NewClientWithRetry("superman", 3, 5*time.Minute)
	// 1
	listQueues, err := client.ListQueue()
	if err != nil {
		log.Fatal(err)
	}
	// 2
	if err := client.Push(listQueues[0], "some random data to store"); err != nil {
		log.Fatal(err)
	}

	data, err := client.Poll(listQueues[1])
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("data from the ", listQueues[1], " |> ", data.Data)
}
