package main

import (
	"context"
	"log"
	"sync"
)

// recives background jobs through a Jobs channel
func backgroundJobs(jobs <-chan func()) {
	for job := range jobs {
		go func(j func()) {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Job panic recovered: %v\n", r)
				}
			}()
			j()
		}(job)
	}
}

// function iterater
func Parallel(jobs []interface{}) func(yield func(int, interface{}) bool) {
	return func(yield func(int, interface{}) bool) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		var wg sync.WaitGroup
		wg.Add(len(jobs))
		for i, job := range jobs {
			go func() {
				defer wg.Done()
				select {
				case <-ctx.Done():
					return
				default:
					if !yield(i, job) {
						cancel()
					}
				}
			}()
		}
	}
}
