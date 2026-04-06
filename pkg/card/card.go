// Package card contains the definition of a card.
package card

import (
	"github.com/google/uuid"
	"mtgsim/pkg/mana"
)

// Card represents a single Magic: The Gathering card.
type Card struct {
	ID           uuid.UUID
	Name         string
	ManaCost     *mana.Cost
	Tapped       bool
	ControllerID int
	// other fields like type, abilities will be added later.
}
