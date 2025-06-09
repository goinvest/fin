// Copyright (c) 2019-2025 The goinvest/fin developers. All rights reserved.
// Project site: https://github.com/goinvest/fin
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE file for the project.

package cf

import (
	"fmt"
	"math"
)

type ToleranceType int

const (
	Relative ToleranceType = 1
	Absolute ToleranceType = 2
)

type IRROptions struct {
	InitialGuess  float64
	Tolerance     float64
	ToleranceType ToleranceType
	MaxIterations int
}

// IRR calculates the Internal Rate of Return (IRR), which is the discount rate
// for which the Net Present Value (NPV) equals zero. The IRR assumes that cash
// flows are reinvested at the IRR, which is why the Modified IRR (MIRR) is
// preferred. The IRR is calculated using Newton's Method (i.e.,
// Newton-Raphson Method) as follows:
//
// NPV = 0 = ∑(CF_n / (1 + IRR)^n) for n=0...N
//
// If the IRR function is called without the optional struct, the defaults will
// be initialGuess = 0.1, tolerance = 1e-8, toleranceType = Absolute, and
// maxIterations = 100.
func IRR(cashflows []float64, opts ...IRROptions) (float64, error) {

	if len(cashflows) < 2 {
		return math.NaN(), fmt.Errorf("need at least two cash flows")
	}

	// Default options.
	initialGuess := 0.1
	tolerance := 1e-8
	toleranceType := Absolute
	maxIterations := 100

	// Override default options if provided.
	if len(opts) > 0 && opts[0].InitialGuess != 0.0 {
		initialGuess = opts[0].InitialGuess
	}
	if len(opts) > 0 && opts[0].Tolerance != 0.0 {
		tolerance = opts[0].Tolerance
	}
	if len(opts) > 0 && opts[0].ToleranceType != 0 {
		toleranceType = opts[0].ToleranceType
	}
	if len(opts) > 0 && opts[0].MaxIterations != 0 {
		maxIterations = opts[0].MaxIterations
	}

	// Calculate the IRR using the Newton-Raphson method.
	rate := initialGuess
	for i := 0; i < maxIterations; i++ {
		f, fdk := 0.0, 0.0
		for i, cf := range cashflows {
			n := float64(i)
			// rate = k = discount rate, which is the IRR
			// f (function) = NPV = ∑n=0-N: CF_n / (1+rate)^n
			// fdk (derivative) = d/dk NPV = ∑n=0-N: -n * CF_n / (1+rate)^(n+1)
			f += cf / math.Pow(1+rate, n)
			fdk -= n * cf / math.Pow(1+rate, n+1)
		}
		if math.Abs(fdk) < 1e-12 {
			return math.NaN(), fmt.Errorf("derivative too close to zero, cannot converge")
		}
		newRate := rate - (f / fdk)

		if toleranceType == Relative && math.Abs(newRate-rate)/rate < tolerance {
			return newRate, nil
		} else if toleranceType == Absolute && math.Abs(newRate-rate) < tolerance {
			return newRate, nil
		}
		rate = newRate
	}

	return math.NaN(), fmt.Errorf("failed to converge after %d iterations", maxIterations)
}

// MIRR calculates the Modified Internal Rate of Return (MIRR), which is the
// discount rate at which the present value of the cash outflows equals the
// discounted future value of cash inflows—the discounted terminal value. Cash
// outflows (negative cashflows), regardless of when they occur, are treated as
// a cost and discounted using the cost of capital (k) to calculate the present
// value. Cash inflows (positive cashflows) are reinvested at the cost of
// capital (k), so cash inflows are compounded using the cost of capital (k) to
// calculate the terminal value.
//
// MIRR = [Future Value Cash Inflows / Present Value Cash Outflows]^(1/n) - 1
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
