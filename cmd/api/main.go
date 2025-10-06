package main

import (
	"log"
	"net/http"
)

func main() {
	config := NewConfig()
	container := NewContainer()

	r := SetupRouter(container.UserHandler)

	log.Printf("Server starting on :%s", config.Server.Port)
	if err := http.ListenAndServe(":"+config.Server.Port, r); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
