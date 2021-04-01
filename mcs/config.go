// Copyright (c) 2019-2020 The goinvest/fin developers. All rights reserved.
// Project site: https://github.com/goinvest/fin
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE file for the project.

package mcs

import (
	"encoding/json"
	"io/ioutil"
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
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return c, err
	}
	return Parse(b)
}

// Parse converts the byte slice into a Config struct.
func Parse(b []byte) (Config, error) {
	var c Config
	err := json.Unmarshal(b, &c)
	if err != nil {
		return c, err
	}
	var aux struct {
		Name        string `json:"name"`
		StartPeriod int    `json:"startPeriod"`
		EndPeriod   int    `json:"endPeriod"`
		Sims        int    `json:"sims"`
		GrowthRates []gr   `json:"growthRates"`
		Cashflows   []cf   `json:"cashflows"`
	}

	c.Name = aux.Name
	c.StartPeriod = aux.StartPeriod
	c.EndPeriod = aux.EndPeriod
	c.Sims = aux.Sims

	// Create each growth randomizer.
	for _, g := range aux.GrowthRates {
		switch g.Dist.Type {
		case "tri":

		}
	}

	return c, nil
}

type cf struct {
	Name   string `json:"name"`
	Dir    string `json:"dir"`
	Apply  string `json:"apply"`
	Dist   dist   `json:"dist"`
	Growth string `json:"growth"`
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
