// Package stack manages the game's stack and spell objects.
package stack

import (
	"mtgsim/pkg/card"
	"mtgsim/pkg/mana"
)

// Spell represents a card that has been cast and is currently on the stack.
type Spell struct {
	Card *card.Card
	// Choices made during casting
	XValue      int
	ChosenModes []string // Placeholder for specific modes chosen
	Targets     []string // Placeholder for specific targets chosen (e.g., creature IDs)

	// Final calculated cost after all modifications (e.g., cost reductions, alternative costs)
	FinalCost *mana.Cost
}

// Stack represents the game's stack where spells and abilities await resolution.
type Stack struct {
	Spells []*Spell
}

// Push adds a spell to the top of the stack.
func (s *Stack) Push(spell *Spell) {
	s.Spells = append(s.Spells, spell)
}

// Pop removes and returns the top spell from the stack.
// Returns nil if the stack is empty.
func (s *Stack) Pop() *Spell {
	if len(s.Spells) == 0 {
		return nil
	}
	spell := s.Spells[len(s.Spells)-1]
	s.Spells = s.Spells[:len(s.Spells)-1]
	return spell
}
