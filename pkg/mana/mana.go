// Package mana handles mana, mana pools, and costs.
package mana

// Type represents a type of mana.
type Type int

const (
	White Type = iota
	Blue
	Black
	Red
	Green
	Colorless
	// numManaTypes is the total number of mana types.
	numManaTypes
)

// String returns the string representation of a mana type.
func (t Type) String() string {
	switch t {
	case White:
		return "W"
	case Blue:
		return "U"
	case Black:
		return "B"
	case Red:
		return "R"
	case Green:
		return "G"
	case Colorless:
		return "C"
	default:
		return ""
	}
}

// Pool represents a collection of mana.
type Pool struct {
	Amounts [numManaTypes]int
}

// Add adds mana of a given type to the pool.
func (p *Pool) Add(manaType Type, amount int) {
	if manaType < numManaTypes {
		p.Amounts[manaType] += amount
	}
}

// Total returns the total amount of mana in the pool.
func (p *Pool) Total() int {
	total := 0
	for _, amount := range p.Amounts {
		total += amount
	}
	return total
}

// Cost represents the mana cost of a card or ability.
type Cost struct {
	Colored [numManaTypes - 1]int
	Generic int
}

// Payment represents the mana selected to pay a cost.
type Payment struct {
	Amounts [numManaTypes]int
}

// CanPay checks if the pool has enough mana to pay the cost.
func (p *Pool) CanPay(cost Cost) bool {
	// Check colored costs.
	for i, c := range cost.Colored {
		if p.Amounts[i] < c {
			return false
		}
	}

	// Check generic cost.
	remainingMana := 0
	for i, amount := range p.Amounts {
		if i < len(cost.Colored) {
			remainingMana += amount - cost.Colored[i]
		} else {
			remainingMana += amount
		}
	}

	return remainingMana >= cost.Generic
}

// Pay deducts the cost from the pool given a specific payment.
// Returns true if payment was successful, false otherwise.
func (p *Pool) Pay(cost Cost, payment Payment) bool {
	// 1. Check if the payment covers the cost.
	// a. Colored costs
	for i, c := range cost.Colored {
		if payment.Amounts[i] < c {
			return false // Payment doesn't cover colored cost.
		}
	}
	// b. Total cost
	totalPayment := 0
	for _, amount := range payment.Amounts {
		totalPayment += amount
	}
	totalCost := cost.Generic
	for _, c := range cost.Colored {
		totalCost += c
	}
	if totalPayment != totalCost {
		return false // Payment doesn't match total cost.
	}

	// 2. Check if the player has the mana for the payment.
	for i, amount := range payment.Amounts {
		if p.Amounts[i] < amount {
			return false // Not enough mana in the pool.
		}
	}

	// 3. Deduct the payment from the pool.
	for i, amount := range payment.Amounts {
		p.Amounts[i] -= amount
	}

	return true
}
