package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	wg.Add(1)
	/* this kicks off receiveOrders and an asynchronous task.
	You also need to pass the address WaitGroup in because we don't know when the asynchronous work is
	done until receiveOrders finishes*/
	go receiveOrders(&wg)

	wg.Wait()

	fmt.Println(orders)
}

// receiveOrders needs the wait group passed in
func receiveOrders(wg *sync.WaitGroup) {
	for _, rawOrder := range rawOrders {
		var newOrder order
		/*in side the json package in the standard library, there is a package called encoding/json.
		This allows you to convert Go object to and from representations.
		The function ".Unmarshal" takes a JSON representation and converts it to a Go object.
		You'll need to pass in the raw order as a []byte (byte slice)*/
		err := json.Unmarshal([]byte(rawOrder), &newOrder)
		if err != nil {
			log.Print(err)
			continue
		}
		orders = append(orders, newOrder)
	}
	wg.Done()
}

var rawOrders = []string{
	`{"productCode": 1111, "quantity": 5, "status": 1}`,
	`{"productCode": 2222, "quantity": 42.3, "status": 1}`,
	`{"productCode": 3333, "quantity": 19, "status": 1}`,
	`{"productCode": 4444, "quantity": 8, "status": 1}`,
}