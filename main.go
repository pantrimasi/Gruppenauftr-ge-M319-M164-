package main

import (
	"fmt"
	"log"
	"os"

	"codera-battle/combat"
	"codera-battle/db"
	"codera-battle/dragon"
	"codera-battle/hero/arkan"
	"codera-battle/hero/druide"
	krieger "codera-battle/hero/funktions-krieger"
	"codera-battle/internal"
)

func main() {
	// Datenbankverbindung (Daten-Druide – Angelos)
	database, err := db.Connect()
	if err != nil {
		fmt.Printf("[DB] Warnung: %v\n[DB] Starte ohne Datenbankverbindung.\n", err)
	} else {
		fmt.Println("[DB] Verbindung zur codera-Datenbank erfolgreich.")

		if err := database.AutoMigrate(&db.Role{}, &db.Equipment{}, &db.Skill{}, &db.Hero{}); err != nil {
			log.Printf("[DB] AutoMigrate Warnung: %v", err)
		}

		if err := db.Seed(database); err != nil {
			log.Printf("[DB] Seed Warnung: %v", err)
		}

		// M164 – Lazar: Query-Ergebnisse validieren
		heroes, err := db.LoadHeroes(database)
		if err != nil {
			log.Printf("[DB] LoadHeroes Warnung: %v", err)
		} else {
			fmt.Println("\n[DB] Helden in der Datenbank:")
			fmt.Printf("  %-24s  %-10s  %5s  %6s  %7s  %5s\n",
				"Name", "Rolle", "MaxHP", "Attack", "Defense", "Speed")
			fmt.Println("  --------------------------------------------------------------")
			for _, h := range heroes {
				fmt.Printf("  %-24s  %-10s  %5d  %6d  %7d  %5d\n",
					h.Name, h.Role.Name, h.MaxHP, h.Attack, h.Defense, h.Speed)
			}
		}
		fmt.Println()
	}

	helden := buildHeroes()
	entropyDragon := dragon.New()

	combat.CombatLoop(helden, entropyDragon)
	os.Exit(0)
}

// buildHeroes assembles the full hero roster for combat.
func buildHeroes() []internal.Combatant {
	return []internal.Combatant{
		krieger.New(),                  // Lazar        (Funktions-Krieger*in)
		druide.NewDruide(),             // Angelos       (Daten-Druide)
		arkan.NewHero("Masato Kuster"), // Masato        (Arkan-Dokumentar*in)
	}
}
