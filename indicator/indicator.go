package indicator

import (
	"fmt"
	sdk "github.com/TinkoffCreditSystems/invest-openapi-go-sdk"
	"github.com/sdcoffey/big"
	"github.com/sdcoffey/techan"
	"time"
)

const BUY = "BUY"
const SELL = "SELL"
const NEUTRAL = "neutral"

func NewSeries(candles []sdk.Candle) *techan.TimeSeries {
	ts := techan.NewTimeSeries()
	for _, datum := range candles {
		period := techan.NewTimePeriod(time.Unix(datum.TS.Unix(), 0), time.Hour*24)

		candle := techan.NewCandle(period)
		candle.OpenPrice = big.NewDecimal(datum.OpenPrice)
		candle.ClosePrice = big.NewDecimal(datum.ClosePrice)
		candle.MaxPrice = big.NewDecimal(datum.HighPrice)
		candle.MinPrice = big.NewDecimal(datum.LowPrice)

		ts.AddCandle(candle)
	}
	return ts
}

func Macd(ts *techan.TimeSeries) (string, string) {
	closePriceIndicator := techan.NewClosePriceIndicator(ts)
	macd := techan.NewMACDIndicator(closePriceIndicator, 12, 26)
	macdHistogram := techan.NewMACDHistogramIndicator(macd, 9)
	entryMACDConstant := techan.NewConstantIndicator(0)
	exitMACDConstant := techan.NewConstantIndicator(0)

	entryRule := techan.NewCrossUpIndicatorRule(macdHistogram, entryMACDConstant)

	exitRule := techan.NewCrossDownIndicatorRule(macdHistogram, exitMACDConstant)

	record := techan.NewTradingRecord()

	strategy := techan.RuleStrategy{
		UnstablePeriod: 10, // Period before which ShouldEnter and ShouldExit will always return false
		EntryRule:      entryRule,
		ExitRule:       exitRule,
	}

	i := len(ts.Candles) - 1
	value := macd.Calculate(i).FormattedString(2)
	if strategy.ShouldEnter(i, record) {
		return BUY, value
	} else if strategy.ShouldExit(i, record) {
		return SELL, value
	}
	return NEUTRAL, value
}

func Rsi(ts *techan.TimeSeries) (string, string) {
	closePriceIndicator := techan.NewClosePriceIndicator(ts)
	indicator := techan.NewRelativeStrengthIndexIndicator(closePriceIndicator, 14)
	entryRSIConstant := techan.NewConstantIndicator(40)
	exitRCIConstant := techan.NewConstantIndicator(60)

	entryRule := techan.NewCrossUpIndicatorRule(indicator, entryRSIConstant)

	exitRule := techan.NewCrossDownIndicatorRule(indicator, exitRCIConstant)

	record := techan.NewTradingRecord()

	strategy := techan.RuleStrategy{
		UnstablePeriod: 10, // Period before which ShouldEnter and ShouldExit will always return false
		EntryRule:      entryRule,
		ExitRule:       exitRule,
	}

	i := len(ts.Candles) - 1
	value := indicator.Calculate(i).FormattedString(2)
	if strategy.ShouldEnter(i, record) {
		return BUY, value
	} else if strategy.ShouldExit(i, record) {
		return SELL, value
	}
	return NEUTRAL, value
}

func Arron(ts *techan.TimeSeries) (string, string) {
	aroonUp := techan.NewAroonUpIndicator(techan.NewHighPriceIndicator(ts), 25)
	aroonDown := techan.NewAroonDownIndicator(techan.NewLowPriceIndicator(ts), 25)

	entryRule := techan.NewCrossUpIndicatorRule(aroonDown, aroonUp)
	exitRule := techan.NewCrossDownIndicatorRule(aroonUp, aroonDown)

	record := techan.NewTradingRecord()

	strategy := techan.RuleStrategy{
		UnstablePeriod: 10, // Period before which ShouldEnter and ShouldExit will always return false
		EntryRule:      entryRule,
		ExitRule:       exitRule,
	}

	i := len(ts.Candles) - 1
	value := fmt.Sprintf("%f", aroonUp.Calculate(i).Float()-aroonDown.Calculate(i).Float())
	if strategy.ShouldEnter(i, record) {
		return BUY, value
	} else if strategy.ShouldExit(i, record) {
		return SELL, value
	}
	return NEUTRAL, value
}
