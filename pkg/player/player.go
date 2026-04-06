// Package player contains structures and functions related to players.
package player

import "mtgsim/pkg/mana"

// Player represents a single player in the game.
type Player struct {
	ID       int
	Life     int
	ManaPool *mana.Pool
	// other fields like hand, library, graveyard will be added later.
}
