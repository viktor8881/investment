package tinkoffMarket

import (
	"context"
	sdk "github.com/TinkoffCreditSystems/invest-openapi-go-sdk"
	"math/rand"
	"time"
)

type Client struct {
	restClient *sdk.SandboxRestClient
}

var client Client

func CreateClient(token string) (Client, error) {
	rand.Seed(time.Now().UnixNano()) // инициируем Seed рандома для функции requestID
	client = Client{sdk.NewSandboxRestClient(token)}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.restClient.Register(ctx, sdk.AccountTinkoff)
	if err != nil {
		return client, err
		//log.Fatalln(errorHandle(err))
	}
	return client, nil
}

func (client *Client) GetDataByFigiFor4Month(figi string) ([]sdk.Candle, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	candles, err := client.restClient.Candles(ctx, time.Now().AddDate(0, -4, 0), time.Now(), sdk.CandleInterval1Day, figi)
	if err != nil {
		return nil, err
		//log.Fatalln(err)
	}
	return candles, nil
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// Генерируем уникальный ID для запроса
func requestID() string {
	b := make([]rune, 12)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return string(b)
}

func errorHandle(err error) error {
	if err == nil {
		return nil
	}

	if tradingErr, ok := err.(sdk.TradingError); ok {
		if tradingErr.InvalidTokenSpace() {
			tradingErr.Hint = "Do you use sandbox token in production environment or vise verse?"
			return tradingErr
		}
	}

	return err
}
