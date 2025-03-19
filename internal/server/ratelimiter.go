package server

import (
	"sync"
	"time"
)

type RateLimiter struct {
	attempts map[string][]time.Time
	window   time.Duration
	limit    int
	mu       sync.RWMutex
}

func NewRateLimiter(window time.Duration, limit int) *RateLimiter {
	return &RateLimiter{
		attempts: make(map[string][]time.Time),
		window:   window,
		limit:    limit,
	}
}

func (rl *RateLimiter) Allow(addr string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-rl.window)

	// Remove tentativas antigas
	attempts := rl.attempts[addr]
	valid := attempts[:0]
	for _, t := range attempts {
		if t.After(windowStart) {
			valid = append(valid, t)
		}
	}
	rl.attempts[addr] = valid

	// Verifica se excedeu o limite
	if len(valid) >= rl.limit {
		return false
	}

	// Adiciona nova tentativa
	rl.attempts[addr] = append(rl.attempts[addr], now)
	return true
}

func (rl *RateLimiter) Cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-rl.window)

	for addr, attempts := range rl.attempts {
		valid := attempts[:0]
		for _, t := range attempts {
			if t.After(windowStart) {
				valid = append(valid, t)
			}
		}
		if len(valid) == 0 {
			delete(rl.attempts, addr)
		} else {
			rl.attempts[addr] = valid
		}
	}
}
