package main

import (
	"log"
	"net/http"

	"openfeed/internal/api"
)

func main() {

	http.Handle("/", http.FileServer(http.Dir("./web")))

	http.HandleFunc("/api/channel/", api.Channel)

	http.HandleFunc("/api/status", api.Status)

	http.HandleFunc("/api/download", api.Download)

	log.Println("OpenFeed started on :8080")

	log.Fatal(http.ListenAndServe(":8080", nil))

}
