package main

import (
	"github.com/dev/falcon-agent/internal/config"
	"github.com/dev/falcon-agent/internal/service"
	"github.com/dev/falcon-agent/pkg/logger"
	"github.com/dev/falcon-agent/pkg/ui"
)

func main() {
	// Inicializa a configuração
	cfg := config.New()

	// Inicializa o logger
	log, err := logger.New(cfg.LogPath)
	if err != nil {
		panic(err)
	}

	log.Info("Iniciando Falcon Agent...")

	// Coleta informações da máquina
	machineInfo, err := service.CollectMachineInfo()
	if err != nil {
		log.Error("Erro ao coletar informações da máquina: %v", err)
		panic(err)
	}

	// Inicia a interface gráfica
	app := ui.New(machineInfo)
	app.Run()
}
