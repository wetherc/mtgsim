// Package game contains the core game state and rules engine.
package game

import (
	"mtgsim/pkg/player"
	"mtgsim/pkg/turn"
)

// Game represents the entire state of a single game of Magic: The Gathering.
type Game struct {
	Players []*player.Player
	Turn    *turn.Turn
}
