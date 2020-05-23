package chart

import (
	sdk "github.com/TinkoffCreditSystems/invest-openapi-go-sdk"
	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const DIR_MODE = 755
const DEFAULT_PATH = "no-data.png"

func CreateChartByCandels(candels []sdk.Candle, prefixPathImg string) string {
	if len(candels) == 0 {
		return prefixPathImg + "/" + DEFAULT_PATH
	}

	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	currentTime := time.Now()
	dirs := []string{prefixPathImg, strconv.Itoa(currentTime.Year()), strconv.Itoa(int(currentTime.Month()))}
	absPathImg := dir + string(os.PathSeparator) + strings.Join(dirs, string(os.PathSeparator)) + string(os.PathSeparator)
	if _, err := os.Stat(absPathImg); os.IsNotExist(err) {
		os.MkdirAll(absPathImg, DIR_MODE)
	}
	fileName := candels[len(candels)-1].FIGI + "_" + candels[len(candels)-1].TS.Format("20060102") + ".png"
	f, _ := os.Create(absPathImg + fileName)
	defer f.Close()

	xv, yv := xyValues(candels)

	priceSeries := chart.TimeSeries{
		Name: "SPY",
		Style: chart.Style{
			StrokeColor: chart.GetDefaultColor(0),
		},
		XValues: xv,
		YValues: yv,
	}

	smaSeries := chart.SMASeries{
		Name: "SPY - SMA",
		Style: chart.Style{
			StrokeColor:     drawing.ColorRed,
			StrokeDashArray: []float64{5.0, 5.0},
		},
		InnerSeries: priceSeries,
	}

	bbSeries := &chart.BollingerBandsSeries{
		Name: "SPY - Bol. Bands",
		Style: chart.Style{
			StrokeColor: drawing.ColorFromHex("efefef"),
			FillColor:   drawing.ColorFromHex("efefef").WithAlpha(64),
		},
		InnerSeries: priceSeries,
	}

	graph := chart.Chart{
		XAxis: chart.XAxis{
			TickPosition: chart.TickPositionBetweenTicks,
		},
		Series: []chart.Series{
			bbSeries,
			priceSeries,
			smaSeries,
		},
	}
	graph.Render(chart.PNG, f)
	return strings.Join(dirs, "/") + "/" + fileName
}

func xyValues(candels []sdk.Candle) ([]time.Time, []float64) {
	var xvalue []time.Time
	var yvalue []float64
	for _, candle := range candels {
		xvalue = append(xvalue, candle.TS)
		yvalue = append(yvalue, candle.OpenPrice)
	}
	return xvalue, yvalue
}
