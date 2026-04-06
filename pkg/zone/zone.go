// Package zone manages the game's zones (e.g., library, graveyard).
package zone

import (
	"errors"
	"math/rand"
	"mtgsim/pkg/card"
	"time"
)

// Zone represents a collection of cards.
type Zone struct {
	Cards []*card.Card
}

// NewZone creates a new, empty zone.
func NewZone() *Zone {
	return &Zone{
		Cards: []*card.Card{},
	}
}

// Add adds a card to the top of the zone (e.g., top of library).
func (z *Zone) Add(c *card.Card) {
	z.Cards = append(z.Cards, c)
}

// Remove removes a specific card from the zone.
func (z *Zone) Remove(c *card.Card) error {
	for i, cardInZone := range z.Cards {
		if cardInZone.ID == c.ID {
			z.Cards = append(z.Cards[:i], z.Cards[i+1:]...)
			return nil
		}
	}
	return errors.New("card not found in zone")
}

// Draw removes and returns the top card from the zone.
func (z *Zone) Draw() (*card.Card, error) {
	if len(z.Cards) == 0 {
		return nil, errors.New("zone is empty")
	}
	card := z.Cards[len(z.Cards)-1]
	z.Cards = z.Cards[:len(z.Cards)-1]
	return card, nil
}

// Shuffle randomizes the order of cards in the zone.
func (z *Zone) Shuffle() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(z.Cards), func(i, j int) {
		z.Cards[i], z.Cards[j] = z.Cards[j], z.Cards[i]
	})
}
