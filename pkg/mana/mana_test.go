package mana

import (
	"reflect"
	"testing"
)

func TestPool_Add(t *testing.T) {
	tests := []struct {
		name     string
		initial  Pool
		manaType Type
		amount   int
		expected Pool
	}{
		{
			name:     "Add White Mana",
			initial:  Pool{},
			manaType: White,
			amount:   1,
			expected: Pool{Amounts: [NumManaTypes]int{White: 1}},
		},
		{
			name:     "Add Multiple Mana Types",
			initial:  Pool{Amounts: [NumManaTypes]int{White: 1}},
			manaType: Blue,
			amount:   2,
			expected: Pool{Amounts: [NumManaTypes]int{White: 1, Blue: 2}},
		},
		{
			name:     "Add Zero Amount",
			initial:  Pool{Amounts: [NumManaTypes]int{White: 5}},
			manaType: White,
			amount:   0,
			expected: Pool{Amounts: [NumManaTypes]int{White: 5}},
		},
		{
			name:     "Add Colorless Mana",
			initial:  Pool{},
			manaType: Colorless,
			amount:   3,
			expected: Pool{Amounts: [NumManaTypes]int{Colorless: 3}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tt.initial
			p.Add(tt.manaType, tt.amount)
			if !reflect.DeepEqual(p, tt.expected) {
				t.Errorf("Pool.Add() = %v, want %v", p, tt.expected)
			}
		})
	}
}

func TestPool_Total(t *testing.T) {
	tests := []struct {
		name     string
		pool     Pool
		expected int
	}{
		{
			name:     "Empty Pool",
			pool:     Pool{},
			expected: 0,
		},
		{
			name:     "Single Mana Type",
			pool:     Pool{Amounts: [NumManaTypes]int{White: 3}},
			expected: 3,
		},
		{
			name:     "Multiple Mana Types",
			pool:     Pool{Amounts: [NumManaTypes]int{White: 1, Blue: 2, Colorless: 3}},
			expected: 6,
		},
		{
			name:     "Negative Mana (should not happen in game, but testing logic)",
			pool:     Pool{Amounts: [NumManaTypes]int{White: 5, Black: -2}},
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pool.Total(); got != tt.expected {
				t.Errorf("Pool.Total() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestPool_CanPay(t *testing.T) {
	tests := []struct {
		name     string
		pool     Pool
		cost     Cost
		expected bool
	}{
		{
			name:     "Can Pay Exact Colored",
			pool:     Pool{Amounts: [NumManaTypes]int{White: 1, Blue: 1}},
			cost:     Cost{Colored: [NumManaTypes - 1]int{White: 1, Blue: 1}},
			expected: true,
		},
		{
			name:     "Can Pay More Colored",
			pool:     Pool{Amounts: [NumManaTypes]int{White: 2, Blue: 1}},
			cost:     Cost{Colored: [NumManaTypes - 1]int{White: 1}},
			expected: true,
		},
		{
			name:     "Cannot Pay Colored",
			pool:     Pool{Amounts: [NumManaTypes]int{White: 0, Blue: 1}},
			cost:     Cost{Colored: [NumManaTypes - 1]int{White: 1}},
			expected: false,
		},
		{
			name:     "Can Pay Exact Generic with Colorless",
			pool:     Pool{Amounts: [NumManaTypes]int{Colorless: 2}},
			cost:     Cost{Generic: 2},
			expected: true,
		},
		{
			name:     "Can Pay Exact Generic with Colored",
			pool:     Pool{Amounts: [NumManaTypes]int{White: 2}},
			cost:     Cost{Generic: 2},
			expected: true,
		},
		{
			name:     "Can Pay Exact Generic with Mixed",
			pool:     Pool{Amounts: [NumManaTypes]int{White: 1, Colorless: 1}},
			cost:     Cost{Generic: 2},
			expected: true,
		},
		{
			name:     "Cannot Pay Generic",
			pool:     Pool{Amounts: [NumManaTypes]int{White: 1, Colorless: 0}},
			cost:     Cost{Generic: 2},
			expected: false,
		},
		{
			name:     "Complex Cost Can Pay",
			pool:     Pool{Amounts: [NumManaTypes]int{White: 2, Blue: 1, Green: 1, Colorless: 3}},
			cost:     Cost{Colored: [NumManaTypes - 1]int{White: 1, Blue: 1}, Generic: 2},
			expected: true,
		},
		{
			name:     "Complex Cost Cannot Pay (Colored)",
			pool:     Pool{Amounts: [NumManaTypes]int{White: 0, Blue: 1, Green: 1, Colorless: 3}},
			cost:     Cost{Colored: [NumManaTypes - 1]int{White: 1, Blue: 1}, Generic: 2},
			expected: false,
		},
		{
			name:     "Complex Cost Cannot Pay (Generic)",
			pool:     Pool{Amounts: [NumManaTypes]int{White: 2, Blue: 1, Green: 1, Colorless: 0}},
			cost:     Cost{Colored: [NumManaTypes - 1]int{White: 1, Blue: 1}, Generic: 3},
			expected: false,
		},
		{
			name:     "Empty Pool, Empty Cost",
			pool:     Pool{},
			cost:     Cost{},
			expected: true,
		},
		{
			name:     "Pool with some mana, Empty Cost",
			pool:     Pool{Amounts: [NumManaTypes]int{White: 1}},
			cost:     Cost{},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pool.CanPay(tt.cost); got != tt.expected {
				t.Errorf("Pool.CanPay() for cost %v and pool %v = %v, want %v", tt.cost, tt.pool, got, tt.expected)
			}
		})
	}
}

func TestPool_Pay(t *testing.T) {
	tests := []struct {
		name        string
		initialPool Pool
		cost        Cost
		payment     Payment
		expectedErr bool
		expectedPool Pool
	}{
		{
			name:        "Successful Colored Payment",
			initialPool: Pool{Amounts: [NumManaTypes]int{White: 1, Blue: 1}},
			cost:        Cost{Colored: [NumManaTypes - 1]int{White: 1}},
			payment:     Payment{Amounts: [NumManaTypes]int{White: 1}},
			expectedErr: false,
			expectedPool: Pool{Amounts: [NumManaTypes]int{White: 0, Blue: 1}},
		},
		{
			name:        "Successful Generic Payment with White",
			initialPool: Pool{Amounts: [NumManaTypes]int{White: 2}},
			cost:        Cost{Generic: 2},
			payment:     Payment{Amounts: [NumManaTypes]int{White: 2}},
			expectedErr: false,
			expectedPool: Pool{Amounts: [NumManaTypes]int{White: 0}},
		},
		{
			name:        "Successful Generic Payment with Colorless",
			initialPool: Pool{Amounts: [NumManaTypes]int{Colorless: 3}},
			cost:        Cost{Generic: 3},
			payment:     Payment{Amounts: [NumManaTypes]int{Colorless: 3}},
			expectedErr: false,
			expectedPool: Pool{Amounts: [NumManaTypes]int{Colorless: 0}},
		},
		{
			name:        "Successful Mixed Payment",
			initialPool: Pool{Amounts: [NumManaTypes]int{White: 2, Blue: 2, Colorless: 2}},
			cost:        Cost{Colored: [NumManaTypes - 1]int{White: 1}, Generic: 3},
			payment:     Payment{Amounts: [NumManaTypes]int{White: 1, Blue: 1, Colorless: 2}},
			expectedErr: false,
			expectedPool: Pool{Amounts: [NumManaTypes]int{White: 1, Blue: 1, Colorless: 0}},
		},
		{
			name:        "Failure: Not Enough Colored Mana in Pool for Payment",
			initialPool: Pool{Amounts: [NumManaTypes]int{White: 0, Blue: 1}},
			cost:        Cost{Colored: [NumManaTypes - 1]int{White: 1}},
			payment:     Payment{Amounts: [NumManaTypes]int{White: 1}},
			expectedErr: true,
			expectedPool: Pool{Amounts: [NumManaTypes]int{White: 0, Blue: 1}}, // Pool should not change
		},
		{
			name:        "Failure: Payment Doesn't Cover Colored Cost",
			initialPool: Pool{Amounts: [NumManaTypes]int{White: 1, Blue: 1}},
			cost:        Cost{Colored: [NumManaTypes - 1]int{White: 1, Blue: 1}},
			payment:     Payment{Amounts: [NumManaTypes]int{White: 1}}, // Missing Blue
			expectedErr: true,
			expectedPool: Pool{Amounts: [NumManaTypes]int{White: 1, Blue: 1}}, // Pool should not change
		},
		{
			name:        "Failure: Payment Doesn't Match Total Cost",
			initialPool: Pool{Amounts: [NumManaTypes]int{White: 5}},
			cost:        Cost{Generic: 2},
			payment:     Payment{Amounts: [NumManaTypes]int{White: 3}}, // Overpaying
			expectedErr: true,
			expectedPool: Pool{Amounts: [NumManaTypes]int{White: 5}}, // Pool should not change
		},
		{
			name:        "Failure: Empty Pool, Non-empty Cost",
			initialPool: Pool{},
			cost:        Cost{Generic: 1},
			payment:     Payment{Amounts: [NumManaTypes]int{Colorless: 1}},
			expectedErr: true,
			expectedPool: Pool{}, // Pool should not change
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tt.initialPool
			got := p.Pay(tt.cost, tt.payment)
			if got != !tt.expectedErr { // if expectedErr is true, got should be false
				t.Errorf("Pool.Pay() for cost %v, payment %v and initial pool %v = %v, want %v", tt.cost, tt.payment, tt.initialPool, got, !tt.expectedErr)
			}
			if !reflect.DeepEqual(p, tt.expectedPool) {
				t.Errorf("Pool.Pay() pool after operation = %v, want %v", p, tt.expectedPool)
			}
		})
	}
}

func TestManaType_String(t *testing.T) {
	tests := []struct {
		name     string
		manaType Type
		expected string
	}{
		{name: "White", manaType: White, expected: "W"},
		{name: "Blue", manaType: Blue, expected: "U"},
		{name: "Black", manaType: Black, expected: "B"},
		{name: "Red", manaType: Red, expected: "R"},
		{name: "Green", manaType: Green, expected: "G"},
		{name: "Colorless", manaType: Colorless, expected: "C"},
		{name: "Invalid", manaType: Type(99), expected: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.manaType.String(); got != tt.expected {
				t.Errorf("Type.String() for %v = %v, want %v", tt.manaType, got, tt.expected)
			}
		})
	}
}
