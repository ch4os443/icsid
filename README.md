# ICSID

Sistema legítimo de gerenciamento remoto multiplataforma desenvolvido em Go.

## Características

- Comunicação SSH segura
- Interface web responsiva
- Multiplataforma (Windows, Linux, macOS)
- Execução sem privilégios administrativos
- Gerenciamento de sessões
- Monitoramento em tempo real

## Requisitos

- Go 1.21 ou superior
- OpenSSL para certificados
- Navegador web moderno

## Instalação

```bash
# Clone o repositório
git clone git@github.com:ch4os443/icsid.git

# Entre no diretório
cd icsid

# Instale as dependências
./scripts/install.sh

# Execute os testes
./scripts/test.sh

# Compile o projeto
./scripts/build.sh
```

## Uso

```bash
# Inicie o servidor em modo de desenvolvimento
./scripts/dev.sh

# Acesse a interface web
https://localhost:8443
```

## Segurança

- Comunicação criptografada
- Autenticação segura
- Sem privilégios elevados
- Logs detalhados

## Licença

MIT License - veja o arquivo [LICENSE](LICENSE) para mais detalhes. 