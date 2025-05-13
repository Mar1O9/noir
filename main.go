package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload"
)

func routes() http.Handler {
	r := mux.NewRouter()
	r.Handle("/", getHome())
	return r
}

var (
	wg       sync.WaitGroup
	JobsChan chan func()
)

func main() {
	JobsChan = make(chan func())
	srv := http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: routes(),
	}

	go backgroundJobs(JobsChan)

	go func() {
		log.Printf("LISTENTING on PORT: %v", os.Getenv("PORT"))
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("Listening Failed with error: %v", err)
		}
	}()

	shutdown, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	<-shutdown.Done()

	ctx, ctxClose := context.WithTimeout(context.Background(), 30*time.Second)
	defer ctxClose()

	log.Println("Shutting down the server")
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Shutting down with error: %v", err)
	}
	log.Println("Ending Background Jobs")
	wg.Wait()
	defer close(JobsChan)

	log.Println("Shutdown Complete")
}
