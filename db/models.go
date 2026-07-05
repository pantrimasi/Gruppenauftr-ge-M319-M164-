package db

// Role bildet die Tabelle Rolle ab.
type Role struct {
	ID   uint   `gorm:"column:ID;primaryKey"`
	Name string `gorm:"column:name;size:100;unique;not null"`
}

func (Role) TableName() string {
	return "Rolle"
}

// Equipment bildet die Tabelle Ausrüstung ab.
type Equipment struct {
	ID           uint   `gorm:"column:ID;primaryKey"`
	Name         string `gorm:"column:name;size:100;unique;not null"`
	Type         string `gorm:"column:type;size:100;not null"`
	AttackBonus  int    `gorm:"column:attack_bonus;default:0"`
	DefenseBonus int    `gorm:"column:defense_bonus;default:0"`
	SpeedBonus   int    `gorm:"column:speed_bonus;default:0"`
	HPBonus      int    `gorm:"column:hp_bonus;default:0"`
}

func (Equipment) TableName() string {
	return "Ausrüstung"
}

// Skill bildet die Tabelle Skills ab.
type Skill struct {
	ID          uint    `gorm:"column:ID;primaryKey"`
	RoleID      uint    `gorm:"column:Rolle;not null"`
	Name        string  `gorm:"column:name;size:100;not null"`
	Description string  `gorm:"column:description;size:200"`
	DamageMin   int     `gorm:"column:damage_min;default:0"`
	DamageMax   int     `gorm:"column:damage_max;default:0"`
	HealMin     int     `gorm:"column:heal_min;default:0"`
	HealMax     int     `gorm:"column:heal_max;default:0"`
	Accuracy    float64 `gorm:"column:accuracy;type:numeric(3,2);not null"`
	TargetType  string  `gorm:"column:target_type;size:100;not null"`

	Role Role `gorm:"foreignKey:RoleID;references:ID"`
}

func (Skill) TableName() string {
	return "Skills"
}

// Hero bildet die Tabelle Helden ab.
type Hero struct {
	Name                string `gorm:"column:Name;size:100;primaryKey"`
	RoleID              uint   `gorm:"column:Rolle;not null"`
	MaxHP               int    `gorm:"column:max_hp;not null"`
	CurrentHP           int    `gorm:"column:current_hp;not null"`
	Attack              int    `gorm:"column:attack;not null"`
	Defense             int    `gorm:"column:defense;not null"`
	Speed               int    `gorm:"column:speed;not null"`
	EquippedWeaponID    uint   `gorm:"column:equipped_weapon"`
	EquippedArmorID     uint   `gorm:"column:equipped_armor"`
	EquippedAccessoryID uint   `gorm:"column:equipped_accessory"`

	Role      Role      `gorm:"foreignKey:RoleID;references:ID"`
	Weapon    Equipment `gorm:"foreignKey:EquippedWeaponID;references:ID"`
	Armor     Equipment `gorm:"foreignKey:EquippedArmorID;references:ID"`
	Accessory Equipment `gorm:"foreignKey:EquippedAccessoryID;references:ID"`
}

func (Hero) TableName() string {
	return "Helden"
}
