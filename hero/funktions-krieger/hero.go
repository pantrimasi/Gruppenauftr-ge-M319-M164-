// Package krieger implements the Funktions-Krieger*in hero character for the
// Codera Battle system. The warrior's identity is "Praezision der Funktionen".
//
// # Character Design
//
// Base stats (before equipment): HP 150, ATK 22, DEF 14, SPD 8.
// After applying all three pieces of equipment the effective values are:
// MaxHP 190, ATK 32, DEF 22, SPD 10.
//
// # Concurrency
//
// All HP and buff fields are protected by a sync.Mutex so the struct is safe
// to use from multiple goroutines. The Double Strike feature (see [Krieger.doubleStrike])
// intentionally exercises this by computing two hits in parallel.
//
// # Round-based Buff System
//
//   - Schutzschild adds +5 DEF for the remainder of the current round.
//     The bonus is cleared by [Krieger.OnTurnStart] at the top of the warrior's
//     next turn, after the dragon has already had a chance to attack into it.
//   - Kampfschrei queues +5 ATK (pendingAtkBonus). OnTurnStart promotes the
//     pending value to activeAtkBonus so it applies from the next turn onward.
//
// # Auto-Defense
//
// When the warrior's HP drops below 30 % of MaxHP, the combat loop calls
// no user prompt and selects Schutzschild automatically. This keeps the hero
// alive long enough to land a final decisive blow.
package krieger

import (
	"sync"

	"codera-battle/internal"
)

// equipment holds the stat bonuses granted by one piece of gear.
type equipment struct {
	name       string
	itemType   string // "weapon" | "armor" | "accessory"
	bonusATK   int
	bonusDEF   int
	bonusSPD   int
	bonusMaxHP int
}

// loadout returns the three pieces of gear worn by the Funktions-Krieger.
func loadout() [3]equipment {
	return [3]equipment{
		{name: "Funktions-Schwert", itemType: "weapon", bonusATK: 10},
		{name: "Krieger-Ruestung", itemType: "armor", bonusDEF: 8},
		{name: "Gurt der Ausdauer", itemType: "accessory", bonusSPD: 2, bonusMaxHP: 40},
	}
}

// Krieger is the Funktions-Krieger*in hero.
// Create instances with [New]; do not construct directly.
type Krieger struct {
	mu sync.Mutex

	name      string
	maxHP     int
	currentHP int
	baseStats internal.Stats // stats without equipment or buffs

	// round-based buff tracking
	tempDefBonus    int // active until OnTurnStart resets it (Schutzschild)
	activeAtkBonus  int // applied this turn (promoted from pending)
	pendingAtkBonus int // promoted to active on next OnTurnStart (Kampfschrei)
}

// New returns a fully equipped Lazar (Funktions-Krieger*in) ready for battle.
// Equipment bonuses are baked into the returned struct's effective stats.
func New() *Krieger {
	base := internal.Stats{
		MaxHP:   150,
		Attack:  22,
		Defense: 14,
		Speed:   8,
	}

	maxHP := base.MaxHP
	atk := base.Attack
	def := base.Defense
	spd := base.Speed

	for _, e := range loadout() {
		maxHP += e.bonusMaxHP
		atk += e.bonusATK
		def += e.bonusDEF
		spd += e.bonusSPD
	}

	return &Krieger{
		name:      "Lazar (Funktions-Krieger*in)",
		maxHP:     maxHP,
		currentHP: maxHP,
		baseStats: internal.Stats{
			MaxHP:   maxHP,
			Attack:  atk,
			Defense: def,
			Speed:   spd,
		},
	}
}

// GetName implements [internal.Combatant].
func (k *Krieger) GetName() string { return k.name }

// GetStats implements [internal.Combatant].
// Returns current stats including equipment bonuses and any active round buffs.
func (k *Krieger) GetStats() internal.Stats {
	k.mu.Lock()
	defer k.mu.Unlock()
	return internal.Stats{
		MaxHP:   k.baseStats.MaxHP,
		Attack:  k.baseStats.Attack + k.activeAtkBonus,
		Defense: k.baseStats.Defense + k.tempDefBonus,
		Speed:   k.baseStats.Speed,
	}
}

// GetCurrentHP implements [internal.Combatant].
func (k *Krieger) GetCurrentHP() int {
	k.mu.Lock()
	defer k.mu.Unlock()
	return k.currentHP
}

// SetCurrentHP implements [internal.Combatant]. Clamps hp to [0, MaxHP].
func (k *Krieger) SetCurrentHP(hp int) {
	k.mu.Lock()
	defer k.mu.Unlock()
	switch {
	case hp < 0:
		k.currentHP = 0
	case hp > k.maxHP:
		k.currentHP = k.maxHP
	default:
		k.currentHP = hp
	}
}

// GetMaxHP implements [internal.Combatant].
func (k *Krieger) GetMaxHP() int { return k.maxHP }

// IsAlive implements [internal.Combatant].
func (k *Krieger) IsAlive() bool {
	k.mu.Lock()
	defer k.mu.Unlock()
	return k.currentHP > 0
}

// OnTurnStart implements [internal.SkillUser].
// Must be called at the beginning of every hero turn before any other action.
// It promotes the pending ATK bonus (from a previous Kampfschrei) to active,
// and clears the temporary DEF bonus (from a previous Schutzschild).
func (k *Krieger) OnTurnStart() {
	k.mu.Lock()
	defer k.mu.Unlock()
	k.activeAtkBonus = k.pendingAtkBonus
	k.pendingAtkBonus = 0
	k.tempDefBonus = 0
}

// GetSkills implements [internal.SkillUser].
// Returns the three role-specific skills of the Funktions-Krieger*in.
func (k *Krieger) GetSkills() []internal.Skill {
	return []internal.Skill{
		{
			Name:        "Praeziser Hieb",
			Description: "Kraeftiger physischer Angriff  (18-32 DMG, 80%)",
			DamageMin:   18,
			DamageMax:   32,
			Accuracy:    0.80,
			Target:      internal.TargetSingleEnemy,
		},
		{
			Name:        "Schutzschild",
			Description: "Erhoet eigene DEF um 5 fuer diese Runde  (100%, self)",
			DamageMin:   0,
			DamageMax:   0,
			Accuracy:    1.0,
			Target:      internal.TargetSelf,
		},
		{
			Name:        "Kampfschrei",
			Description: "Schwaecher Angriff + ATK+5 naechste Runde  (8-16 DMG, 90%)",
			DamageMin:   8,
			DamageMax:   16,
			Accuracy:    0.90,
			Target:      internal.TargetSingleEnemy,
		},
	}
}

// ApplyTempDefBonus activates the Schutzschild defense boost for this round.
func (k *Krieger) ApplyTempDefBonus(amount int) {
	k.mu.Lock()
	defer k.mu.Unlock()
	k.tempDefBonus += amount
}

// QueueAtkBonus queues an ATK bonus that becomes active on the next turn (Kampfschrei).
func (k *Krieger) QueueAtkBonus(amount int) {
	k.mu.Lock()
	defer k.mu.Unlock()
	k.pendingAtkBonus += amount
}

// IsAutoDefenseRequired returns true when HP has dropped below 30 % of MaxHP.
// The combat loop uses this to skip the player menu and auto-select Schutzschild.
func (k *Krieger) IsAutoDefenseRequired() bool {
	k.mu.Lock()
	defer k.mu.Unlock()
	return float64(k.currentHP)/float64(k.maxHP) < 0.30
}
