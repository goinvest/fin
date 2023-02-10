// Copyright (c) 2019-2023 The goinvest/fin developers. All rights reserved.
// Project site: https://github.com/goinvest/fin
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE file for the project.

package mcs

import (
	"github.com/goinvest/distuvx"
	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
)

// Fixed is a fixed number.
type Fixed float64

// Randomize a new fixed number.
func (f Fixed) Randomize(src rand.Source) Rander {
	return distuvx.NewFixed(float64(f))
}

// PERT setups up a new PERT distribution.
type PERT []float64

// Randomize sets up a new PERT distribution.
func (p PERT) Randomize(src rand.Source) Rander {
	if len(p) != 3 {
		panic("wrong number of PERT arguments")
	}
	return distuvx.NewPERT(p[0], p[2], p[1], src)
}

// Triangle is a triangle distribution with the values min, mode, max.
type Triangle []float64

// Randomize sets up a new triangle distribution.
func (t Triangle) Randomize(src rand.Source) Rander {
	if len(t) != 3 {
		return distuv.NewTriangle(1, 0, 0, src)
	}
	return distuv.NewTriangle(t[0], t[2], t[1], src)
}

// TriangleOne use a triangle distribution once and then returns the same number
// each time.
type TriangleOne []float64

// Randomize sets up a new one-time only triangle distribution.
func (t TriangleOne) Randomize(src rand.Source) Rander {
	if len(t) != 3 {
		return distuv.NewTriangle(1, 0, 0, src)
	}
	triangle := distuv.NewTriangle(t[0], t[2], t[1], src)
	return distuvx.NewFixed(triangle.Rand())
}

// Uniform is a uniform distribution with the values min and max.
type Uniform []float64

// Randomize sets up a new uniform distribution.
func (u Uniform) Randomize(src rand.Source) Rander {
	if len(u) != 2 {
		panic("wrong number of uniform arguments")
	}
	min, max := u[0], u[1]
	if u[0] > u[1] {
		max, min = min, max
	}
	uni := distuv.Uniform{
		Min: min,
		Max: max,
		Src: src,
	}
	return uni
}
