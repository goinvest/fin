// Copyright (c) 2019-2025 The goinvest/fin developers. All rights reserved.
// Project site: https://github.com/goinvest/fin
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE file for the project.

package mcs

import (
	"encoding/json"
	"reflect"
	"testing"

	"golang.org/x/exp/rand"
)

func TestNewTriangle(t *testing.T) {
	testCases := []struct {
		min  float64
		max  float64
		mode float64
	}{
		{5.0, 10.0, 7.0},
	}
	for _, tc := range testCases {
		got := NewTriangle(tc.min, tc.max, tc.mode)
		if reflect.TypeOf(got).String() != "mcs.Triangle" {
			t.Errorf("wrong type. expected mcs.Triangle. / got %v", reflect.TypeOf(got))
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

func TestTriangleUnmarshalJSON(t *testing.T) {
	testCases := []struct {
		data []byte
		min  float64
		max  float64
		mode float64
	}{
		{[]byte(`{"type":"tri","min":5,"max":10,"mode":7}`), 5.0, 10.0, 7.0},
	}
	for _, tc := range testCases {
		var dist Triangle
		err := json.Unmarshal(tc.data, &dist)
		if err != nil {
			t.Errorf("error unmarshaling to triangle dist: %v", tc.data)
		}
		if reflect.TypeOf(dist).String() != "mcs.Triangle" {
			t.Errorf("wrong type. expected mcs.Triangle. / got %v", reflect.TypeOf(dist))
		}
		if dist.min != tc.min {
			t.Errorf("wrong min. expected %v / got %v", tc.min, dist.min)
		}
		if dist.max != tc.max {
			t.Errorf("wrong max. expected %v / got %v", tc.max, dist.max)
		}
		if dist.mode != tc.mode {
			t.Errorf("wrong mode. expected %v / got %v", tc.mode, dist.mode)
		}
	}
}

func TestTriangleOneUnmarshalJSON(t *testing.T) {
	testCases := []struct {
		data []byte
		min  float64
		max  float64
		mode float64
	}{
		{[]byte(`{"type":"tri_one","min":5,"max":10,"mode":7}`), 5.0, 10.0, 7.0},
	}
	for _, tc := range testCases {
		var dist TriangleOne
		err := json.Unmarshal(tc.data, &dist)
		if err != nil {
			t.Errorf("error unmarshaling to triangleOne dist: %v", tc.data)
		}
		if reflect.TypeOf(dist).String() != "mcs.TriangleOne" {
			t.Errorf("wrong type. expected mcs.TriangleOne. / got %v", reflect.TypeOf(dist))
		}
		if dist.min != tc.min {
			t.Errorf("wrong min. expected %v / got %v", tc.min, dist.min)
		}
		if dist.max != tc.max {
			t.Errorf("wrong max. expected %v / got %v", tc.max, dist.max)
		}
		if dist.mode != tc.mode {
			t.Errorf("wrong mode. expected %v / got %v", tc.mode, dist.mode)
		}
		// Since this is a triangle one distribution, the same value should always
		// be returned.
		r := dist.Randomize(rand.NewSource(1234))
		value1 := r.Rand()
		value2 := r.Rand()
		if value1 != value2 {
			t.Errorf("expected same value each time. got %v and then got %v", value1, value2)
		}
	}
}
