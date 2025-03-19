#!/bin/bash

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${GREEN}🚀 Iniciando ambiente de desenvolvimento do ICSID...${NC}"

# Verifica se o Go está instalado
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Go não está instalado. Por favor, instale o Go 1.21 ou superior.${NC}"
    exit 1
fi

# Verifica se o Air está instalado
if ! command -v air &> /dev/null; then
    echo -e "${YELLOW}📥 Instalando Air para hot-reload...${NC}"
    go install github.com/cosmtrek/air@latest
fi

# Verifica se os certificados existem
if [ ! -f "certs/cert.pem" ] || [ ! -f "certs/key.pem" ]; then
    echo -e "${YELLOW}🔒 Gerando certificados SSL...${NC}"
    mkdir -p certs
    openssl genrsa -out certs/key.pem 2048
    openssl req -x509 -new -nodes \
        -key certs/key.pem \
        -days 365 \
        -out certs/cert.pem \
        -subj "/C=BR/ST=SP/L=SP/O=ICSID/CN=localhost"
fi

# Verifica se o arquivo de configuração existe
if [ ! -f "config.yaml" ]; then
    echo -e "${YELLOW}📝 Criando arquivo de configuração...${NC}"
    cat > config.yaml << EOL
server:
  ssh:
    port: 2222
    host_key: "certs/host_key"
    username: "admin"
    password: "changeme"  # Altere em produção!
  web:
    port: 8443
    cert_file: "certs/cert.pem"
    key_file: "certs/key.pem"
client:
  server_address: "localhost:2222"
  username: "admin"
  password: "changeme"  # Altere em produção!
persistence:
  enabled: true
  method: "registry"  # Windows: registry, Linux/macOS: service
  path: "/usr/local/bin/icsid"  # Caminho para o binário
EOL
fi

# Verifica se o arquivo .air.toml existe
if [ ! -f ".air.toml" ]; then
    echo -e "${YELLOW}📝 Criando configuração do Air...${NC}"
    cat > .air.toml << EOL
root = "."
tmp_dir = "tmp"

[build]
  cmd = "go build -o ./tmp/main ./cmd/icsid"
  bin = "./tmp/main"
  delay = 1000
  exclude_dir = ["assets", "tmp", "vendor"]
  include_ext = ["go", "tpl", "tmpl", "html"]
  exclude_regex = ["_test.go"]
EOL
fi

# Inicia o servidor com hot-reload
echo -e "${GREEN}🚀 Iniciando servidor com hot-reload...${NC}"
air 