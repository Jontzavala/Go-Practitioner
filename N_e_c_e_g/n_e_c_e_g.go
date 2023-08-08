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
	recievedOrdersCh := receiveOrders()
	validateOrdersCh, invalidOrdersCh := validateOrders(recievedOrdersCh)

	wg.Add(1)

	//this will handle messages coming from the validOrdersCh
	go func(validOrderCh <-chan order, invalidOrderCh <-chan invalidOrder) {
		loop:
		for {
			select {
			case order, ok := <-validateOrdersCh:
				if ok {
					fmt.Printf("Valid order received: %v\n", order)
				} else {
					break loop
				}
			case order, ok := <-invalidOrdersCh:
				if ok {
					fmt.Printf("Invalid order recieved: %v. Issue: %v\n", order.order, order.err)
				} else {
					break loop
				}
			}
		}
		wg.Done()

		// to make directional pass the channels in as parameters, then when you receive them you can receive them as directional channels.
	}(validateOrdersCh, invalidOrdersCh)

	wg.Wait()

}

func validateOrders(in <-chan order) (<-chan order, <-chan invalidOrder) {
	out := make(chan order)
	errCh := make(chan invalidOrder, 1)
	go func() {
		for order := range in {

			if order.Quantity <= 0 {
				// error condition
				errCh <- invalidOrder{order: order, err: errors.New("quantity must be greater than zero")}
			} else {
				// success path
				out <- order
			}
		}
		close(out)
		close(errCh)
	}()
	return out, errCh

}

// encapsulating the goroutine
func receiveOrders() chan order {
	out := make(chan order)
	go func() {
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
		close(out)
	}()
	return out
}

var rawOrders = []string{
	`{"productCode": 1111, "quantity": 5, "status": 1}`,
	`{"productCode": 2222, "quantity": 42.3, "status": 1}`,
	`{"productCode": 3333, "quantity": 19, "status": 1}`,
	`{"productCode": 4444, "quantity": 8, "status": 1}`,
}