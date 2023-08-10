package main

import "fmt"

func main() {
	// Buffer allows the message not to be blocked
	ch := make(chan string, 1)

	//creating the channel this way will cause the message to be blocked because it will be waiting on a receiver.
	//ch := make(chan string)

	ch <- "message"

	fmt.Println(<-ch)
}
