// Copyright (c) 2019-2020 The goinvest/fin developers. All rights reserved.
// Project site: https://github.com/goinvest/fin
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE file for the project.

package mcs

import (
	"log"

	"github.com/goinvest/seq"
)

// nrcf models a non-random cashflow.
type nrcf struct {
	outflow    bool
	start      int
	end        int
	dist       Randomizer
	growth     Randomizer
	name       string
	applicable []bool
	grow       []bool
}

// setupNRCFs sets up the non-random cashflows (NRCFs).
func setupNRCFs(start, end int, cashflows []Cashflow) ([]nrcf, error) {
	var nonrandomCFs []nrcf
	for _, cashflow := range cashflows {
		cf, err := newNRCF(start, end, cashflow)
		if err != nil {
			log.Printf("error creating %s cashflow", cashflow.Name)
			return []nrcf{}, err
		}
		applied := 0
		for _, apply := range cf.applicable {
			if apply {
				applied++
			}
		}
		nonrandomCFs = append(nonrandomCFs, cf)
	}
	log.Printf("Done setting up %d non-random cashflows", len(nonrandomCFs))
	return nonrandomCFs, nil
}

// newNRCF creates a new non-random cashflow for a given random distribution
// with another random distribution for the growth rate and the given number of
// periods.
func newNRCF(start, end int, cashflow Cashflow) (nrcf, error) {
	cf := nrcf{
		outflow: cashflow.Outflow,
		start:   start,
		end:     end,
		dist:    cashflow.Dist,
		name:    cashflow.Name,
	}

	// Setup the applicable periods both for the distribution and growth.
	applicable, err := parsePeriods(start, end, cashflow.Periods)
	if err != nil {
		log.Printf("error parsing applicable periods for %s: %s", cashflow.Name, cashflow.Periods)
		return nrcf{}, err
	}
	cf.applicable = applicable

	// If there is no growth, set all grow periods to false.
	if cashflow.Growth == (Growth{}) {
		grow := make([]bool, len(applicable))
		cf.grow = grow
		// Return without growth.
		return cf, nil
	}
	grow, err := parsePeriods(start, end, cashflow.Growth.Periods)
	if err != nil {
		return cf, err
	}
	cf.grow = grow
	cf.growth = cashflow.Growth.Dist
	return cf, nil
}

func parsePeriods(start, end int, s string) ([]bool, error) {
	numPeriods := end - start + 1
	a := make([]bool, numPeriods)

	// Parse the applicable period string into a slice of ints.
	periods, err := seq.Parse(s)
	if err != nil {
		return nil, err
	}

	// Setup the applicable periods.
	for _, period := range periods {
		a[period-start] = true
	}
	return a, nil
}
