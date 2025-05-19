package metrics

import (
	"sync"
	"time"
)

const maxDataPoints = 60 // 5 minutos de histórico (1 ponto a cada 5 segundos)

type MetricPoint struct {
	Timestamp time.Time
	Value     float64
}

type MetricHistory struct {
	mu     sync.RWMutex
	points []MetricPoint
}

func NewMetricHistory() *MetricHistory {
	return &MetricHistory{
		points: make([]MetricPoint, 0, maxDataPoints),
	}
}

func (h *MetricHistory) Add(value float64) {
	h.mu.Lock()
	defer h.mu.Unlock()

	point := MetricPoint{
		Timestamp: time.Now(),
		Value:     value,
	}

	if len(h.points) >= maxDataPoints {
		h.points = h.points[1:] // Remove o ponto mais antigo
	}
	h.points = append(h.points, point)
}

func (h *MetricHistory) GetPoints() []MetricPoint {
	h.mu.RLock()
	defer h.mu.RUnlock()

	points := make([]MetricPoint, len(h.points))
	copy(points, h.points)
	return points
}

type SystemMetrics struct {
	CPUUsage    *MetricHistory
	MemoryUsage *MetricHistory
}

func NewSystemMetrics() *SystemMetrics {
	return &SystemMetrics{
		CPUUsage:    NewMetricHistory(),
		MemoryUsage: NewMetricHistory(),
	}
}
