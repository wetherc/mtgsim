// Package turn manages the turn structure, phases, and steps.
package turn

// Phase represents a phase in a turn.
type Phase string

const (
	BeginningPhase Phase = "Beginning"
	PreCombatMain  Phase = "Pre-Combat Main"
	CombatPhase    Phase = "Combat"
	PostCombatMain Phase = "Post-Combat Main"
	EndPhase       Phase = "End"
)

// Step represents a step within a phase.
type Step string

const (
	// Beginning Phase steps
	UntapStep   Step = "Untap"
	UpkeepStep  Step = "Upkeep"
	DrawStep    Step = "Draw"

	// Combat Phase steps
	BeginCombatStep     Step = "Begin Combat"
	DeclareAttackersStep Step = "Declare Attackers"
	DeclareBlockersStep  Step = "Declare Blockers"
	CombatDamageStep     Step = "Combat Damage"
	EndCombatStep        Step = "End Combat"

	// End Phase steps
	EndStep       Step = "End"
	CleanupStep   Step = "Cleanup"
)

// Turn represents the current turn, including the active phase and step.
type Turn struct {
	CurrentPhase Phase
	CurrentStep  Step
	TurnNumber   int
}
