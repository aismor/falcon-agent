# Falcon Agent

Um agente em Go para monitoramento e gerenciamento de sistemas, compatível com Windows e Linux.

## Requisitos

- Go 1.21 ou superior
- Windows 10/11 ou Linux (Ubuntu 20.04+, CentOS 7+)

## Estrutura do Projeto

```
falcon-agent/
├── cmd/
│   └── falcon-agent/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── handler/
│   ├── service/
│   │   └── agent.go
├── pkg/
│   ├── logger/
│   │   └── logger.go
│   └── utils/
├── go.mod
└── README.md
```

## Instalação

1. Clone o repositório:
```bash
git clone https://github.com/seu-usuario/falcon-agent.git
cd falcon-agent
```

2. Instale as dependências:
```bash
go mod download
```

## Executando

### Linux
```bash
# Executar diretamente
go run main.go

# Ou compilar e executar
go build
./falcon-agent
```

### Windows
```bash
# Executar diretamente
go run main.go

# Ou compilar e executar
go build
.\falcon-agent.exe
```

## Compilando para Outras Plataformas

### Para Windows (a partir do Linux)
```bash
GOOS=windows GOARCH=amd64 go build -o falcon-agent.exe
```

### Para Linux (a partir do Windows)
```bash
GOOS=linux GOARCH=amd64 go build -o falcon-agent
```

## Diretórios de Configuração

### Windows
- Logs: `C:\ProgramData\FalconAgent\logs`
- Dados: `C:\ProgramData\FalconAgent\data`
- Configuração: `C:\ProgramData\FalconAgent\config`

### Linux
- Logs: `/var/log/falcon-agent`
- Dados: `/var/lib/falcon-agent`
- Configuração: `/etc/falcon-agent`

## Funcionalidades

- Coleta de métricas do sistema
- Logging multiplataforma
- Gerenciamento de configuração específica por plataforma
- Tratamento de sinais do sistema operacional
- Graceful shutdown

## Contribuindo

1. Fork o projeto
2. Crie sua branch de feature (`git checkout -b feature/nova-feature`)
3. Commit suas mudanças (`git commit -am 'Adiciona nova feature'`)
4. Push para a branch (`git push origin feature/nova-feature`)
5. Crie um novo Pull Request