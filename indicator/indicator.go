package indicator

import (
	sdk "github.com/TinkoffCreditSystems/invest-openapi-go-sdk"
	"github.com/sdcoffey/big"
	"github.com/sdcoffey/techan"
	"strconv"
	"time"
)

const BUY = "BUY"
const SELL = "SELL"
const KEEP_BUY = "Keep Buy"
const KEEP_SELL = "Keep Sell"
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

func Macd(ts *techan.TimeSeries, window int) (string, string) {
	closePriceIndicator := techan.NewClosePriceIndicator(ts)
	macd := techan.NewMACDIndicator(closePriceIndicator, 12, 26)
	macdHistogram := techan.NewMACDHistogramIndicator(macd, window)
	entryMACDConstant := techan.NewConstantIndicator(0)
	exitMACDConstant := techan.NewConstantIndicator(0)

	entryRule := techan.NewCrossUpIndicatorRule(macdHistogram, entryMACDConstant)
	exitRule := techan.NewCrossDownIndicatorRule(macdHistogram, exitMACDConstant)

	return runStrategy(entryRule, exitRule, window, len(ts.Candles)-1)
}

func Rsi(ts *techan.TimeSeries, window int) (string, string) {
	closePriceIndicator := techan.NewClosePriceIndicator(ts)
	indicator := techan.NewRelativeStrengthIndexIndicator(closePriceIndicator, window)
	entryRSIConstant := techan.NewConstantIndicator(40)
	exitRCIConstant := techan.NewConstantIndicator(60)

	entryRule := techan.NewCrossUpIndicatorRule(indicator, entryRSIConstant)
	exitRule := techan.NewCrossDownIndicatorRule(indicator, exitRCIConstant)

	return runStrategy(entryRule, exitRule, window, len(ts.Candles)-1)
}

func Arron(ts *techan.TimeSeries, window int) (string, string) {
	aroonUp := techan.NewAroonUpIndicator(techan.NewHighPriceIndicator(ts), window)
	aroonDown := techan.NewAroonDownIndicator(techan.NewLowPriceIndicator(ts), window)

	entryRule := techan.NewCrossUpIndicatorRule(aroonDown, aroonUp)
	exitRule := techan.NewCrossDownIndicatorRule(aroonUp, aroonDown)

	return runStrategy(entryRule, exitRule, window, len(ts.Candles)-1)
}

func Cci(ts *techan.TimeSeries, window int) (string, string) {
	indicator := techan.NewCCIIndicator(ts, window)
	entryCCIConstant := techan.NewConstantIndicator(-100)
	exitCCIConstant := techan.NewConstantIndicator(100)

	entryRule := techan.NewCrossUpIndicatorRule(indicator, entryCCIConstant)
	exitRule := techan.NewCrossDownIndicatorRule(indicator, exitCCIConstant)

	return runStrategy(entryRule, exitRule, window, len(ts.Candles)-1)
}

func runStrategy(entryRule, exitRule techan.Rule, window, lastIndex int) (string, string) {
	record := techan.NewTradingRecord()

	strategy := techan.RuleStrategy{
		UnstablePeriod: window, // Period before which ShouldEnter and ShouldExit will always return false
		EntryRule:      entryRule,
		ExitRule:       exitRule,
	}

	if strategy.ShouldEnter(lastIndex, record) {
		return BUY, ""
	} else if strategy.ShouldExit(lastIndex, record) {
		return SELL, ""
	}
	for i := lastIndex - 1; i > window; i-- {
		if strategy.ShouldEnter(i, record) {
			return KEEP_BUY, "should buy " + strconv.Itoa(lastIndex-window) + " day(s) ago"
		} else if strategy.ShouldExit(i, record) {
			return KEEP_SELL, "should sell " + strconv.Itoa(lastIndex-window) + " day(s) ago"
		}
	}
	return NEUTRAL, "recommendation did`t found"
}
