// Copyright (c) 2019-2023 The goinvest/fin developers. All rights reserved.
// Project site: https://github.com/goinvest/fin
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE file for the project.

package main

import (
	"log"

	"github.com/goinvest/fin"
)

func main() {

	cashflows := []float64{-1000, 500, 400, 300, 100}
	irr := fin.IRR(cashflows)
	log.Printf("IRR = %f", irr)
}
