package krieger

import (
	"testing"

	"codera-battle/internal"
)

// ── Construction ─────────────────────────────────────────────────────────────

func TestNew_EffectiveStats(t *testing.T) {
	k := New()
	s := k.GetStats()

	// Base + equipment: MaxHP 150+40=190, ATK 22+10=32, DEF 14+8=22, SPD 8+2=10
	if s.MaxHP != 190 {
		t.Errorf("MaxHP: want 190, got %d", s.MaxHP)
	}
	if s.Attack != 32 {
		t.Errorf("Attack: want 32, got %d", s.Attack)
	}
	if s.Defense != 22 {
		t.Errorf("Defense: want 22, got %d", s.Defense)
	}
	if s.Speed != 10 {
		t.Errorf("Speed: want 10, got %d", s.Speed)
	}
}

func TestNew_StartsAtFullHP(t *testing.T) {
	k := New()
	if k.GetCurrentHP() != k.GetMaxHP() {
		t.Errorf("expected full HP at start: got %d/%d", k.GetCurrentHP(), k.GetMaxHP())
	}
}

// ── HP management ────────────────────────────────────────────────────────────

func TestSetCurrentHP_ClampsToZero(t *testing.T) {
	k := New()
	k.SetCurrentHP(-100)
	if k.GetCurrentHP() != 0 {
		t.Errorf("expected 0, got %d", k.GetCurrentHP())
	}
}

func TestSetCurrentHP_ClampsToMax(t *testing.T) {
	k := New()
	k.SetCurrentHP(k.GetMaxHP() + 999)
	if k.GetCurrentHP() != k.GetMaxHP() {
		t.Errorf("expected %d, got %d", k.GetMaxHP(), k.GetCurrentHP())
	}
}

func TestIsAlive_TrueAtFullHP(t *testing.T) {
	k := New()
	if !k.IsAlive() {
		t.Error("expected IsAlive=true at full HP")
	}
}

func TestIsAlive_FalseAtZeroHP(t *testing.T) {
	k := New()
	k.SetCurrentHP(0)
	if k.IsAlive() {
		t.Error("expected IsAlive=false at 0 HP")
	}
}

// ── Round buff system ────────────────────────────────────────────────────────

func TestSchutzschild_TempDefBonusActiveThisRound(t *testing.T) {
	k := New()
	baseDef := k.GetStats().Defense
	k.ApplyTempDefBonus(5)

	if k.GetStats().Defense != baseDef+5 {
		t.Errorf("DEF after Schutzschild: want %d, got %d", baseDef+5, k.GetStats().Defense)
	}
}

func TestSchutzschild_BonusClearedOnNextTurnStart(t *testing.T) {
	k := New()
	baseDef := k.GetStats().Defense
	k.ApplyTempDefBonus(5)
	k.OnTurnStart() // simulate next turn

	if k.GetStats().Defense != baseDef {
		t.Errorf("DEF after OnTurnStart: want %d (bonus cleared), got %d", baseDef, k.GetStats().Defense)
	}
}

func TestKampfschrei_AtkBonusActiveNextTurn(t *testing.T) {
	k := New()
	baseAtk := k.GetStats().Attack
	k.QueueAtkBonus(5) // queued this turn

	// Not active yet
	if k.GetStats().Attack != baseAtk {
		t.Errorf("ATK should not change before OnTurnStart; want %d, got %d", baseAtk, k.GetStats().Attack)
	}

	k.OnTurnStart() // next turn: promotes pending to active

	if k.GetStats().Attack != baseAtk+5 {
		t.Errorf("ATK after OnTurnStart: want %d, got %d", baseAtk+5, k.GetStats().Attack)
	}
}

func TestKampfschrei_PendingBonusClearedAfterPromotion(t *testing.T) {
	k := New()
	k.QueueAtkBonus(5)
	k.OnTurnStart() // promotes +5 to active
	k.OnTurnStart() // second turn: pending is 0, active resets to 0

	if k.GetStats().Attack != k.baseStats.Attack {
		t.Errorf("ATK should reset after second OnTurnStart; want %d, got %d",
			k.baseStats.Attack, k.GetStats().Attack)
	}
}

// ── Auto-defense threshold ───────────────────────────────────────────────────

func TestAutoDefense_FalseAbove30Pct(t *testing.T) {
	k := New()
	k.SetCurrentHP(k.GetMaxHP()) // 100%
	if k.IsAutoDefenseRequired() {
		t.Error("auto-defense should not trigger at full HP")
	}
}

func TestAutoDefense_TrueBelow30Pct(t *testing.T) {
	k := New()
	k.SetCurrentHP(int(float64(k.GetMaxHP()) * 0.29))
	if !k.IsAutoDefenseRequired() {
		t.Error("auto-defense should trigger below 30% HP")
	}
}

func TestAutoDefense_FalseAtExactly30Pct(t *testing.T) {
	k := New()
	k.SetCurrentHP(int(float64(k.GetMaxHP()) * 0.30))
	// exactly 30% -- threshold is strictly less than 0.30
	if k.IsAutoDefenseRequired() {
		t.Error("auto-defense should not trigger at exactly 30% HP")
	}
}

// ── Skills ───────────────────────────────────────────────────────────────────

func TestGetSkills_ExactlyThree(t *testing.T) {
	k := New()
	if len(k.GetSkills()) != 3 {
		t.Fatalf("expected 3 skills, got %d", len(k.GetSkills()))
	}
}

func TestGetSkills_CorrectTargetTypes(t *testing.T) {
	k := New()
	skills := k.GetSkills()

	cases := []struct {
		name   string
		target internal.TargetType
	}{
		{"Praeziser Hieb", internal.TargetSingleEnemy},
		{"Schutzschild", internal.TargetSelf},
		{"Kampfschrei", internal.TargetSingleEnemy},
	}

	for i, c := range cases {
		if skills[i].Name != c.name {
			t.Errorf("skill[%d] name: want %q, got %q", i, c.name, skills[i].Name)
		}
		if skills[i].Target != c.target {
			t.Errorf("skill[%d] target: want %q, got %q", i, c.target, skills[i].Target)
		}
	}
}

// ── Interface compliance ─────────────────────────────────────────────────────

func TestImplementsCombatant(t *testing.T) {
	var _ internal.Combatant = New()
}

func TestImplementsSkillUser(t *testing.T) {
	var _ internal.SkillUser = New()
}
