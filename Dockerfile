# Usar uma imagem base com GLIBC 2.27
FROM ubuntu:18.04

# Instalar dependências necessárias
RUN apt-get update && apt-get install -y \
    build-essential \
    wget \
    git \
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
    pkg-config \
    && rm -rf /var/lib/apt/lists/*

# Instalar Go
RUN wget https://go.dev/dl/go1.21.6.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf go1.21.6.linux-amd64.tar.gz \
    && rm go1.21.6.linux-amd64.tar.gz

# Configurar variáveis de ambiente
ENV PATH="/usr/local/go/bin:${PATH}"
ENV GOPATH="/go"
ENV PATH="${GOPATH}/bin:${PATH}"

# Criar diretório de trabalho
WORKDIR /app

# Copiar os arquivos do projeto
COPY . .

# Compilar a aplicação
RUN go mod download && \
    go build -ldflags="-s -w" -o falcon-agent cmd/api/main.go

# Criar uma imagem de runtime mais leve
FROM ubuntu:18.04

# Instalar apenas as dependências necessárias para execução
RUN apt-get update && apt-get install -y \
    libgl1-mesa-dev \
    libxcursor-dev \
    libxrandr-dev \
    libxinerama-dev \
    libxi-dev \
    libxxf86vm-dev \
    libxss-dev \
    libusb-1.0-0 \
    && rm -rf /var/lib/apt/lists/*

# Copiar o binário compilado
COPY --from=0 /app/falcon-agent /usr/local/bin/

# Definir o comando padrão
CMD ["falcon-agent"] 