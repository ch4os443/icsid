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

# Executa os testes
echo "Executando testes..."
go test -v ./...

# Verifica a cobertura de código
echo "Verificando cobertura de código..."
go test -coverprofile=coverage.out ./...
coverage=$(go tool cover -func=coverage.out | grep total: | awk '{print $3}')
echo "Cobertura total: $coverage"

# Remove o arquivo de cobertura
rm coverage.out

echo "Testes concluídos!" 