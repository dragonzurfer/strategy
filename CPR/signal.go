package cpr

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/dragonzurfer/revclose"
)

func processCrossedLevels(LevelsWithinRange []revclose.LevelInterface) []CPRLevel {
	crossedLevels := []CPRLevel{}
	for _, levelInterface := range LevelsWithinRange {

		level, _ := levelInterface.(*CPRLevel)
		crossedLevels = append(crossedLevels, *level)

	}
	return crossedLevels
}

func convertToCPRSignal(revSignal *revclose.Signal, levels *[]CPRLevel) Signal {
	cprSignal := initializeCPRSignal(revSignal)

	sortLevelsByPrice(levels)

	startIndex := findStartIndex(revSignal.Value, revSignal.Close, levels)
	// fmt.Println(startIndex, levels)
	if startIndex != math.MinInt32 {
		setTargetLevelsAndPrice(revSignal.Value, startIndex, *levels, &cprSignal)
		setMessage(&cprSignal, levels)
	}

	if len(revSignal.LevelsWithinRange) == 0 {
		setNeutralSignal(&cprSignal)
		cprSignal.Message = "No level has crossed"
		return cprSignal
	}

	if len(cprSignal.TargetLevels) < 1 {
		setNeutralSignal(&cprSignal)
		cprSignal.Message = "No more levels"
		return cprSignal
	}

	if !isRiskRewardRatioOneToOne(cprSignal) {
		setNeutralSignal(&cprSignal)
		cprSignal.Message = "No good RR"
		return cprSignal
	}

	if !minTargetAbailable(cprSignal.TargetPrice, cprSignal.EntryPrice) {
		setNeutralSignal(&cprSignal)
		cprSignal.Message = "Not enough points"
		return cprSignal
	}
	// set cprSignal.Message to Prices it crossed comma seperated appended by the level Description (123,1231,12312) (PREVIOUS DAY HIGH,RESISTANCE, CENTRAL PIVOT)
	return cprSignal
}

func setMessage(cprSignal *Signal, levels *[]CPRLevel) {
	var crossedLevelsStr []string

	for _, level := range cprSignal.CrossedLevels {
		crossedLevelsStr = append(crossedLevelsStr, fmt.Sprintf("%.2f (%s)", level.Price, level.Type))
	}

	cprSignal.Message = strings.Join(crossedLevelsStr, ", ")
}

func initializeCPRSignal(revSignal *revclose.Signal) Signal {
	return Signal{
		Signal:        SignalValue(revSignal.Value),
		EntryPrice:    revSignal.Close,
		StopLossPrice: revSignal.Reversal,
		TargetPrice:   0,
		TargetLevels:  []CPRLevel{},
		CrossedLevels: processCrossedLevels(revSignal.LevelsWithinRange),
	}
}

func sortLevelsByPrice(levels *[]CPRLevel) {
	sort.Slice(*levels, func(i, j int) bool {
		return (*levels)[i].Price < (*levels)[j].Price
	})
}

func findStartIndex(value revclose.SignalValue, closePrice float64, levels *[]CPRLevel) int {

	levelsLen := len(*levels)
	if value == revclose.Sell {

		for i, level := range *levels {
			if closePrice < level.Price {
				return i
			}
		}
		return levelsLen
	}
	if value == revclose.Buy {
		for i := levelsLen - 1; i >= 0; i-- {
			if closePrice > (*levels)[i].Price && value == revclose.Buy {
				return i
			}
		}
		return -1
	}
	return math.MinInt32
}

func setTargetLevelsAndPrice(value revclose.SignalValue, startIndex int, levels []CPRLevel, cprSignal *Signal) {
	if value == revclose.Buy {
		cprSignal.TargetLevels = levels[startIndex+1:]
	} else if value == revclose.Sell {
		cprSignal.TargetLevels = levels[:startIndex]
	}
	length := len(cprSignal.TargetLevels)
	if length > 0 && value == revclose.Sell {
		cprSignal.TargetPrice = cprSignal.TargetLevels[length-1].Price
	}
	if length > 0 && value == revclose.Buy {
		cprSignal.TargetPrice = cprSignal.TargetLevels[0].Price
	}
}

func isRiskRewardRatioOneToOne(cprSignal Signal) bool {
	return math.Abs(cprSignal.StopLossPrice-cprSignal.EntryPrice) <= math.Abs(cprSignal.TargetPrice-cprSignal.EntryPrice)
}

func minTargetAbailable(targetPrice, entryPrice float64) bool {
	if targetPrice < entryPrice {
		return targetPrice <= entryPrice*(1-MIN_POINTS_MULTIPLIER)
	} else {
		return targetPrice >= entryPrice*(1+MIN_POINTS_MULTIPLIER)
	}
}

func setNeutralSignal(cprSignal *Signal) {
	cprSignal.Signal = Neutral
	cprSignal.TargetPrice = 0
	cprSignal.TargetLevels = []CPRLevel{}
}
