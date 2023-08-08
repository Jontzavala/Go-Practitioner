package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()
		for range time.Tick(500*time.Millisecond) {
			// you have to check and see if the context was canceled
			if ctx.Err() != nil {
				log.Println(ctx.Err())
				return
			}
			fmt.Println("tick!")
		}
	}(ctx)

	time.Sleep(2 * time.Second)
	cancel()

	wg.Wait()
}