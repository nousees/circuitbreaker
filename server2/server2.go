package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
)

func main() {
	http.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
		if rand.Float32() < 0.75 {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "Success response from server 2")
	})

	addr := ":8081"
	log.Printf("Starting server 2 on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server 2 failed: %v", err)
	}
}
