// Copyright (c) 2019-2025 The goinvest/fin developers. All rights reserved.
// Project site: https://github.com/goinvest/fin
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE file for the project.

package mcs

import (
	"testing"
)

func TestMinSimsPerCPU(t *testing.T) {
	testCases := []struct {
		sims     int
		cpus     int
		expected []int
	}{
		{1, 1, []int{1}},
		{2, 2, []int{1, 1}},
		{11, 2, []int{6, 5}},
		{11, 3, []int{4, 4, 3}},
		{20, 3, []int{7, 7, 6}},
		{22, 3, []int{8, 7, 7}},
		{20, 4, []int{5, 5, 5, 5}},
		{22, 4, []int{6, 6, 5, 5}},
	}
	for _, tc := range testCases {
		got := calcSimsPerCPU(tc.sims, tc.cpus)
		if len(got) != len(tc.expected) {
			t.Errorf("wrong length int slice received; got = %v / expected = %v", got, tc.expected)
		}
		for i := range got {
			if got[i] != tc.expected[i] {
				t.Errorf(
					"wrong num sims for cpu %v; got = %v / expected = %v",
					i,
					got[i],
					tc.expected[i],
				)
			}
		}
	}

}
