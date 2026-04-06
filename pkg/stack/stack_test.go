package stack

import (
	"mtgsim/pkg/card"
	"mtgsim/pkg/mana"
	"reflect"
	"testing"

	"github.com/google/uuid"
)

func TestSpell_DetermineTotalCost(t *testing.T) {
	// Create dummy Card objects for testing
	cardWithCost := &card.Card{
		ID:   uuid.New(),
		Name: "Test Card 1",
		ManaCost: &mana.Cost{
			Colored: [mana.NumManaTypes - 1]int{mana.Blue: 1, mana.Red: 1},
			Generic: 2,
		},
	}

	cardWithXCost := &card.Card{
		ID:   uuid.New(),
		Name: "Test Card 2 (X Cost)",
		ManaCost: &mana.Cost{
			Colored: [mana.NumManaTypes - 1]int{mana.Green: 1},
			Generic: mana.XValuePlaceholder, // Assume XValuePlaceholder is a constant representing X
		},
	}

	cardNoCost := &card.Card{
		ID:       uuid.New(),
		Name:     "Test Card 3 (No Cost)",
		ManaCost: nil,
	}

	tests := []struct {
		name         string
		spell        Spell
		expectedCost mana.Cost
	}{
		{
			name: "Fixed Mana Cost",
			spell: Spell{
				Card: cardWithCost,
			},
			expectedCost: mana.Cost{
				Colored: [mana.NumManaTypes - 1]int{mana.Blue: 1, mana.Red: 1},
				Generic: 2,
			},
		},
		{
			name: "X Mana Cost with XValue = 0",
			spell: Spell{
				Card: cardWithXCost,
				XValue: 0,
			},
			expectedCost: mana.Cost{
				Colored: [mana.NumManaTypes - 1]int{mana.Green: 1},
				Generic: 0,
			},
		},
		{
			name: "X Mana Cost with XValue = 3",
			spell: Spell{
				Card: cardWithXCost,
				XValue: 3,
			},
			expectedCost: mana.Cost{
				Colored: [mana.NumManaTypes - 1]int{mana.Green: 1},
				Generic: 3,
			},
		},
		{
			name: "No Mana Cost",
			spell: Spell{
				Card: cardNoCost,
			},
			expectedCost: mana.Cost{}, // Should result in an empty cost
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.spell
			gotCost := s.DetermineTotalCost()

			if !reflect.DeepEqual(*gotCost, tt.expectedCost) {
				t.Errorf("DetermineTotalCost() got cost %v, want %v", *gotCost, tt.expectedCost)
			}
			if !reflect.DeepEqual(*s.FinalCost, tt.expectedCost) {
				t.Errorf("DetermineTotalCost() spell.FinalCost got %v, want %v", *s.FinalCost, tt.expectedCost)
			}
		})
	}
}

func TestStack_PushPop(t *testing.T) {
	s := &Stack{}

	spell1 := &Spell{Card: &card.Card{Name: "Spell 1"}}
	spell2 := &Spell{Card: &card.Card{Name: "Spell 2"}}

	s.Push(spell1)
	if len(s.Spells) != 1 || s.Spells[0] != spell1 {
		t.Errorf("Push failed, stack: %v", s.Spells)
	}

	s.Push(spell2)
	if len(s.Spells) != 2 || s.Spells[1] != spell2 {
		t.Errorf("Push failed, stack: %v", s.Spells)
	}

	poppedSpell := s.Pop()
	if poppedSpell != spell2 {
		t.Errorf("Pop got %v, want %v", poppedSpell, spell2)
	}
	if len(s.Spells) != 1 {
		t.Errorf("Stack length after pop got %v, want %v", len(s.Spells), 1)
	}

	poppedSpell = s.Pop()
	if poppedSpell != spell1 {
		t.Errorf("Pop got %v, want %v", poppedSpell, spell1)
	}
	if len(s.Spells) != 0 {
		t.Errorf("Stack length after pop got %v, want %v", len(s.Spells), 0)
	}

	poppedSpell = s.Pop()
	if poppedSpell != nil {
		t.Errorf("Pop from empty stack got %v, want nil", poppedSpell)
	}
}
