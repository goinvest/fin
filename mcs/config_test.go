// Copyright (c) 2023 The goinvest/fin developers. All rights reserved.
// Project site: https://github.com/goinvest/fin
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE file for the project.

package mcs

import (
	"fmt"
	"testing"
)

func TestParseFile(t *testing.T) {
	testCases := []struct {
		filename string
		expected Config
	}{
		{
			"sample_config.json",
			Config{
				Name:        "Buy",
				StartPeriod: 1,
				EndPeriod:   48,
				Sims:        10,
				Cashflows: []Cashflow{
					{
						Name:      "Revenues",
						IsOutflow: false,
						Periods:   "1-48",
						Dist:      Triangle{50, 100, 70},
						Growth: Growth{
							Periods: "13,25,37",
							Dist:    Triangle{-0.15, 0.35, 0.15},
							Name:    "revenueGrowth",
						},
					},
					{
						Name:      "Variable Expenses",
						IsOutflow: true,
						Periods:   "1-48",
						Dist:      PERT{30, 65, 42},
						Growth: Growth{
							Periods: "13,25,37",
							Dist:    Triangle{0.0, 0.03, 0.01},
							Name:    "expenseGrowth",
						},
					},
					{
						Name:      "Fixed Expenses",
						IsOutflow: true,
						Periods:   "25-48",
						Dist:      Fixed(5),
					},
				},
			},
		},
	}
	for _, tc := range testCases {
		cfg, err := ParseFile(fmt.Sprintf("./testdata/%s", tc.filename))
		if err != nil {
			t.Errorf("error parsing json file: %s", err)
		}
		// Test config general information.
		if cfg.Name != tc.expected.Name {
			t.Errorf("got %v for name but wanted %v", cfg.Name, tc.expected.Name)
		}
		if cfg.StartPeriod != tc.expected.StartPeriod {
			t.Errorf("got %v for start period but wanted %v", cfg.StartPeriod, tc.expected.StartPeriod)
		}
		if cfg.EndPeriod != tc.expected.EndPeriod {
			t.Errorf("got %v for end period but wanted %v", cfg.EndPeriod, tc.expected.EndPeriod)
		}
		if cfg.Sims != tc.expected.Sims {
			t.Errorf("got %v for num of sims but wanted %v", cfg.Sims, tc.expected.Sims)
		}
		// Test first cashflow.
		if cfg.Cashflows[0].Name != tc.expected.Cashflows[0].Name {
			t.Errorf("got %v for name of first cashflow %v", cfg.Cashflows[0].Name, tc.expected.Cashflows[0].Name)
		}
		if cfg.Cashflows[0].IsOutflow != tc.expected.Cashflows[0].IsOutflow {
			t.Errorf("got %v for isOutflow of first cashflow; expected %v", cfg.Cashflows[0].IsOutflow, tc.expected.Cashflows[0].IsOutflow)
		}
		if cfg.Cashflows[0].Periods != tc.expected.Cashflows[0].Periods {
			t.Errorf("got %v for periods of first cashflow; expected %v", cfg.Cashflows[0].Periods, tc.expected.Cashflows[0].Periods)
		}
		switch v := cfg.Cashflows[0].Dist.(type) {
		case Triangle:
		default:
			t.Errorf("got %v / expected Triangle", v)
		}
		// Test second cashflow.
		if cfg.Cashflows[1].Name != tc.expected.Cashflows[1].Name {
			t.Errorf("got %v for name of second cashflow %v", cfg.Cashflows[1].Name, tc.expected.Cashflows[1].Name)
		}
		if cfg.Cashflows[1].IsOutflow != tc.expected.Cashflows[1].IsOutflow {
			t.Errorf("got %v for isOutflow of second cashflow %v", cfg.Cashflows[1].IsOutflow, tc.expected.Cashflows[1].IsOutflow)
		}
		if cfg.Cashflows[1].Periods != tc.expected.Cashflows[1].Periods {
			t.Errorf("got %v for periods of first cashflow; expected %v", cfg.Cashflows[1].Periods, tc.expected.Cashflows[1].Periods)
		}
		switch v := cfg.Cashflows[1].Dist.(type) {
		case PERT:
		default:
			t.Errorf("got %v / expected PERT", v)
		}
	}
}
