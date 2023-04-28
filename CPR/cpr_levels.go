package cpr

import (
	"math"

	"github.com/dragonzurfer/revclose"
)

func ConverToRevCloseLevelsInterface(cprLevels []CPRLevel) []revclose.LevelInterface {
	levelInterfaces := make([]revclose.LevelInterface, len(cprLevels))
	for i, level := range cprLevels {
		new_level := level
		levelInterfaces[i] = &new_level
	}
	return levelInterfaces
}

func GetCPRLevels(previousDayCandle, currentDay5MinCandles *CPRCandles) []CPRLevel {
	var cprLevels []CPRLevel

	// Calculate Pivot Points
	extendedPreviousDayCandle := &ExtendedCPRCandles{*previousDayCandle}
	extendedCurrentDay5MinCandles := &ExtendedCPRCandles{*currentDay5MinCandles}

	pivotPoints := calculatePivotPoints(extendedPreviousDayCandle)
	cprLevels = append(cprLevels, pivotPoints...)

	// Calculate init balance high low
	initBalanceLevels := getInitBalanceLevels(extendedCurrentDay5MinCandles)
	cprLevels = append(cprLevels, initBalanceLevels...)

	// Calculate Previous day high low
	previousDayHighLowLevels := getPreviousDayHighLowLevels(extendedPreviousDayCandle)
	cprLevels = append(cprLevels, previousDayHighLowLevels...)
	return cprLevels
}

func getPreviousDayHighLowLevels(candles *ExtendedCPRCandles) []CPRLevel {
	var previousDayHighLowLevels []CPRLevel

	previousDayHigh, previousDayLow, _ := candles.GetHighLowClose()

	// Create CPRLevel for Previous Day High and append it to previousDayHighLowLevels
	prevDayHighLevel := CPRLevel{
		Price: previousDayHigh,
		Type:  PDH,
	}
	previousDayHighLowLevels = append(previousDayHighLowLevels, prevDayHighLevel)

	// Create CPRLevel for Previous Day Low and append it to previousDayHighLowLevels
	prevDayLowLevel := CPRLevel{
		Price: previousDayLow,
		Type:  PDL,
	}
	previousDayHighLowLevels = append(previousDayHighLowLevels, prevDayLowLevel)

	return previousDayHighLowLevels
}

func getInitBalanceLevels(candles *ExtendedCPRCandles) []CPRLevel {
	var initBalanceLevels []CPRLevel

	// Check if the first hour exists
	if candles.GetCandlesLength() >= 12 {
		initialBalanceHigh, initialBalanceLow := candles.GetCandle(0).GetHigh(), candles.GetCandle(0).GetLow()

		for i := 1; i < 12; i++ {
			candle := candles.GetCandle(i)
			initialBalanceHigh = math.Max(initialBalanceHigh, candle.GetHigh())
			initialBalanceLow = math.Min(initialBalanceLow, candle.GetLow())
		}

		// Create CPRLevel for Initial Balance High and append it to initBalanceLevels
		initBalanceHighLevel := CPRLevel{
			Price: initialBalanceHigh,
			Type:  InitBalanceHigh,
		}
		initBalanceLevels = append(initBalanceLevels, initBalanceHighLevel)

		// Create CPRLevel for Initial Balance Low and append it to initBalanceLevels
		initBalanceLowLevel := CPRLevel{
			Price: initialBalanceLow,
			Type:  InitBalanceLow,
		}
		initBalanceLevels = append(initBalanceLevels, initBalanceLowLevel)
	}
	return initBalanceLevels
}

func calculatePivotPoints(candles *ExtendedCPRCandles) []CPRLevel {
	high, low, close := candles.GetHighLowClose()
	pivot := (high + low + close) / 3

	// Calculate bottom and top pivot
	bottomPivot := (high + low) / 2
	topPivot := (2 * pivot) - bottomPivot

	// Resistance levels
	r1 := 2*pivot - low
	r2 := pivot + high - low
	r3 := r1 + high - low
	r4 := r3 + r2 - r1

	// Support levels
	s1 := 2*pivot - high
	s2 := pivot - high + low
	s3 := s1 - high + low
	s4 := s3 + s2 - s1

	levels := []CPRLevel{
		{Price: pivot, Type: CentralPivot},
		{Price: bottomPivot, Type: BottomPivot},
		{Price: topPivot, Type: TopPivot},
		{Price: r1, Type: Resistance},
		{Price: r2, Type: Resistance},
		{Price: r3, Type: Resistance},
		{Price: r4, Type: Resistance},
		{Price: s1, Type: Support},
		{Price: s2, Type: Support},
		{Price: s3, Type: Support},
		{Price: s4, Type: Support},
	}

	if levels[1].Price > levels[2].Price {
		temp := levels[1].Price
		levels[1].Price = levels[2].Price
		levels[2].Price = temp
	}
	return levels
}
