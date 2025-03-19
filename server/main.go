package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Dimashey/sigma-go/server/items"
)

var logChannel = make(chan string, 100)

func startLogger() {
	for msg := range logChannel {
		log.Println(msg)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	time.Sleep(2 * time.Second)

	w.WriteHeader(http.StatusOK)

	_, err := w.Write([]byte("Hello, User!"))
	if err != nil {
		log.Printf("Error writing response: %v\n", err)
	}

	logChannel <- fmt.Sprintf("Handled request %s in %v", r.URL.Path, time.Since(start))
}

func main() {
	go startLogger()

	itemsTransport := items.NewTransport()

	http.HandleFunc("/", handler)


	http.HandleFunc("/v1/items", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			itemsTransport.Create(w, r)
		case http.MethodGet:
      itemsTransport.GetMany(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/v1/items/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
      itemsTransport.GetOne(w, r)
		case http.MethodPut:
      itemsTransport.Update(w, r)
		case http.MethodDelete:
      itemsTransport.Delete(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	port := ":8080"
	log.Printf("Starting server on %s\n", port)

	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalf("Server failed: %v\n", err)
	}
}
