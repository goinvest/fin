// Copyright (c) 2019-2020 The goinvest/fin developers. All rights reserved.
// Project site: https://github.com/goinvest/fin
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE file for the project.

package fin

import (
	"log"
	"math"
)

// IRR calculates the Internal Rate of Return (IRR), which is the discount rate
// for which the Net Present Value (NPV) equals zero. The IRR assumes that cash
// flows are reinvested at the IRR, which is why the Modified IRR (MIRR) is
// preferred.
//
// NPV = 0 = ∑(CF_t / (1 + IRR)^t) for t=0...n
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

// MIRR calcualtes the Modified Internal Rate of Return (MIRR). Cash outflows
// (COF, negative cashflows), regardless of when they occur, are treated as a
// cost and discounted using the cost of capital (k). Cash inflows (CIF,
// positive cashflows) are treated as part of the terminal value, which are
// compounded using the cost of capital (k).
//
// PV Costs = PV Terminal Value
// ∑(COFt / (1+k)^t) = (∑ [CIF_t * (1+k)^(n-t)]) / (1 + MIRR)^n
func MIRR(cashflows []float64, k float64) float64 {
	pvCosts, tv := 0.0, 0.0
	n := float64(len(cashflows) - 1)
	for i, cf := range cashflows {
		t := float64(i)
		if cf > 0 {
			// Cash inflows (CIF)
			tv += cf * math.Pow(1+k, n-t)
		} else {
			// Cash outflow (COF)
			pvCosts -= cf / math.Pow(1+k, t)
		}
	}
	log.Printf("PV costs = %f / TV = %f", pvCosts, tv)
	return math.Pow(tv/pvCosts, 1/n) - 1
}

// DiscountedPaybackPeriod calculates the expected number of periods required
// to recover the original investment using the given discount rate (k). If the
// investment never pays back, then NaN is returned.
func DiscountedPaybackPeriod(cashflows []float64, k float64) float64 {
	cumulative := 0.0
	for i, cf := range cashflows {
		discountedCF := cf / math.Pow(1+k, float64(i))
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
// discount rate (k). The initial cashflow is not discounted.
//
// NPV = ∑(CF_t / (1+k)^t) for t=0...n
func NPV(cashflows []float64, k float64) float64 {
	npv := 0.0
	for i, cf := range cashflows {
		npv += cf / math.Pow(1+k, float64(i))
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
