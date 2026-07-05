package arkan

import "codera-battle/internal"

// fake combatant for testing
type fakeCombatant struct {
	hp    int
	maxHP int
	alive bool
}

func (f *fakeCombatant) GetName() string       { return "test" }
func (f *fakeCombatant) GetStats() internal.Stats { return internal.Stats{MaxHP: f.maxHP} }
func (f *fakeCombatant) GetCurrentHP() int     { return f.hp }
func (f *fakeCombatant) SetCurrentHP(hp int)   { f.hp = hp }
func (f *fakeCombatant) GetMaxHP() int         { return f.maxHP }
func (f *fakeCombatant) IsAlive() bool         { return f.alive }

// test heal choice
func TestChooseActionHealsWhenCritical(t *testing.T) {
	hero := NewHero("Masato")
	allies := []internal.Combatant{&fakeCombatant{hp: 20, maxHP: 100, alive: true}}

	skill := hero.ChooseAction(allies)

	if skill.Name != "Klaerende Annotation" {
		t.Errorf("expected heal skill, got %s", skill.Name)
	}
}

// test attack choice
func TestChooseActionAttacksWhenHealthy(t *testing.T) {
	hero := NewHero("Masato")
	allies := []internal.Combatant{&fakeCombatant{hp: 90, maxHP: 100, alive: true}}

	skill := hero.ChooseAction(allies)

	if skill.Name == "Klaerende Annotation" {
		t.Errorf("expected attack skill, got heal")
	}
}
