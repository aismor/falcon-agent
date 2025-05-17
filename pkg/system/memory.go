package system

import (
	"fmt"
	"log"

	"github.com/jaypipes/ghw"
)

// MemoryInfo representa as informações de memória do sistema
type MemoryInfo struct {
	Size         string `json:"size"`
	Manufacturer string `json:"manufacturer"`
	SerialNumber string `json:"serial_number"`
}

// GetMemoryInfo retorna as informações de memória do sistema
func GetMemoryInfo() (*MemoryInfo, error) {
	memory, err := ghw.Memory()
	if err != nil {
		return nil, fmt.Errorf("erro ao obter informações de memória: %v", err)
	}

	info := &MemoryInfo{}

	// Log do número total de módulos encontrados
	log.Printf("Número de módulos de memória encontrados: %d", len(memory.Modules))

	// Se houver módulos de memória, pegar informações do primeiro módulo
	if len(memory.Modules) > 0 {
		module := memory.Modules[0]

		// Log das informações brutas do módulo
		log.Printf("Informações do módulo: %+v", module)

		// Converter tamanho para GB
		sizeGB := float64(module.SizeBytes) / (1024 * 1024 * 1024)
		info.Size = fmt.Sprintf("%.0f GB", sizeGB)

		info.Manufacturer = module.Vendor
		info.SerialNumber = module.SerialNumber
	} else {
		log.Printf("Nenhum módulo de memória encontrado")
	}

	return info, nil
}
