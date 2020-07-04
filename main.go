package main

import (
	"fmt"
	"github.com/sadlil/gologger"
	"github.com/sdcoffey/techan"
	"investment/chart"
	"investment/configuration"
	"investment/email"
	"investment/indicator"
	"investment/tinkoffMarket"
	"net/http"
	"os"
)

const INDICATOR_MACD = 1
const INDICATOR_RSI = 2
const INDICATOR_AROON = 3
const INDICATOR_CCI = 4

func ping(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "pong")
}

type Indicator struct {
	IName              string
	IDescr             string
	IRecomendation     string
	ILastRecomandation string
}

type investment struct {
	CurrencyName string
	PathImg      string
	Indicators   []*Indicator
}

func newIndicator(name int8, timeSeries *techan.TimeSeries) *Indicator {
	ind := Indicator{}
	switch name {
	case INDICATOR_MACD:
		ind.IName = "MACD"
		ind.IDescr = "MACD Histogram - рекомендации к покупке/продаже"
		ind.IRecomendation, ind.ILastRecomandation = indicator.Macd(timeSeries, 9)
	case INDICATOR_RSI:
		ind.IName = "RSI"
		ind.IDescr = "Relative Strength Index - показывает возможный разворот"
		ind.IRecomendation, ind.ILastRecomandation = indicator.Rsi(timeSeries, 14)
	case INDICATOR_AROON:
		ind.IName = "Aroon"
		ind.IDescr = "Aroon - - рекомендации к покупке/продаже"
		ind.IRecomendation, ind.ILastRecomandation = indicator.Arron(timeSeries, 25)
	case INDICATOR_CCI:
		ind.IName = "CCI"
		ind.IDescr = "CCI - - рекомендации к покупке/продаже"
		ind.IRecomendation, ind.ILastRecomandation = indicator.Cci(timeSeries, 40)
	}
	return &ind
}

func sendEmail(w http.ResponseWriter, req *http.Request) {
	config, err := configuration.LoadConfiguration("config/config.json")
	if err != nil {
		fmt.Fprintf(w, "error read configuration:"+err.Error())
	}
	logger := gologger.GetLogger(gologger.FILE, config.LogFile)

	client, err := tinkoffMarket.CreateClient(config.Sdk.Token)
	if err != nil {
		logger.Log(err.Error())
	}

	var investments []investment
	for _, favorit := range config.Sdk.Favorits {
		candles, err := client.GetDataByFigiFor4Month(favorit.FIGI)
		if err != nil {
			logger.Log(err.Error())
		}

		fileName := chart.CreateChartByCandels(candles, config.PathImg)
		var investment = investment{CurrencyName: favorit.ExchangeName, PathImg: config.AbsPathChart + fileName}
		if len(candles) != 0 {
			series := indicator.NewSeries(candles)
			investment.Indicators = append(investment.Indicators, newIndicator(INDICATOR_MACD, series))
			investment.Indicators = append(investment.Indicators, newIndicator(INDICATOR_RSI, series))
			investment.Indicators = append(investment.Indicators, newIndicator(INDICATOR_AROON, series))
			investment.Indicators = append(investment.Indicators, newIndicator(INDICATOR_CCI, series))
		}
		investments = append(investments, investment)
	}

	Investments := struct {
		Investments []investment
	}{
		Investments: investments,
	}
	email.SetSetting(config.Email.Host, config.Email.Port, config.Email.UseSmtpAuth, config.Email.Username, config.Email.Password, config.Email.To, config.Email.From, "summary netangels")
	_, err = email.SendAnalytics("template/email-layout.html", Investments)
	if err != nil {
		logger.Log(err.Error())
	}
	fmt.Fprintf(w, "email has sent")
}

func main() {
	http.HandleFunc("/ping", ping)
	http.HandleFunc("/email", sendEmail)

	APP_IP := os.Getenv("APP_IP")
	APP_PORT := os.Getenv("APP_PORT")
	if APP_IP == "" {
		APP_IP = "127.0.0.1"
	}
	if APP_PORT == "" {
		APP_PORT = "8080"
	}
	http.ListenAndServe(APP_IP+":"+APP_PORT, nil)
}
