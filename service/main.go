package main

import (
	"fmt"
	"imgnheap/app"
	"log"
	"net/http"
	"time"
)

func main() {
	port := 8080
	router := app.RegisterRouter()

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	log.Printf("listening on port %d...\n", port)
	log.Fatal(server.ListenAndServe())
}
