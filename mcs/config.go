// Copyright (c) 2019-2023 The goinvest/fin developers. All rights reserved.
// Project site: https://github.com/goinvest/fin
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE file for the project.

package mcs

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config models the net cashflows configuration information.
type Config struct {
	Name        string     `json:"name"`
	StartPeriod int        `json:"startPeriod"`
	EndPeriod   int        `json:"endPeriod"`
	Sims        int        `json:"sims"`
	Cashflows   []Cashflow `json:"cashflows"`
}

// ParseFile parses the JSON configuration file into a Config struct.
func ParseFile(filename string) (Config, error) {
	var c Config
	b, err := os.ReadFile(filename)
	if err != nil {
		return c, err
	}
	return Parse(b)
}

// Parse converts the byte slice into a Config struct.
func Parse(b []byte) (Config, error) {
	var c Config
	var aux struct {
		Name        string `json:"name"`
		StartPeriod int    `json:"startPeriod"`
		EndPeriod   int    `json:"endPeriod"`
		Sims        int    `json:"sims"`
		GrowthRates []gr   `json:"growthRates"`
		Cashflows   []cf   `json:"cashflows"`
	}

	err := json.Unmarshal(b, &aux)
	if err != nil {
		return c, err
	}
	c.Name = aux.Name
	c.StartPeriod = aux.StartPeriod
	c.EndPeriod = aux.EndPeriod
	c.Sims = aux.Sims

	// Create each growth rate.
	growthRates := make(map[string]Growth)
	for _, gr := range aux.GrowthRates {
		thisGrowthRate := Growth{
			Name:    gr.Name,
			Periods: gr.Apply,
		}
		switch gr.Dist.Type {
		case "tri":
			// FIXME(mdr): I'm hard coding the min, max, mode, which is wrong.
			thisGrowthRate.Dist = Triangle{1.0, 10.0, 5.0}
		case "pert":
			// FIXME(mdr): I'm hard coding the min, max, mode, which is wrong.
			thisGrowthRate.Dist = PERT{1.0, 10.0, 5.0}
		case "fixed":
			thisGrowthRate.Dist = Fixed(1.0)
		default:
			// FIXME(mdr): I'm missing other distribution types.
			return c, fmt.Errorf("bad distribution type %v in growth rate %v", gr.Dist.Type, gr.Name)
		}
		growthRates[gr.Name] = thisGrowthRate
	}

	// Create each cashflow.
	c.Cashflows = make([]Cashflow, len(aux.Cashflows))
	for i, cf := range aux.Cashflows {
		thisCashflow := Cashflow{
			Name:      cf.Name,
			IsOutflow: cf.IsOutflow,
			Periods:   cf.Apply,
		}

		// Determine distribution type for this cashflow.
		switch cf.Dist.Type {
		case "tri":
			// FIXME(mdr): I'm hard coding the min, max, mode, which is wrong.
			thisCashflow.Dist = Triangle{1.0, 10.0, 5.0}
		case "pert":
			// FIXME(mdr): I'm hard coding the min, max, mode, which is wrong.
			thisCashflow.Dist = PERT{1.0, 10.0, 5.0}
		case "fixed":
			thisCashflow.Dist = Fixed(1.0)
		default:
			// FIXME(mdr): I'm missing other distribution types.
			return c, fmt.Errorf("bad distribution type %v in cashflow %v", cf.Dist.Type, cf.Name)
		}

		// Apply growth rate to this cash flow.

		c.Cashflows[i] = thisCashflow
	}

	return c, nil
}

type cf struct {
	Name      string `json:"name"`
	IsOutflow bool   `json:"outflow"`
	Apply     string `json:"apply"`
	Dist      dist   `json:"dist"`
	Growth    string `json:"growth"`
}

type gr struct {
	Name  string `json:"name"`
	Apply string `json:"apply"`
	Dist  dist   `json:"dist"`
}

type dist struct {
	Type string                 `json:"type"`
	X    map[string]interface{} `json:"-"`
}
