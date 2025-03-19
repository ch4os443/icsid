package server

import (
	"testing"
	"time"
)

func TestRateLimiter(t *testing.T) {
	// Cria um rate limiter com 5 tentativas em 5 minutos
	limiter := NewRateLimiter(5*time.Minute, 5)

	// Testa tentativas permitidas
	for i := 0; i < 5; i++ {
		if !limiter.Allow("127.0.0.1") {
			t.Errorf("Tentativa %d não deveria ser bloqueada", i+1)
		}
	}

	// Testa tentativa após o limite
	if limiter.Allow("127.0.0.1") {
		t.Error("Tentativa após o limite deveria ser bloqueada")
	}

	// Testa tentativa de outro IP
	if !limiter.Allow("192.168.1.1") {
		t.Error("Tentativa de outro IP deveria ser permitida")
	}
}

func TestRateLimiterCleanup(t *testing.T) {
	// Cria um rate limiter com 5 tentativas em 1 segundo
	limiter := NewRateLimiter(1*time.Second, 5)

	// Faz algumas tentativas
	for i := 0; i < 5; i++ {
		if !limiter.Allow("127.0.0.1") {
			t.Errorf("Tentativa %d não deveria ser bloqueada", i+1)
		}
	}

	// Aguarda o tempo de expiração
	time.Sleep(2 * time.Second)

	// Limpa as tentativas antigas
	limiter.Cleanup()

	// Verifica se novas tentativas são permitidas
	if !limiter.Allow("127.0.0.1") {
		t.Error("Nova tentativa deveria ser permitida após limpeza")
	}
}

func TestRateLimiterConcurrent(t *testing.T) {
	// Cria um rate limiter com 5 tentativas em 5 minutos
	limiter := NewRateLimiter(5*time.Minute, 5)

	// Testa tentativas concorrentes
	done := make(chan bool)
	concurrent := 10

	for i := 0; i < concurrent; i++ {
		go func(id int) {
			allowed := limiter.Allow("127.0.0.1")
			if id < 5 && !allowed {
				t.Errorf("Tentativa %d deveria ser permitida", id+1)
			} else if id >= 5 && allowed {
				t.Errorf("Tentativa %d deveria ser bloqueada", id+1)
			}
			done <- true
		}(i)
	}

	// Aguarda todas as goroutines terminarem
	for i := 0; i < concurrent; i++ {
		<-done
	}
}

func TestRateLimiterDifferentIPs(t *testing.T) {
	// Cria um rate limiter com 5 tentativas em 5 minutos
	limiter := NewRateLimiter(5*time.Minute, 5)

	// Testa diferentes IPs
	ips := []string{"127.0.0.1", "192.168.1.1", "10.0.0.1"}
	for _, ip := range ips {
		// Faz 5 tentativas para cada IP
		for i := 0; i < 5; i++ {
			if !limiter.Allow(ip) {
				t.Errorf("Tentativa %d para IP %s não deveria ser bloqueada", i+1, ip)
			}
		}

		// Verifica se a sexta tentativa é bloqueada
		if limiter.Allow(ip) {
			t.Errorf("Tentativa após o limite para IP %s deveria ser bloqueada", ip)
		}
	}
}

func TestRateLimiterReset(t *testing.T) {
	// Cria um rate limiter com 5 tentativas em 1 segundo
	limiter := NewRateLimiter(1*time.Second, 5)

	// Faz 5 tentativas
	for i := 0; i < 5; i++ {
		if !limiter.Allow("127.0.0.1") {
			t.Errorf("Tentativa %d não deveria ser bloqueada", i+1)
		}
	}

	// Aguarda o tempo de expiração
	time.Sleep(2 * time.Second)

	// Limpa as tentativas antigas
	limiter.Cleanup()

	// Faz mais 5 tentativas
	for i := 0; i < 5; i++ {
		if !limiter.Allow("127.0.0.1") {
			t.Errorf("Tentativa %d após reset não deveria ser bloqueada", i+1)
		}
	}

	// Verifica se a sexta tentativa é bloqueada
	if limiter.Allow("127.0.0.1") {
		t.Error("Tentativa após o limite após reset deveria ser bloqueada")
	}
}
