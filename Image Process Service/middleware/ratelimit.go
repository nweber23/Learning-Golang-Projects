package middleware

import (
	"net/http"
	"sync"
	"time"
)

type RateLimiter struct {
	mu         sync.RWMutex
	limits     map[string]*userLimit
	maxPerHour int
}

type userLimit struct {
	count     int
	resetTime time.Time
}

func NewRateLimiter(maxPerHour int) *RateLimiter {
	return &RateLimiter{
		limits:     make(map[string]*userLimit),
		maxPerHour: maxPerHour,
	}
}

func (rl *RateLimiter) RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := GetUserIDFromContext(r)
		if userID == "" {
			next.ServeHTTP(w, r)
			return
		}
		if !rl.AllowRequest(userID) {
			JSONError(w, http.StatusTooManyRequests, "Rate limit exceeded: max 100 transformations per hour")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (rl *RateLimiter) AllowRequest(userID string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	now := time.Now()
	limit, exists := rl.limits[userID]
	if !exists {
		rl.limits[userID] = &userLimit{
			count:     1,
			resetTime: now.Add(time.Hour),
		}
		return true
	}
	if limit.count >= rl.maxPerHour {
		return false
	}
	limit.count++
	return true
}
