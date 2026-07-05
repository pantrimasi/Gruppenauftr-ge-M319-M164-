# Codera Battle

*Ein rundenbasiertes Kampfsystem in Go, entwickelt als Abschlussauftrag für die Module M319 und M164.*

## Überblick

Codera Battle ist ein Kommandozeilen Spiel, bei dem eine Gruppe von Helden gemeinsam gegen den Entropie Drachen antritt. Jedes Gruppenmitglied hat einen eigenen Helden Charakter mit individuellen Stats, Ausrüstung und Skills implementiert. Die Helden werden aus einer PostgreSQL Datenbank geladen, der Kampf selbst läuft rundenbasiert ab und die Reihenfolge der Züge richtet sich nach dem Speed Wert jedes Kämpfers. Schaden, Genauigkeit und kritische Treffer werden über Zufallszahlen berechnet.

Der Drache ist vollständig vorgegeben und wird nicht verändert, die Helden Charaktere sowie die Kampf Logik wurden von der Gruppe selbst entwickelt.

## Rollen und Zuständigkeiten

| Person | Rolle | Zuständigkeit |
|---|---|---|
| Masato | Arkan-Dokumentar | Codequalität, C4 Diagramme, Linter, eigener Charakter |
| Angelos | Daten-Druide | GORM Modelle, Datenbankanbindung, eigener Charakter |
| Lazar | Funktions-Krieger | Kampf Loop, Goroutines, eigener Charakter |

## Installation

Voraussetzung ist eine installierte Go Version sowie eine laufende PostgreSQL Instanz, entweder lokal oder über Docker.

Repository klonen:

    git clone https://github.com/AngelosDaroukakis/Gruppenauftr-ge-M319-M164-.git
    cd Gruppenauftr-ge-M319-M164-

Abhängigkeiten installieren:

    go mod tidy

## Konfiguration

Die Datei .env-example im Root kopieren und in .env umbenennen, danach die eigenen Datenbank Zugangsdaten eintragen.

    DB_HOST=localhost
    DB_PORT=5432
    DB_USER=codera
    DB_PASSWORD=codera
    DB_NAME=codera

## Programm starten

    go run main.go

Das Programm lädt beim Start automatisch die Helden aus der Datenbank, initialisiert den Drachen und startet den Kampf in der Kommandozeile.

## Tests ausführen

    go test ./...

## Linter Setup

Jede Person im Team installiert den Linter lokal auf dem eigenen Rechner.

Modul Support aktivieren:

    set GO111MODULE=on

Linter installieren:

    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

Linter ausführen, im Root des Projekts:

    golangci-lint run

Die genauen Regeln sind in der Datei .golangci.yml hinterlegt.

## Branching Strategie

Wir arbeiten nach GitHub Flow, erweitert um einen zusätzlichen develop Branch. Jede Person hat ihren eigenen Feature Branch, fertige Arbeit wird über einen Pull Request zuerst nach develop gemerged. Erst wenn develop stabil ist, geht es gemeinsam nach main. Direktes Pushen auf main oder develop ist über Branch Protection Rules gesperrt, jeder Merge braucht mindestens eine Freigabe.

## Clean Code

Wer am Projekt mitarbeitet, folgt diesen Regeln:

- Immer auf dem eigenen Feature Branch arbeiten
- Commit Nachrichten nach Conventional Commits, zum Beispiel feat, fix, docs
- Jeder Commit auf main oder develop muss mit go build ./... fehlerfrei kompilieren
- Vor jedem Pull Request den Linter laufen lassen
- Clean Code Regeln aus .opencode/rules.md beachten

## Projektstruktur

    main.go            Startpunkt des Programms
    combat/            Kampf Loop und Regeln
    dragon/            Vorgegebener Drache, wird nicht verändert
    hero/              Helden Charaktere, ein Paket pro Rolle
    internal/          Combatant Interface und gemeinsame Typen
    db/                GORM Modelle, Datenbankverbindung und Seeds

## Lizenz und Credits

Dieses Projekt entstand im Rahmen der Ausbildung als Schulauftrag und dient ausschliesslich Lernzwecken.

Entwickelt von Masato, Angelos und Lazar.