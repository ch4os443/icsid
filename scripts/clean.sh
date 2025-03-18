#!/bin/bash

echo "Limpando arquivos temporários e de build..."

# Remove arquivos de build
rm -rf build/
rm -rf dist/

# Remove arquivos temporários
rm -rf tmp/
rm -f *.out
rm -f *.test
rm -f coverage.out

# Remove arquivos de cache do Go
go clean -cache -modcache

echo "Limpeza concluída!" 