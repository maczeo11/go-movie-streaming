package middleware

import (
	"math"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter is an in-memory token bucket keyed by client IP.
// A request costs one token; tokens refill at `rate` per second up to
// `burst`. Clients that run out get a 429 until the bucket refills.
type RateLimiter struct {
	mu      sync.Mutex
	buckets map[string]*bucket
	rate    float64
	burst   float64
	window  time.Duration
}

type bucket struct {
	tokens     float64
	lastRefill time.Time
}

func NewRateLimiter(rate, burst float64, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		buckets: make(map[string]*bucket),
		rate:    rate,
		burst:   burst,
		window:  window,
	}
	go rl.prune()
	return rl
}

func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	b, ok := rl.buckets[key]
	if !ok {
		b = &bucket{tokens: rl.burst, lastRefill: now}
		rl.buckets[key] = b
	}

	elapsed := now.Sub(b.lastRefill).Seconds()
	b.tokens = math.Min(rl.burst, b.tokens+elapsed*rl.rate)
	b.lastRefill = now

	if b.tokens < 1 {
		return false
	}
	b.tokens--
	return true
}

// prune drops buckets for clients that have been idle longer than the
// window, so the map doesn't grow without bound.
func (rl *RateLimiter) prune() {
	for {
		time.Sleep(rl.window)
		rl.mu.Lock()
		for key, b := range rl.buckets {
			if time.Since(b.lastRefill) > rl.window {
				delete(rl.buckets, key)
			}
		}
		rl.mu.Unlock()
	}
}

func RateLimit(limiter *RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !limiter.Allow(c.ClientIP()) {
			c.Header("Retry-After", "60")
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests, please slow down",
			})
			return
		}
		c.Next()
	}
}
