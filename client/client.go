package main

import (
	"circuit_breaker/apiclient"
	"circuit_breaker/config"
	"context"
	"log"
	"time"
)

func main() {
	cfg := config.Load()

	client1 := apiclient.NewApiClient("http://server1:8080", cfg)
	client2 := apiclient.NewApiClient("http://server2:8081", cfg)

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
