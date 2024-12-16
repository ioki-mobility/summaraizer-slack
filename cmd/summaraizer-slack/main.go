package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/ioki-mobility/summaraizer-slack/api"
)

func main() {
	http.HandleFunc("/", api.IndexHandler)
	http.HandleFunc("/healthz", api.HealthzHandler)
	http.HandleFunc("/ready", api.ReadyHandler)

	port := parsePort()
	address := fmt.Sprintf(":%s", port)
	if err := http.ListenAndServe(address, nil); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}

func parsePort() string {
	port := flag.String("port", "8080", "Port to listen on. Default is 8080")
	flag.Parse()

	return *port
}
