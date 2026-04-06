package mana_test

import (
	"mtgsim/pkg/mana"
	"reflect"
	"testing"
)

func TestPay(t *testing.T) {
	tests := []struct {
		name          string
		initialPool   mana.Pool
		cost          mana.Cost
		payment       mana.Payment
		expectedSuccess bool
		expectedPool  mana.Pool
	}{
		{
			name:        "Success: Exact payment for colored cost",
			initialPool: mana.Pool{Amounts: [mana.NumManaTypes]int{mana.Red: 1}},
			cost:        mana.Cost{Colored: [mana.NumManaTypes - 1]int{mana.Red: 1}},
			payment:     mana.Payment{Amounts: [mana.NumManaTypes]int{mana.Red: 1}},
			expectedSuccess: true,
			expectedPool:  mana.Pool{Amounts: [mana.NumManaTypes]int{mana.Red: 0}},
		},
		{
			name:        "Success: Payment for generic cost",
			initialPool: mana.Pool{Amounts: [mana.NumManaTypes]int{mana.Colorless: 2}},
			cost:        mana.Cost{Generic: 2},
			payment:     mana.Payment{Amounts: [mana.NumManaTypes]int{mana.Colorless: 2}},
			expectedSuccess: true,
			expectedPool:  mana.Pool{Amounts: [mana.NumManaTypes]int{}},
		},
		{
			name:        "Success: Using colored mana for generic cost",
			initialPool: mana.Pool{Amounts: [mana.NumManaTypes]int{mana.Blue: 2}},
			cost:        mana.Cost{Generic: 2},
			payment:     mana.Payment{Amounts: [mana.NumManaTypes]int{mana.Blue: 2}},
			expectedSuccess: true,
			expectedPool:  mana.Pool{Amounts: [mana.NumManaTypes]int{}},
		},
		{
			name:        "Success: Mixed cost with exact payment",
			initialPool: mana.Pool{Amounts: [mana.NumManaTypes]int{mana.Red: 1, mana.Colorless: 2}},
			cost:        mana.Cost{Colored: [mana.NumManaTypes - 1]int{mana.Red: 1}, Generic: 2},
			payment:     mana.Payment{Amounts: [mana.NumManaTypes]int{mana.Red: 1, mana.Colorless: 2}},
			expectedSuccess: true,
			expectedPool:  mana.Pool{Amounts: [mana.NumManaTypes]int{}},
		},
		{
			name:        "Success: Mixed cost using other colored mana for generic part",
			initialPool: mana.Pool{Amounts: [mana.NumManaTypes]int{mana.Red: 1, mana.Blue: 2}},
			cost:        mana.Cost{Colored: [mana.NumManaTypes - 1]int{mana.Red: 1}, Generic: 2},
			payment:     mana.Payment{Amounts: [mana.NumManaTypes]int{mana.Red: 1, mana.Blue: 2}},
			expectedSuccess: true,
			expectedPool:  mana.Pool{Amounts: [mana.NumManaTypes]int{}},
		},
		{
			name:        "Failure: Not enough mana in pool",
			initialPool: mana.Pool{Amounts: [mana.NumManaTypes]int{mana.Red: 1}},
			cost:        mana.Cost{Colored: [mana.NumManaTypes - 1]int{mana.Red: 2}},
			payment:     mana.Payment{Amounts: [mana.NumManaTypes]int{mana.Red: 2}},
			expectedSuccess: false,
			expectedPool:  mana.Pool{Amounts: [mana.NumManaTypes]int{mana.Red: 1}}, // Pool unchanged
		},
		{
			name:        "Failure: Incorrect color paid",
			initialPool: mana.Pool{Amounts: [mana.NumManaTypes]int{mana.Blue: 1}},
			cost:        mana.Cost{Colored: [mana.NumManaTypes - 1]int{mana.Red: 1}},
			payment:     mana.Payment{Amounts: [mana.NumManaTypes]int{mana.Blue: 1}},
			expectedSuccess: false,
			expectedPool:  mana.Pool{Amounts: [mana.NumManaTypes]int{mana.Blue: 1}}, // Pool unchanged
		},
		{
			name:        "Failure: Payment does not match total cost",
			initialPool: mana.Pool{Amounts: [mana.NumManaTypes]int{mana.Red: 2}},
			cost:        mana.Cost{Colored: [mana.NumManaTypes - 1]int{mana.Red: 1}},
			payment:     mana.Payment{Amounts: [mana.NumManaTypes]int{mana.Red: 2}}, // Paying 2 for a cost of 1
			expectedSuccess: false,
			expectedPool:  mana.Pool{Amounts: [mana.NumManaTypes]int{mana.Red: 2}}, // Pool unchanged
		},
		{
			name:        "Success: Zero cost spell with zero payment",
			initialPool: mana.Pool{Amounts: [mana.NumManaTypes]int{mana.Red: 1}},
			cost:        mana.Cost{},
			payment:     mana.Payment{},
			expectedSuccess: true,
			expectedPool:  mana.Pool{Amounts: [mana.NumManaTypes]int{mana.Red: 1}},
		},
		{
			name:        "Failure: Zero cost spell with non-zero payment",
			initialPool: mana.Pool{Amounts: [mana.NumManaTypes]int{mana.Red: 1}},
			cost:        mana.Cost{},
			payment:     mana.Payment{Amounts: [mana.NumManaTypes]int{mana.Red: 1}},
			expectedSuccess: false,
			expectedPool:  mana.Pool{Amounts: [mana.NumManaTypes]int{mana.Red: 1}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool := tt.initialPool // Make a copy to avoid modifying the test case
			success := pool.Pay(tt.cost, tt.payment)

			if success != tt.expectedSuccess {
				t.Errorf("Pay() success = %v, want %v", success, tt.expectedSuccess)
			}

			if !reflect.DeepEqual(pool, tt.expectedPool) {
				t.Errorf("Pay() final pool = %v, want %v", pool, tt.expectedPool)
			}
		})
	}
}
