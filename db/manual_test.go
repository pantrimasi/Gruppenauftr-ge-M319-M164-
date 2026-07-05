package db

import (
	"fmt"
	"os"
	"testing"
)

// set test env
func setTestEnv() {
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "codera")
	os.Setenv("DB_PASSWORD", "codera")
	os.Setenv("DB_NAME", "codera")
}

// manual check load heroes
func TestManualLoadHeroes(t *testing.T) {
	setTestEnv()

	database, err := Connect()
	if err != nil {
		t.Fatal(err)
	}

	heroes, err := LoadHeroes(database)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(heroes)
}

// manual check load skills
func TestManualLoadSkillsByRole(t *testing.T) {
	setTestEnv()

	database, err := Connect()
	if err != nil {
		t.Fatal(err)
	}

	skills, err := LoadSkillsByRole(database, "krieger")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(skills)
}
