# Clean Code Regeln für Codera Battle

## Allgemein

Jede Funktion macht genau eine Sache.
Namen sind aussagekräftig und auf Englisch.
Keine verschachtelten Strukturen, lieber früh zurückkehren (early return).
Kommentare nur wo nötig, guter Code erklärt sich selbst.
Jede exportierte Funktion, jeder Struct und jedes Paket hat einen Godoc Kommentar, der mit dem Namen beginnt.

## Fehlerbehandlung

Keine ignorierten Errors, jeder Error wird geprüft oder weitergegeben.
Panic nur bei nicht behebbaren Fehlern, zum Beispiel beim Programmstart wenn die Datenbank nicht erreichbar ist.

## Struktur

Ein Paket pro Rolle unter hero/.
Öffentliche Funktionen und Structs beginnen mit Grossbuchstaben, interne Hilfsfunktionen bleiben klein geschrieben.
Keine Business Logik in main.go, main.go startet nur das Programm.

## Formatierung

Code wird immer mit gofmt formatiert bevor er committet wird.
golangci-lint run läuft ohne Fehler bevor ein Pull Request erstellt wird.

## Commits

Commit Nachrichten folgen Conventional Commits, also feat, fix, docs, ci und ähnliche Prefixe.
Jeder Commit auf main oder develop kompiliert fehlerfrei mit go build ./....