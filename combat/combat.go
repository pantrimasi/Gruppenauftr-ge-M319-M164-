// Package combat contains the main battle loop and all combat mechanics.
//
// Responsibilities of the Funktions-Krieger*in (Lazar):
//   - [CombatLoop]: round-based battle management
//   - [processHeroTurn]: CLI interaction, auto-defense, skill execution
//   - [doubleStrike]: parallel damage via Goroutines + Mutex (optional feature)
//   - [processDragonTurn]: dragon AI using ChooseAction, rage bonus, mutex safety
package combat

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"

	"codera-battle/dragon"
	"codera-battle/internal"
)

// ── ANSI colors ──────────────────────────────────────────────────────────────

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorBold   = "\033[1m"
	colorDim    = "\033[2m"
)

func red(s string) string    { return colorRed + s + colorReset }
func green(s string) string  { return colorGreen + s + colorReset }
func yellow(s string) string { return colorYellow + s + colorReset }
func cyan(s string) string   { return colorCyan + s + colorReset }
func bold(s string) string   { return colorBold + s + colorReset }
func dim(s string) string    { return colorDim + s + colorReset }

// ── Layout constants ─────────────────────────────────────────────────────────

const (
	width  = 56
	barLen = 22
)

var (
	divider = "+" + strings.Repeat("=", width) + "+"
	line    = "+" + strings.Repeat("-", width) + "+"
)

// ── Types ─────────────────────────────────────────────────────────────────────

// ActionResult records everything that happened during one combat action.
type ActionResult struct {
	ActorName  string
	SkillName  string
	TargetName string
	Damage     int
	Healing    int
	IsCrit     bool
	IsMiss     bool
	IsAOE      bool
}

// CombatantInfo wraps a Combatant with its role in the initiative order.
type CombatantInfo struct {
	Combatant internal.Combatant
	IsDragon  bool
}

// CalculateDamage computes final damage including stat scaling, defense,
// accuracy, and critical hit chance.
//
// This function is provided by the assignment. Do not modify.
func CalculateDamage(baseMin, baseMax int, attackerStat, defenderDef int, accuracy float64) (int, bool, bool) {
	if rand.Float64() > accuracy {
		return 0, false, true
	}

	baseDamage := rand.Intn(baseMax-baseMin+1) + baseMin
	attackMultiplier := 1.0 + float64(attackerStat)/20.0
	defenseReduction := 1.0 - float64(defenderDef)/100.0
	if defenseReduction < 0.1 {
		defenseReduction = 0.1
	}

	finalDamage := int(float64(baseDamage) * attackMultiplier * defenseReduction)
	if finalDamage < 1 {
		finalDamage = 1
	}

	isCrit := rand.Float64() < 0.10
	if isCrit {
		finalDamage = int(float64(finalDamage) * 1.5)
	}

	return finalDamage, isCrit, false
}

// ── Entry point ───────────────────────────────────────────────────────────────

// CombatLoop runs the round-based battle until the dragon dies or all heroes fall.
func CombatLoop(heroes []internal.Combatant, d *dragon.EntropyDragon) {
	printIntroScreen(heroes, d)

	participants := buildInitiativeOrder(heroes, d)
	round := 1
	prevDragonEnraged := false

	for {
		if !d.IsAlive() {
			clearScreen()
			printHeader()
			printFullStatus(heroes, d, round)
			fmt.Println()
			printVictoryScreen(heroes)
			return
		}
		if allDead(heroes) {
			clearScreen()
			printHeader()
			printFullStatus(heroes, d, round)
			fmt.Println()
			printDefeatScreen()
			return
		}

		// Announce rage activation once
		if d.IsEnraged && !prevDragonEnraged {
			prevDragonEnraged = true
			printRageWarning(d)
		}

		for _, p := range participants {
			if !p.Combatant.IsAlive() {
				continue
			}

			var result ActionResult
			if p.IsDragon {
				result = processDragonTurn(d, heroes, round)
			} else {
				result = processHeroTurn(p.Combatant, heroes, d, round)
			}

			// Check for hero deaths after each action
			announceDeaths(heroes)

			// TODO: Logging jeder Kampfaktion (Code-Kleriker*in)
			_ = result

			if !d.IsAlive() || allDead(heroes) {
				break
			}
		}

		round++
	}
}

// buildInitiativeOrder sorts all fighters by Speed descending.
// On a tie the hero goes before the dragon.
func buildInitiativeOrder(heroes []internal.Combatant, d *dragon.EntropyDragon) []CombatantInfo {
	participants := make([]CombatantInfo, 0, len(heroes)+1)
	for _, h := range heroes {
		if h.IsAlive() {
			participants = append(participants, CombatantInfo{Combatant: h})
		}
	}
	participants = append(participants, CombatantInfo{Combatant: d, IsDragon: true})

	sort.SliceStable(participants, func(i, j int) bool {
		si := participants[i].Combatant.GetStats().Speed
		sj := participants[j].Combatant.GetStats().Speed
		if si == sj {
			return !participants[i].IsDragon && participants[j].IsDragon
		}
		return si > sj
	})

	return participants
}

// ── Hero turn ─────────────────────────────────────────────────────────────────

// processHeroTurn renders the full UI, handles auto-defense and skill selection.
func processHeroTurn(hero internal.Combatant, allies []internal.Combatant, d *dragon.EntropyDragon, round int) ActionResult {
	su, hasSkills := hero.(internal.SkillUser)
	if hasSkills {
		su.OnTurnStart()
	}

	clearScreen()
	printHeader()
	printFullStatus(allies, d, round)
	fmt.Printf("\n  %s\n", dim(separator()))
	fmt.Printf("  %s\n", bold("Zug: "+hero.GetName()))
	fmt.Printf("  %s\n", dim(separator()))

	if !hasSkills {
		result := autoAttack(hero, d)
		fmt.Println()
		printActionResult(result)
		waitEnter()
		return result
	}

	skills := su.GetSkills()

	// Auto-defense when HP is critical
	if autoDefenser, ok := hero.(interface{ IsAutoDefenseRequired() bool }); ok {
		if autoDefenser.IsAutoDefenseRequired() {
			fmt.Printf("\n  %s\n", yellow("[AUTO-VERTEIDIGUNG] HP unter 30% — Schutzschild aktiviert!"))
			result := executeSkill(hero, skills[1], d, allies)
			fmt.Println()
			printActionResult(result)
			waitEnter()
			return result
		}
	}

	printSkillMenu(skills)
	choice := readInt(1, len(skills))
	result := executeSkill(hero, skills[choice-1], d, allies)

	clearScreen()
	printHeader()
	printFullStatus(allies, d, round)
	fmt.Println()
	printActionResult(result)
	waitEnter()

	return result
}

// executeSkill applies the chosen skill effect.
func executeSkill(hero internal.Combatant, skill internal.Skill, d *dragon.EntropyDragon, allies []internal.Combatant) ActionResult {
	name := hero.GetName()
	stats := hero.GetStats()

	switch skill.Target {
	case internal.TargetSelf:
		if defender, ok := hero.(interface{ ApplyTempDefBonus(int) }); ok {
			defender.ApplyTempDefBonus(5)
		}
		return ActionResult{ActorName: name, SkillName: skill.Name, TargetName: name}

	case internal.TargetSingleEnemy:
		if skill.Name == "Praeziser Hieb" && float64(d.GetCurrentHP())/float64(d.GetMaxHP()) < 0.30 {
			return doubleStrike(hero, skill, d)
		}

		dmg, isCrit, isMiss := CalculateDamage(skill.DamageMin, skill.DamageMax, stats.Attack, d.Defense, skill.Accuracy)
		if isMiss {
			return ActionResult{ActorName: name, SkillName: skill.Name, TargetName: d.GetName(), IsMiss: true}
		}
		d.TakeDamage(dmg)

		if skill.Name == "Kampfschrei" {
			if queuer, ok := hero.(interface{ QueueAtkBonus(int) }); ok {
				queuer.QueueAtkBonus(5)
			}
		}
		return ActionResult{ActorName: name, SkillName: skill.Name, TargetName: d.GetName(), Damage: dmg, IsCrit: isCrit}

	default:
		return autoAttack(hero, d)
	}
}

// doubleStrike fires two goroutine-parallel hits when dragon is below 30% HP.
//
// sync.WaitGroup waits for both goroutines to finish calculating damage before
// any HP is modified. sync.Mutex serialises the writes to the dragon's HP.
func doubleStrike(hero internal.Combatant, skill internal.Skill, d *dragon.EntropyDragon) ActionResult {
	name := hero.GetName()
	stats := hero.GetStats()

	type hit struct {
		dmg    int
		isCrit bool
		isMiss bool
	}

	results := make([]hit, 2)
	var wg sync.WaitGroup
	var mu sync.Mutex

	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			dmg, isCrit, isMiss := CalculateDamage(skill.DamageMin, skill.DamageMax, stats.Attack, d.Defense, skill.Accuracy)
			mu.Lock()
			results[idx] = hit{dmg: dmg, isCrit: isCrit, isMiss: isMiss}
			mu.Unlock()
		}(i)
	}
	wg.Wait()

	totalDamage := 0
	for _, h := range results {
		if !h.isMiss {
			mu.Lock()
			d.TakeDamage(h.dmg)
			totalDamage += h.dmg
			mu.Unlock()
		}
	}

	return ActionResult{
		ActorName:  name,
		SkillName:  "Doppel-Angriff (Praeziser Hieb x2)",
		TargetName: d.GetName(),
		Damage:     totalDamage,
	}
}

// autoAttack is the fallback for heroes without a skill list.
func autoAttack(hero internal.Combatant, d *dragon.EntropyDragon) ActionResult {
	stats := hero.GetStats()
	dmg, isCrit, isMiss := CalculateDamage(10, 20, stats.Attack, d.Defense, 0.85)
	if isMiss {
		return ActionResult{ActorName: hero.GetName(), SkillName: "Angriff", TargetName: d.GetName(), IsMiss: true}
	}
	d.TakeDamage(dmg)
	return ActionResult{ActorName: hero.GetName(), SkillName: "Angriff", TargetName: d.GetName(), Damage: dmg, IsCrit: isCrit}
}

// ── Dragon turn ───────────────────────────────────────────────────────────────

func processDragonTurn(d *dragon.EntropyDragon, heroes []internal.Combatant, round int) ActionResult {
	skill, targetIdx := d.ChooseAction(len(heroes))
	name := d.GetName()
	effectiveAttack := d.GetEffectiveAttack()

	clearScreen()
	printHeader()
	printFullStatus(heroes, d, round)
	fmt.Printf("\n  %s\n", dim(separator()))
	fmt.Printf("  %s\n", bold(red("DRACHEN-ZUG: "+skill.Name)))
	fmt.Printf("  %s\n", dim(separator()))

	var result ActionResult

	switch {
	case skill.Healing > 0:
		d.Heal(skill.Healing)
		fmt.Printf("\n  %s\n", green(fmt.Sprintf("Der Drache heilt sich um %d HP!", skill.Healing)))
		result = ActionResult{ActorName: name, SkillName: skill.Name, TargetName: name, Healing: skill.Healing}

	case skill.IsAOE:
		fmt.Println()
		totalDamage := 0
		for _, h := range heroes {
			if !h.IsAlive() {
				continue
			}
			dmg, isCrit, isMiss := CalculateDamage(skill.DamageMin, skill.DamageMax, effectiveAttack, h.GetStats().Defense, skill.Accuracy)
			if isMiss {
				fmt.Printf("  -> %-28s %s\n", h.GetName(), dim("VERFEHLT"))
				continue
			}
			h.SetCurrentHP(h.GetCurrentHP() - dmg)
			totalDamage += dmg
			fmt.Printf("  -> %-28s %s\n", h.GetName(), red(fmt.Sprintf("-%d HP%s", dmg, critSuffix(isCrit))))
		}
		result = ActionResult{ActorName: name, SkillName: skill.Name, TargetName: "alle Helden", Damage: totalDamage, IsAOE: true}

	default:
		target := liveTarget(heroes, targetIdx)
		if target == nil {
			result = ActionResult{ActorName: name, SkillName: skill.Name, IsMiss: true}
			break
		}
		fmt.Println()
		dmg, isCrit, isMiss := CalculateDamage(skill.DamageMin, skill.DamageMax, effectiveAttack, target.GetStats().Defense, skill.Accuracy)
		if isMiss {
			fmt.Printf("  -> %-28s %s\n", target.GetName(), dim("VERFEHLT"))
			result = ActionResult{ActorName: name, SkillName: skill.Name, TargetName: target.GetName(), IsMiss: true}
		} else {
			target.SetCurrentHP(target.GetCurrentHP() - dmg)
			fmt.Printf("  -> %-28s %s\n", target.GetName(), red(fmt.Sprintf("-%d HP%s", dmg, critSuffix(isCrit))))
			result = ActionResult{ActorName: name, SkillName: skill.Name, TargetName: target.GetName(), Damage: dmg, IsCrit: isCrit}
		}
	}

	fmt.Printf("\n  %s\n", dim("[Enter] weiter..."))
	waitEnter()
	return result
}

// ── UI screens ────────────────────────────────────────────────────────────────

func printIntroScreen(heroes []internal.Combatant, d *dragon.EntropyDragon) {
	clearScreen()
	fmt.Println(divider)
	fmt.Printf("|  %-*s  |\n", width-4, bold("CODERA BATTLE"))
	fmt.Printf("|  %-*s  |\n", width-4, "Der finale Kampf gegen den Entropie-Drachen")
	fmt.Println(divider)

	fmt.Printf("\n  %s\n\n", bold("EUER TEAM:"))
	fmt.Printf("  %-24s  %5s  %6s  %7s  %5s\n", "Name", "HP", "ATK", "DEF", "SPD")
	fmt.Printf("  %s\n", strings.Repeat("-", 52))
	for _, h := range heroes {
		s := h.GetStats()
		fmt.Printf("  %-24s  %5d  %6d  %7d  %5d\n",
			cyan(h.GetName()), s.MaxHP, s.Attack, s.Defense, s.Speed)
	}

	fmt.Printf("\n  %s\n\n", bold("GEGNER:"))
	fmt.Printf("  %-24s  %5d HP\n", red(d.GetName()), d.GetMaxHP())
	fmt.Printf("  %s\n\n", dim("Der Drache rasiert bei unter 50% HP (Schaden x1.5)"))

	fmt.Printf("  %s\n\n", yellow("Bereit? Viel Erfolg!"))
	fmt.Printf("  %s", dim("[Enter] Kampf starten..."))
	waitEnter()
}

func printRageWarning(d *dragon.EntropyDragon) {
	clearScreen()
	fmt.Println(divider)
	fmt.Printf("|  %-*s  |\n", width-4, red(bold("!!! ACHTUNG !!!")))
	fmt.Println(divider)
	fmt.Printf("\n  %s\n", red(bold(d.GetName()+" ist jetzt RASEND!")))
	fmt.Printf("  %s\n\n", yellow("Schaden erhöht um 50% — alle Aktionen zählen jetzt doppelt!"))
	fmt.Printf("  %s", dim("[Enter] weiter..."))
	waitEnter()
}

// announceDeaths prints a single notification if any hero just fell (HP reached 0 this turn).
// We track this by checking if IsAlive() is false but we haven't announced them yet using
// a package-level set of already-announced names.
var announcedDeaths = map[string]bool{}

func announceDeaths(heroes []internal.Combatant) {
	for _, h := range heroes {
		if !h.IsAlive() && !announcedDeaths[h.GetName()] {
			announcedDeaths[h.GetName()] = true
			clearScreen()
			fmt.Println(line)
			fmt.Printf("| %-*s |\n", width-2, red(bold("  "+h.GetName()+" ist gefallen!")))
			fmt.Println(line)
			fmt.Printf("\n  %s\n\n", dim("[Enter] weiter..."))
			waitEnter()
		}
	}
}

func printVictoryScreen(heroes []internal.Combatant) {
	fmt.Println(divider)
	fmt.Printf("|  %-*s  |\n", width-4, green(bold("*** SIEG! ***")))
	fmt.Printf("|  %-*s  |\n", width-4, "Der Entropie-Drache wurde besiegt!")
	fmt.Println(divider)
	fmt.Println()
	fmt.Println(line)
	fmt.Printf("| %-*s |\n", width-2, bold("  Kampfergebnis"))
	fmt.Println(line)
	for _, h := range heroes {
		if h.IsAlive() {
			pct := int(float64(h.GetCurrentHP()) / float64(h.GetMaxHP()) * 100)
			fmt.Printf("| %-*s |\n", width-2,
				fmt.Sprintf("  %s  %-22s  %d/%d HP  (%d%%)",
					green("UEBERLEBT"), h.GetName(), h.GetCurrentHP(), h.GetMaxHP(), pct))
		} else {
			fmt.Printf("| %-*s |\n", width-2,
				fmt.Sprintf("  %s   %s", red("GEFALLEN "), dim(h.GetName())))
		}
	}
	fmt.Println(line)
	fmt.Println()
}

func printDefeatScreen() {
	fmt.Println(divider)
	fmt.Printf("|  %-*s  |\n", width-4, red(bold("*** NIEDERLAGE ***")))
	fmt.Printf("|  %-*s  |\n", width-4, "Alle Helden sind gefallen...")
	fmt.Println(divider)
	fmt.Println()
}

// ── Status panel ──────────────────────────────────────────────────────────────

func printHeader() {
	fmt.Println(divider)
	fmt.Printf("|  %-*s  |\n", width-4, bold("CODERA BATTLE")+" — Kampf gegen den Entropie-Drachen")
	fmt.Println(divider)
}

// printFullStatus renders the complete battlefield: dragon + all heroes with HP bars.
func printFullStatus(heroes []internal.Combatant, d *dragon.EntropyDragon, round int) {
	rageTag := ""
	if d.IsEnraged {
		rageTag = "  " + yellow("RAGE")
	}

	fmt.Printf("\n  %s\n", dim(fmt.Sprintf("Runde %d", round)))
	fmt.Printf("  %s\n", dim(separator()))

	fmt.Printf("  %s\n", bold("GEGNER"))
	fmt.Printf("  %-28s  %s  %s%s\n",
		red(d.GetName()), coloredHPBar(d.GetCurrentHP(), d.GetMaxHP()),
		hpText(d.GetCurrentHP(), d.GetMaxHP()), rageTag)

	fmt.Printf("\n  %s\n", bold("TEAM"))
	for _, h := range heroes {
		tag := ""
		name := h.GetName()
		if !h.IsAlive() {
			tag = "  " + dim("[GEFALLEN]")
			name = dim(name)
		} else if float64(h.GetCurrentHP())/float64(h.GetMaxHP()) < 0.30 {
			tag = "  " + red("[KRITISCH!]")
			name = yellow(name)
		}
		fmt.Printf("  %-28s  %s  %s%s\n",
			name, coloredHPBar(h.GetCurrentHP(), h.GetMaxHP()),
			hpText(h.GetCurrentHP(), h.GetMaxHP()), tag)
	}
	fmt.Printf("  %s\n", dim(separator()))
}

// ── Skill menu ────────────────────────────────────────────────────────────────

func printSkillMenu(skills []internal.Skill) {
	fmt.Printf("\n  %s\n", bold("Aktionen:"))
	for i, s := range skills {
		dmgInfo := ""
		if s.DamageMin > 0 || s.DamageMax > 0 {
			dmgInfo = dim(fmt.Sprintf("  [%d-%d Schaden, %.0f%% Trefferchance]", s.DamageMin, s.DamageMax, s.Accuracy*100))
		} else if s.Healing > 0 {
			dmgInfo = dim(fmt.Sprintf("  [+%d Heilung]", s.Healing))
		}
		fmt.Printf("  %s  %-22s  %s%s\n",
			cyan(fmt.Sprintf("[%d]", i+1)), s.Name, dim(s.Description), dmgInfo)
	}
	fmt.Printf("\n  %s", bold("Deine Wahl: "))
}

// ── Action result ─────────────────────────────────────────────────────────────

func printActionResult(r ActionResult) {
	fmt.Printf("  %s\n", dim(separator()))
	switch {
	case r.IsMiss:
		fmt.Printf("  %s — %s  %s\n", r.ActorName, r.SkillName, yellow("VERFEHLT!"))
	case r.Healing > 0:
		fmt.Printf("  %s — %s  %s\n", r.ActorName, r.SkillName, green(fmt.Sprintf("+%d HP", r.Healing)))
	case r.SkillName == "Schutzschild":
		fmt.Printf("  %s — %s  %s\n", r.ActorName, r.SkillName, cyan("DEF +5 diese Runde"))
	case r.SkillName == "Kampfschrei":
		fmt.Printf("  %s — %s  %s  %s\n", r.ActorName, r.SkillName,
			red(fmt.Sprintf("-%d HP%s", r.Damage, critSuffix(r.IsCrit))), cyan("(ATK +5 naechste Runde)"))
	case strings.Contains(r.SkillName, "Doppel"):
		fmt.Printf("  %s — %s  %s\n", r.ActorName, r.SkillName,
			red(bold(fmt.Sprintf("-%d HP (2 Treffer!)", r.Damage))))
	default:
		fmt.Printf("  %s — %s  %s\n", r.ActorName, r.SkillName,
			red(fmt.Sprintf("-%d HP%s", r.Damage, critSuffix(r.IsCrit))))
	}
	fmt.Printf("  %s\n", dim(separator()))
	fmt.Printf("\n  %s", dim("[Enter] weiter..."))
}

// ── HP bar helpers ────────────────────────────────────────────────────────────

// coloredHPBar returns a color-coded HP bar: green > 60%, yellow > 30%, red <= 30%.
func coloredHPBar(current, max int) string {
	bar := hpBarRaw(current, max)
	ratio := 0.0
	if max > 0 {
		ratio = float64(current) / float64(max)
	}
	switch {
	case ratio > 0.60:
		return green(bar)
	case ratio > 0.30:
		return yellow(bar)
	default:
		return red(bar)
	}
}

func hpBarRaw(current, max int) string {
	if max <= 0 {
		return "[" + strings.Repeat(" ", barLen) + "]"
	}
	filled := int(float64(current) / float64(max) * float64(barLen))
	if filled < 0 {
		filled = 0
	}
	if filled > barLen {
		filled = barLen
	}
	return "[" + strings.Repeat("#", filled) + strings.Repeat(" ", barLen-filled) + "]"
}

func hpText(current, max int) string {
	return dim(fmt.Sprintf("%d/%d HP", current, max))
}

func separator() string {
	return strings.Repeat("-", width)
}

// ── Input / screen helpers ────────────────────────────────────────────────────

// clearScreen clears the terminal using ANSI escape codes.
func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

// waitEnter blocks until the player presses Enter.
func waitEnter() {
	sc := bufio.NewScanner(os.Stdin)
	sc.Scan()
}

func readInt(min, max int) int {
	sc := bufio.NewScanner(os.Stdin)
	for {
		if sc.Scan() {
			n, err := strconv.Atoi(strings.TrimSpace(sc.Text()))
			if err == nil && n >= min && n <= max {
				return n
			}
		}
		fmt.Printf("  %s", yellow(fmt.Sprintf("Bitte %d-%d eingeben: ", min, max)))
	}
}

// ── Combat helpers ────────────────────────────────────────────────────────────

func allDead(heroes []internal.Combatant) bool {
	for _, h := range heroes {
		if h.IsAlive() {
			return false
		}
	}
	return true
}

func liveTarget(heroes []internal.Combatant, idx int) internal.Combatant {
	if idx >= 0 && idx < len(heroes) && heroes[idx].IsAlive() {
		return heroes[idx]
	}
	for _, h := range heroes {
		if h.IsAlive() {
			return h
		}
	}
	return nil
}

func critSuffix(isCrit bool) string {
	if isCrit {
		return yellow("  ★ KRITISCH!")
	}
	return ""
}
