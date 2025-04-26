package middleware

import (
	"math"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

type TokenBucket struct {
	mu             sync.Mutex
	tokens         float64
	maxTokens      float64
	refillRate     float64
	lastRefillTime time.Time
}

func NewTokenBucket(maxTokens float64, refillRate float64) *TokenBucket {
	return &TokenBucket{
		tokens:         maxTokens,
		maxTokens:      maxTokens,
		refillRate:     refillRate,
		lastRefillTime: time.Now(),
	}
}

func (tb *TokenBucket) refill() {
	now := time.Now()
	duration := now.Sub(tb.lastRefillTime)
	tokensToAdd := tb.refillRate * duration.Seconds()
	tb.tokens = math.Min(tb.tokens+tokensToAdd, tb.maxTokens)
	tb.lastRefillTime = now
}


func (tb *TokenBucket) Request(tokens float64) bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.refill()
	if tokens <= tb.tokens {
		tb.tokens -= tokens
		return true
	}
	return false
}


var globalTokenBucket = NewTokenBucket(1000, 10) 

func RateLimiter() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if !globalTokenBucket.Request(1) {
			return c.Status(fiber.StatusTooManyRequests).SendString("Too many requests")
		}

		return c.Next()
	}
}
