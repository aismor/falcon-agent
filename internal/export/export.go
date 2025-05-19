package export

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/dev/falcon-agent/internal/metrics"
)

type Exporter interface {
	Export(data interface{}, filename string) error
}

type CSVExporter struct{}

func (e *CSVExporter) Export(data interface{}, filename string) error {
	metrics, ok := data.(*metrics.SystemMetrics)
	if !ok {
		return fmt.Errorf("dados inválidos para exportação CSV")
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Escreve o cabeçalho
	header := []string{"Timestamp", "CPU Usage (%)", "Memory Usage (%)"}
	if err := writer.Write(header); err != nil {
		return err
	}

	cpuPoints := metrics.CPUUsage.GetPoints()
	memPoints := metrics.MemoryUsage.GetPoints()

	// Escreve os dados
	for i := 0; i < len(cpuPoints); i++ {
		row := []string{
			cpuPoints[i].Timestamp.Format(time.RFC3339),
			fmt.Sprintf("%.2f", cpuPoints[i].Value),
			fmt.Sprintf("%.2f", memPoints[i].Value),
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

type JSONExporter struct{}

func (e *JSONExporter) Export(data interface{}, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

func ExportData(data interface{}, format string, baseDir string) error {
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return err
	}

	timestamp := time.Now().Format("2006-01-02_15-04-05")
	var exporter Exporter
	var filename string

	switch format {
	case "csv":
		exporter = &CSVExporter{}
		filename = filepath.Join(baseDir, fmt.Sprintf("metrics_%s.csv", timestamp))
	case "json":
		exporter = &JSONExporter{}
		filename = filepath.Join(baseDir, fmt.Sprintf("metrics_%s.json", timestamp))
	default:
		return fmt.Errorf("formato de exportação não suportado: %s", format)
	}

	return exporter.Export(data, filename)
}
