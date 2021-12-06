// Copyright (c) 2019-2022 The goinvest/fin developers. All rights reserved.
// Project site: https://github.com/goinvest/fin
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE file for the project.

package fin

// CurrentRatio measures a copmany's current assets against its current
// liabilities.
//
// Source: https://financialintelligencebook.com
func CurrentRatio(currentAssets, currentLiabilities float64) float64 {
	return currentAssets / currentLiabilities
}

// QuickRatio also known as the acid test, assesses how quickly a company could
// pay off their short-term debt without waiting to sell inventory.
//
// Source: https://financialintelligencebook.com
func QuickRatio(currentAssets, currentLiabilities, inventory float64) float64 {
	return (currentAssets - inventory) / currentLiabilities
}
