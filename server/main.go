package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
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

  http.HandleFunc("/", handler)

  port := ":8080"
  log.Printf("Starting server on %s\n", port)

  err := http.ListenAndServe(port, nil)

  if err != nil {
    log.Fatalf("Server failed: %v\n", err)
  }
}
