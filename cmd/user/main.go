package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("User Service starting on :8081...")
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("User Service OK"))
	})
	log.Fatal(http.ListenAndServe(":8081", nil))
}
