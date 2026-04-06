// Package player contains structures and functions related to players.
package player

import (
	"mtgsim/pkg/mana"
)

// Player represents a single player in the game.
type Player struct {
	ID       int
	Life     int
	ManaPool *mana.Pool
}

// NewPlayer creates a new player with a given starting life total.
func NewPlayer(id int, startingLife int) *Player {
	return &Player{
		ID:       id,
		Life:     startingLife,
		ManaPool: &mana.Pool{},
	}
}
