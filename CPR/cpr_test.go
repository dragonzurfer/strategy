package cpr_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/dragonzurfer/revclose"
	cpr "github.com/dragonzurfer/strategy/CPR"
)

type TestCandle struct {
	Open  float64 `json:"open"`
	High  float64 `json:"high"`
	Low   float64 `json:"low"`
	Close float64 `json:"close"`
}
type TestSignal struct {
	Signal        cpr.SignalValue `json:"signal"`
	EntryPrice    float64         `json:"entry_price"`
	StopLossPrice float64         `json:"stop_loss_price"`
	TargetPrice   float64         `json:"target_price"`
	TargetLevels  []cpr.CPRLevel  `json:"target_levels"`
	CrossedLevels []cpr.CPRLevel  `json:"crossed_levels"`
	Message       string          `json:"message"`
}

type TestCase struct {
	MinPointsPercent      float64      `json:"min_points_percent"`
	PreviousDayCandle     []TestCandle `json:"previous_day_candle"`
	CurrentDay5MinCandles []TestCandle `json:"current_day_5_min_candles"`
	ExpectedSignal        TestSignal   `json:"expected_signal"`
}

func (ts *TestSignal) ToSignal() cpr.Signal {
	return cpr.Signal{
		Signal:        ts.Signal,
		EntryPrice:    ts.EntryPrice,
		StopLossPrice: ts.StopLossPrice,
		TargetPrice:   ts.TargetPrice,
		TargetLevels:  ts.TargetLevels,
		CrossedLevels: ts.CrossedLevels,
		Message:       ts.Message,
	}
}

func (tc *TestCandle) GetOpen() float64 {
	return tc.Open
}

func (tc *TestCandle) GetHigh() float64 {
	return tc.High
}

func (tc *TestCandle) GetLow() float64 {
	return tc.Low
}

func (tc *TestCandle) GetClose() float64 {
	return tc.Close
}

func (tc *TestCandle) GetOHLC() (float64, float64, float64, float64) {
	return tc.Open, tc.High, tc.Low, tc.Close
}

type TestCandles []TestCandle

func (tc TestCandles) GetCandle(index int) revclose.RevCandle {
	return &tc[index]
}

func (tc TestCandles) GetCandlesLength() int {
	return len(tc)
}

const tolerance = 0.7 // You can adjust this value to your needs

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) <= tolerance
}

func signalsAlmostEqual(tc string, expected, actual cpr.Signal) error {
	fmt.Errorf("Error on test:%v\n" + tc)
	if expected.Signal != actual.Signal {
		return fmt.Errorf("Signal: expected %v, got %v", expected.Signal, actual.Signal)
	}

	if !almostEqual(expected.EntryPrice, actual.EntryPrice) {
		return fmt.Errorf("EntryPrice: expected %f, got %f", expected.EntryPrice, actual.EntryPrice)
	}

	if !almostEqual(expected.StopLossPrice, actual.StopLossPrice) {
		return fmt.Errorf("StopLossPrice: expected %f, got %f", expected.StopLossPrice, actual.StopLossPrice)
	}

	if !almostEqual(expected.TargetPrice, actual.TargetPrice) {
		return fmt.Errorf("TargetPrice: expected %f, got %f", expected.TargetPrice, actual.TargetPrice)
	}

	// Check TargetLevels
	if len(expected.TargetLevels) != len(actual.TargetLevels) {
		return fmt.Errorf("TargetLevels length: expected %d, got %d", len(expected.TargetLevels), len(actual.TargetLevels))
	}
	for i := range expected.TargetLevels {
		if expected.TargetLevels[i].Type != actual.TargetLevels[i].Type ||
			!almostEqual(expected.TargetLevels[i].Price, actual.TargetLevels[i].Price) {
			return fmt.Errorf("TargetLevels[%d]: expected {Price: %f, Type: %v}, got {Price: %f, Type: %v}",
				i, expected.TargetLevels[i].Price, expected.TargetLevels[i].Type, actual.TargetLevels[i].Price, actual.TargetLevels[i].Type)
		}
	}

	// Check CrossedLevels
	if len(expected.CrossedLevels) != len(actual.CrossedLevels) {
		return fmt.Errorf("CrossedLevels length: expected %d, got %d", len(expected.CrossedLevels), len(actual.CrossedLevels))
	}
	for i := range expected.CrossedLevels {
		if expected.CrossedLevels[i].Type != actual.CrossedLevels[i].Type ||
			!almostEqual(expected.CrossedLevels[i].Price, actual.CrossedLevels[i].Price) {
			return fmt.Errorf("CrossedLevels[%d]: expected {Price: %f, Type: %v}, got {Price: %f, Type: %v}",
				i, expected.CrossedLevels[i].Price, expected.CrossedLevels[i].Type, actual.CrossedLevels[i].Price, actual.CrossedLevels[i].Type)
		}
	}

	if expected.Message != actual.Message {
		switch tc {
		case "testcase8.json":
			if expected.Message != actual.Message {
				return fmt.Errorf("expected:%v \nactual:%v\n", expected.Message, actual.Message)
			}
		}
	}

	return nil
}

func TestGetCPRSignal(t *testing.T) {
	testCases := []string{
		"testcase1.json",
		"testcase2.json",
		"testcase3.json",
		"testcase4.json",
		"testcase5.json",
		"testcase6.json",
		"testcase7.json",
		"testcase8.json",
		"testcase9.json",
		"testcase10.json",
		"testcase11.json",
		"testcase12.json",
	}

	for _, tc := range testCases {
		wd, err := os.Getwd()
		if err != nil {
			t.Fatalf("Error getting working directory: %v", err)
		}

		data, err := ioutil.ReadFile(filepath.Join(wd, "testcases", tc))

		if err != nil {
			t.Fatalf("Error reading %s: %v", tc, err)
		}

		var testCase TestCase
		err = json.Unmarshal(data, &testCase)
		if err != nil {
			t.Fatalf("Error unmarshalling JSON data in %s: %v", tc, err)
		}

		t.Run(fmt.Sprintf("TestCase_%s", tc), func(t *testing.T) {
			prevDayCandles := TestCandles(testCase.PreviousDayCandle)
			currDayCandles := TestCandles(testCase.CurrentDay5MinCandles)

			actualSignal := cpr.GetCPRSignal(0.05, testCase.MinPointsPercent, prevDayCandles, currDayCandles)
			expectedSignal := testCase.ExpectedSignal.ToSignal()

			// Sort the TargetLevels in both expected and actual signals
			sort.Slice(expectedSignal.TargetLevels, func(i, j int) bool {
				return expectedSignal.TargetLevels[i].Price < expectedSignal.TargetLevels[j].Price
			})
			sort.Slice(actualSignal.TargetLevels, func(i, j int) bool {
				return actualSignal.TargetLevels[i].Price < actualSignal.TargetLevels[j].Price
			})

			sort.Slice(expectedSignal.CrossedLevels, func(i, j int) bool {
				return expectedSignal.CrossedLevels[i].Price < expectedSignal.CrossedLevels[j].Price
			})
			sort.Slice(actualSignal.CrossedLevels, func(i, j int) bool {
				return actualSignal.CrossedLevels[i].Price < actualSignal.CrossedLevels[j].Price
			})

			if err := signalsAlmostEqual(tc, expectedSignal, actualSignal); err != nil {
				t.Errorf("Signals not equal: %v", err)
				t.Errorf("Expected: %+v\n Actual:%+v\n", expectedSignal, actualSignal)
			} else {
				fmt.Println(actualSignal.Message)
			}
		})
	}
}
