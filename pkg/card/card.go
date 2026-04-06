// Package card contains the definition of a card.
package card

import "github.com/google/uuid"

// Card represents a single Magic: The Gathering card.
type Card struct {
	ID   uuid.UUID
	Name string
	// other fields like cost, type, abilities will be added later.
}
