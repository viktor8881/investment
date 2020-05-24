package main

import (
	"fmt"
	"github.com/sadlil/gologger"
	"investment/chart"
	"investment/configuration"
	"investment/email"
	"investment/tinkoffMarket"
	"net/http"
	"os"
)

const DEFAULT_CHART = "no-data.png"

func ping(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "pong")
}

type investment struct {
	CurrencyName string
	PathImg      string
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

		fileName := DEFAULT_CHART
		if len(candles) != 0 {
			fileName = chart.CreateChartByCandels(candles, config.PathImg)
		}
		var investment = investment{CurrencyName: favorit.ExchangeName, PathImg: config.AbsPathChart + fileName}
		investments = append(investments, investment)
	}

	Investments := struct {
		Investments []investment
	}{
		Investments: investments,
	}
	email.SetSetting(config.Email.Host, config.Email.Port, config.Email.Username, config.Email.Password, config.Email.To, config.Email.From, "summary")
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
