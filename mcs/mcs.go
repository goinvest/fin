// Copyright (c) 2019-2020 The goinvest/fin developers. All rights reserved.
// Project site: https://github.com/goinvest/fin
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE file for the project.

package mcs

import (
	"fmt"
	"log"

	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"

	"github.com/goinvest/distuvx"
)

// Rander is the interface for the Rand method.
type Rander interface {
	Rand() float64
}

// Valuer provides a value given a period.
type Valuer interface {
	Value(period int) float64
}

// Setuper is the interface for the Setup method.
type Setuper interface {
	Setup(src rand.Source) Rander
}

// CashflowSetup models the setup information needed for a cashflow for the
// Monte Carlo simulation.
type CashflowSetup struct {
	Outflow     bool
	StartPeriod int
	EndPeriod   int
	Random      Setuper
	Growth      Setuper
	Name        string
}

// Cashflow models a cash flow.
type Cashflow struct {
	outflow           bool
	value             Rander
	annualGrowthRate  Rander
	Name              string
	applicablePeriods []float64
	growthRates       []float64
}

// SetupCFs sets up the cashflows.
func SetupCFs(setups []CashflowSetup, seed uint64) ([]Cashflow, error) {
	// FIXME(mdr): I should probably search for the maximum period and make sure
	// that each cashflow is the same length.
	var cashflows []Cashflow
	for _, setup := range setups {
		cf, err := NewCF(setup, seed)
		if err != nil {
			log.Fatalf("Error creating cash flow: %s", err)
		}
		cashflows = append(cashflows, cf)
	}
	return cashflows, nil
}

// NewCF creates a new cash flow for a given random distribution with another
// random distribution for the annual growth rate and the given number of
// months. The start and end months are one (1) based. The months are the total
// number of months in the entire simulation.
func NewCF(cfg CashflowSetup, seed uint64) (Cashflow, error) {
	src := rand.New(rand.NewSource(seed))
	if cfg.StartPeriod > cfg.EndPeriod {
		return Cashflow{}, fmt.Errorf(
			"start period %d must be less than or equal to end period %d", cfg.StartPeriod, cfg.EndPeriod)
	}
	applicablePeriods := make([]float64, cfg.EndPeriod)
	gr := cfg.Growth.Setup(src)
	growthRates := make([]float64, cfg.EndPeriod)
	growthRates[0] = 1.0
	for i := 1; i < len(growthRates); i++ {
		if i%12 == 0 {
			growthRates[i] = growthRates[i-1] * (1 + gr.Rand())
		} else {
			growthRates[i] = growthRates[i-1]
		}
	}
	for i := 0; i < len(applicablePeriods); i++ {
		applicablePeriods[i] = 0.0
		if i >= cfg.StartPeriod-1 && i < cfg.EndPeriod {
			applicablePeriods[i] = 1.0
		}
	}
	// log.Printf("%s applicable months = %v", name, applicablePeriods)
	cf := Cashflow{
		outflow:           cfg.Outflow,
		value:             cfg.Random.Setup(src),
		annualGrowthRate:  gr,
		applicablePeriods: applicablePeriods,
		growthRates:       growthRates,
		Name:              cfg.Name,
	}
	return cf, nil
}

// RandomizeGrowthRates randomizes the growth rates.
func (cf *Cashflow) RandomizeGrowthRates() {
	cf.growthRates[0] = 1.0
	for i := 1; i < len(cf.growthRates); i++ {
		if i%12 == 0 {
			cf.growthRates[i] = cf.growthRates[i-1] * (1 + cf.annualGrowthRate.Rand())
		} else {
			cf.growthRates[i] = cf.growthRates[i-1]
		}
	}
}

// Value returns a random number for the given period.
func (cf *Cashflow) Value(period int) float64 {
	applicable := 0.0
	if period > 0 && period <= len(cf.applicablePeriods) {
		applicable = cf.applicablePeriods[period-1]
	}
	growthRate := 1.0
	if period > 0 && period <= len(cf.growthRates) {
		growthRate = cf.growthRates[period-1]
	}
	return cf.value.Rand() * applicable * growthRate
}

// NetCashflows calculates the net cashflows, cash inflows, and cash outflows
// for a given number of simulations, number of periods, cashflow
// distributions, and random source.
func NetCashflows(sims, cpus, periods int, seed uint64, setups []CashflowSetup) ([]float64, []float64, []float64) {
	simsPerCPU := sims / cpus
	leftovers := sims - cpus*simsPerCPU
	cpuSims := make([]int, cpus)
	for i := 0; i < cpus; i++ {
		if i < leftovers {
			cpuSims[i] = simsPerCPU + 1
		} else {
			cpuSims[i] = simsPerCPU
		}
	}

	ch := make(chan inOutflow, cpus)

	// Start the simulation in a goroutine for each CPU.
	for cpu := 0; cpu < cpus; cpu++ {
		cpuSeed := seed + uint64(cpu*100)
		go sim(cpuSims[cpu], periods, cpuSeed, ch, setups)
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

func sim(sims, periods int, seed uint64, ch chan inOutflow, setups []CashflowSetup) {
	cfs, err := SetupCFs(setups, seed)
	if err != nil {
		log.Printf("error: %s", err)
	}
	// Loop through each simulation
	for sim := 0; sim < sims; sim++ {
		// Reset all growth rates
		for _, cf := range cfs {
			cf.RandomizeGrowthRates()
		}
		netCashflow := 0.0
		netInflow := 0.0
		netOutflow := 0.0
		// Loop through each month
		for period := 1; period <= periods; period++ {
			// Sum each cash flow.
			periodInflows := 0.0
			periodOutlfows := 0.0
			for _, cf := range cfs {
				val := cf.Value(period)
				if cf.outflow {
					periodOutlfows += val
				} else {
					periodInflows += val
				}
			}
			netInflow += periodInflows
			netOutflow += periodOutlfows
			netCashflow += periodInflows - periodOutlfows
		}
		ch <- inOutflow{
			in:  netInflow,
			out: netOutflow,
		}
	}
}

// NoGrowth is used to model no growth of a distribution.
var NoGrowth = Fixed(0.0)

// Triangle is a triangle distribution with the values min, mode, max.
type Triangle []float64

// Setup sets up a new triangle distribution.
func (t Triangle) Setup(src rand.Source) Rander {
	if len(t) != 3 {
		return distuv.NewTriangle(1, 0, 0, src)
	}
	return distuv.NewTriangle(t[0], t[2], t[1], src)
}

// Fixed is a fixed number.
type Fixed float64

// Setup a new fixed number.
func (f Fixed) Setup(src rand.Source) Rander {
	return distuvx.NewFixed(float64(f))
}

// TriangleOne use a triangle distribution once and then returns the same number
// each time.
type TriangleOne []float64

// Setup sets up a new one-time only triangle distribution.
func (t TriangleOne) Setup(src rand.Source) Rander {
	if len(t) != 3 {
		return distuv.NewTriangle(1, 0, 0, src)
	}
	triangle := distuv.NewTriangle(t[0], t[2], t[1], src)
	return distuvx.NewFixed(triangle.Rand())
}

// Uniform is a uniform distribution with the values min and max.
type Uniform []float64

// Setup sets up a new uniform distribution.
func (u Uniform) Setup(src rand.Source) Rander {
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

// PERT setups up a new PERT distribution.
type PERT []float64

// Setup sets up a new triangle distribution.
func (p PERT) Setup(src rand.Source) Rander {
	if len(p) != 3 {
		panic("wrong number of PERT arguments")
	}
	return distuvx.NewPERT(p[0], p[2], p[1], src)
}
