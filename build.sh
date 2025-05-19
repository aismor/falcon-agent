#!/bin/bash

# Construir a imagem Docker
docker build -t falcon-agent-builder .

# Criar um container temporário e copiar o binário
docker create --name temp-container falcon-agent-builder
docker cp temp-container:/usr/local/bin/falcon-agent ./falcon-agent
docker rm temp-container

echo "Build concluído! O binário está disponível em ./falcon-agent" 