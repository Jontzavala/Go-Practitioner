package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	// This makes the channel and the type of message the channel will be working with is order.
	var recievedOrdersCh = make(chan order)
	var validateOrdersCh = make(chan order)
	var invalidOrdersCh = make(chan invalidOrder)
	/* this kicks off receiveOrders and an asynchronous task.
	You also need to pass the address WaitGroup in because we don't know when the asynchronous work is
	done until receiveOrders finishes*/
	go receiveOrders(recievedOrdersCh)
	go validateOrders(recievedOrdersCh, validateOrdersCh, invalidOrdersCh)

	wg.Add(1)

	//this will handle messages coming from the validOrdersCh
	go func(validOrderCh <-chan order, invalidOrderCh <-chan invalidOrder){
		select {
		case order := <- validateOrdersCh:
			fmt.Printf("Valid order received: %v\n", order)
		case order := <- invalidOrdersCh:
			fmt.Printf("Invalid order recieved: %v. Issue: %v\n", order.order, order.err)
		}
		wg.Done()

		// to make directional pass the channels in as parameters, then when you receive them you can receive them as directional channels.
	}(validateOrdersCh, invalidOrdersCh)

	wg.Wait()

}

// This is going to receive its orders from the receiveOrders goroutine. Using a channel.
func validateOrders(in <-chan order, out chan<- order, errCh chan<- invalidOrder) {
	order := <- in
	if order.Quantity <= 0 {
		// error condition
		errCh <- invalidOrder{order: order, err: errors.New("quantity must be greater than zero")}
	} else {
		// success path
		out <- order
	}
}

// receiveOrders needs the wait group passed in
func receiveOrders(out chan<- order) {
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
		out <- newOrder
	}
}

var rawOrders = []string{
	`{"productCode": 1111, "quantity": 5, "status": 1}`,
	`{"productCode": 2222, "quantity": 42.3, "status": 1}`,
	`{"productCode": 3333, "quantity": 19, "status": 1}`,
	`{"productCode": 4444, "quantity": 8, "status": 1}`,
}