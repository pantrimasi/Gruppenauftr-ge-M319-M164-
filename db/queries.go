package db

import "gorm.io/gorm"

// LoadHeroes lädt alle Helden mit Rolle und Ausrüstung aus der Datenbank.
func LoadHeroes(database *gorm.DB) ([]Hero, error) {
	var heroes []Hero

	err := database.
		Preload("Role").
		Preload("Weapon").
		Preload("Armor").
		Preload("Accessory").
		Find(&heroes).Error

	return heroes, err
}

// LoadSkillsByRole lädt alle Skills von einer bestimmten Rolle.
func LoadSkillsByRole(database *gorm.DB, roleName string) ([]Skill, error) {
	var role Role
	if err := database.Where("name = ?", roleName).First(&role).Error; err != nil {
		return nil, err
	}

	var skills []Skill
	err := database.Where("\"Rolle\" = ?", role.ID).Find(&skills).Error
	return skills, err
}
