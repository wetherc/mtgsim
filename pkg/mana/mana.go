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
)

// Pool represents a collection of mana.
type Pool struct {
	// White mana
	W int
	// Blue mana
	U int
	// Black mana
	B int
	// Red mana
	R int
	// Green mana
	G int
	// Colorless mana
	C int
}

// Add adds mana of a given type to the pool.
func (p *Pool) Add(manaType Type, amount int) {
	switch manaType {
	case White:
		p.W += amount
	case Blue:
		p.U += amount
	case Black:
		p.B += amount
	case Red:
		p.R += amount
	case Green:
		p.G += amount
	case Colorless:
		p.C += amount
	}
}
