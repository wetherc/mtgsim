// Package game contains the core game state and rules engine.
package game

import (
	"errors"
	"mtgsim/pkg/card"
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

// InitiateCasting initiates the process of casting a spell.
// It moves the card from the player's hand and creates a Spell object.
// It does NOT yet handle cost calculation, payment, or putting the spell onto the stack.
func (g *Game) InitiateCasting(player *player.Player, cardToCast *card.Card, choices *CastChoices) (*stack.Spell, error) {
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

	// 2. Remove card from hand.
	player.Hand = append(player.Hand[:cardIndex], player.Hand[cardIndex+1:]...)

	// 3. Create a new Spell object.
	spell := &stack.Spell{
		Card:        cardToCast,
		XValue:      choices.XValue,
		ChosenModes: choices.ChosenModes,
		Targets:     choices.Targets,
		FinalCost:   cardToCast.ManaCost, // Initial cost, will be modified later
	}

	return spell, nil
}
