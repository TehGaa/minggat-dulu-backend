package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Booking Service starting on :8082...")
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Booking Service OK"))
	})
	log.Fatal(http.ListenAndServe(":8082", nil))
}
