package main

import (
	"fmt"
	"log"

	"codera-battle/db"
	"codera-battle/hero/druide"
)

func main() {
	database, err := db.Connect()
	if err != nil {
		log.Fatal("Datenbankverbindung fehlgeschlagen: ", err)
	}

	err = database.AutoMigrate(
		&db.Role{},
		&db.Equipment{},
		&db.Skill{},
		&db.Hero{},
	)
	if err != nil {
		log.Fatal("AutoMigration fehlgeschlagen: ", err)
	}

	if err := db.Seed(database); err != nil {
		log.Fatal("Seed-Daten konnten nicht eingefügt werden: ", err)
	}

	heroes, err := db.LoadHeroes(database)
	if err != nil {
		log.Fatal("Helden konnten nicht geladen werden: ", err)
	}

	fmt.Println("Datenbank verbunden.")
	fmt.Println("Geladene Helden:")
	for _, hero := range heroes {
		fmt.Printf("- %s (%s) HP: %d/%d\n", hero.Name, hero.Role.Name, hero.CurrentHP, hero.MaxHP)
	}

	myHero := druide.NewDruide()
	fmt.Printf("Eigener Charakter im Code: %s HP: %d/%d\n", myHero.GetName(), myHero.GetCurrentHP(), myHero.GetMaxHP())
}
