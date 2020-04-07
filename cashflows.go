// Copyright (c) 2019-2020 The goinvest/fin developers. All rights reserved.
// Project site: https://github.com/goinvest/fin
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE file for the project.

package fin

import (
	"math"
)

// IRR calculates the internal rate of return.
func IRR(cashflows []float64) float64 {
	// The IRR is calculated using the Newton-Raphson method.
	// TODO: Should tolerance be a variadic argument to the IRR function?
	const relError = 1e-8
	const maxIterations = 100
	k0, k1 := 1.0, 0.0
	for i := 0; i < maxIterations; i++ {
		k0 = k1
		f, fdk := 0.0, 0.0
		for i, cf := range cashflows {
			t := float64(i)
			// k = discount rate, which is the IRR
			// f = NPV = ∑t=0-n: CF_t * (1+k)^-t
			// fdk = d/dk NPV = ∑t=0-n: -t * CF_t * (1+k)^(-t-1)
			f += cf * math.Pow(1+k0, -t)
			fdk -= t * cf * math.Pow(1+k0, -t-1)
		}
		k1 = k0 - (f / fdk)
		if math.Abs(k1-k0)/k0 < relError {
			return k1
		}
	}
	return math.NaN()
}

// DiscountedPaybackPeriod calculates the expected number of periods required
// to recover the original investment using the given discount rate. If the
// investment never pays back, then NaN is returned.
func DiscountedPaybackPeriod(cashflows []float64, discountRate float64) float64 {
	cumulative := 0.0
	for i, cf := range cashflows {
		discountedCF := cf / math.Pow(1+discountRate, float64(i))
		if cumulative+discountedCF >= 0.0 {
			return float64(i-1) - cumulative/discountedCF
		}
		cumulative += discountedCF
	}
	return math.NaN()
}

// PaybackPeriod calculates the expected number of periods required to recover
// the original investment. If the investment never pays back, then NaN is
// returned.
func PaybackPeriod(cashflows []float64) float64 {
	cumulative := 0.0
	for i, cf := range cashflows {
		if cumulative+cf >= 0.0 {
			return float64(i-1) - cumulative/cf
		}
		cumulative += cf
	}
	return math.NaN()
}

// NPV calculates the Net Present Value (NPV) for the cashflows based on the
// discount rate. The initial cashflow is not discounted.
//
// NPV = ∑t=0-n: CF_t / (1+k)^t
func NPV(cashflows []float64, discountRate float64) float64 {
	npv := 0.0
	for i, cf := range cashflows {
		npv += cf / math.Pow(1+discountRate, float64(i))
	}
	return npv
}

// NCF calcualates the Net Cash Flows (NCF) for the cashflows given per period.
func NCF(cashflows []float64) float64 {
	sum := 0.0
	for _, cf := range cashflows {
		sum += cf
	}
	return sum
}
