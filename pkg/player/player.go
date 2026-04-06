// Package player contains structures and functions related to players.
package player

import (
	"mtgsim/pkg/card"
	"mtgsim/pkg/mana"
)

// Player represents a single player in the game.
type Player struct {
	ID       int
	Life     int
	ManaPool *mana.Pool
	Hand     []*card.Card
	// other fields like library, graveyard will be added later.
}
