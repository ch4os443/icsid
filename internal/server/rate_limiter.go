package server

import (
	"sync"
	"time"
)

type RateLimiter struct {
	attempts map[string][]time.Time
	mu       sync.RWMutex
	window   time.Duration
	maxTries int
}

func NewRateLimiter(window time.Duration, maxTries int) *RateLimiter {
	return &RateLimiter{
		attempts: make(map[string][]time.Time),
		window:   window,
		maxTries: maxTries,
	}
}

func (rl *RateLimiter) Allow(addr string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	attempts := rl.attempts[addr]

	// Remove tentativas antigas
	valid := attempts[:0]
	for _, t := range attempts {
		if now.Sub(t) <= rl.window {
			valid = append(valid, t)
		}
	}

	// Verifica se excedeu o limite
	if len(valid) >= rl.maxTries {
		return false
	}

	// Adiciona nova tentativa
	rl.attempts[addr] = append(valid, now)
	return true
}

func (rl *RateLimiter) Cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	for addr, attempts := range rl.attempts {
		valid := attempts[:0]
		for _, t := range attempts {
			if now.Sub(t) <= rl.window {
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
