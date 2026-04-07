package tokenbucket

import (
	"testing"
	"time"
)

// Test 1 — bucket starts full
func TestBucketStartsFull(t *testing.T) {
	tb := NewTokenBucket(5, 1)
	for i := 0; i < 5; i++ {
		if !tb.Allow() {
			t.Errorf("expected request %d to be allowed", i+1)
		}
	}
}

// Test 2 — bucket denies after empty
func TestBucketDeniesWhenEmpty(t *testing.T) {
	tb := NewTokenBucket(3, 0) // refillRate=0, no refill
	tb.Allow()
	tb.Allow()
	tb.Allow()

	if tb.Allow() {
		t.Error("expected request to be denied when bucket empty")
	}
}

// Test 3 — bucket refills over time
func TestBucketRefills(t *testing.T) {
	tb := NewTokenBucket(5, 2) // 2 tokens per second
	tb.tokens = 0              // drain manually

	time.Sleep(1 * time.Second) // wait 1 second → 2 tokens added

	if !tb.Allow() {
		t.Error("expected request to be allowed after refill")
	}
}

// Test 4 — bucket never exceeds capacity
func TestBucketNeverExceedsCapacity(t *testing.T) {
	tb := NewTokenBucket(3, 100) // very fast refill
	time.Sleep(1 * time.Second)
	tb.refill()

	if tb.tokens > float64(tb.capacity) {
		t.Errorf("tokens %f exceeded capacity %d", tb.tokens, tb.capacity)
	}
}

// Test 5 — concurrent requests are thread safe
func TestConcurrentRequests(t *testing.T) {
	tb := NewTokenBucket(100, 50)
	done := make(chan bool)

	for i := 0; i < 50; i++ {
		go func() {
			tb.Allow()
			done <- true
		}()
	}

	for i := 0; i < 50; i++ {
		<-done
	}
	// no race condition → test passes ✅
}
