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
	states       []turnState
	currentIndex int
}

type turnState struct {
	phase Phase
	step  Step
}

var turnOrder = []turnState{
	{BeginningPhase, UntapStep},
	{BeginningPhase, UpkeepStep},
	{BeginningPhase, DrawStep},
	{PreCombatMain, ""},
	{CombatPhase, BeginCombatStep},
	{CombatPhase, DeclareAttackersStep},
	{CombatPhase, DeclareBlockersStep},
	{CombatPhase, CombatDamageStep},
	{CombatPhase, EndCombatStep},
	{PostCombatMain, ""},
	{EndPhase, EndStep},
	{EndPhase, CleanupStep},
}

// NewTurn creates a new Turn manager.
func NewTurn() *Turn {
	return &Turn{
		TurnNumber:   1,
		CurrentPhase: turnOrder[0].phase,
		CurrentStep:  turnOrder[0].step,
		states:       turnOrder,
		currentIndex: 0,
	}
}

// Next advances the turn to the next state.
func (t *Turn) Next() {
	t.currentIndex++
	if t.currentIndex >= len(t.states) {
		t.currentIndex = 0
		t.TurnNumber++
	}

	nextState := t.states[t.currentIndex]
	t.CurrentPhase = nextState.phase
	t.CurrentStep = nextState.step
}
