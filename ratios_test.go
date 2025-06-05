// Copyright (c) 2019-2025 The goinvest/fin developers. All rights reserved.
// Project site: https://github.com/goinvest/fin
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE file for the project.

package fin

import (
	"fmt"
	"math"
	"testing"
)

func TestReturnOnEquity(t *testing.T) {
	testCases := []struct {
		netProfit float64
		equity    float64
		want      float64
	}{
		{10.0, 2.0, 5.0},
		{250.0, 125.0, 2.0},
	}
	for i, test := range testCases {
		name := fmt.Sprintf("return_on_equity_%d", i)
		t.Run(name, func(t *testing.T) {
			got := ReturnOnEquity(test.netProfit, test.equity)
			assertFloat64(t, name, got, test.want, 0.0001)
		})
	}
}

func TestDebtToAssets(t *testing.T) {
	testCases := []struct {
		debt   float64
		assets float64
		want   float64
	}{
		{10.0, 2.0, 5.0},
		{250.0, 125.0, 2.0},
	}
	for i, test := range testCases {
		name := fmt.Sprintf("debt_to_assets_%d", i)
		t.Run(name, func(t *testing.T) {
			got := DebtToAssets(test.debt, test.assets)
			assertFloat64(t, name, got, test.want, 0.0001)
		})
	}
}

func TestDebtToEquity(t *testing.T) {
	testCases := []struct {
		debt   float64
		equity float64
		want   float64
	}{
		{10.0, 2.0, 5.0},
		{250.0, 125.0, 2.0},
	}
	for i, test := range testCases {
		name := fmt.Sprintf("debt_to_equity_%d", i)
		t.Run(name, func(t *testing.T) {
			got := DebtToEquity(test.debt, test.equity)
			assertFloat64(t, name, got, test.want, 0.0001)
		})
	}
}

func TestLiabilitiesToAssets(t *testing.T) {
	testCases := []struct {
		liabilities float64
		assets      float64
		want        float64
	}{
		{10.0, 2.0, 5.0},
		{250.0, 125.0, 2.0},
	}
	for i, test := range testCases {
		name := fmt.Sprintf("liabilities_to_assets_%d", i)
		t.Run(name, func(t *testing.T) {
			got := LiabilitiesToAssets(test.liabilities, test.assets)
			assertFloat64(t, name, got, test.want, 0.0001)
		})
	}
}

func TestLiabilitiesToEquity(t *testing.T) {
	testCases := []struct {
		liabilities float64
		equity      float64
		want        float64
	}{
		{10.0, 2.0, 5.0},
		{250.0, 125.0, 2.0},
	}
	for i, test := range testCases {
		name := fmt.Sprintf("liabilities_to_equity_%d", i)
		t.Run(name, func(t *testing.T) {
			got := LiabilitiesToEquity(test.liabilities, test.equity)
			assertFloat64(t, name, got, test.want, 0.0001)
		})
	}
}

func TestEquityMultiplier(t *testing.T) {
	testCases := []struct {
		assets float64
		equity float64
		want   float64
	}{
		{10.0, 2.0, 5.0},
		{250.0, 125.0, 2.0},
	}
	for i, test := range testCases {
		name := fmt.Sprintf("equity_multiplier_%d", i)
		t.Run(name, func(t *testing.T) {
			got := EquityMultiplier(test.assets, test.equity)
			assertFloat64(t, name, got, test.want, 0.0001)
		})
	}
}

func TestTimesInterestEarned(t *testing.T) {
	testCases := []struct {
		ebit     float64
		interest float64
		want     float64
	}{
		{10.0, 2.0, 5.0},
		{250.0, 125.0, 2.0},
	}
	for i, test := range testCases {
		name := fmt.Sprintf("time_interest_earned_%d", i)
		t.Run(name, func(t *testing.T) {
			got := TimesInterestEarned(test.ebit, test.interest)
			assertFloat64(t, name, got, test.want, 0.0001)
		})
	}
}

func TestInterestCoverage(t *testing.T) {
	testCases := []struct {
		operatingProfit float64
		interest        float64
		want            float64
	}{
		{10.0, 2.0, 5.0},
		{250.0, 125.0, 2.0},
	}
	for i, test := range testCases {
		name := fmt.Sprintf("interest_coverage_%d", i)
		t.Run(name, func(t *testing.T) {
			got := InterestCoverage(test.operatingProfit, test.interest)
			assertFloat64(t, name, got, test.want, 0.0001)
		})
	}
}

func TestCurrentRatio(t *testing.T) {
	testCases := []struct {
		currentAssets      float64
		currentLiabilities float64
		want               float64
	}{
		{10.0, 2.0, 5.0},
		{250.0, 125.0, 2.0},
	}
	for i, test := range testCases {
		name := fmt.Sprintf("current_ratio_%d", i)
		t.Run(name, func(t *testing.T) {
			got := CurrentRatio(test.currentAssets, test.currentLiabilities)
			assertFloat64(t, name, got, test.want, 0.0001)
		})
	}
}

func TestQuickRatio(t *testing.T) {
	testCases := []struct {
		currentAssets      float64
		currentLiabilities float64
		inventory          float64
		want               float64
	}{
		{10.0, 2.0, 4.0, 3.0},
		{250.0, 100.0, 50.0, 2.0},
	}
	for i, test := range testCases {
		name := fmt.Sprintf("quick_ratio_%d", i)
		t.Run(name, func(t *testing.T) {
			got := QuickRatio(test.currentAssets, test.currentLiabilities, test.inventory)
			assertFloat64(t, name, got, test.want, 0.0001)
		})
	}
}

func TestDaysInInventory(t *testing.T) {
	testCases := []struct {
		inventory float64
		cogs      float64
		days      int
		want      float64
	}{
		{64.0, 720.0, 360, 32.0},
		{125.0, 250.0, 360, 180.0},
	}
	for i, test := range testCases {
		name := fmt.Sprintf("dii_%d", i)
		t.Run(name, func(t *testing.T) {
			got := DaysInInventory(test.inventory, test.cogs, test.days)
			assertFloat64(t, name, got, test.want, 0.0001)
		})
	}
}

func TestDaysInInventory360(t *testing.T) {
	testCases := []struct {
		inventory float64
		cogs      float64
		want      float64
	}{
		{64.0, 720.0, 32.0},
		{125.0, 250.0, 180.0},
	}
	for i, test := range testCases {
		name := fmt.Sprintf("dii360_%d", i)
		t.Run(name, func(t *testing.T) {
			got := DaysInInventory360(test.inventory, test.cogs)
			assertFloat64(t, name, got, test.want, 0.0001)
		})
	}
}

func TestInventoryTurns(t *testing.T) {
	testCases := []struct {
		inventory float64
		cogs      float64
		want      float64
	}{
		{64.0, 720.0, 11.25},
		{125.0, 250.0, 2.0},
	}
	for i, test := range testCases {
		name := fmt.Sprintf("inventory_turns_%d", i)
		t.Run(name, func(t *testing.T) {
			got := InventoryTurns(test.inventory, test.cogs)
			assertFloat64(t, name, got, test.want, 0.0001)
		})
	}
}

func TestDaysSalesOutstanding(t *testing.T) {
	testCases := []struct {
		receivables float64
		revenues    float64
		days        int
		want        float64
	}{
		{10.0, 200.0, 360, 18.0},
		{250.0, 500.0, 360, 180.0},
	}
	for i, test := range testCases {
		name := fmt.Sprintf("days_sales_outstanding_%d", i)
		t.Run(name, func(t *testing.T) {
			got := DaysSalesOutstanding(test.receivables, test.revenues, test.days)
			assertFloat64(t, name, got, test.want, 0.0001)
		})
	}
}

func TestDaysSalesOutstanding360(t *testing.T) {
	testCases := []struct {
		receivables float64
		revenues    float64
		want        float64
	}{
		{10.0, 200.0, 18.0},
		{250.0, 500.0, 180.0},
	}
	for i, test := range testCases {
		name := fmt.Sprintf("days_sales_outstanding_360_%d", i)
		t.Run(name, func(t *testing.T) {
			got := DaysSalesOutstanding360(test.receivables, test.revenues)
			assertFloat64(t, name, got, test.want, 0.0001)
		})
	}
}

func TestDaysPayableOutstanding(t *testing.T) {
	testCases := []struct {
		payables float64
		cogs     float64
		days     int
		want     float64
	}{
		{10.0, 20.0, 360, 180.0},
		{250.0, 1000.0, 360, 90.0},
	}
	for i, test := range testCases {
		name := fmt.Sprintf("days_payable_outstanding_%d", i)
		t.Run(name, func(t *testing.T) {
			got := DaysPayableOutstanding(test.payables, test.cogs, test.days)
			assertFloat64(t, name, got, test.want, 0.0001)
		})
	}
}

func TestDaysPayableOutstanding360(t *testing.T) {
	testCases := []struct {
		payables float64
		cogs     float64
		want     float64
	}{
		{10.0, 20.0, 180.0},
		{250.0, 1000.0, 90.0},
	}
	for i, test := range testCases {
		name := fmt.Sprintf("days_payable_outstanding_360_%d", i)
		t.Run(name, func(t *testing.T) {
			got := DaysPayableOutstanding360(test.payables, test.cogs)
			assertFloat64(t, name, got, test.want, 0.0001)
		})
	}
}

func TestPPETurnover(t *testing.T) {
	testCases := []struct {
		revenue float64
		ppe     float64
		want    float64
	}{
		{10.0, 2.0, 5.0},
		{250.0, 125.0, 2.0},
	}
	for i, test := range testCases {
		name := fmt.Sprintf("ppe_turnover_%d", i)
		t.Run(name, func(t *testing.T) {
			got := PPETurnover(test.revenue, test.ppe)
			assertFloat64(t, name, got, test.want, 0.0001)
		})
	}
}

func TestTotalAssetTurnover(t *testing.T) {
	testCases := []struct {
		revenue float64
		assets  float64
		want    float64
	}{
		{10.0, 2.0, 5.0},
		{250.0, 125.0, 2.0},
	}
	for i, test := range testCases {
		name := fmt.Sprintf("total_asset_turnover_%d", i)
		t.Run(name, func(t *testing.T) {
			got := TotalAssetTurnover(test.revenue, test.assets)
			assertFloat64(t, name, got, test.want, 0.0001)
		})
	}
}

func assertFloat64(t *testing.T, label string, got, want, tolerance float64) {
	if diff := math.Abs(want - got); diff >= tolerance {
		t.Errorf("\t got = %f %s\n\t\t\twant = %f", got, label, want)
	}
}
