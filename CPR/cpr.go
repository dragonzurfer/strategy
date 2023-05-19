package cpr

import (
	"github.com/dragonzurfer/revclose"
)

var MIN_POINTS_MULTIPLIER float64
var MIN_POINTS_MULTIPLIER_SL float64
var MAX_POINT_MULTIPLIER_SL float64

type SignalValue string

const (
	Buy     SignalValue = "BUY"
	Sell    SignalValue = "SELL"
	Neutral SignalValue = "NEUTRAL"
)

type CPRCandles interface {
	GetCandle(int) revclose.RevCandle
	GetCandlesLength() int
}

type Signal struct {
	Signal        SignalValue
	EntryPrice    float64
	StopLossPrice float64
	TargetPrice   float64
	TargetLevels  []CPRLevel
	CrossedLevels []CPRLevel
	Message       string
}

func GetCPRSignal(minPointsStopLossPercent float64, minPointsPercent float64, previousDayCandle, currendDay5MinCandles CPRCandles) Signal {
	MIN_POINTS_MULTIPLIER = minPointsPercent / 100
	MIN_POINTS_MULTIPLIER_SL = minPointsStopLossPercent / 100
	MAX_POINT_MULTIPLIER_SL = MIN_POINTS_MULTIPLIER * 1.2
	var signal Signal
	currenDayCandlesLength := currendDay5MinCandles.GetCandlesLength()
	closingPrice := currendDay5MinCandles.GetCandle(currenDayCandlesLength - 1).GetClose()
	levels := GetCPRLevels(&previousDayCandle, &currendDay5MinCandles)
	revCloseLevelInterface := ConverToRevCloseLevelsInterface(levels)
	revCloseSignal := revclose.GetSignal(currendDay5MinCandles, revCloseLevelInterface, closingPrice)
	signal = convertToCPRSignal(&revCloseSignal, &levels)

	return signal
}
