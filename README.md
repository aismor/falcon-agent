# Falcon Agent

O Falcon Agent é uma aplicação desktop desenvolvida em Go que monitora e gerencia dispositivos USB, fornecendo uma interface gráfica moderna e intuitiva.

## Características

- Interface gráfica moderna usando Fyne
- Monitoramento de dispositivos USB em tempo real
- Suporte a múltiplos dispositivos
- Minimização para bandeja do sistema
- Design responsivo e compacto

## Requisitos do Sistema

- Sistema operacional Linux
- GLIBC 2.27 ou superior
- Dependências do sistema:
  - libgl1-mesa-dev
  - libxcursor-dev
  - libxrandr-dev
  - libxinerama-dev
  - libxi-dev
  - libxxf86vm-dev
  - libxss-dev
  - libusb-1.0-0

## Compilação

### Método 1: Compilação Local

1. Instale o Go 1.21 ou superior
2. Instale as dependências do sistema:
```bash
sudo apt-get update && sudo apt-get install -y \
    build-essential \
    libgl1-mesa-dev \
    xorg-dev \
    libxcursor-dev \
    libxrandr-dev \
    libxinerama-dev \
    libxi-dev \
    libxxf86vm-dev \
    libxss-dev \
    libglib2.0-dev \
    libusb-1.0-0-dev \
    pkg-config
```
3. Clone o repositório:
```bash
git clone https://github.com/seu-usuario/falcon-agent.git
cd falcon-agent
```
4. Compile o projeto:
```bash
go mod download
go build -o falcon-agent cmd/api/main.go
```

### Método 2: Compilação com Docker (Recomendado)

1. Certifique-se de ter o Docker instalado
2. Clone o repositório:
```bash
git clone https://github.com/seu-usuario/falcon-agent.git
cd falcon-agent
```
3. Execute o script de build:
```bash
chmod +x build.sh
./build.sh
```

O binário compilado será gerado na raiz do projeto como `falcon-agent`.

## Execução

Para executar a aplicação:

```bash
./falcon-agent
```

A aplicação será iniciada e aparecerá na bandeja do sistema. Você pode:
- Clicar no ícone para mostrar/esconder a janela principal
- Usar o menu de contexto para sair da aplicação

## Desenvolvimento

### Estrutura do Projeto

```
falcon-agent/
├── cmd/
│   └── api/
│       └── main.go      # Ponto de entrada da aplicação
├── pkg/
│   ├── ui/
│   │   └── app.go       # Lógica da interface gráfica
│   └── usb/
│       └── monitor.go   # Monitoramento de dispositivos USB
├── Dockerfile           # Configuração para build com Docker
├── build.sh            # Script de build
└── README.md
```

### Contribuindo

1. Faça um fork do projeto
2. Crie uma branch para sua feature (`git checkout -b feature/nova-feature`)
3. Commit suas mudanças (`git commit -am 'Adiciona nova feature'`)
4. Push para a branch (`git push origin feature/nova-feature`)
5. Crie um Pull Request

## Licença

Este projeto está licenciado sob a licença MIT - veja o arquivo [LICENSE](LICENSE) para detalhes.