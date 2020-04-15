// Copyright (c) 2019-2020 The goinvest/fin developers. All rights reserved.
// Project site: https://github.com/goinvest/fin
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE file for the project.

package mcs

import (
	"log"
	"math"
	"testing"
)

const tolerance = 0.000001

func TestCalcPeriods(t *testing.T) {
	testCases := []struct {
		start    int
		end      int
		given    string
		expected []float64
	}{
		{1, 4, "1-4", []float64{1, 1, 1, 1}},
	}
	for _, tc := range testCases {
		periods, err := calcPeriods(tc.start, tc.end, tc.given)
		if err != nil {
			log.Fatalf("error: %s", err)
		}
		for i, period := range periods {
			if !almostEqual(tc.expected[i], period) {
				t.Errorf("expected = %f / got = %f", tc.expected[i], period)
			}
		}
	}
}

func almostEqual(f1, f2 float64) bool {
	return math.Abs(f1-f2) < tolerance
}
