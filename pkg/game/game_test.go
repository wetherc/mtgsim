package game

import (
	"mtgsim/pkg/card"
	"mtgsim/pkg/mana"
	"mtgsim/pkg/player"
	"reflect"
	"testing"

	"github.com/google/uuid"
)

// Helper function to create a basic player with a mana pool
func newTestPlayer(id int, life int, pool mana.Pool) *player.Player {
	return &player.Player{
		ID:       id,
		Life:     life,
		ManaPool: &pool,
		Hand:     []*card.Card{},
	}
}

// Helper function to create a basic card
func newTestCard(name string, cost mana.Cost) *card.Card {
	return &card.Card{
		ID:       uuid.New(),
		Name:     name,
		ManaCost: &cost,
	}
}

func TestNewGame(t *testing.T) {
	game := NewGame()
	if game == nil {
		t.Error("NewGame returned nil")
	}
	if game.Stack == nil {
		t.Error("NewGame created game with nil Stack")
	}
	if len(game.Stack.Spells) != 0 {
		t.Errorf("NewGame created game with non-empty Stack, got %d spells", len(game.Stack.Spells))
	}
}

func TestCastSpell(t *testing.T) {
	// Test Cases
	tests := []struct {
		name         string
		initialPlayer *player.Player // Use initialPlayer to define a fresh player state for each test
		card         *card.Card
		choices      *CastChoices
		payment      mana.Payment
		expectedErr  string
		expectedHand int
		expectedPool mana.Pool
		expectedStack int
	}{
		{
			name:        "Success: Cast a simple spell",
			initialPlayer: newTestPlayer(1, 20, mana.Pool{Amounts: [mana.NumManaTypes]int{mana.White: 5, mana.Red: 5, mana.Colorless: 5}}),
			card:        newTestCard("Lightning Bolt", mana.Cost{Colored: [mana.NumManaTypes - 1]int{mana.Red: 1}, Generic: 1}),
			choices:     &CastChoices{},
			payment:     mana.Payment{Amounts: [mana.NumManaTypes]int{mana.Red: 1, mana.Colorless: 1}},
			expectedErr: "",
			expectedHand: 0,
			expectedPool: mana.Pool{Amounts: [mana.NumManaTypes]int{mana.White: 5, mana.Red: 4, mana.Colorless: 4}},
			expectedStack: 1,
		},
		{
			name:        "Failure: Card not in hand",
			initialPlayer: newTestPlayer(1, 20, mana.Pool{Amounts: [mana.NumManaTypes]int{mana.White: 5, mana.Red: 5, mana.Colorless: 5}}),
			card:        newTestCard("Non-existent Card", mana.Cost{Generic: 1}),
			choices:     &CastChoices{},
			payment:     mana.Payment{Amounts: [mana.NumManaTypes]int{mana.Colorless: 1}},
			expectedErr: "card not in hand",
			expectedHand: 0, // Hand size should remain unchanged as card was not there
			expectedPool: mana.Pool{Amounts: [mana.NumManaTypes]int{mana.White: 5, mana.Red: 5, mana.Colorless: 5}},
			expectedStack: 0,
		},
		{
			name:        "Failure: Not enough mana (payment fails)",
			initialPlayer: newTestPlayer(2, 20, mana.Pool{Amounts: [mana.NumManaTypes]int{mana.White: 0}}), // Player with no mana
			card:        newTestCard("Expensive Spell", mana.Cost{Colored: [mana.NumManaTypes - 1]int{mana.White: 1}, Generic: 1}),
			choices:     &CastChoices{},
			payment:     mana.Payment{Amounts: [mana.NumManaTypes]int{mana.White: 1, mana.Colorless: 1}},
			expectedErr: "failed to pay mana cost",
			expectedHand: 0, // Card remains in hand
			expectedPool: mana.Pool{Amounts: [mana.NumManaTypes]int{mana.White: 0}},
			expectedStack: 0,
		},
		{
			name:        "Success: Cast spell with X cost (X=2)",
			initialPlayer: newTestPlayer(1, 20, mana.Pool{Amounts: [mana.NumManaTypes]int{mana.White: 5, mana.Red: 5, mana.Colorless: 5}}),
			card:        newTestCard("X Spell", mana.Cost{Generic: mana.XValuePlaceholder}),
			choices:     &CastChoices{XValue: 2},
			payment:     mana.Payment{Amounts: [mana.NumManaTypes]int{mana.White: 1, mana.Red: 1, mana.Colorless: 0}}, // Using 2 generic mana
			expectedErr: "",
			expectedHand: 0,
			expectedPool: mana.Pool{Amounts: [mana.NumManaTypes]int{mana.White: 4, mana.Red: 4, mana.Colorless: 5}}, // 1W, 1R paid for X, 5 colorless remains
			expectedStack: 1,
		},
		{
			name:        "Failure: Cast spell with X cost, insufficient mana",
			initialPlayer: newTestPlayer(3, 20, mana.Pool{Amounts: [mana.NumManaTypes]int{mana.Colorless: 1}}), // Only 1 mana
			card:        newTestCard("X Spell 2", mana.Cost{Generic: mana.XValuePlaceholder}),
			choices:     &CastChoices{XValue: 2},
			payment:     mana.Payment{Amounts: [mana.NumManaTypes]int{mana.Colorless: 2}},
			expectedErr: "failed to pay mana cost",
			expectedHand: 0, // Card remains in hand
			expectedPool: mana.Pool{Amounts: [mana.NumManaTypes]int{mana.Colorless: 1}},
			expectedStack: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clone player for tests that modify player state
			p := &player.Player{
				ID:       tt.initialPlayer.ID,
				Life:     tt.initialPlayer.Life,
				ManaPool: &mana.Pool{Amounts: tt.initialPlayer.ManaPool.Amounts},
				Hand:     make([]*card.Card, len(tt.initialPlayer.Hand)),
			}
			copy(p.Hand, tt.initialPlayer.Hand)

			// Add card to player's hand if it's supposed to be there for casting
			if tt.card != nil && tt.expectedErr != "card not in hand" {
				p.Hand = append(p.Hand, tt.card)
			}

			g := NewGame() // Each test gets a fresh game
			g.Players = append(g.Players, p)

			// initialStackLen is always 0 for a fresh game
			initialStackLen := 0

			spell, err := g.CastSpell(p, tt.card, tt.choices, tt.payment)

			if tt.expectedErr != "" {
				if err == nil || err.Error() != tt.expectedErr {
					t.Errorf("CastSpell() got error %v, want %v", err, tt.expectedErr)
				}
				// For failure cases, hand, pool, and stack should remain unchanged (except for card removal if not in hand)
				if !reflect.DeepEqual(*p.ManaPool, tt.expectedPool) {
					t.Errorf("CastSpell() failed payment, player mana pool got %v, want %v", *p.ManaPool, tt.expectedPool)
				}
				if len(g.Stack.Spells) != initialStackLen {
					t.Errorf("CastSpell() failed payment, stack size got %v, want %v", len(g.Stack.Spells), initialStackLen)
				}
				// If the card was in hand but payment failed, it should still be in hand
				if tt.expectedErr != "card not in hand" && len(p.Hand) != 1 {
					t.Errorf("CastSpell() failed payment, card not returned to hand. Hand size got %v, want %v", len(p.Hand), 1)
				}
				// If card was not in hand, the hand should remain empty
				if tt.expectedErr == "card not in hand" && len(p.Hand) != 0 {
					t.Errorf("CastSpell() failed payment, hand size got %v, want %v", len(p.Hand), 0)
				}

			} else {
				if err != nil {
					t.Errorf("CastSpell() got unexpected error: %v", err)
				}
				if len(p.Hand) != tt.expectedHand {
					t.Errorf("CastSpell() player hand got %v, want %v", len(p.Hand), tt.expectedHand)
				}
				if !reflect.DeepEqual(*p.ManaPool, tt.expectedPool) {
					t.Errorf("CastSpell() player mana pool got %v, want %v", *p.ManaPool, tt.expectedPool)
				}
				if len(g.Stack.Spells) != tt.expectedStack {
					t.Errorf("CastSpell() stack size got %v, want %v", len(g.Stack.Spells), tt.expectedStack)
				}
				if g.Stack.Spells[len(g.Stack.Spells)-1] != spell {
					t.Errorf("CastSpell() last spell on stack got %v, want %v", g.Stack.Spells[len(g.Stack.Spells)-1].Card.Name, spell.Card.Name)
				}
			}
		})
	}
}
