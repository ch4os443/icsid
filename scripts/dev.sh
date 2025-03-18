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

# Instala as dependências
echo "Instalando dependências..."
go mod download

# Gera os certificados SSL se não existirem
if [ ! -f "certs/server.crt" ] || [ ! -f "certs/server.key" ]; then
    echo "Gerando certificados SSL..."
    ./scripts/generate_cert.sh
fi

# Executa o programa em modo de desenvolvimento
echo "Iniciando o programa em modo de desenvolvimento..."
go run cmd/icsid/main.go --debug 