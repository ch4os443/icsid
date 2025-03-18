#!/bin/bash

# Verifica se o Go está instalado
if ! command -v go &> /dev/null; then
    echo "Go não está instalado. Por favor, instale o Go 1.21 ou superior."
    exit 1
fi

# Verifica a versão do Go
go_version=$(go version | cut -d' ' -f3 | cut -d'.' -f2)
if [ "$go_version" -lt 21 ]; then
    echo "Versão do Go muito antiga. Por favor, instale o Go 1.21 ou superior."
    exit 1
fi

# Gera os certificados SSL se não existirem
if [ ! -f "certs/cert.pem" ] || [ ! -f "certs/key.pem" ]; then
    echo "Gerando certificados SSL..."
    ./scripts/generate_cert.sh
fi

# Instala as dependências
echo "Instalando dependências..."
go mod download

# Compila o projeto
echo "Compilando o projeto..."
GOOS=windows GOARCH=amd64 go build -o icsid.exe ./cmd/icsid
GOOS=linux GOARCH=amd64 go build -o icsid_linux ./cmd/icsid
GOOS=darwin GOARCH=amd64 go build -o icsid_darwin ./cmd/icsid

echo "Compilação concluída!" 