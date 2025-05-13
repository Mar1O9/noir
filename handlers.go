package main

import (
	"log"
	"net/http"
	"time"
)

func getHome() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("HOME"))
	})
}

// an example for how to start a background task in the handler
func getExample() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//
		wg.Add(1)
		task := func() {
			defer wg.Done()
			// Simulate long-running task
			time.Sleep(2 * time.Second)
			log.Println("Completed: Task from HTTP request")
		}

		// Send job to channel (non-blocking if worker is ready)
		JobsChan <- task

		w.Write([]byte("example!!"))
	})
}
