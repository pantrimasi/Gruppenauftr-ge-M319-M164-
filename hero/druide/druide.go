package druide

import "codera-battle/internal"

// Druide ist mein eigener Held für die Rolle Daten-Druide.
type Druide struct {
	Name      string
	Stats     internal.Stats
	CurrentHP int
}

// NewDruide erstellt den Daten-Druiden Angelos Daroukakis.
func NewDruide() *Druide {
	return &Druide{
		Name: "Angelos Daroukakis",
		Stats: internal.Stats{
			MaxHP:   100,
			Attack:  14,
			Defense: 10,
			Speed:   16,
		},
		CurrentHP: 100,
	}
}

func (d *Druide) GetName() string {
	return d.Name
}

func (d *Druide) GetStats() internal.Stats {
	return d.Stats
}

func (d *Druide) GetCurrentHP() int {
	return d.CurrentHP
}

func (d *Druide) SetCurrentHP(hp int) {
	if hp < 0 {
		hp = 0
	}
	d.CurrentHP = hp
}

func (d *Druide) GetMaxHP() int {
	return d.Stats.MaxHP
}

func (d *Druide) IsAlive() bool {
	return d.CurrentHP > 0
}
