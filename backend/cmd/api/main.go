package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func habitHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello from /habit")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/habit", habitHandler)

	service := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Println("service started at localhost:8080")
	log.Fatal(service.ListenAndServe())
}
