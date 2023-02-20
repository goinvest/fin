// Copyright (c) 2019-2023 The goinvest/fin developers. All rights reserved.
// Project site: https://github.com/goinvest/fin
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE file for the project.

package mcs

import (
	"encoding/json"
	"fmt"

	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
)

// Uniform is a uniform distribution with the values min and max.
type Uniform struct {
	min float64
	max float64
}

// NewUniform constructs a new Uniform distribution without a random source
// using the given min and max values.
func NewUniform(min, max float64) Uniform {
	checkUniformArguments(min, max)
	return Uniform{min, max}
}

// Randomize sets up a new uniform distribution.
func (u Uniform) Randomize(src rand.Source) Rander {
	return distuv.Uniform{
		Min: u.min,
		Max: u.max,
		Src: src,
	}
}

// UnmarshalJSON unmarshals the given JSON data into a PERT MCS distribution.
func (u *Uniform) UnmarshalJSON(data []byte) error {
	var aux struct {
		Type string  `json:"type"`
		Min  float64 `json:"min"`
		Max  float64 `json:"max"`
	}
	err := json.Unmarshal(data, &aux)
	if err != nil {
		return err
	}
	if aux.Type != "uniform" {
		return fmt.Errorf("instead of uniform type distribution found %s", aux.Type)
	}
	*u = NewUniform(aux.Min, aux.Max)
	return nil
}

func checkUniformArguments(min, max float64) {
	if min > max {
		panic("uniform constraint of min <= max violated")
	}
}
