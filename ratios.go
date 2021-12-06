// Copyright (c) 2019-2022 The goinvest/fin developers. All rights reserved.
// Project site: https://github.com/goinvest/fin
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE file for the project.

package fin

///////////////////////////////////////////////////////////////////////////////
//
// Profitability Ratios
//
///////////////////////////////////////////////////////////////////////////////

// ReturnOnEquity tells what percentage of profit is made for every dollar of
// equity invested in the company.
//
// Source: https://financialintelligencebook.com
func ReturnOnEquity(netProfit, equity float64) float64 {
	return netProfit / equity
}

///////////////////////////////////////////////////////////////////////////////
//
// Leverage Ratios
//
///////////////////////////////////////////////////////////////////////////////

// DebtToAssets measures the ratio of all debt (short-term debt and long-term
// debt, but not other liabilities) to all assets.
func DebtToAssets(debt, assets float64) float64 {
	return debt / assets
}

// DebtToEquity measures the ratio of all debt (short-term debt and long-term
// debt, but not other liabilities) to shareholders' equity.
func DebtToEquity(debt, equity float64) float64 {
	return debt / equity
}

// LiabilitiesToAssets measures the ratio of total liabilities to total assets.
func LiabilitiesToAssets(liabilities, assets float64) float64 {
	return liabilities / assets
}

// LiabilitiesToEquity measures the ratio of total liabilities to total
// shareholders' equity.
func LiabilitiesToEquity(liabilities, equity float64) float64 {
	return liabilities / equity
}

// EquityMultiplier ratio is the factor by which the return on assets is
// multiplied to determine the return on equity. It is defined as the ratio of
// total assets to common equity. [Source: Financial Management: Theory &
// Practice, 16th ed.]
func EquityMultiplier(assets, equity float64) float64 {
	return assets / equity
}

// TimesInterestEarned ratio, also called the interest coverage ratio, is
// determined by dividing earnings before interest and taxes (EBIT) by the
// interest expense. [Source: Financial Management: Theory & Practice, 16th
// ed.]
func TimesInterestEarned(ebit, interest float64) float64 {
	return ebit / interest
}

// InterestCoverage calculates the ratio of many times it can pay its annual
// interest charges.
//
// Source: https://financialintelligencebook.com
func InterestCoverage(operatingProfit, interest float64) float64 {
	return operatingProfit / interest
}

///////////////////////////////////////////////////////////////////////////////
//
// Liquidity Ratios
//
///////////////////////////////////////////////////////////////////////////////

// CurrentRatio measures a copmany's current assets against its current
// liabilities.
//
// Source: https://financialintelligencebook.com
func CurrentRatio(currentAssets, currentLiabilities float64) float64 {
	return currentAssets / currentLiabilities
}

// QuickRatio, also known as the acid test, assesses how quickly a company
// could pay off their short-term debt without waiting to sell inventory.
//
// Source: https://financialintelligencebook.com
func QuickRatio(currentAssets, currentLiabilities, inventory float64) float64 {
	return (currentAssets - inventory) / currentLiabilities
}

///////////////////////////////////////////////////////////////////////////////
//
// Efficiency Ratios
//
///////////////////////////////////////////////////////////////////////////////

// DaysInInventory measures the number of days inventory stays in the system.
//
// Source: https://financialintelligencebook.com
func DaysInInventory(inventory, cogs float64, days int) float64 {
	return inventory / (cogs / float64(days))
}

// DaysInInventory360 measures the number of days inventory stays in the system
// assuming 360 days for the year.
//
// Source: https://financialintelligencebook.com
func DaysInInventory360(inventory, cogs float64) float64 {
	return DaysInInventory(inventory, cogs, 360)
}

// InventoryTurns measures how many times inventory turns over in a year.
//
// Source: https://financialintelligencebook.com
func InventoryTurns(inventory, cogs float64) float64 {
	return cogs / inventory
}

// DaysSalesOutstanding measures the average time it takes to collect the
// cash from sales.
//
// Source: https://financialintelligencebook.com
func DaysSalesOutstanding(receivables, revenues float64, days int) float64 {
	return receivables / (revenues / float64(days))
}

// DaysSalesOutstanding360 measures the average time it takes to collect the
// cash from sales assuming 360 days for the year.
//
// Source: https://financialintelligencebook.com
func DaysSalesOutstanding360(receivables, revenues float64) float64 {
	return DaysSalesOutstanding(receivables, revenues, 360)
}

// DaysPayableOutstanding shows the average number of days it take a company to
// pay its own outstanding invoices.
//
// Source: https://financialintelligencebook.com
func DaysPayableOutstanding(payables, cogs float64, days int) float64 {
	return payables / (cogs / float64(days))
}

// DaysPayableOutstanding360 shows the average number of days it take a company to
// pay its own outstanding invoices assuming 360 days for the year.
//
// Source: https://financialintelligencebook.com
func DaysPayableOutstanding360(payables, cogs float64) float64 {
	return DaysPayableOutstanding(payables, cogs, 360)
}

// PPETurnover calculates the amount of revenue generated per dollar invested
// in Property, Plant, and Equipment (PPE).
//
// Source: https://financialintelligencebook.com
func PPETurnover(revenue, ppe float64) float64 {
	return revenue / ppe
}

// TotalAssetTurnover calculates the amount of revenue generated per dollar
// invested in all assets.
//
// Source: https://financialintelligencebook.com
func TotalAssetTurnover(revenue, assets float64) float64 {
	return revenue / assets
}
