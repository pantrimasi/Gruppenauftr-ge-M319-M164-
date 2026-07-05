package arkan

import "github.com/codera/battle/internal"

// Equipment represents an item a hero can wear.
type Equipment struct {
	Name         string
	Type         string
	AttackBonus  int
	DefenseBonus int
	SpeedBonus   int
	HPBonus      int
}

// Skill represents a usable ability.
type Skill struct {
	Name        string
	MinDamage   int
	MaxDamage   int
	Healing     int
	Accuracy    float64
	Target      string
	Description string
}

// Hero represents the Arkan-Dokumentar character.
type Hero struct {
	Name      string
	Stats     internal.Stats
	CurrentHP int
	Weapon    Equipment
	Armor     Equipment
	Accessory Equipment
	Skills    []Skill
}

// NewHero creates a new Arkan-Dokumentar hero.
func NewHero(name string) *Hero {
	// base stats
	stats := internal.Stats{
		MaxHP:   120,
		Attack:  18,
		Defense: 8,
		Speed:   14,
	}

	weapon := Equipment{Name: "Pergament-Stab", Type: "weapon", AttackBonus: 8}
	armor := Equipment{Name: "Runen-Gewand", Type: "armor", DefenseBonus: 5}
	accessory := Equipment{Name: "Tintenfass-Amulett", Type: "accessory", SpeedBonus: 3, HPBonus: 20}

	skills := []Skill{
		{Name: "Runen-Geschoss", MinDamage: 12, MaxDamage: 24, Accuracy: 0.90, Target: "single_enemy", Description: "Magischer Runenangriff"},
		{Name: "Arkaner Bann", MinDamage: 8, MaxDamage: 16, Accuracy: 0.85, Target: "all_enemies", Description: "Schwacher Flaechenangriff"},
		{Name: "Klaerende Annotation", Healing: 25, Accuracy: 1.0, Target: "single_ally", Description: "Heilt einen Verbuendeten"},
	}

	// apply equipment bonus
	totalHP := stats.MaxHP + accessory.HPBonus

	return &Hero{
		Name:      name,
		Stats:     stats,
		CurrentHP: totalHP,
		Weapon:    weapon,
		Armor:     armor,
		Accessory: accessory,
		Skills:    skills,
	}
}

// GetName returns the hero name.
func (h *Hero) GetName() string {
	return h.Name
}

// GetStats returns the hero stats.
func (h *Hero) GetStats() internal.Stats {
	return h.Stats
}

// GetCurrentHP returns current HP.
func (h *Hero) GetCurrentHP() int {
	return h.CurrentHP
}

// SetCurrentHP sets current HP.
func (h *Hero) SetCurrentHP(hp int) {
	h.CurrentHP = hp
}

// GetMaxHP returns max HP including equipment bonus.
func (h *Hero) GetMaxHP() int {
	return h.Stats.MaxHP + h.Accessory.HPBonus
}

// IsAlive checks if hero is still alive.
func (h *Hero) IsAlive() bool {
	return h.CurrentHP > 0
}

// ChooseAction decides which skill to use based on ally HP
func (h *Hero) ChooseAction(allies []internal.Combatant) Skill {
	weakest := findWeakestAlly(allies)

	// heal check
	if weakest != nil && isCritical(weakest) {
		return h.Skills[2]
	}

	return h.Skills[0]
}

// findWeakestAlly returns the ally with lowest HP
func findWeakestAlly(allies []internal.Combatant) internal.Combatant {
	var weakest internal.Combatant
	lowest := -1

	for _, ally := range allies {
		if !ally.IsAlive() {
			continue
		}
		if lowest == -1 || ally.GetCurrentHP() < lowest {
			lowest = ally.GetCurrentHP()
			weakest = ally
		}
	}

	return weakest
}

// isCritical checks if HP is below 30 percent
func isCritical(c internal.Combatant) bool {
	threshold := float64(c.GetMaxHP()) * 0.3
	return float64(c.GetCurrentHP()) < threshold
}
