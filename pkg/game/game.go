// Package game contains the core game state and rules engine.
package game

import (
	"errors"
	"mtgsim/pkg/card"
	"mtgsim/pkg/mana"
	"mtgsim/pkg/player"
	"mtgsim/pkg/stack"
	"mtgsim/pkg/turn"
)

// CastChoices encapsulates player choices made during the casting of a spell.
type CastChoices struct {
	XValue      int
	ChosenModes []string // Placeholder for specific modes chosen
	Targets     []string // Placeholder for specific targets chosen (e.g., creature IDs)
}

// Game represents the entire state of a single game of Magic: The Gathering.
type Game struct {
	Players []*player.Player
	Turn    *turn.Turn
	Stack   *stack.Stack // The game's stack
}

// NewGame creates and initializes a new game.
func NewGame() *Game {
	return &Game{
		Stack: &stack.Stack{},
	}
}

// PaySpellCost attempts to pay the mana cost of a spell.
func (g *Game) PaySpellCost(player *player.Player, spell *stack.Spell, payment mana.Payment) error {
	spell.DetermineTotalCost() // Ensure FinalCost is calculated

	if !player.ManaPool.Pay(*spell.FinalCost, payment) {
		return errors.New("failed to pay mana cost")
	}
	return nil
}

// CastSpell orchestrates the full process of casting a spell.
func (g *Game) CastSpell(player *player.Player, cardToCast *card.Card, choices *CastChoices, payment mana.Payment) (*stack.Spell, error) {
	// 1. Check if card is in hand.
	cardInHand := false
	cardIndex := -1
	for i, c := range player.Hand {
		if c.ID == cardToCast.ID {
			cardInHand = true
			cardIndex = i
			break
		}
	}
	if !cardInHand {
		return nil, errors.New("card not in hand")
	}

	// 2. Create a new Spell object.
	spell := &stack.Spell{
		Card:        cardToCast,
		XValue:      choices.XValue,
		ChosenModes: choices.ChosenModes,
		Targets:     choices.Targets,
		FinalCost:   cardToCast.ManaCost, // Initial cost, will be modified by DetermineTotalCost
	}

	// 3. Determine total cost.
	spell.DetermineTotalCost()

	// 4. Pay spell cost.
	err := g.PaySpellCost(player, spell, payment)
	if err != nil {
		return nil, err
	}

	// 5. Remove card from hand (only after successful payment).
	player.Hand = append(player.Hand[:cardIndex], player.Hand[cardIndex+1:]...)

	// 6. Put spell on stack.
	g.Stack.Push(spell)

	return spell, nil
}
