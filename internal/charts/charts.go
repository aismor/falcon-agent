package charts

import (
	"image"
	"image/color"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

type TimeValue struct {
	Time  time.Time
	Value float64
}

func CreateLineChart(data []TimeValue, title, yLabel string) (image.Image, error) {
	if len(data) == 0 {
		return nil, nil
	}

	p := plot.New()

	p.Title.Text = title
	p.X.Label.Text = "Tempo"
	p.Y.Label.Text = yLabel

	// Configurar o estilo do gráfico
	p.BackgroundColor = color.White
	p.Title.TextStyle.Color = color.Black
	p.X.Label.TextStyle.Color = color.Black
	p.Y.Label.TextStyle.Color = color.Black
	p.X.Color = color.Black
	p.Y.Color = color.Black

	// Criar os pontos para o gráfico
	pts := make(plotter.XYs, len(data))
	minTime := data[0].Time
	for i, d := range data {
		pts[i].X = d.Time.Sub(minTime).Seconds()
		pts[i].Y = d.Value
	}

	// Criar a linha
	line, err := plotter.NewLine(pts)
	if err != nil {
		return nil, err
	}
	line.Color = color.RGBA{R: 0, G: 0, B: 255, A: 255}
	line.Width = vg.Points(2)

	// Adicionar a linha ao gráfico
	p.Add(line)

	// Adicionar grade
	p.Add(plotter.NewGrid())

	// Criar a imagem
	width, height := vg.Points(400), vg.Points(200)
	c := vgimg.New(width, height)

	// Desenhar o gráfico
	p.Draw(draw.New(c))

	// Retornar a imagem final
	return c.Image(), nil
}

func CreateCPUChart(data []TimeValue) (image.Image, error) {
	return CreateLineChart(data, "Uso de CPU", "Porcentagem (%)")
}

func CreateMemoryChart(data []TimeValue) (image.Image, error) {
	return CreateLineChart(data, "Uso de Memória", "Porcentagem (%)")
}
