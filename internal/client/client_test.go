package client

import (
	"testing"

	"github.com/icsid/icsid/internal/config"
)

func TestClientCreation(t *testing.T) {
	cfg := &config.Config{}
	cfg.Client.ServerAddress = "localhost:2222"
	cfg.Client.Username = "test"
	cfg.Client.Password = "test"

	client := New(cfg)
	if client == nil {
		t.Fatal("Cliente é nil")
	}

	if client.config != cfg {
		t.Error("Configuração não foi corretamente atribuída")
	}
}

func TestClientInstall(t *testing.T) {
	cfg := &config.Config{}
	cfg.Persistence.Enabled = true
	cfg.Persistence.Method = "registry"

	client := New(cfg)
	if client == nil {
		t.Fatal("Cliente é nil")
	}

	// Testa a instalação
	err := client.Install()
	if err != nil {
		t.Errorf("Erro ao instalar cliente: %v", err)
	}
}

func TestClientConnect(t *testing.T) {
	cfg := &config.Config{}
	cfg.Client.ServerAddress = "localhost:2222"
	cfg.Client.Username = "test"
	cfg.Client.Password = "test"

	client := New(cfg)
	if client == nil {
		t.Fatal("Cliente é nil")
	}

	// Testa a conexão
	err := client.Connect()
	if err != nil {
		t.Errorf("Erro ao conectar cliente: %v", err)
	}

	// Fecha a conexão
	if err := client.Close(); err != nil {
		t.Errorf("Erro ao fechar cliente: %v", err)
	}
}
