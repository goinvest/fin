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
	err = json.Unmarshal(b, &c)
	return c, err
}

// UnmarshalJSON unmarshals the given JSON byte slice into a Config struct.
func (c *Config) UnmarshalJSON(b []byte) error {
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
		return err
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

		grDistType, err := determineDistType(gr.Dist)
		if err != nil {
			return err
		}

		switch grDistType {
		case "tri":
			var t Triangle
			err := json.Unmarshal(gr.Dist, &t)
			if err != nil {
				return err
			}
			thisGrowthRate.Dist = t
		case "pert":
			var p PERT
			err := json.Unmarshal(gr.Dist, &p)
			if err != nil {
				return err
			}
			thisGrowthRate.Dist = p
		case "fixed":
			var f Fixed
			err := json.Unmarshal(gr.Dist, &f)
			if err != nil {
				return err
			}
			thisGrowthRate.Dist = f
		default:
			// FIXME(mdr): I'm missing other distribution types.
			return fmt.Errorf("bad distribution type %v in growth rate %v", grDistType, gr.Name)
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

		cfDistType, err := determineDistType(cf.Dist)
		if err != nil {
			return err
		}

		// Determine distribution type for this cashflow.
		switch cfDistType {
		case "tri":
			var t Triangle
			err := json.Unmarshal(cf.Dist, &t)
			if err != nil {
				return err
			}
			thisCashflow.Dist = t
		case "pert":
			var p PERT
			err := json.Unmarshal(cf.Dist, &p)
			if err != nil {
				return err
			}
			thisCashflow.Dist = p
		case "fixed":
			var f Fixed
			err := json.Unmarshal(cf.Dist, &f)
			if err != nil {
				return err
			}
			thisCashflow.Dist = f
		default:
			// FIXME(mdr): I'm missing other distribution types.
			return fmt.Errorf("bad distribution type %v in cashflow %v", cfDistType, cf.Name)
		}

		// Apply growth rate to this cash flow.

		c.Cashflows[i] = thisCashflow
	}

	return nil
}

type cf struct {
	Name      string          `json:"name"`
	IsOutflow bool            `json:"outflow"`
	Apply     string          `json:"apply"`
	Dist      json.RawMessage `json:"dist"`
	Growth    string          `json:"growth"`
}

type gr struct {
	Name  string          `json:"name"`
	Apply string          `json:"apply"`
	Dist  json.RawMessage `json:"dist"`
}

func determineDistType(data []byte) (string, error) {
	var aux struct {
		Type string `json:"type"`
	}
	err := json.Unmarshal(data, &aux)
	if err != nil {
		return "", err
	}
	return aux.Type, nil
}
