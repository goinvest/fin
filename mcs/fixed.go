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

// Fixed is a fixed number.
type Fixed float64

// NewFixed returns a Fixed distribution.
func NewFixed(num float64) Fixed {
	return Fixed(num)
}

// Randomize a new fixed number.
func (f Fixed) Randomize(src rand.Source) Rander {
	return distuvx.NewFixed(float64(f))
}

// String implements the io.Stringer interface.
func (f Fixed) String() string {
	return fmt.Sprintf("Fixed value of %f", f)
}

// UnmarshalJSON unmarshals the given JSON data into a Fixed MCS distribution.
func (f *Fixed) UnmarshalJSON(data []byte) error {
	var aux struct {
		Type  string  `json:"type"`
		Value float64 `json:"val"`
	}
	err := json.Unmarshal(data, &aux)
	if err != nil {
		return err
	}
	if aux.Type != "fixed" {
		return fmt.Errorf("instead of fixed type distribution found %s", aux.Type)
	}
	*f = Fixed(aux.Value)
	return nil
}
