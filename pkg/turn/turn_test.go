package turn_test

import (
	"mtgsim/pkg/turn"
	"testing"
)

func TestTurnProgression(t *testing.T) {
	gameTurn := turn.NewTurn()

	tests := []struct {
		expectedPhase turn.Phase
		expectedStep  turn.Step
	}{
		{turn.BeginningPhase, turn.UpkeepStep},
		{turn.BeginningPhase, turn.DrawStep},
		{turn.PreCombatMain, ""},
		{turn.CombatPhase, turn.BeginCombatStep},
		{turn.CombatPhase, turn.DeclareAttackersStep},
		{turn.CombatPhase, turn.DeclareBlockersStep},
		{turn.CombatPhase, turn.CombatDamageStep},
		{turn.CombatPhase, turn.EndCombatStep},
		{turn.PostCombatMain, ""},
		{turn.EndPhase, turn.EndStep},
		{turn.EndPhase, turn.CleanupStep},
	}

	for _, tc := range tests {
		gameTurn.Next()
		if gameTurn.CurrentPhase != tc.expectedPhase {
			t.Errorf("Expected phase %s, got %s", tc.expectedPhase, gameTurn.CurrentPhase)
		}
		if gameTurn.CurrentStep != tc.expectedStep {
			t.Errorf("Expected step %s, got %s", tc.expectedStep, gameTurn.CurrentStep)
		}
	}

	// Test wrapping to a new turn
	gameTurn.Next()
	if gameTurn.TurnNumber != 2 {
		t.Errorf("Expected turn number 2, got %d", gameTurn.TurnNumber)
	}
	if gameTurn.CurrentPhase != turn.BeginningPhase {
		t.Errorf("Expected phase %s, got %s", turn.BeginningPhase, gameTurn.CurrentPhase)
	}
	if gameTurn.CurrentStep != turn.UntapStep {
		t.Errorf("Expected step %s, got %s", turn.UntapStep, gameTurn.CurrentStep)
	}
}
