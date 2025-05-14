package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
)

func main() {
	http.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
		if rand.Float32() < 0.5 {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "Success response from server 1")
	})

	addr := ":8080"
	log.Printf("Starting server 1 on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server 1 failed: %v", err)
	}
}
