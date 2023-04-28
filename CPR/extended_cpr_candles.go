package cpr

import "math"

type ExtendedCPRCandles struct {
	CPRCandles
}

func (ec *ExtendedCPRCandles) GetHighLowClose() (high, low, close float64) {
	length := ec.GetCandlesLength()
	high, low = ec.GetCandle(0).GetHigh(), ec.GetCandle(0).GetLow()

	for i := 1; i < length; i++ {
		candle := ec.GetCandle(i)
		high = math.Max(high, candle.GetHigh())
		low = math.Min(low, candle.GetLow())
	}

	close = ec.GetCandle(length - 1).GetClose()

	return high, low, close
}
