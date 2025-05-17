package service

import (
	"runtime"
	"time"

	"github.com/dev/falcon-agent/internal/config"
	"github.com/dev/falcon-agent/pkg/logger"
)

// Agent representa o serviço principal do agente
type Agent struct {
	config *config.Config
	logger logger.Logger
}

// New cria uma nova instância do agente
func New(cfg *config.Config, log logger.Logger) *Agent {
	return &Agent{
		config: cfg,
		logger: log,
	}
}

// Start inicia o agente
func (a *Agent) Start() error {
	a.logger.Info("Iniciando Falcon Agent na plataforma: %s", a.config.Platform)

	// Inicia o loop principal do agente
	go a.mainLoop()

	return nil
}

// Stop para o agente
func (a *Agent) Stop() error {
	a.logger.Info("Parando Falcon Agent...")
	return nil
}

// mainLoop é o loop principal do agente
func (a *Agent) mainLoop() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			a.collectMetrics()
		}
	}
}

// collectMetrics coleta métricas do sistema
func (a *Agent) collectMetrics() {
	// Exemplo de coleta de métricas básicas
	a.logger.Info("Coletando métricas do sistema...")
	a.logger.Info("Sistema Operacional: %s", runtime.GOOS)
	a.logger.Info("Arquitetura: %s", runtime.GOARCH)
	a.logger.Info("Número de CPUs: %d", runtime.NumCPU())
}
