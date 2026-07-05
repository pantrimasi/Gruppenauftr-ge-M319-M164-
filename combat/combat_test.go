package combat

import (
	"testing"

	"codera-battle/dragon"
	"codera-battle/internal"
)

// ── CalculateDamage ──────────────────────────────────────────────────────────

func TestCalculateDamage_AlwaysHitsAtAccuracy1(t *testing.T) {
	for i := 0; i < 200; i++ {
		dmg, _, isMiss := CalculateDamage(10, 20, 0, 0, 1.0)
		if isMiss {
			t.Fatal("expected hit with accuracy 1.0, got miss")
		}
		if dmg < 1 {
			t.Fatalf("damage must be >= 1, got %d", dmg)
		}
	}
}

func TestCalculateDamage_AlwaysMissesAtAccuracy0(t *testing.T) {
	for i := 0; i < 200; i++ {
		_, _, isMiss := CalculateDamage(10, 20, 0, 0, 0.0)
		if !isMiss {
			t.Fatal("expected miss with accuracy 0.0, got hit")
		}
	}
}

func TestCalculateDamage_HighDefenseReducesDamage(t *testing.T) {
	lowTotal, highTotal := 0, 0
	for i := 0; i < 500; i++ {
		d1, _, m1 := CalculateDamage(20, 20, 10, 0, 1.0)
		d2, _, m2 := CalculateDamage(20, 20, 10, 80, 1.0)
		if !m1 {
			lowTotal += d1
		}
		if !m2 {
			highTotal += d2
		}
	}
	if lowTotal <= highTotal {
		t.Errorf("higher defense should reduce total damage: lowDef=%d highDef=%d", lowTotal, highTotal)
	}
}

func TestCalculateDamage_MinimumOneDamage(t *testing.T) {
	for i := 0; i < 500; i++ {
		dmg, _, isMiss := CalculateDamage(1, 1, 0, 99, 1.0)
		if !isMiss && dmg < 1 {
			t.Errorf("minimum damage must be 1, got %d", dmg)
		}
	}
}

// ── Initiative order ─────────────────────────────────────────────────────────

func TestBuildInitiativeOrder_DescendingSpeed(t *testing.T) {
	h1 := &stub{name: "slow", speed: 4}
	h2 := &stub{name: "fast", speed: 20}
	d := dragon.New() // speed 14

	order := buildInitiativeOrder([]internal.Combatant{h1, h2}, d)

	for i := 1; i < len(order); i++ {
		prev := order[i-1].Combatant.GetStats().Speed
		curr := order[i].Combatant.GetStats().Speed
		if curr > prev {
			t.Errorf("order not descending at index %d: %d > %d", i, curr, prev)
		}
	}
}

func TestBuildInitiativeOrder_HeroBeforeDragonOnTie(t *testing.T) {
	h := &stub{name: "hero", speed: 14, hp: 100} // same speed as dragon
	d := dragon.New()

	order := buildInitiativeOrder([]internal.Combatant{h}, d)

	if len(order) < 2 {
		t.Fatal("expected at least 2 participants")
	}
	if order[0].IsDragon {
		t.Error("hero should go before dragon on speed tie")
	}
}

// ── allDead ──────────────────────────────────────────────────────────────────

func TestAllDead_False(t *testing.T) {
	heroes := []internal.Combatant{&stub{hp: 0}, &stub{hp: 50}}
	if allDead(heroes) {
		t.Error("allDead should be false when one hero is alive")
	}
}

func TestAllDead_True(t *testing.T) {
	heroes := []internal.Combatant{&stub{hp: 0}, &stub{hp: 0}}
	if !allDead(heroes) {
		t.Error("allDead should be true when all heroes are dead")
	}
}

// ── liveTarget ───────────────────────────────────────────────────────────────

func TestLiveTarget_ReturnsIndexedTarget(t *testing.T) {
	h0 := &stub{name: "h0", hp: 10}
	h1 := &stub{name: "h1", hp: 10}
	got := liveTarget([]internal.Combatant{h0, h1}, 1)
	if got != h1 {
		t.Errorf("expected h1, got %v", got)
	}
}

func TestLiveTarget_FallsBackWhenIndexedIsDead(t *testing.T) {
	h0 := &stub{name: "dead", hp: 0}
	h1 := &stub{name: "alive", hp: 10}
	got := liveTarget([]internal.Combatant{h0, h1}, 0)
	if got != h1 {
		t.Errorf("expected fallback to h1, got %v", got)
	}
}

func TestLiveTarget_NilWhenNoneAlive(t *testing.T) {
	got := liveTarget([]internal.Combatant{&stub{hp: 0}}, 0)
	if got != nil {
		t.Errorf("expected nil, got %v", got)
	}
}

// ── stub ─────────────────────────────────────────────────────────────────────

type stub struct {
	name  string
	hp    int
	speed int
}

func (s *stub) GetName() string          { return s.name }
func (s *stub) GetStats() internal.Stats { return internal.Stats{MaxHP: s.hp, Speed: s.speed} }
func (s *stub) GetCurrentHP() int        { return s.hp }
func (s *stub) SetCurrentHP(hp int)      { s.hp = hp }
func (s *stub) GetMaxHP() int            { return s.hp }
func (s *stub) IsAlive() bool            { return s.hp > 0 }
