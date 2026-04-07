package main

import (
	"fmt"
	"rate-limiter/tokenbucket"
	"time"
)

func main() {
	tb := tokenbucket.NewTokenBucket(3, 1) // 3 tokens, 1/sec refill

	for i := 1; i <= 6; i++ {
		if tb.Allow() {
			fmt.Printf("Request %d → ✅ allowed\n", i)
		} else {
			fmt.Printf("Request %d → ❌ denied\n", i)
		}
		time.Sleep(300 * time.Millisecond)
	}
}
