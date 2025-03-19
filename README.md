# ICSID - Sistema de Controle Remoto

ICSID é um sistema de controle remoto desenvolvido em Go que oferece uma interface web moderna e segura para gerenciar conexões SSH.

## Características

- Interface web responsiva com Tailwind CSS
- Autenticação segura
- Monitoramento em tempo real
- Execução de comandos remotos
- Informações detalhadas do sistema
- Rate limiting para prevenir ataques
- Logging completo
- Suporte multi-plataforma

## Requisitos

- Go 1.21 ou superior
- OpenSSL
- Git

## Instalação

1. Clone o repositório:
```bash
git clone https://github.com/ch4os443/icsid.git
cd icsid
```

2. Execute o script de instalação:
```bash
./scripts/install.sh
```

3. Configure o arquivo `config.yaml`:
```yaml
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
```

## Uso

1. Inicie o servidor:
```bash
./icsid
```

2. Acesse a interface web:
```
https://localhost:8443
```

3. Use as credenciais configuradas para fazer login

## Funcionalidades

### Interface Web

- Lista de sessões ativas
- Informações detalhadas do sistema
- Execução de comandos remotos
- Monitoramento em tempo real
- Estatísticas de uso

### SSH

- Conexão segura via SSH
- Autenticação por usuário/senha
- Rate limiting para prevenir ataques
- Logging completo de comandos

## Segurança

- Comunicação criptografada (SSH + HTTPS)
- Autenticação em duas camadas
- Rate limiting
- Logging de eventos de segurança
- Senhas nunca armazenadas em texto puro

## Desenvolvimento

1. Instale as dependências:
```bash
go mod download
```

2. Execute os testes:
```bash
go test ./...
```

3. Compile para diferentes plataformas:
```bash
./scripts/build.sh
```

## Contribuição

1. Faça um fork do projeto
2. Crie uma branch para sua feature
3. Faça commit das mudanças
4. Push para a branch
5. Abra um Pull Request

## Licença

Este projeto está licenciado sob a licença MIT - veja o arquivo [LICENSE](LICENSE) para detalhes. 