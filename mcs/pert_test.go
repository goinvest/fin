// Copyright (c) 2023 The goinvest/fin developers. All rights reserved.
// Project site: https://github.com/goinvest/fin
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE file for the project.

package mcs

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestNewPERT(t *testing.T) {
	testCases := []struct {
		min  float64
		max  float64
		mode float64
	}{
		{5.0, 10.0, 7.0},
	}
	for _, tc := range testCases {
		got := NewPERT(tc.min, tc.max, tc.mode)
		if reflect.TypeOf(got).String() != "mcs.PERT" {
			t.Errorf("wrong type. expected pert / got %v", reflect.TypeOf(got))
		}
		if got.min != tc.min {
			t.Errorf("wrong min. expected %v / got %v", tc.min, got.min)
		}
		if got.max != tc.max {
			t.Errorf("wrong max. expected %v / got %v", tc.max, got.max)
		}
		if got.mode != tc.mode {
			t.Errorf("wrong mode. expected %v / got %v", tc.mode, got.mode)
		}
	}
}

func TestPERTUnmarshalJSON(t *testing.T) {
	testCases := []struct {
		data     []byte
		expected PERT
	}{
		{
			[]byte(`{"type":"pert","min":5,"max":10,"mode":7}`),
			PERT{5, 10, 7},
		},
	}
	for _, tc := range testCases {
		var p PERT
		err := json.Unmarshal(tc.data, &p)
		if err != nil {
			t.Errorf("error unmarshaling to pert dist: %v", tc.data)
		}
	}
}
