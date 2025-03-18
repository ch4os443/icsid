package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/icsid/icsid/internal/config"
	"github.com/icsid/icsid/internal/server"
)

var (
	configPath = flag.String("config", "config.yaml", "Caminho para o arquivo de configuração")
)

func main() {
	flag.Parse()

	// Carrega configurações
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Erro ao carregar configurações: %v", err)
	}

	// Inicializa o servidor
	srv, err := server.New(cfg)
	if err != nil {
		log.Fatalf("Erro ao inicializar servidor: %v", err)
	}

	// Canal para sinais de interrupção
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Inicia o servidor em uma goroutine
	go func() {
		if err := srv.Start(); err != nil {
			log.Printf("Erro ao iniciar servidor: %v", err)
		}
	}()

	// Aguarda sinal de interrupção
	<-sigChan

	// Encerra o servidor graciosamente
	if err := srv.Stop(); err != nil {
		log.Printf("Erro ao encerrar servidor: %v", err)
	}
} 