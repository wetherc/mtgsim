package game

import (
	"fmt"
	"mtgsim/pkg/card"
	"mtgsim/pkg/mana"
	"mtgsim/pkg/player"
	"testing"

	"github.com/google/uuid"
)

// Helper function to create a basic card
func newTestCard(name string, cost mana.Cost) *card.Card {
	return &card.Card{
		ID:       uuid.New(),
		Name:     name,
		ManaCost: &cost,
	}
}

func TestNewGame(t *testing.T) {
	playerIDs := []int{1, 2}
	startingLife := 20
	game, err := NewGame(playerIDs, startingLife)
	if err != nil {
		t.Fatalf("NewGame() returned an unexpected error: %v", err)
	}

	if game == nil {
		t.Fatal("NewGame returned nil")
	}
	if game.Stack == nil {
		t.Error("NewGame created game with nil Stack")
	}
	if len(game.Players) != 2 {
		t.Errorf("Expected 2 players, got %d", len(game.Players))
	}
	if game.ActivePlayer.ID != 1 {
		t.Errorf("Expected active player to be P1, got P%d", game.ActivePlayer.ID)
	}
	if game.Zones[fmt.Sprintf("p%d_hand", 1)] == nil {
		t.Error("P1 hand zone not initialized")
	}
	if game.Zones["battlefield"] == nil {
		t.Error("Battlefield zone not initialized")
	}
}

func TestPassPriority(t *testing.T) {
	game, _ := NewGame([]int{1, 2}, 20)

	if game.PriorityPlayer.ID != 1 {
		t.Fatalf("Expected priority player to be P1, got P%d", game.PriorityPlayer.ID)
	}

	game.PassPriority()
	if game.PriorityPlayer.ID != 2 {
		t.Errorf("Expected priority player to be P2, got P%d", game.PriorityPlayer.ID)
	}

	game.PassPriority()
	if game.PriorityPlayer.ID != 1 {
		t.Errorf("Expected priority player to be P1 after wrapping, got P%d", game.PriorityPlayer.ID)
	}
}

func TestCheckState(t *testing.T) {
	t.Run("Advance step when stack is empty and all pass", func(t *testing.T) {
		game, _ := NewGame([]int{1, 2}, 20)
		initialPhase := game.Turn.CurrentPhase
		initialStep := game.Turn.CurrentStep

		// All players pass
		game.PassPriority()
		game.PassPriority()

		game.CheckState()

		if game.Turn.CurrentPhase == initialPhase && game.Turn.CurrentStep == initialStep {
			t.Error("Turn did not advance after all players passed on empty stack")
		}
		if game.consecutivePasses != 0 {
			t.Errorf("Expected consecutivePasses to be 0 after advancing step, got %d", game.consecutivePasses)
		}
		if game.PriorityPlayer.ID != game.ActivePlayer.ID {
			t.Errorf("Expected active player to have priority after advancing step, got P%d", game.PriorityPlayer.ID)
		}
	})

	t.Run("Resolve stack when not empty and all pass", func(t *testing.T) {
		game, _ := NewGame([]int{1, 2}, 20)
		p1 := game.Players[0]
		testCard := newTestCard("Test Spell", mana.Cost{})
		
		// P1 casts a spell
		handZone := game.Zones[fmt.Sprintf("p%d_hand", p1.ID)]
		handZone.Add(testCard)
		p1.ManaPool.Add(mana.Colorless, 0) // Zero cost
		game.CastSpell(p1, testCard, &CastChoices{}, mana.Payment{})

		// All players pass
		game.PassPriority()
		game.PassPriority()

		game.CheckState()

		if len(game.Stack.Spells) != 0 {
			t.Errorf("Expected stack to be empty after resolution, got %d spells", len(game.Stack.Spells))
		}
		battlefield := game.Zones["battlefield"]
		if len(battlefield.Cards) != 1 || battlefield.Cards[0].ID != testCard.ID {
			t.Error("Spell did not resolve to the battlefield correctly")
		}
		if game.consecutivePasses != 0 {
			t.Errorf("Expected consecutivePasses to be 0 after resolution, got %d", game.consecutivePasses)
		}
		if game.PriorityPlayer.ID != game.ActivePlayer.ID {
			t.Errorf("Expected active player to have priority after resolution, got P%d", game.PriorityPlayer.ID)
		}
	})
}

func TestPriorityAndPassing(t *testing.T) {
	game, _ := NewGame([]int{1, 2}, 20)
	p1 := game.Players[0]
	p2 := game.Players[1]
	testCard := newTestCard("Test Spell", mana.Cost{})

	// Setup: P1 has a card and mana
	p1.ManaPool.Add(mana.Colorless, 1)
	handZoneP1 := game.Zones[fmt.Sprintf("p%d_hand", p1.ID)]
	handZoneP1.Add(testCard)

	// Action: P1 casts a spell
	_, err := game.CastSpell(p1, testCard, &CastChoices{}, mana.Payment{})
	if err != nil {
		t.Fatalf("CastSpell failed unexpectedly: %v", err)
	}

	// Validation 1: P1 holds priority after casting
	if game.PriorityPlayer.ID != p1.ID {
		t.Errorf("Expected P1 to have priority after casting, but P%d does", game.PriorityPlayer.ID)
	}
	if game.consecutivePasses != 0 {
		t.Errorf("Expected consecutivePasses to be 0 after casting, but got %d", game.consecutivePasses)
	}

	// Action: P1 passes priority
	game.PassPriority()

	// Validation 2: P2 gets priority
	if game.PriorityPlayer.ID != p2.ID {
		t.Errorf("Expected P2 to have priority after P1 passes, but P%d does", game.PriorityPlayer.ID)
	}
	if game.consecutivePasses != 1 {
		t.Errorf("Expected consecutivePasses to be 1 after P1 passes, but got %d", game.consecutivePasses)
	}

	// Action: P2 passes priority
	game.PassPriority()

	// Validation 3: Priority wraps back to P1
	if game.PriorityPlayer.ID != p1.ID {
		t.Errorf("Expected P1 to have priority after P2 passes, but P%d does", game.PriorityPlayer.ID)
	}
	if game.consecutivePasses != 2 {
		t.Errorf("Expected consecutivePasses to be 2 after P2 passes, but got %d", game.consecutivePasses)
	}
}


func TestCastSpell(t *testing.T) {
	testCard := newTestCard("Lightning Bolt", mana.Cost{Colored: [mana.NumManaTypes-1]int{mana.Red: 1}})
	
	tests := []struct {
		name          string
		setup         func(g *Game) *player.Player // Setup returns the player who will cast
		cardToCast    *card.Card
		choices       *CastChoices
		payment       mana.Payment
		expectedErr   string
		validate      func(t *testing.T, g *Game, p *player.Player)
	}{
		{
			name: "Success: Cast a simple spell",
			setup: func(g *Game) *player.Player {
				p := g.Players[0]
				p.ManaPool.Add(mana.Red, 1)
				handZone := g.Zones[fmt.Sprintf("p%d_hand", p.ID)]
				handZone.Add(testCard)
				return p
			},
			cardToCast:  testCard,
			choices:     &CastChoices{},
			payment:     mana.Payment{Amounts: [mana.NumManaTypes]int{mana.Red: 1}},
			expectedErr: "",
			validate: func(t *testing.T, g *Game, p *player.Player) {
				if len(g.Stack.Spells) != 1 {
					t.Errorf("Expected 1 spell on the stack, got %d", len(g.Stack.Spells))
				}
				if g.Stack.Spells[0].Card.ID != testCard.ID {
					t.Error("Wrong card on stack")
				}
				handZone := g.Zones[fmt.Sprintf("p%d_hand", p.ID)]
				if len(handZone.Cards) != 0 {
					t.Errorf("Expected hand to be empty, got %d cards", len(handZone.Cards))
				}
				if p.ManaPool.Amounts[mana.Red] != 0 {
					t.Errorf("Expected mana pool to have 0 red mana, got %d", p.ManaPool.Amounts[mana.Red])
				}
			},
		},
		{
			name: "Failure: Card not in hand",
			setup: func(g *Game) *player.Player {
				p := g.Players[0]
				p.ManaPool.Add(mana.Red, 1)
				// Card is not added to hand
				return p
			},
			cardToCast:  testCard,
			choices:     &CastChoices{},
			payment:     mana.Payment{Amounts: [mana.NumManaTypes]int{mana.Red: 1}},
			expectedErr: "card not in hand",
			validate: func(t *testing.T, g *Game, p *player.Player) {
				if len(g.Stack.Spells) != 0 {
					t.Errorf("Expected stack to be empty, got %d", len(g.Stack.Spells))
				}
				if p.ManaPool.Amounts[mana.Red] != 1 {
					t.Errorf("Expected mana pool to be unchanged, got %d red mana", p.ManaPool.Amounts[mana.Red])
				}
			},
		},
		{
			name: "Failure: Not enough mana",
			setup: func(g *Game) *player.Player {
				p := g.Players[0]
				// No mana added to pool
				handZone := g.Zones[fmt.Sprintf("p%d_hand", p.ID)]
				handZone.Add(testCard)
				return p
			},
			cardToCast:  testCard,
			choices:     &CastChoices{},
			payment:     mana.Payment{Amounts: [mana.NumManaTypes]int{mana.Red: 1}},
			expectedErr: "failed to pay mana cost",
			validate: func(t *testing.T, g *Game, p *player.Player) {
				handZone := g.Zones[fmt.Sprintf("p%d_hand", p.ID)]
				if len(handZone.Cards) != 1 {
					t.Errorf("Expected card to remain in hand, got %d cards", len(handZone.Cards))
				}
				if len(g.Stack.Spells) != 0 {
					t.Errorf("Expected stack to be empty, got %d", len(g.Stack.Spells))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game, _ := NewGame([]int{1}, 20)
			player := tt.setup(game)

			_, err := game.CastSpell(player, tt.cardToCast, tt.choices, tt.payment)

			if tt.expectedErr != "" {
				if err == nil || err.Error() != tt.expectedErr {
					t.Errorf("CastSpell() error = %v, wantErr %v", err, tt.expectedErr)
				}
			} else if err != nil {
				t.Errorf("CastSpell() unexpected error = %v", err)
			}
			
			tt.validate(t, game, player)
		})
	}
}
