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

// TargetType beschreibt, wen ein Skill treffen kann.
type TargetType string

const (
	TargetSingleEnemy TargetType = "single_enemy"
	TargetAllEnemies  TargetType = "all_enemies"
	TargetSelf        TargetType = "self"
	TargetSingleAlly  TargetType = "single_ally"
	TargetAllAllies   TargetType = "all_allies"
)

// Skill beschreibt eine Kampffähigkeit.
type Skill struct {
	Name        string
	Description string
	DamageMin   int
	DamageMax   int
	Healing     int
	Accuracy    float64
	Target      TargetType
}

// SkillUser wird von Helden implementiert, die aktive Skills besitzen.
type SkillUser interface {
	GetSkills() []Skill
	OnTurnStart()
}
