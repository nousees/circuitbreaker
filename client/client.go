package main

import (
	"circuit_breaker/apiclient"
	"context"
	"log"
	"time"
)

func main() {
	client1 := apiclient.NewApiClient("http://localhost:8080", 2, 3)
	client2 := apiclient.NewApiClient("http://localhost:8081", 2, 3)

	ctx := context.Background()

	for i := 0; i < 5; i++ {
		resp, err := client1.Call(ctx, "/api/data")
		if err != nil {
			log.Printf("Client 1 error: %v", err)
		} else {
			log.Printf("Client 1 response: %s", resp)
		}
		time.Sleep(1 * time.Second)
	}

	for i := 0; i < 5; i++ {
		resp, err := client2.Call(ctx, "/api/data")
		if err != nil {
			log.Printf("Client 2 error: %v", err)
		} else {
			log.Printf("Client 2 response: %s", resp)
		}
		time.Sleep(1 * time.Second)
	}
}
