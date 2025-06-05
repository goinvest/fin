// Copyright (c) 2019-2025 The goinvest/fin developers. All rights reserved.
// Project site: https://github.com/goinvest/fin
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE file for the project.

package mcs

import (
	"log"
	"testing"
)

func TestParsePeriods(t *testing.T) {
	testCases := []struct {
		start    int
		end      int
		given    string
		expected []bool
	}{
		{1, 4, "2-3", []bool{false, true, true, false}},
		{0, 5, "1-3", []bool{false, true, true, true, false, false}},
		{0, 5, "", []bool{false, false, false, false, false, false}},
	}
	for _, tc := range testCases {
		periods, err := parsePeriods(tc.start, tc.end, tc.given)
		if err != nil {
			log.Fatalf("error parsing periods in test: %s", err)
		}
		for i, period := range periods {
			if tc.expected[i] != period {
				t.Errorf("idx %d: expected = %t / got = %t", i, tc.expected[i], period)
			}
		}
	}
}
