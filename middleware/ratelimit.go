package middleware

import (
	"context"
	"math"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type bucket struct {
	tokens     float64
	lastRefill time.Time
	mu         sync.Mutex
}

type Limiter struct {
	capacity     int
	refillPerSec float64
	mu           sync.RWMutex
	buckets      map[string]*bucket
}

func NewRateLimiter(capacity int, refillPerSec float64) *Limiter {
	buckets := make(map[string]*bucket)
	return &Limiter{
		capacity:     capacity,
		refillPerSec: refillPerSec,
		buckets:      buckets,
	}
}

func (l *Limiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr
		}
		b := l.getBucket(ip)
		b.mu.Lock()
		now := time.Now()
		elapsed := now.Sub(b.lastRefill)
		b.tokens += elapsed.Seconds() * l.refillPerSec
		b.lastRefill = now
		if b.tokens > float64(l.capacity) {
			b.tokens = float64(l.capacity)
		}

		allowed := b.tokens >= 1
		var tokenSnap float64
		if allowed {
			b.tokens -= 1
		} else {
			tokenSnap = b.tokens
		}

		b.mu.Unlock()

		if allowed {
			next.ServeHTTP(w, r)
			return
		}
		retryAfter := int(math.Ceil((1.0 - tokenSnap) / l.refillPerSec))
		if retryAfter < 1 {
			retryAfter = 1
		}
		w.Header().Set("Retry-After", strconv.Itoa(retryAfter))
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = w.Write([]byte(`{"error":"too many requests"}`))
	})
}

func (l *Limiter) getBucket(ip string) *bucket {
	l.mu.RLock()
	b, ok := l.buckets[ip]
	l.mu.RUnlock()
	if ok {
		return b
	}

	l.mu.Lock()
	defer l.mu.Unlock()
	if b, ok := l.buckets[ip]; ok {
		return b
	}
	b = &bucket{
		tokens:     float64(l.capacity),
		lastRefill: time.Now(),
	}
	l.buckets[ip] = b
	return b
}

func (l *Limiter) StartEvictor(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			l.evict(60 * time.Second)
		}
	}
}

func (l *Limiter) evict(idleAfter time.Duration) {
	now := time.Now()

	var staleKeys []string
	l.mu.RLock()
	for k, b := range l.buckets {
		b.mu.Lock()
		if now.Sub(b.lastRefill) > idleAfter {
			staleKeys = append(staleKeys, k)
		}
		b.mu.Unlock()
	}
	l.mu.RUnlock()

	l.mu.Lock()
	for _, k := range staleKeys {
		l.buckets[k].mu.Lock()
		if b, ok := l.buckets[k]; ok && now.Sub(b.lastRefill) > idleAfter {
			delete(l.buckets, k)
		}
		l.buckets[k].mu.Unlock()
	}
	l.mu.Unlock()

}
