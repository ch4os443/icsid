#!/bin/bash

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${GREEN}🚀 Iniciando backup do ICSID...${NC}"

# Verifica se o diretório de backup existe
BACKUP_DIR="backups"
if [ ! -d "$BACKUP_DIR" ]; then
    echo -e "${YELLOW}📁 Criando diretório de backup...${NC}"
    mkdir -p "$BACKUP_DIR"
fi

# Gera timestamp para o backup
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="$BACKUP_DIR/icsid_backup_$TIMESTAMP.tar.gz"

# Lista de arquivos e diretórios para backup
BACKUP_ITEMS=(
    "config.yaml"
    "certs"
    "logs"
    "data"
)

# Verifica se os itens existem antes de fazer backup
echo -e "${YELLOW}📦 Verificando arquivos para backup...${NC}"
for item in "${BACKUP_ITEMS[@]}"; do
    if [ ! -e "$item" ]; then
        echo -e "${YELLOW}⚠️  Item não encontrado: $item${NC}"
    else
        echo -e "${GREEN}✅ Item encontrado: $item${NC}"
    fi
done

# Cria o backup
echo -e "${YELLOW}💾 Criando backup...${NC}"
tar -czf "$BACKUP_FILE" "${BACKUP_ITEMS[@]}"

# Verifica se o backup foi criado com sucesso
if [ -f "$BACKUP_FILE" ]; then
    echo -e "${GREEN}✅ Backup criado com sucesso: $BACKUP_FILE${NC}"
    echo -e "Tamanho do backup: $(du -h "$BACKUP_FILE" | cut -f1)"
else
    echo -e "${RED}❌ Erro ao criar backup${NC}"
    exit 1
fi

# Remove backups antigos (mantém os últimos 5)
echo -e "${YELLOW}🧹 Limpando backups antigos...${NC}"
cd "$BACKUP_DIR" || exit
ls -t | tail -n +6 | xargs -r rm
cd ..

echo -e "${GREEN}✨ Backup concluído com sucesso!${NC}" 