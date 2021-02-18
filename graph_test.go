package main

import (
	"fmt"
	"github.com/wcharczuk/go-chart/v2"
	"github.com/wcharczuk/go-chart/v2/drawing"
	"math"
	"math/rand"
	"os"
	"testing"
	"time"
)

func TestChart(t *testing.T) {
	n := 60
	ratings := make([]float64, n)
	days := make([]time.Time, n)
	kds := make([]float64, n)
	now := time.Now().AddDate(0, 0, -n)
	minValue := 10000.0
	maxValue := 0.0
	ratings[0] = 1500
	days[0] = now
	kds[0] = 1.1
	rand.Seed(7)
	for i := 1; i < n; i++ {
		days[i] = now.AddDate(0, 0, i)
		sign := 1
		if rand.Intn(100) < 50 {
			sign = -1
		}
		ratings[i] = math.Max(0, ratings[i-1]+float64(rand.Intn(200)*sign))
		kds[i] = math.Max(0, kds[i-1]+float64(rand.Intn(80)*sign)/1000)
		if ratings[i] < minValue {
			minValue = ratings[i]
		}
		if ratings[i] > maxValue {
			maxValue = ratings[i]
		}
	}
	ratingSeries := chart.TimeSeries{
		Name: "Rating",
		Style: chart.Style{
			StrokeColor: drawing.ColorRed,               // will supercede defaults
			FillColor:   drawing.ColorRed.WithAlpha(64), // will supercede defaults
		},
		YAxis:   chart.YAxisPrimary,
		XValues: days,
		YValues: ratings,
	}
	graph := chart.Chart{
		YAxis: chart.YAxis{
			Range: &chart.ContinuousRange{
				Min: minValue - 0.5*minValue,
				Max: maxValue + 0.3*maxValue,
			},
		},
		XAxis: chart.XAxis{
			TickPosition: chart.TickPositionBetweenTicks,
			ValueFormatter: func(v interface{}) string {
				typed := v.(float64)
				typedDate := chart.TimeFromFloat64(typed)
				return fmt.Sprintf("%d-%d-%d", typedDate.Month(), typedDate.Day(), typedDate.Year())
			},
		},
		Series: []chart.Series{
			ratingSeries,
			chart.TimeSeries{
				Name: "K/D",
				Style: chart.Style{
					StrokeColor:     drawing.ColorBlue,
					StrokeDashArray: []float64{5.0, 5.0},
				},
				YAxis:   chart.YAxisSecondary,
				XValues: days,
				YValues: kds,
			},
			chart.SMASeries{
				Name: "Rating - SMA",
				Style: chart.Style{
					StrokeColor:     drawing.ColorBlack,
					StrokeDashArray: []float64{5.0, 5.0},
				},
				InnerSeries: ratingSeries,
			},
		},
	}

	graph.Elements = []chart.Renderable{
		chart.LegendLeft(&graph),
	}

	file, err := os.Create("test-image.png")
	if err != nil {
		t.FailNow()
	}
	err = graph.Render(chart.PNG, file)
	if err != nil {
		t.FailNow()
	}
}
