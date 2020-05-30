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

func CreateChartByCandels(candles []sdk.Candle, PathImg string) string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	currentTime := time.Now()
	dirs := []string{strconv.Itoa(currentTime.Year()), strconv.Itoa(int(currentTime.Month()))}
	absPathImg := dir + string(os.PathSeparator) + PathImg + string(os.PathSeparator) + strings.Join(dirs, string(os.PathSeparator)) + string(os.PathSeparator)
	if _, err := os.Stat(absPathImg); os.IsNotExist(err) {
		os.MkdirAll(absPathImg, 0755)
	}
	// dd_HHii
	fileName := candles[len(candles)-1].FIGI + "_" + candles[len(candles)-1].TS.Format("02_1504") + ".png"
	f, _ := os.Create(absPathImg + fileName)
	defer f.Close()
	os.Chmod(absPathImg+fileName, 0444)

	xv, yv := xyValues(candles)

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
