#!/bin/bash

# Cria o diretório para os certificados
mkdir -p certs

# Gera a chave privada
openssl genrsa -out certs/key.pem 2048

# Gera o certificado auto-assinado
openssl req -x509 -new -nodes \
    -key certs/key.pem \
    -days 365 \
    -out certs/cert.pem \
    -subj "/C=BR/ST=SP/L=SP/O=ICSID/CN=localhost"

# Ajusta as permissões
chmod 600 certs/key.pem
chmod 644 certs/cert.pem

echo "Certificados gerados com sucesso!" 