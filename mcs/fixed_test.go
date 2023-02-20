// Copyright (c) 2023 The goinvest/fin developers. All rights reserved.
// Project site: https://github.com/goinvest/fin
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE file for the project.

package mcs

import (
	"encoding/json"
	"testing"
)

func TestNewFixed(t *testing.T) {
	testCases := []struct {
		value    float64
		expected Fixed
	}{
		{5.0, Fixed(5.0)},
		{25.0, Fixed(25.0)},
	}
	for _, tc := range testCases {
		f := NewFixed(tc.value)
		if f != tc.expected {
			t.Errorf("error creating new fixed dist: got %v / expected %v", f, tc.expected)
		}
	}
}

func TestFixedUnmarshalJSON(t *testing.T) {
	testCases := []struct {
		data     []byte
		expected Fixed
	}{
		{
			[]byte(`{"type":"fixed","val":5}`),
			Fixed(5),
		},
		{
			[]byte(`{"type":"fixed","val":25.0}`),
			Fixed(25.0),
		},
	}
	for _, tc := range testCases {
		var f Fixed
		err := json.Unmarshal(tc.data, &f)
		if err != nil {
			t.Errorf("error unmarshaling to fixed dist: %v", tc.data)
		}
		if f != tc.expected {
			t.Errorf("error unmarshaling fixed dist: got %v / expected %v", f, tc.expected)
		}
	}
}
