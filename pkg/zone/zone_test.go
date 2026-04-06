package zone_test

import (
	"mtgsim/pkg/card"
	"mtgsim/pkg/zone"
	"testing"

	"github.com/google/uuid"
)

func TestZone(t *testing.T) {
	z := zone.NewZone()
	c1 := &card.Card{ID: uuid.New()}
	c2 := &card.Card{ID: uuid.New()}

	// Test Add
	z.Add(c1)
	z.Add(c2)
	if len(z.Cards) != 2 {
		t.Fatalf("Expected 2 cards in zone, got %d", len(z.Cards))
	}
	if z.Cards[1].ID != c2.ID {
		t.Errorf("Expected c2 to be the last card added")
	}

	// Test Draw
	drawnCard, err := z.Draw()
	if err != nil {
		t.Fatalf("Unexpected error drawing card: %v", err)
	}
	if drawnCard.ID != c2.ID {
		t.Errorf("Expected to draw c2, got %s", drawnCard.ID)
	}
	if len(z.Cards) != 1 {
		t.Fatalf("Expected 1 card in zone after drawing, got %d", len(z.Cards))
	}

	// Test Remove
	err = z.Remove(c1)
	if err != nil {
		t.Fatalf("Unexpected error removing card: %v", err)
	}
	if len(z.Cards) != 0 {
		t.Fatalf("Expected 0 cards in zone after removing, got %d", len(z.Cards))
	}

	// Test Remove non-existent
	err = z.Remove(&card.Card{ID: uuid.New()})
	if err == nil {
		t.Error("Expected an error when removing a non-existent card")
	}

	// Test Draw from empty
	_, err = z.Draw()
	if err == nil {
		t.Error("Expected an error when drawing from an empty zone")
	}
}

func TestShuffle(t *testing.T) {
	z := zone.NewZone()
	for i := 0; i < 100; i++ {
		z.Add(&card.Card{ID: uuid.New()})
	}

	originalOrder := make([]*card.Card, len(z.Cards))
	copy(originalOrder, z.Cards)

	z.Shuffle()

	if len(originalOrder) != len(z.Cards) {
		t.Fatalf("Shuffle changed the number of cards")
	}

	isSameOrder := true
	for i := range originalOrder {
		if originalOrder[i].ID != z.Cards[i].ID {
			isSameOrder = false
			break
		}
	}
	if isSameOrder {
		t.Errorf("Shuffle did not change the order of cards")
	}
}
