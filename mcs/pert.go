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
)

// PERT models a PERT distribution with the values min, max, mode but without a
// random source.
type PERT struct {
	min  float64
	max  float64
	mode float64
}

// NewPERT constructs a new PERT distribution without a random source using the
// given min, max, and mode. Constraints are min < max and min ≤ mode ≤ max.
func NewPERT(min, max, mode float64) PERT {
	checkPERTArguments(min, max, mode)
	return PERT{
		min:  min,
		max:  max,
		mode: mode,
	}
}

// Randomize sets up a new PERT distribution using the given random source.
func (p PERT) Randomize(src rand.Source) Rander {
	return distuvx.NewPERT(p.min, p.max, p.mode, src)
}

// UnmarshalJSON unmarshals the given JSON data into a PERT MCS distribution.
func (p *PERT) UnmarshalJSON(data []byte) error {
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
	if aux.Type != "pert" {
		return fmt.Errorf("instead of pert type distribution found %s", aux.Type)
	}
	*p = NewPERT(aux.Min, aux.Max, aux.Mode)
	return nil
}

// PERTOne models a PERT distribution once with the values min, max, mode, and
// then returns the same number each time. PERTOne does not have a random source.
type PERTOne struct {
	min  float64
	max  float64
	mode float64
}

// NewPERTOne constructs a new PERT distribution without a random source using the
// given min, max, and mode. Constraints are min < max and min ≤ mode ≤ max.
// Once the random number has been generated using the PERT distribution, the
// same number will be returned each time.
func NewPERTOne(min, max, mode float64) PERTOne {
	checkPERTArguments(min, max, mode)
	return PERTOne{
		min:  min,
		max:  max,
		mode: mode,
	}
}

// UnmarshalJSON unmarshals the given JSON data into a PERT MCS distribution.
func (p *PERTOne) UnmarshalJSON(data []byte) error {
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
	if aux.Type != "pert_one" {
		return fmt.Errorf("instead of pert_one type distribution found %s", aux.Type)
	}
	*p = NewPERTOne(aux.Min, aux.Max, aux.Mode)
	return nil
}

// Randomize sets up a new one-time only PERT distribution.
func (p PERTOne) Randomize(src rand.Source) Rander {
	return distuvx.NewPERTOne(p.min, p.max, p.mode, src)
}

func checkPERTArguments(min, max, mode float64) {
	if min >= max {
		panic("pert constraint of min < max violated")
	}
	if mode < min {
		panic("pert constraint of min <= mode violated")
	}
	if mode > max {
		panic("pert constraint of mode <= max violated")
	}
}
