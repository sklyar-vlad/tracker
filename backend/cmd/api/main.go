package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world!"))
	})
	log.Println("service started at localhost:8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", mux))
}
