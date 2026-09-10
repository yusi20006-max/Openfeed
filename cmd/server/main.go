package main

import (
	"log"
	"net/http"
	"os"

	"openfeed/internal/api"
)

const defaultPort = "7006"

func main() {

	http.Handle("/", http.FileServer(http.Dir("./web")))

	http.HandleFunc("/api/channel/", api.Channel)

	http.HandleFunc("/api/status", api.Status)

	http.HandleFunc("/api/download", api.Download)

	http.HandleFunc("/api/image", api.Image)

	port := os.Getenv("OPENFEED_PORT")
	if port == "" {
		port = defaultPort
	}
	addr := ":" + port

	log.Printf("OpenFeed started on %s", addr)

	log.Fatal(http.ListenAndServe(addr, nil))

}
