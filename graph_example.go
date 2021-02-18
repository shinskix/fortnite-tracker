package main

import (
	"encoding/json"
	chart "github.com/wcharczuk/go-chart/v2"
	"github.com/wcharczuk/go-chart/v2/drawing"
	"log"
	"os"
	"time"
)

func mainTest() {
	info := &PlayerInfo{}
	fd, err := os.Open("./stats.json")
	if err != nil {
		log.Fatal(err)
	}
	stat, err := fd.Stat()
	if err != nil {
		log.Fatal(err)
	}
	bytes := make([]byte, stat.Size())
	fd.Read(bytes)
	json.Unmarshal(bytes, info)

	ratings := make([]float64, 0)

	profitStyle := chart.Style{
		FillColor:   drawing.ColorFromHex("13c158"),
		StrokeColor: drawing.ColorFromHex("13c158"),
		StrokeWidth: 0,
	}

	lossStyle := chart.Style{
		FillColor:   drawing.ColorFromHex("c11313"),
		StrokeColor: drawing.ColorFromHex("c11313"),
		StrokeWidth: 0,
	}

	ratingChanges2 := make([]chart.Value, 0)
	kds := make([]float64, 0)
	for _, matchStats := range info.RecentMatchesStats {
		if matchStats.GameMode != "p10" {
			continue
		}
		ratings = append(ratings, float64(matchStats.TrnRating))
		var style chart.Style
		if matchStats.TrnRatingChange < 0 {
			style = lossStyle
		} else {
			style = profitStyle
		}
		ratingChanges2 = append(ratingChanges2, chart.Value{
			Label: time.Time(matchStats.DateCollected).Format("Mon Jan 2 15:04:05 MST 2006"),
			Value: float64(matchStats.TrnRatingChange),
			Style: style,
		})
		kds = append(kds, float64(matchStats.Kills)/float64(matchStats.Matches))
	}

	/*
		PercentChangeSeries
	*/
	//graph := chart.Chart{
	//	XAxis: chart.XAxis{
	//		TickPosition: chart.TickPositionBetweenTicks,
	//		ValueFormatter: func(v interface{}) string {
	//			typed := v.(float64)
	//			typedDate := chart.TimeFromFloat64(typed)
	//			return fmt.Sprintf("%d-%d-%d", typedDate.Month(), typedDate.Day(), typedDate.Year())
	//		},
	//	},
	//	Series: []chart.Series{
	//		chart.PercentChangeSeries{
	//			Style: chart.Style{
	//				StrokeColor: drawing.ColorRed,               // will supercede defaults
	//				FillColor:   drawing.ColorRed.WithAlpha(64), // will supercede defaults
	//			},
	//			YAxis: chart.YAxisPrimary,
	//			InnerSeries: chart.TimeSeries{
	//				XValues: dates,
	//				YValues: kds,
	//			},
	//		},
	//		chart.PercentChangeSeries{
	//			Name: ""
	//			Style: chart.Style{
	//				StrokeColor: drawing.ColorRed,               // will supercede defaults
	//				FillColor:   drawing.ColorRed.WithAlpha(64), // will supercede defaults
	//			},
	//			YAxis: chart.YAxisPrimary,
	//			InnerSeries: chart.TimeSeries{
	//				Style: chart.Style{
	//					StrokeColor: drawing.ColorBlue,               // will supercede defaults
	//					FillColor:   drawing.ColorBlue.WithAlpha(64), // will supercede defaults
	//				},
	//				YAxis:   chart.YAxisSecondary,
	//				XValues: dates,
	//				YValues: ratings,
	//			},
	//		},
	//	},
	//}

	//graph := chart.Chart{
	//	XAxis: chart.XAxis{
	//		TickPosition: chart.TickPositionBetweenTicks,
	//		ValueFormatter: func(v interface{}) string {
	//			typed := v.(float64)
	//			typedDate := chart.TimeFromFloat64(typed)
	//			return fmt.Sprintf("%d-%d-%d", typedDate.Month(), typedDate.Day(), typedDate.Year())
	//		},
	//	},
	//	Series: []chart.Series{
	//		chart.TimeSeries{
	//			Name: "k/d",
	//			Style: chart.Style{
	//				StrokeColor: drawing.ColorRed,               // will supercede defaults
	//				FillColor:   drawing.ColorRed.WithAlpha(64), // will supercede defaults
	//			},
	//			XValues: dates,
	//			YValues: kds,
	//		},
	//		chart.TimeSeries{
	//			Name: "rating",
	//			Style: chart.Style{
	//				StrokeColor: drawing.ColorBlue,               // will supercede defaults
	//				FillColor:   drawing.ColorBlue.WithAlpha(64), // will supercede defaults
	//			},
	//			YAxis:   chart.YAxisSecondary,
	//			XValues: dates,
	//			YValues: ratings,
	//		},
	//	},
	//}

	graph := chart.BarChart{
		Title: "TRN Rating Change (" + info.Name + ")",
		Background: chart.Style{
			Padding: chart.Box{
				Top: 40,
			},
		},
		Height:       512,
		BarWidth:     60,
		UseBaseValue: true,
		BaseValue:    0.0,
		Bars:         ratingChanges2,
	}

	file, err := os.Create("test-image.png")
	if err != nil {
		log.Fatal(err)
	}
	err = graph.Render(chart.PNG, file)
	if err != nil {
		log.Fatalln(err)
	}
}
