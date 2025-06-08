// Copyright (c) 2019-2025 The goinvest/fin developers. All rights reserved.
// Project site: https://github.com/goinvest/fin
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE file for the project.

package cf

import (
	"math"
	"testing"
)

const tolerance = 0.000001

func TestMIRR(t *testing.T) {
	testCases := []struct {
		cashflows     []float64
		costOfCapital float64
		expected      float64
	}{
		{[]float64{-1000, 500, 400, 300, 100}, 0.10, 0.121063},
		{[]float64{-1000, 100, 300, 400, 600}, 0.10, 0.113281},
		{[]float64{1000, 100, 300, 400, 600}, 0.10, math.Inf(1)},
	}
	for _, tc := range testCases {
		mirr := MIRR(tc.cashflows, tc.costOfCapital)
		if !math.IsInf(mirr, 1) && !almostEqual(tc.expected, mirr) {
			t.Errorf("MIRR calculated = %f, expected = %f", mirr, tc.expected)
		} else if math.IsInf(mirr, 1) && !math.IsInf(tc.expected, 1) {
			t.Errorf("MIRR = +Inf, expected = %f", tc.expected)
		}
	}
}

func TestIRR(t *testing.T) {
	testCases := []struct {
		cashflows []float64
		expected  float64
		options   IRROptions
	}{
		{[]float64{-1000, 500, 400, 300, 100}, 0.144888, IRROptions{0.0, 1e-5, 10}},
		{[]float64{-1000, 100, 300, 400, 600}, 0.117906, IRROptions{0.0, 1e-5, 10}},
	}
	for _, tc := range testCases {
		irr := IRR(tc.cashflows, tc.options)
		if !math.IsNaN(irr) && !almostEqual(tc.expected, irr) {
			t.Errorf("IRR calculated = %f, expected = %f", irr, tc.expected)
		} else if math.IsNaN(irr) && !math.IsNaN(tc.expected) {
			t.Errorf("IRR = NaN, expected = %f", tc.expected)
		}
	}
}

func TestIRRWithOptions(t *testing.T) {
	testCases := []struct {
		cashflows []float64
		expected  float64
	}{
		{[]float64{-1000, 500, 400, 300, 100}, 0.144888},
		{[]float64{-1000, 100, 300, 400, 600}, 0.117906},
		{[]float64{1000, 100, 300, 400, 600}, math.NaN()},
	}
	for _, tc := range testCases {
		irr := IRR(tc.cashflows)
		if !math.IsNaN(irr) && !almostEqual(tc.expected, irr) {
			t.Errorf("IRR calculated = %f, expected = %f", irr, tc.expected)
		} else if math.IsNaN(irr) && !math.IsNaN(tc.expected) {
			t.Errorf("IRR = NaN, expected = %f", tc.expected)
		}
	}
}

func TestDiscountedPaybackPeriod(t *testing.T) {
	testCases := []struct {
		cashflows    []float64
		discountRate float64
		expected     float64
	}{
		{[]float64{-1000, 500, 400, 300, 100}, 0.10, 2.9533333},
		{[]float64{-1000, 100, 300, 400, 600}, 0.10, 3.8800000},
		{[]float64{-1000, -100, -300, -400, -600}, 0.10, math.NaN()},
	}
	for _, tc := range testCases {
		paybackPeriod := DiscountedPaybackPeriod(tc.cashflows, tc.discountRate)
		if !math.IsNaN(paybackPeriod) && !almostEqual(tc.expected, paybackPeriod) {
			t.Errorf("Payback Period calculated = %f, expected = %f", paybackPeriod, tc.expected)
		} else if math.IsNaN(paybackPeriod) && !math.IsNaN(tc.expected) {
			t.Errorf("Payback Period = NaN, expected = %f", tc.expected)
		}
	}
}

func TestPaybackPeriod(t *testing.T) {
	testCases := []struct {
		cashflows []float64
		expected  float64
	}{
		{[]float64{-1000, 500, 400, 300, 100}, 2.3333333},
		{[]float64{-1000, 100, 300, 400, 600}, 3.3333333},
		{[]float64{-1000, -100, -300, -400, -600}, math.NaN()},
	}
	for _, tc := range testCases {
		paybackPeriod := PaybackPeriod(tc.cashflows)
		if !math.IsNaN(paybackPeriod) && !almostEqual(tc.expected, paybackPeriod) {
			t.Errorf("Payback Period calculated = %f, expected = %f", paybackPeriod, tc.expected)
		} else if math.IsNaN(paybackPeriod) && !math.IsNaN(tc.expected) {
			t.Errorf("Payback Period = NaN, expected = %f", tc.expected)
		}
	}
}

func TestNPV(t *testing.T) {
	testCases := []struct {
		cashflows    []float64
		discountRate float64
		expected     float64
	}{
		{[]float64{-1000, 500, 400, 300, 100}, 0.10, 78.819753},
		{[]float64{-3000, 1300, 1300, 1300}, 0.08, 350.226083},
	}
	for _, tc := range testCases {
		npv := NPV(tc.cashflows, tc.discountRate)
		if !almostEqual(tc.expected, npv) {
			t.Errorf("NPV calculated = %f, expected = %f", npv, tc.expected)
		}
	}
}

func TestNCF(t *testing.T) {
	testCases := []struct {
		cashflows []float64
		expected  float64
	}{
		{[]float64{-1000, 500, 400, 300, 100}, 300},
		{[]float64{-1000, 100, 300, 400, 600}, 400},
	}
	for _, tc := range testCases {
		ncf := NCF(tc.cashflows)
		if !almostEqual(tc.expected, ncf) {
			t.Errorf("NCF calculated = %f, expected = %f", ncf, tc.expected)
		}
	}
}

func almostEqual(f1, f2 float64) bool {
	return math.Abs(f1-f2) < tolerance
}
