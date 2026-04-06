// Package game contains the core game state and rules engine.
package game

import (
	"errors"
	"fmt"
	"mtgsim/pkg/api"
	"mtgsim/pkg/card"
	"mtgsim/pkg/mana"
	"mtgsim/pkg/player"
	"mtgsim/pkg/stack"
	"mtgsim/pkg/turn"
	"mtgsim/pkg/zone"
)

// CastChoices encapsulates player choices made during the casting of a spell.
type CastChoices struct {
	XValue      int
	ChosenModes []string // Placeholder for specific modes chosen
	Targets     []string // Placeholder for specific targets chosen (e.g., creature IDs)
}

// Game represents the entire state of a single game of Magic: The Gathering.
type Game struct {
	Players           []*player.Player
	ActivePlayer      *player.Player
	PriorityPlayer    *player.Player
	Turn              *turn.Turn
	Stack             *stack.Stack
	Zones             map[string]*zone.Zone
	consecutivePasses int
}

// NewGame creates and initializes a new game for a set of players.
func NewGame(playerIDs []int, startingLife int) (*Game, error) {
	if len(playerIDs) == 0 {
		return nil, errors.New("must have at least one player")
	}

	g := &Game{
		Stack:             stack.NewStack(),
		Turn:              turn.NewTurn(),
		Zones:             make(map[string]*zone.Zone),
		consecutivePasses: 0,
	}

	for _, id := range playerIDs {
		p := player.NewPlayer(id, startingLife)
		g.Players = append(g.Players, p)

		g.Zones[fmt.Sprintf("p%d_library", id)] = zone.NewZone()
		g.Zones[fmt.Sprintf("p%d_hand", id)] = zone.NewZone()
		g.Zones[fmt.Sprintf("p%d_graveyard", id)] = zone.NewZone()
	}

	g.ActivePlayer = g.Players[0]
	g.PriorityPlayer = g.Players[0]

	g.Zones["stack"] = zone.NewZone()
	g.Zones["exile"] = zone.NewZone()
	g.Zones["battlefield"] = zone.NewZone()

	return g, nil
}

// PaySpellCost attempts to pay the mana cost of a spell.
func (g *Game) PaySpellCost(player *player.Player, spell *stack.Spell, payment mana.Payment) error {
	spell.DetermineTotalCost() // Ensure FinalCost is calculated

	if !player.ManaPool.Pay(*spell.FinalCost, payment) {
		return errors.New("failed to pay mana cost")
	}
	return nil
}

// CheckState evaluates the game state to determine if the stack should resolve or the turn should advance.
func (g *Game) CheckState() {
	if g.consecutivePasses < len(g.Players) {
		// Not all players have passed priority yet.
		return
	}

	if len(g.Stack.Spells) > 0 {
		// Stack is not empty, resolve the top item.
		g.resolve(g.Stack.Pop())
	} else {
		// Stack is empty, advance to the next step/phase.
		g.Turn.Next()
		g.handleStepBasedActions()
	}

	// After resolution or advancing the turn, reset passes and give priority to the active player.
	g.consecutivePasses = 0
	g.PriorityPlayer = g.ActivePlayer
}

// PerformAction is the main entry point for a player to perform an action.
func (g *Game) PerformAction(playerID int, action *api.Action) error {
	if g.PriorityPlayer.ID != playerID {
		return fmt.Errorf("player %d does not have priority", playerID)
	}

	switch action.Type {
	case api.ActionType_PASS_PRIORITY:
		g.PassPriority()
	default:
		return fmt.Errorf("unsupported action type: %s", action.Type)
	}

	// After any action, the game state should be checked.
	g.CheckState()

	return nil
}

// GetValidTargets returns a list of valid targets for a given card.
// This is a placeholder for the complex targeting logic to be implemented.
func (g *Game) GetValidTargets(c *card.Card) []*card.Card {
	// For now, no cards have targets.
	return []*card.Card{}
}

// resolve handles the effect of a spell resolving.
// This is a placeholder and will become much more complex.
func (g *Game) resolve(spell *stack.Spell) {
	// For now, we assume all spells are permanents that go to the battlefield.
	spell.Card.ControllerID = spell.CasterID
	battlefield := g.Zones["battlefield"]
	battlefield.Add(spell.Card)
}

// handleStepBasedActions executes any state-based actions or turn-based actions for the current step.
func (g *Game) handleStepBasedActions() {
	switch g.Turn.CurrentStep {
	case turn.UntapStep:
		g.untapPermanents()
	}
}

// untapPermanents handles the untap action for the active player.
func (g *Game) untapPermanents() {
	battlefield := g.Zones["battlefield"]
	for _, card := range battlefield.Cards {
		if card.ControllerID == g.ActivePlayer.ID && g.canUntap(card) {
			card.Tapped = false
		}
	}
}

// canUntap checks if a given permanent is allowed to untap.
// This is the hook for effects that prevent untapping (e.g., stun counters).
func (g *Game) canUntap(c *card.Card) bool {
	// For now, all permanents can untap.
	return true
}

// PassPriority passes priority to the next player in turn order.
func (g *Game) PassPriority() {
	g.consecutivePasses++

	currentIndex := -1
	for i, p := range g.Players {
		if p.ID == g.PriorityPlayer.ID {
			currentIndex = i
			break
		}
	}

	if currentIndex == -1 {
		// Should not happen, but as a fallback, give priority to the active player.
		g.PriorityPlayer = g.ActivePlayer
		return
	}

	nextIndex := (currentIndex + 1) % len(g.Players)
	g.PriorityPlayer = g.Players[nextIndex]
}

// CastSpell orchestrates the full process of casting a spell.
func (g *Game) CastSpell(p *player.Player, cardToCast *card.Card, choices *CastChoices, payment mana.Payment) (*stack.Spell, error) {
	// 1. Check if card is in hand and get the hand zone.
	handZoneName := fmt.Sprintf("p%d_hand", p.ID)
	handZone, ok := g.Zones[handZoneName]
	if !ok {
		return nil, fmt.Errorf("hand zone not found for player %d", p.ID)
	}

	cardInHand := false
	for _, c := range handZone.Cards {
		if c.ID == cardToCast.ID {
			cardInHand = true
			break
		}
	}
	if !cardInHand {
		return nil, errors.New("card not in hand")
	}

	// 2. Create a new Spell object.
	spell := &stack.Spell{
		Card:        cardToCast,
		CasterID:    p.ID,
		XValue:      choices.XValue,
		ChosenModes: choices.ChosenModes,
		Targets:     choices.Targets,
		FinalCost:   cardToCast.ManaCost, // Initial cost, will be modified by DetermineTotalCost
	}

	// 3. Determine total cost.
	spell.DetermineTotalCost()

	// 4. Pay spell cost.
	err := g.PaySpellCost(p, spell, payment)
	if err != nil {
		return nil, err
	}

	// 5. Remove card from hand (only after successful payment).
	if err := handZone.Remove(cardToCast); err != nil {
		// This should ideally not happen if the card was found, but handle it defensively.
		return nil, fmt.Errorf("failed to remove card from hand after payment: %w", err)
	}

	// 6. Put spell on stack.
	g.Stack.Push(spell)

	// 7. Active player gets priority after casting a spell.
	g.PriorityPlayer = p
	g.consecutivePasses = 0

	return spell, nil
}
