package internal

// Stats speichert die Kampfwerte von einem Helden oder Gegner.
type Stats struct {
	MaxHP   int
	Attack  int
	Defense int
	Speed   int
}

// Combatant beschreibt alles, was am Kampf teilnehmen kann.
type Combatant interface {
	GetName() string
	GetStats() Stats
	GetCurrentHP() int
	SetCurrentHP(hp int)
	GetMaxHP() int
	IsAlive() bool
}
