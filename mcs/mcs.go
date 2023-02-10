// Copyright (c) 2019-2023 The goinvest/fin developers. All rights reserved.
// Project site: https://github.com/goinvest/fin
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE file for the project.

package mcs

import (
	"log"

	"golang.org/x/exp/rand"
)

// Rander is the interface for the Rand method.
type Rander interface {
	Rand() float64
}

// Randomizer is the interface for the Setup method.
type Randomizer interface {
	Randomize(src rand.Source) Rander
}

// Cashflow models the setup information needed for a cashflow for the
// Monte Carlo simulation.
type Cashflow struct {
	Outflow bool
	Periods string
	Dist    Randomizer
	Growth  Growth
	Name    string
}

// Growth models the setup information for a Growth rate.
type Growth struct {
	Periods string
	Dist    Randomizer
	Name    string
}

// NetCashflows calculates the net cashflows, cash inflows, and cash outflows
// for a given number of simulations, number of periods, cashflow
// distributions, and random source.
func NetCashflows(cfg Config, seed uint64, cpus int) ([]float64, []float64, []float64) {
	sims := cfg.Sims
	start := cfg.StartPeriod
	end := cfg.EndPeriod
	cfs := cfg.Cashflows

	// FIXME(mdr): Should probably return a slice of slices instead of three
	// []float64. Should also include error in return.

	nonrandomCFs, err := setupNRCFs(start, end, cfs)
	if err != nil {
		log.Printf("error setting up the non-random cashflows")
		return []float64{}, []float64{}, []float64{}
	}

	// Start the simulations in a goroutine for each CPU.
	simsPerCPU := calcSimsPerCPU(sims, cpus)
	ch := make(chan inOutflow, cpus)
	for cpu := 0; cpu < cpus; cpu++ {
		// Don't use the same seed for each CPU, but we still want reproducible
		// results if a non-random seed is provided.
		cpuSeed := seed * uint64(cpu*5000000)
		go simulate(cpu, simsPerCPU[cpu], start, end, cpuSeed, ch, nonrandomCFs)
	}

	// Assemble the results
	var netCashflows []float64
	var netOutflows []float64
	var netInflows []float64
	var result inOutflow
	for i := 0; i < sims; i++ {
		result = <-ch
		netCashflows = append(netCashflows, result.in-result.out)
		netInflows = append(netInflows, result.in)
		netOutflows = append(netOutflows, result.out)
	}

	return netCashflows, netInflows, netOutflows
}

type inOutflow struct {
	in  float64
	out float64
}

func simulate(cpu, sims, start, end int, seed uint64, ch chan inOutflow, nrcfs []nrcf) {
	periods := end - start + 1
	// Setup each random cashflow.
	src := rand.New(rand.NewSource(seed))
	randomCFs, err := setupRCFs(cpu, src, nrcfs)
	if err != nil {
		log.Printf("error: %s", err)
	}

	// Loop through each simulation
	for sim := 0; sim < sims; sim++ {
		netInflows := 0.0
		netOutflows := 0.0
		// Loop through each period using the period index not the period number
		// (i.e., the index always starts at 0, whereas the period could start at
		// 0, 1, or a different number).
		for i := 0; i < periods; i++ {
			// Loop through each cashflow summing the outflows and inflows
			// separately. Since the call to value(i) has the side effect of changing
			// the randomCFs[j], we need to iterate using the index instead of using
			// the value from range.
			for j := 0; j < len(randomCFs); j++ {
				val := randomCFs[j].value(i)
				// if i == 0 {
				// 	log.Printf("%s [%d] on CPU %d = %f", randomCFs[j].name, i+start, cpu, val)
				// }
				// if j == 1 {
				// 	log.Printf("%s [%d] on CPU %d = %f", rcf.name, i+start, cpu, val)
				// }
				// if i == 0 && j == 0 {
				// 	log.Printf("%s [%d] on CPU %d = %f", rcf.name, i+start, cpu, val)
				// }
				if randomCFs[j].outflow {
					netOutflows += val
				} else {
					netInflows += val
				}
			}
		}
		// Send the inflow and outflow for this simulation to the channel.
		ch <- inOutflow{
			in:  netInflows,
			out: netOutflows,
		}
	}
}

func calcSimsPerCPU(sims, cpus int) []int {
	minSimsPerCPU := sims / cpus
	leftovers := sims - cpus*minSimsPerCPU
	simsPerCPU := make([]int, cpus)
	for i := 0; i < cpus; i++ {
		if i < leftovers {
			simsPerCPU[i] = minSimsPerCPU + 1
		} else {
			simsPerCPU[i] = minSimsPerCPU
		}
	}
	return simsPerCPU
}
