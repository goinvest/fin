// Copyright (c) 2019-2023 The goinvest/fin developers. All rights reserved.
// Project site: https://github.com/goinvest/fin
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE file for the project.

package mcs

import (
	"golang.org/x/exp/rand"
)

// rcf models a random cashflow.
type rcf struct {
	outflow        bool
	start          int
	end            int
	dist           Rander
	growth         Rander
	name           string
	lastGrowthRate float64
	applicable     []bool
	grow           []bool
}

// setupRCFs sets up the random cashflows (RCFs).
func setupRCFs(cpu int, src rand.Source, nrcfs []nrcf) ([]rcf, error) {
	var rcfs []rcf
	for _, cashflow := range nrcfs {
		cf, err := newRCF(src, cashflow)
		if err != nil {
			return []rcf{}, err
		}
		applied := 0
		for _, apply := range cf.applicable {
			if apply {
				applied++
			}
		}
		rcfs = append(rcfs, cf)
	}
	return rcfs, nil
}

// newRCF creates a new random cashflow for a given random distribution with
// another random distribution for the growth rate and the given number of
// periods.
func newRCF(src rand.Source, nonrandomCF nrcf) (rcf, error) {

	// Setup the non-random cashflow.
	cf := rcf{
		outflow:        nonrandomCF.outflow,
		start:          nonrandomCF.start,
		end:            nonrandomCF.end,
		dist:           nonrandomCF.dist.Randomize(src),
		name:           nonrandomCF.name,
		lastGrowthRate: 1.0,
		applicable:     nonrandomCF.applicable,
		grow:           nonrandomCF.grow,
	}
	if nonrandomCF.growth == nil {
		return cf, nil
	}
	cf.growth = nonrandomCF.growth.Randomize(src)
	return cf, nil
}

// value returns a random number for the period index (as opposed for the given
// period).
func (cf *rcf) value(i int) float64 {
	// If we're grabbing the value from the first period, reset the growth rate.
	if i == 0 {
		cf.lastGrowthRate = 1.0
	}
	if !cf.applicable[i] {
		return 0.0
	}
	gr := cf.lastGrowthRate
	if cf.grow[i] {
		gr *= (1.0 + cf.growth.Rand())
	}
	cf.lastGrowthRate = gr
	return cf.dist.Rand() * cf.lastGrowthRate
}
