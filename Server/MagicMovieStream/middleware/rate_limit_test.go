package middleware

import (
	"testing"
	"time"
)

func TestRateLimiterAllowsBurstThenBlocks(t *testing.T) {
	rl := NewRateLimiter(1, 3, time.Second)

	for i := 0; i < 3; i++ {
		if !rl.Allow("1.2.3.4") {
			t.Fatalf("request %d should have been allowed (burst)", i+1)
		}
	}
	if rl.Allow("1.2.3.4") {
		t.Fatal("request beyond the burst should be blocked")
	}
	if !rl.Allow("5.6.7.8") {
		t.Fatal("a different client should still be allowed")
	}
}

func TestRateLimiterRefills(t *testing.T) {
	rl := NewRateLimiter(2, 1, time.Second)

	if !rl.Allow("10.0.0.1") {
		t.Fatal("first request should be allowed")
	}
	if rl.Allow("10.0.0.1") {
		t.Fatal("second request should be blocked, burst is 1")
	}

	time.Sleep(600 * time.Millisecond)
	if !rl.Allow("10.0.0.1") {
		t.Fatal("bucket should have refilled one token after 0.5s at rate 2/s")
	}
}

func TestRateLimiterIsolation(t *testing.T) {
	rl := NewRateLimiter(1, 1, time.Minute)

	rl.Allow("a")
	if !rl.Allow("b") {
		t.Fatal("clients should not share buckets")
	}
}
