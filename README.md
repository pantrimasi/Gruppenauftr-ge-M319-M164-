# Codera Battle – Gruppenauftrag (M319 / M164)

## Codera Battle – rundenbasiertes Kampfsystem mit Datenbankanbindung

## Kurzbeschreibung
*Codera Battle ist ein Schulprojekt, das die Module M319 (Programmieren in Go) und M164 (Datenbanken) kombiniert. Entwickelt wird ein CLI-basiertes, rundenbasiertes Kampfsystem mit PostgreSQL-Datenbank via Docker.*

---

## 1. Projektübersicht

Die Gruppe entwickelt gemeinsam ein rundenbasiertes Kampfspiel, das in der Kommandozeile läuft. Jeder Charakter besitzt eigene Stats, Ausrüstung und Skills und wird in einer PostgreSQL-Datenbank gespeichert, die lokal über Docker betrieben wird.

Alle Spieler treten gemeinsam gegen den Entropie-Drachen an, der fix vorgegeben ist und nicht verändert werden darf. Die Zugreihenfolge richtet sich nach dem Speed-Wert. Schaden, Trefferchance und kritische Treffer werden zufällig berechnet.

---

## 2. Rollen und Aufgaben

### Masato – Arkan-Dokumentar
Verantwortlich für die technische Struktur und Qualität des gesamten Projekts.

Aufgaben:
- Einrichtung des Git-Repositories
- Definition der Branching-Strategie
- Sicherstellen, dass alle den Linter verwenden
- Erstellung der C4-Architekturdiagramme
- Dokumentation des gesamten Systems
- Überprüfung von Codequalität und Struktur
- Implementierung eines eigenen Magier-Charakters

---

### Angelos – Daten-Druide
Verantwortlich für die gesamte Datenbank- und Go-Datenebene.

Aufgaben:
- Erstellung der GORM-Modelle
- Aufbau der Datenbankverbindung in Go
- Entwicklung der Seed-Daten
- Setup eines PostgreSQL-Containers via Docker
- Erstellung und Dokumentation des ERD
- Sicherstellen, dass andere die Datenbank reproduzieren können

---

### Lazar – Funktions-Krieger
Verantwortlich für die Kernlogik des Spiels.

Aufgaben:
- Implementierung des Kampf-Systems (Combat Loop)
- Steuerung der Spiel-Logik im CLI
- Integration von Zufallssystemen (Damage, Hit-Chance, Crits)
- Einsatz von Goroutines für parallele Abläufe
- Absicherung gemeinsamer Ressourcen mit Mutex
- Testen der Datenbankverbindung
- Validierung der Datenbankabfragen

---

## 3. Spielsystem (M319)

Das Kampfsystem basiert auf rundenbasierten Kämpfen in der Kommandozeile.

Wichtige Mechaniken:
- Speed-Wert bestimmt Zugreihenfolge
- RNG steuert Schaden, Treffer und kritische Treffer
- Jeder Spieler hat einen eigenen Helden als separates Go-Paket
- Der Entropie-Drache ist der finale Gegner und fix vorgegeben

---

## 4. Datenbanksystem (M164)

Jede Person betreibt eine lokale PostgreSQL-Instanz via Docker.

Gemeinsam werden folgende Schritte umgesetzt:
- Erstellung eines ERD (Entity Relationship Diagram)
- Umsetzung in ein relationales Schema
- Erstellung der Tabellen via DDL
- Befüllung mit Seed-Daten
- Durchführung von JOIN- und Filterabfragen
- Export und Reimport als SQL-Dump

---

## 5. Abgabeanforderungen

Die Abgabe erfolgt über ein Git-Repository mit vollständiger Historie.

Erforderlich sind:
- Vollständige Git-History mit nachvollziehbaren Beiträgen
- C4-Diagramme
- Gruppendokumentation
- Unit Tests
- Godoc-Dokumentation
- SQL-Skripte (DDL, Seed, Queries, Dump)
- Saubere Projektstruktur

Fehlende oder nicht nachvollziehbare Beiträge können zu Notenabzug für die gesamte Gruppe führen.

## 6. Linter Setup

### Voraussetzung

Go muss installiert sein und im PATH verfügbar sein.

### Installation

Modul Support aktivieren:

    set GO111MODULE=on

Linter installieren:

    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

Installation prüfen:

    golangci-lint --version

### Verwendung

Im Root des Projekts ausführen, dort wo auch die go.mod liegt:

    golangci-lint run

Der Linter zeigt danach alle Findings basierend auf der Konfiguration in .golangci.yml an.

### Häufiger Fehler

Falls beim Installieren die Meldung "modules disabled by GO111MODULE=off" erscheint, muss GO111MODULE zuerst wie oben beschrieben auf on gesetzt werden, bevor der Install Befehl erneut ausgeführt wird.