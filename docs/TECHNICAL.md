# Documentação Técnica - ICSID

## Arquitetura

O ICSID é um sistema de gerenciamento remoto multiplataforma que utiliza SSH para comunicação segura e uma interface web para controle. A arquitetura é composta pelos seguintes componentes:

### 1. Servidor SSH
- Implementado usando a biblioteca `github.com/gliderlabs/ssh`
- Gerencia conexões SSH seguras
- Autenticação por usuário e senha
- Geração automática de chaves RSA
- Suporte a múltiplas sessões simultâneas

### 2. Servidor Web
- Implementado usando `net/http` e `github.com/gorilla/mux`
- Interface web responsiva usando Bootstrap 5
- APIs RESTful para gerenciamento de sessões
- Comunicação em tempo real via WebSocket
- Certificados SSL auto-assinados para HTTPS

### 3. Gerenciador de Sessões
- Gerenciamento thread-safe de sessões SSH
- Armazenamento de histórico de comandos
- Coleta de informações do sistema
- Execução remota de comandos

## Segurança

### Comunicação
- Criptografia SSH para comunicação remota
- HTTPS para interface web
- Certificados SSL auto-assinados
- Autenticação por usuário e senha

### Persistência
- Armazenamento seguro de chaves SSH
- Histórico de comandos em memória
- Sem armazenamento de dados sensíveis

### Permissões
- Execução sem privilégios de administrador
- Permissões mínimas necessárias
- Isolamento de processos

## Desenvolvimento

### Requisitos
- Go 1.21 ou superior
- Node.js 18+ (para interface web)
- OpenSSL (para certificados)

### Estrutura do Projeto
```
.
├── cmd/                    # Ponto de entrada
├── internal/              # Código interno
│   ├── config/           # Configurações
│   ├── server/           # Servidores SSH e Web
│   ├── client/           # Cliente
│   └── models/           # Modelos
├── pkg/                   # Pacotes reutilizáveis
├── web/                   # Interface web
├── scripts/              # Scripts de automação
├── tests/                # Testes
└── docs/                 # Documentação
```

### Scripts
- `scripts/generate_cert.sh`: Gera certificados SSL
- `scripts/build.sh`: Compila o projeto
- `scripts/test.sh`: Executa testes

### CI/CD
- GitHub Actions para automação
- Testes automatizados
- Build multiplataforma
- Backup automático

## Uso

### Compilação
```bash
./scripts/build.sh
```

### Execução
```bash
./icsid.exe -config config.yaml
```

### Interface Web
- Acesse `https://localhost:8443`
- Credenciais padrão: admin/changeme
- Altere as credenciais em produção

### Conexão SSH
```bash
ssh -p 2222 admin@localhost
```

## Manutenção

### Logs
- Logs de sistema no diretório `logs/`
- Rotação automática de logs
- Níveis de log configuráveis

### Backup
- Backup automático a cada commit
- Armazenamento de binários
- Histórico de versões

### Monitoramento
- Status de sessões
- Uso de recursos
- Histórico de comandos

## Contribuição

1. Fork o repositório
2. Crie uma branch para sua feature
3. Faça commit das mudanças
4. Push para a branch
5. Crie um Pull Request

## Licença

Este projeto está sob a licença MIT. Veja o arquivo LICENSE para mais detalhes. 