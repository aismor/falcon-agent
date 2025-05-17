package main

import (
	"fmt"
	"time"

	"github.com/dev/falcon-agent/internal/config"
	"github.com/dev/falcon-agent/internal/service"
	"github.com/dev/falcon-agent/pkg/logger"
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

	// Loop principal
	for {
		// Coleta informações da máquina
		machineInfo, err := service.CollectMachineInfo()
		if err != nil {
			log.Error("Erro ao coletar informações da máquina: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		// Limpa a tela
		fmt.Print("\033[H\033[2J")

		// Exibe as informações
		fmt.Printf("=== Falcon Agent - Informações do Sistema ===\n\n")
		fmt.Printf("Sistema Operacional: %s\n", machineInfo.OS)
		fmt.Printf("Hostname: %s\n\n", machineInfo.Hostname)

		fmt.Printf("Processador:\n")
		fmt.Printf("  Modelo: %s\n", machineInfo.Processor.Model)
		fmt.Printf("  Núcleos: %d\n", machineInfo.Processor.Cores)
		fmt.Printf("  Threads: %d\n", machineInfo.Processor.Threads)
		fmt.Printf("  Frequência: %.2f GHz\n\n", machineInfo.Processor.FrequencyGHz)

		fmt.Printf("BIOS:\n")
		fmt.Printf("  Fabricante: %s\n", machineInfo.BIOS.Vendor)
		fmt.Printf("  Versão: %s\n", machineInfo.BIOS.Version)
		fmt.Printf("  Data: %s\n\n", machineInfo.BIOS.ReleaseDate)

		fmt.Printf("Memória:\n")
		for _, mem := range machineInfo.Memory {
			fmt.Printf("  Slot: %s\n", mem.Slot)
			fmt.Printf("  Capacidade: %d MB\n", mem.SizeMB)
			fmt.Printf("  Fabricante: %s\n", mem.Manufacturer)
			fmt.Printf("  Serial: %s\n\n", mem.SerialNumber)
		}

		fmt.Printf("Discos Rígidos:\n")
		for _, hd := range machineInfo.HDs {
			fmt.Printf("  Modelo: %s\n", hd.Model)
			fmt.Printf("  Serial: %s\n", hd.Serial)
			fmt.Printf("  Tamanho: %d GB\n\n", hd.SizeGB)
		}

		fmt.Printf("Placa-mãe Serial: %s\n", machineInfo.MotherboardSN)
		fmt.Printf("Serial Number: %s\n", machineInfo.SerialNumber)

		fmt.Printf("\nAtualizado em: %s\n", time.Now().Format("02/01/2006 15:04:05"))
		fmt.Printf("Pressione Ctrl+C para sair...\n")

		// Aguarda 5 segundos antes da próxima atualização
		time.Sleep(5 * time.Second)
	}
}
