// Copyright (c) 2019-2023 The goinvest/fin developers. All rights reserved.
// Project site: https://github.com/goinvest/fin
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE file for the project.

package mcs

import (
	"encoding/json"
	"fmt"

	"github.com/goinvest/distuvx"
	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
)

// Triangle is a triangle distribution with the values min, max, mode.
type Triangle struct {
	min  float64
	max  float64
	mode float64
}

// NewTriangle returns a new triangle distribution with the values min, max,
// and mode.
func NewTriangle(min, max, mode float64) Triangle {
	checkTriangleParameters(min, max, mode)
	return Triangle{min, max, mode}
}

// Randomize sets up a new triangle distribution.
func (t Triangle) Randomize(src rand.Source) Rander {
	return distuv.NewTriangle(t.min, t.max, t.max, src)
}

// UnmarshalJSON unmarshals the given JSON data into a PERT MCS distribution.
func (t *Triangle) UnmarshalJSON(data []byte) error {
	var aux struct {
		Type string  `json:"type"`
		Min  float64 `json:"min"`
		Max  float64 `json:"max"`
		Mode float64 `json:"mode"`
	}
	err := json.Unmarshal(data, &aux)
	if err != nil {
		return err
	}
	if aux.Type != "tri" {
		return fmt.Errorf("instead of triangle type distribution found %s", aux.Type)
	}
	*t = NewTriangle(aux.Min, aux.Max, aux.Mode)
	return nil
}

// TriangleOne use a triangle distribution once and then returns the same number
// each time.
type TriangleOne struct {
	min  float64
	max  float64
	mode float64
}

func NewTriangleOne(min, max, mode float64) TriangleOne {
	checkTriangleParameters(min, max, mode)
	return TriangleOne{min, max, mode}
}

// Randomize sets up a new one-time only triangle distribution.
func (t TriangleOne) Randomize(src rand.Source) Rander {
	triangle := distuv.NewTriangle(t.min, t.max, t.mode, src)
	return distuvx.NewFixed(triangle.Rand())
}

// UnmarshalJSON unmarshals the given JSON data into a PERT MCS distribution.
func (t *TriangleOne) UnmarshalJSON(data []byte) error {
	var aux struct {
		Type string  `json:"type"`
		Min  float64 `json:"min"`
		Max  float64 `json:"max"`
		Mode float64 `json:"mode"`
	}
	err := json.Unmarshal(data, &aux)
	if err != nil {
		return err
	}
	if aux.Type != "tri_one" {
		return fmt.Errorf("instead of tri_one type distribution found %s", aux.Type)
	}
	*t = NewTriangleOne(aux.Min, aux.Max, aux.Mode)
	return nil
}

func checkTriangleParameters(min, max, mode float64) {
	if min >= max {
		panic("triangle constraint of min < max violated")
	}
	if mode < min {
		panic("triangle constraint of min <= mode violated")
	}
	if mode > max {
		panic("triangle constraint of mode <= max violated")
	}
}
