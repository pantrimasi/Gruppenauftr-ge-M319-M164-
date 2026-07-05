package db

import "gorm.io/gorm"

// Seed fügt die Rollen, Ausrüstung, Skills und Helden ein.
// FirstOrCreate verhindert doppelte Einträge beim mehrmaligen Starten.
func Seed(database *gorm.DB) error {
	roles := []Role{
		{Name: "arkan"},
		{Name: "druide"},
		{Name: "krieger"},
	}

	for _, role := range roles {
		if err := database.FirstOrCreate(&role, Role{Name: role.Name}).Error; err != nil {
			return err
		}
	}

	equipment := []Equipment{
		{Name: "Pergament-Stab", Type: "weapon", AttackBonus: 8},
		{Name: "Runen-Gewand", Type: "armor", DefenseBonus: 5},
		{Name: "Tintenfass-Amulett", Type: "accessory", SpeedBonus: 3, HPBonus: 20},

		{Name: "Transformations-Kristall", Type: "weapon", AttackBonus: 6},
		{Name: "Datenstrom-Mantel", Type: "armor", DefenseBonus: 4},
		{Name: "Schema-Ring", Type: "accessory", SpeedBonus: 5, HPBonus: 10},

		{Name: "Funktions-Schwert", Type: "weapon", AttackBonus: 10},
		{Name: "Krieger-Rüstung", Type: "armor", DefenseBonus: 8},
		{Name: "Gurt der Ausdauer", Type: "accessory", SpeedBonus: 2, HPBonus: 40},
	}

	for _, item := range equipment {
		if err := database.FirstOrCreate(&item, Equipment{Name: item.Name}).Error; err != nil {
			return err
		}
	}

	var arkan Role
	var druide Role
	var krieger Role

	if err := database.Where("name = ?", "arkan").First(&arkan).Error; err != nil {
		return err
	}
	if err := database.Where("name = ?", "druide").First(&druide).Error; err != nil {
		return err
	}
	if err := database.Where("name = ?", "krieger").First(&krieger).Error; err != nil {
		return err
	}

	skills := []Skill{
		{RoleID: arkan.ID, Name: "Runen-Geschoss", Description: "Magischer Runenangriff", DamageMin: 12, DamageMax: 24, Accuracy: 0.90, TargetType: "single_enemy"},
		{RoleID: arkan.ID, Name: "Arkaner Bann", Description: "Schwacher Flächenangriff", DamageMin: 8, DamageMax: 16, Accuracy: 0.85, TargetType: "all_enemies"},
		{RoleID: arkan.ID, Name: "Klärende Annotation", Description: "Heilt einen Verbündeten", HealMin: 15, HealMax: 25, Accuracy: 1.00, TargetType: "single_ally"},

		{RoleID: druide.ID, Name: "Datenklinge", Description: "Transformierte Daten als Klinge", DamageMin: 10, DamageMax: 20, Accuracy: 0.85, TargetType: "single_enemy"},
		{RoleID: druide.ID, Name: "Strukturwandel", Description: "Hoher Schaden mit niedriger Genauigkeit", DamageMin: 14, DamageMax: 28, Accuracy: 0.70, TargetType: "single_enemy"},
		{RoleID: druide.ID, Name: "Transformative Regeneration", Description: "Heilt sich selbst", HealMin: 12, HealMax: 20, Accuracy: 1.00, TargetType: "self"},

		{RoleID: krieger.ID, Name: "Präziser Hieb", Description: "Kräftiger physischer Angriff", DamageMin: 18, DamageMax: 32, Accuracy: 0.80, TargetType: "single_enemy"},
		{RoleID: krieger.ID, Name: "Schutzschild", Description: "Erhöht eigene Defense um 5", Accuracy: 1.00, TargetType: "self"},
		{RoleID: krieger.ID, Name: "Kampfschrei", Description: "Angriff mit möglichem Buff", DamageMin: 8, DamageMax: 16, Accuracy: 0.90, TargetType: "single_enemy"},
	}

	for _, skill := range skills {
		if err := database.FirstOrCreate(&skill, Skill{Name: skill.Name}).Error; err != nil {
			return err
		}
	}

	var pergamentStab, runenGewand, tintenfass Equipment
	var kristall, mantel, ring Equipment
	var schwert, ruestung, gurt Equipment

	if err := database.Where("name = ?", "Pergament-Stab").First(&pergamentStab).Error; err != nil {
		return err
	}
	if err := database.Where("name = ?", "Runen-Gewand").First(&runenGewand).Error; err != nil {
		return err
	}
	if err := database.Where("name = ?", "Tintenfass-Amulett").First(&tintenfass).Error; err != nil {
		return err
	}

	if err := database.Where("name = ?", "Transformations-Kristall").First(&kristall).Error; err != nil {
		return err
	}
	if err := database.Where("name = ?", "Datenstrom-Mantel").First(&mantel).Error; err != nil {
		return err
	}
	if err := database.Where("name = ?", "Schema-Ring").First(&ring).Error; err != nil {
		return err
	}

	if err := database.Where("name = ?", "Funktions-Schwert").First(&schwert).Error; err != nil {
		return err
	}
	if err := database.Where("name = ?", "Krieger-Rüstung").First(&ruestung).Error; err != nil {
		return err
	}
	if err := database.Where("name = ?", "Gurt der Ausdauer").First(&gurt).Error; err != nil {
		return err
	}

	heroes := []Hero{
		{
			Name:                "Masato Kuster",
			RoleID:              arkan.ID,
			MaxHP:               120,
			CurrentHP:           120,
			Attack:              18,
			Defense:             8,
			Speed:               14,
			EquippedWeaponID:    pergamentStab.ID,
			EquippedArmorID:     runenGewand.ID,
			EquippedAccessoryID: tintenfass.ID,
		},
		{
			Name:                "Angelos Daroukakis",
			RoleID:              druide.ID,
			MaxHP:               100,
			CurrentHP:           100,
			Attack:              14,
			Defense:             10,
			Speed:               16,
			EquippedWeaponID:    kristall.ID,
			EquippedArmorID:     mantel.ID,
			EquippedAccessoryID: ring.ID,
		},
		{
			Name:                "Lazar Zajic",
			RoleID:              krieger.ID,
			MaxHP:               150,
			CurrentHP:           150,
			Attack:              22,
			Defense:             14,
			Speed:               8,
			EquippedWeaponID:    schwert.ID,
			EquippedArmorID:     ruestung.ID,
			EquippedAccessoryID: gurt.ID,
		},
	}

	for _, hero := range heroes {
		if err := database.FirstOrCreate(&hero, Hero{Name: hero.Name}).Error; err != nil {
			return err
		}
	}

	return nil
}
