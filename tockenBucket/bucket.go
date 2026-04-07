package tokenbucket

import (
	"sync"
	"time"
)

type TokenBucket struct {
	capacity   int        // max tokens bucket can hold
	tokens     float64    // current tokens
	refillRate float64    // tokens added per second
	lastRefill time.Time  // last refill timestamp
	mu         sync.Mutex // thread safety
}

func NewTokenBucket(capacity int, refillRate float64) *TokenBucket {
	return &TokenBucket{
		capacity:   capacity,
		tokens:     float64(capacity), // start full
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

// refill tokens based on elapsed time
func (tb *TokenBucket) refill() {
	now := time.Now()
	elapsed := now.Sub(tb.lastRefill).Seconds()
	tb.tokens += elapsed * tb.refillRate

	// cap at capacity
	if tb.tokens > float64(tb.capacity) {
		tb.tokens = float64(tb.capacity)
	}
	tb.lastRefill = now
}

// Allow checks if request can proceed
func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.refill()

	if tb.tokens >= 1 {
		tb.tokens-- // consume one token
		return true // request allowed ✅
	}
	return false // request denied ❌
}
