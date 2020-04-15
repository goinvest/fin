// Copyright (c) 2019-2020 The goinvest/fin developers. All rights reserved.
// Project site: https://github.com/goinvest/fin
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE file for the project.

package mcs

import (
	"log"

	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"

	"github.com/goinvest/distuvx"
	"github.com/goinvest/seq"
)

// Rander is the interface for the Rand method.
type Rander interface {
	Rand() float64
}

// Valuer provides a value given a period.
type Valuer interface {
	Value(period int) float64
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

// cashflow models a cash flow.
type cashflow struct {
	outflow           bool
	start             int
	end               int
	value             Rander
	growthRate        Rander
	Name              string
	applicablePeriods []float64
	growthPeriods     []int
	growthRates       []float64
}

// setupCFs sets up the cashflows.
func setupCFs(start, end int, seed uint64, setups []Cashflow) ([]cashflow, error) {
	periods := end - start + 1
	cashflows := make([]cashflow, periods)
	for _, setup := range setups {
		log.Printf("%s CF / Period string: %s", setup.Name, setup.Periods)
		cf, err := newCF(start, end, seed, setup)
		if err != nil {
			log.Fatalf("error creating %s cash flow: %s", setup.Name, err)
		}
		cashflows = append(cashflows, cf)
	}
	return cashflows, nil
}

// newCF creates a new cash flow for a given random distribution with another
// random distribution for the growth rate and the given number of periods.
func newCF(start, end int, seed uint64, cfg Cashflow) (cashflow, error) {
	// Seed the random numbers
	src := rand.New(rand.NewSource(seed))

	// Setup the applicable periods.
	applicablePeriods, err := calcPeriods(start, end, cfg.Periods)
	if err != nil {
		return cashflow{}, err
	}

	// Setup the growth rates if applicable.
	growthPeriods, err := seq.Parse(cfg.Growth.Periods)
	if err != nil {
		return cashflow{}, err
	}

	var growthRates []float64
	if cfg.Growth == (Growth{}) {
		growthRates, err = calcGrowthPeriods(start, end, src, cfg.Growth)
		if err != nil {
			return cashflow{}, err
		}
	}

	cf := cashflow{
		outflow:           cfg.Outflow,
		start:             start,
		end:               end,
		value:             cfg.Dist.Randomize(src),
		growthRate:        cfg.Growth.Dist.Randomize(src),
		Name:              cfg.Name,
		applicablePeriods: applicablePeriods,
		growthPeriods:     growthPeriods,
		growthRates:       growthRates,
	}
	return cf, nil
}

// RandomizeGrowthRates randomizes the growth rates.
func (cf *cashflow) RandomizeGrowthRates() {
	numPeriods := cf.end - cf.start + 1
	lastGrowthRate := 1.0
	for i := 0; i < numPeriods; i++ {
		for _, growthPeriod := range cf.growthPeriods {
			if cf.start+i == growthPeriod {
				lastGrowthRate *= 1 + cf.growthRate.Rand()
			}
		}
		cf.growthRates[i] = lastGrowthRate
	}
}

// Value returns a random number for the given period.
func (cf *cashflow) Value(period int) float64 {
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
func NetCashflows(sims, cpus, start, end int, seed uint64, cfs []Cashflow) ([]float64, []float64, []float64) {

	// FIXME: Need to setup each cashflow except for the rand.Source, since those
	// need to be created within each goroutine.

	// Start the simulations in a goroutine for each CPU.
	simsPerCPU := calcSimsPerCPU(sims, cpus)
	ch := make(chan inOutflow, cpus)
	for cpu := 0; cpu < cpus; cpu++ {
		cpuSeed := seed + uint64(cpu*100)
		go sim(simsPerCPU[cpu], start, end, cpuSeed, ch, cfs)
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

type inOutflow struct {
	in  float64
	out float64
}

func sim(sims, start, end int, seed uint64, ch chan inOutflow, setups []Cashflow) {
	// Setup each cashflow.
	cfs, err := setupCFs(start, end, seed, setups)
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
		// Loop through each period
		for period := start; period <= end; period++ {
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

// Randomize sets up a new triangle distribution.
func (t Triangle) Randomize(src rand.Source) Rander {
	if len(t) != 3 {
		return distuv.NewTriangle(1, 0, 0, src)
	}
	return distuv.NewTriangle(t[0], t[2], t[1], src)
}

// Fixed is a fixed number.
type Fixed float64

// Randomize a new fixed number.
func (f Fixed) Randomize(src rand.Source) Rander {
	return distuvx.NewFixed(float64(f))
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

// PERT setups up a new PERT distribution.
type PERT []float64

// Randomize sets up a new triangle distribution.
func (p PERT) Randomize(src rand.Source) Rander {
	if len(p) != 3 {
		panic("wrong number of PERT arguments")
	}
	return distuvx.NewPERT(p[0], p[2], p[1], src)
}

func calcPeriods(start, end int, s string) ([]float64, error) {
	numPeriods := end - start + 1
	applicablePeriods := make([]float64, numPeriods)

	// Parse the applicable periods for the cashflow.
	periods, err := seq.Parse(s)
	if err != nil {
		return nil, err
	}

	// Setup the applicable periods for this cashflow.
	for _, period := range periods {
		applicablePeriods[period-start] = 1.0
	}
	return applicablePeriods, nil
}

func calcGrowthPeriods(start, end int, src rand.Source, g Growth) ([]float64, error) {
	numPeriods := end - start + 1
	growthRates := make([]float64, numPeriods)

	// Parse the applicable periods for the growth.
	periods, err := seq.Parse(g.Periods)
	if err != nil {
		return nil, err
	}

	gr := g.Dist.Randomize(src)
	lastGrowthRate := 1.0
	for i := 0; i < numPeriods; i++ {
		for _, period := range periods {
			if start+i == period {
				lastGrowthRate *= 1 + gr.Rand()
			}
		}
		growthRates[i] = lastGrowthRate
	}
	return growthRates, nil
}
