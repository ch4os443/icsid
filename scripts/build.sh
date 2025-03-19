#!/bin/bash

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${GREEN}🚀 Iniciando compilação do ICSID...${NC}"

# Verifica se o Go está instalado
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Go não está instalado. Por favor, instale o Go 1.21 ou superior.${NC}"
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

# Cria diretório de build
echo -e "${YELLOW}📁 Criando diretório de build...${NC}"
mkdir -p build

# Compila para diferentes plataformas
echo -e "${YELLOW}🔨 Compilando para diferentes plataformas...${NC}"

# Windows
echo -e "${YELLOW}Windows...${NC}"
GOOS=windows GOARCH=amd64 go build -o build/icsid.exe ./cmd/icsid

# Linux
echo -e "${YELLOW}Linux...${NC}"
GOOS=linux GOARCH=amd64 go build -o build/icsid_linux ./cmd/icsid
GOOS=linux GOARCH=386 go build -o build/icsid_linux_386 ./cmd/icsid

# macOS
echo -e "${YELLOW}macOS...${NC}"
GOOS=darwin GOARCH=amd64 go build -o build/icsid_mac ./cmd/icsid
GOOS=darwin GOARCH=arm64 go build -o build/icsid_mac_arm64 ./cmd/icsid

# Cria arquivos de distribuição
echo -e "${YELLOW}📦 Criando arquivos de distribuição...${NC}"

# Windows
zip -j build/icsid_windows.zip build/icsid.exe config.yaml README.md LICENSE

# Linux
tar -czf build/icsid_linux.tar.gz build/icsid_linux config.yaml README.md LICENSE
tar -czf build/icsid_linux_386.tar.gz build/icsid_linux_386 config.yaml README.md LICENSE

# macOS
tar -czf build/icsid_mac.tar.gz build/icsid_mac config.yaml README.md LICENSE
tar -czf build/icsid_mac_arm64.tar.gz build/icsid_mac_arm64 config.yaml README.md LICENSE

# Limpa arquivos temporários
echo -e "${YELLOW}🧹 Limpando arquivos temporários...${NC}"
rm build/icsid.exe
rm build/icsid_linux
rm build/icsid_linux_386
rm build/icsid_mac
rm build/icsid_mac_arm64

echo -e "${GREEN}✅ Compilação concluída com sucesso!${NC}"
echo -e "\nArquivos gerados em build/:"
ls -lh build/ 