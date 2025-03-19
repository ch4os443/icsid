#!/bin/bash

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${GREEN}🚀 Iniciando instalação do ICSID...${NC}"

# Verifica se o Go está instalado
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Go não está instalado. Por favor, instale o Go 1.21 ou superior.${NC}"
    exit 1
fi

# Verifica se o OpenSSL está instalado
if ! command -v openssl &> /dev/null; then
    echo -e "${RED}❌ OpenSSL não está instalado. Por favor, instale o OpenSSL.${NC}"
    exit 1
fi

# Verifica se o Git está instalado
if ! command -v git &> /dev/null; then
    echo -e "${RED}❌ Git não está instalado. Por favor, instale o Git.${NC}"
    exit 1
fi

# Cria diretórios necessários
echo -e "${YELLOW}📁 Criando estrutura de diretórios...${NC}"
mkdir -p cmd/icsid
mkdir -p internal/{config,server,client}
mkdir -p web/static
mkdir -p scripts
mkdir -p certs
mkdir -p logs

# Inicializa o módulo Go
echo -e "${YELLOW}📦 Inicializando módulo Go...${NC}"
go mod init github.com/ch4os443/icsid

# Instala dependências
echo -e "${YELLOW}📥 Instalando dependências...${NC}"
go get github.com/gliderlabs/ssh@latest
go get github.com/gorilla/mux@latest
go get gopkg.in/yaml.v3@latest

# Gera certificados SSL
echo -e "${YELLOW}🔒 Gerando certificados SSL...${NC}"
openssl genrsa -out certs/key.pem 2048
openssl req -x509 -new -nodes \
    -key certs/key.pem \
    -days 365 \
    -out certs/cert.pem \
    -subj "/C=BR/ST=SP/L=SP/O=ICSID/CN=localhost"

# Compila o projeto
echo -e "${YELLOW}🔨 Compilando projeto...${NC}"
go build -o icsid ./cmd/icsid

# Configura permissões
echo -e "${YELLOW}🔑 Configurando permissões...${NC}"
chmod +x icsid
chmod +x scripts/*.sh

echo -e "${GREEN}✅ Instalação concluída com sucesso!${NC}"
echo -e "\nPróximos passos:"
echo -e "1. Configure o arquivo config.yaml"
echo -e "2. Execute o servidor: ./icsid"
echo -e "3. Acesse a interface web: https://localhost:8443"