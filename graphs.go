package main

import (
	"image/color"
	"math/rand"

	"golang.org/x/image/colornames"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

// GetColor get color from palette by index
func GetColor(i int) color.Color {
	i += 25

	if i >= len(colornames.Names) {
		i = rand.Intn(len(colornames.Names))
	}
	return colornames.Map[colornames.Names[i]]
}

// MakeGraph make plot graph
func MakeGraph(ds *DataSet) (p *plot.Plot, err error) {
	p, err = plot.New()

	if err != nil {
		return p, err
	}

	p.Y.Min = 0
	p.Y.Max = 100
	plotter.DefaultLineStyle.Width = vg.Points(1)
	plotter.DefaultGlyphStyle.Radius = vg.Points(2)
	p.X.Label.Text = "Time"

	p.X.Tick.Marker = plot.TimeTicks{Format: "2006-01-02 15:04:05"}

	for i, mv := range ds.Metrics {
		total := len(mv.Values)
		x1 := make([]float64, len(mv.Values))
		y1 := make([]float64, len(mv.Values))

		for i := 0; i < total; i++ {
			y1[i] = mv.Values[i].Value
			x1[i] = float64(mv.Values[i].Date.Unix())
		}
		data := xy{x1, y1}

		scatter, _ := plotter.NewLine(data)
		scatter.Color = GetColor(i)
		p.Add(scatter)
	}
	return p, err
}

type xy struct {
	x []float64
	y []float64
}

func (d xy) Len() int {
	return len(d.x)
}

func (d xy) XY(i int) (x, y float64) {
	x = d.x[i]
	y = d.y[i]
	return
}
