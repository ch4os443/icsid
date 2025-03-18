# ICSID - Sistema de Gerenciamento Remoto

Sistema legítimo de gerenciamento remoto multiplataforma desenvolvido em Go.

## Estrutura do Projeto

```
.
├── cmd/                    # Ponto de entrada da aplicação
├── internal/              # Código interno do projeto
│   ├── config/           # Configurações
│   ├── server/           # Servidor SSH e Web
│   ├── client/           # Cliente
│   └── models/           # Modelos de dados
├── pkg/                   # Pacotes reutilizáveis
├── web/                   # Interface web
├── scripts/              # Scripts de automação
├── tests/                # Testes
└── docs/                 # Documentação
```

## Requisitos

- Go 1.21 ou superior
- Node.js 18+ (para interface web)
- Docker (opcional, para desenvolvimento)

## Desenvolvimento

1. Clone o repositório
2. Instale as dependências:
   ```bash
   go mod download
   ```
3. Execute os testes:
   ```bash
   go test ./...
   ```
4. Compile o projeto:
   ```bash
   go build -o icsid.exe ./cmd/icsid
   ```

## CI/CD

O projeto utiliza GitHub Actions para:
- Testes automatizados
- Build e deploy
- Verificação de segurança
- Backup automático

## Backup

O sistema de backup automático é executado a cada commit, armazenando:
- Código fonte
- Configurações
- Documentação
- Artefatos de build

## Segurança

Este é um software legítimo que segue as melhores práticas de segurança:
- Comunicação criptografada via SSH
- Autenticação segura
- Sem privilégios de administrador necessários
- Conformidade com políticas de segurança do Windows 