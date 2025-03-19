// main.go
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func execCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func main() {
	// 1. Criar diretório base
	projectDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("❌ Erro ao obter diretório atual: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n🚀 Iniciando instalação do ICSID v1.0.0\n")
	fmt.Printf("📂 Diretório do projeto: %s\n\n", projectDir)

	// 2. Verificar requisitos
	fmt.Println("📌 Verificando requisitos...")

	if err := execCommand("go", "version"); err != nil {
		fmt.Printf("❌ Go não está instalado: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ Go instalado")

	if err := execCommand("git", "--version"); err != nil {
		fmt.Printf("❌ Git não está instalado: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ Git instalado")

	// 3. Criar estrutura de diretórios
	fmt.Println("\n📌 Criando estrutura de diretórios...")
	dirs := []string{
		"cmd/icsid",
		"internal/config",
		"internal/server",
		"internal/client",
		"web/static",
		"scripts",
		"certs",
		"logs",
	}

	for _, dir := range dirs {
		path := filepath.Join(projectDir, dir)
		if err := os.MkdirAll(path, 0755); err != nil {
			fmt.Printf("❌ Erro ao criar diretório %s: %v\n", dir, err)
			os.Exit(1)
		}
		fmt.Printf("✅ Criado diretório: %s\n", dir)
	}

	// 4. Inicializar módulo Go
	fmt.Println("\n📌 Inicializando módulo Go...")

	// Verifica se o go.mod já existe
	goModPath := filepath.Join(projectDir, "go.mod")
	if _, err := os.Stat(goModPath); err == nil {
		// Se o arquivo existe, tenta atualizar o módulo
		fmt.Println("ℹ️  Módulo Go já existe, atualizando...")
		if err := execCommand("go", "mod", "tidy"); err != nil {
			fmt.Printf("❌ Erro ao atualizar módulo Go: %v\n", err)
			os.Exit(1)
		}
	} else {
		// Se não existe, cria um novo
		if err := execCommand("go", "mod", "init", "github.com/ch4os443/icsid"); err != nil {
			fmt.Printf("❌ Erro ao inicializar módulo Go: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Println("✅ Módulo Go configurado")

	// 5. Instalar dependências
	fmt.Println("\n📌 Instalando dependências...")
	deps := []string{
		"github.com/gliderlabs/ssh@latest",
		"github.com/gorilla/mux@latest",
		"gopkg.in/yaml.v3@latest",
	}

	for _, dep := range deps {
		if err := execCommand("go", "get", dep); err != nil {
			fmt.Printf("❌ Erro ao instalar %s: %v\n", dep, err)
			os.Exit(1)
		}
		fmt.Printf("✅ Instalado: %s\n", dep)
	}

	// 6. Gerar certificados SSL
	fmt.Println("\n📌 Gerando certificados SSL...")
	certsDir := filepath.Join(projectDir, "certs")
	if err := execCommand("openssl", "genrsa", "-out", filepath.Join(certsDir, "key.pem"), "2048"); err != nil {
		fmt.Printf("❌ Erro ao gerar chave privada: %v\n", err)
		os.Exit(1)
	}

	if err := execCommand("openssl", "req", "-x509", "-new", "-nodes",
		"-key", filepath.Join(certsDir, "key.pem"),
		"-days", "365",
		"-out", filepath.Join(certsDir, "cert.pem"),
		"-subj", "/C=BR/ST=SP/L=SP/O=ICSID/CN=localhost"); err != nil {
		fmt.Printf("❌ Erro ao gerar certificado: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ Certificados SSL gerados")

	// 7. Compilar projeto
	fmt.Println("\n📌 Compilando projeto...")
	if runtime.GOOS == "windows" {
		if err := execCommand("go", "build", "-o", "icsid.exe", "./cmd/icsid"); err != nil {
			fmt.Printf("❌ Erro ao compilar projeto: %v\n", err)
			os.Exit(1)
		}
	} else {
		if err := execCommand("go", "build", "-o", "icsid", "./cmd/icsid"); err != nil {
			fmt.Printf("❌ Erro ao compilar projeto: %v\n", err)
			os.Exit(1)
		}
	}
	fmt.Println("✅ Projeto compilado")

	fmt.Println("\n✨ Instalação concluída com sucesso!")
	fmt.Println("\nPróximos passos:")
	fmt.Println("1. Configure o arquivo config.yaml")
	fmt.Println("2. Execute o servidor:")
	if runtime.GOOS == "windows" {
		fmt.Println("   .\\icsid.exe")
	} else {
		fmt.Println("   ./icsid")
	}
	fmt.Println("3. Acesse a interface web: https://localhost:8443")
}
