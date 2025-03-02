package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {
	// initializeStuff() // making all of that stuff global was good-enough
	fmt.Println("\nRick's first Sheet Music Learning App (Treble Clef)")
	fmt.Println("Type the note letter (or 'o' for outliers, 'q' to quit).")

	for {
		note := NewRandomNote(stats) // get the next note

		isCorrect, shouldQuit, outlierAdded := Quiz(note, stats, &outliers)
		if shouldQuit {
			fmt.Printf("Final Score: %d/%d (%.1f%%)\n", correct, total, float64(correct)/float64(total)*100)
			PrintStats(stats)
			printOutliers(outliers)
			fmt.Println("Goodbye.")
			break
		}
		if isCorrect {
			correct++
		}
		total++
		percent := float64(correct) / float64(total) * 100
		fmt.Printf("Score: %d/%d (%.1f%%)\n", correct, total, percent)

		// PrintStats(stats)
		if outlierAdded {
			fmt.Printf("%sOutlier added to ledger.%s\n", Red, Reset)
		}
		fmt.Println("---")
		time.Sleep(1 * time.Second)
	}
}

// NewRandomNote generates a random note
func NewRandomNote(stats map[string]NoteStats) Note {
	pitches := []string{"A5", "G5", "F5", "E5", "D5", "C5", "B4", "A4", "G4", "F4", "E4", "D4", "C4"} // +2 ledger notes
	// 13 pitches in the slice
	r := rand.Intn(13)
	return Note{Pitch: pitches[r]}
}

// Quiz runs a single question, tracks time, and returns result
func Quiz(note Note, stats map[string]NoteStats, outliers *[]Outlier) (bool, bool, bool) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("What note is this? (e.g., 'C', 'o' for outliers, 's' for stats, or 'q' to quit)")
	fmt.Println(DrawStaff(note))
	fmt.Print("Your answer: ")

	start := time.Now()
	answer, _ := reader.ReadString('\n')
	elapsedMs := time.Since(start).Milliseconds()
	elapsedSec := float64(elapsedMs) / 1000.0
	answer = strings.TrimSpace(answer)

	if strings.ToLower(answer) == "q" {
		return false, true, false
	}
	if strings.ToLower(answer) == "o" {
		printOutliers(*outliers)
		return false, false, false // Continue, no score change, no outlier
	}
	if strings.ToLower(answer) == "s" {
		PrintStats(stats)
		return false, false, false // Continue, no score change, no outlier
	}

	pitch := note.Pitch
	if strings.Contains(pitch, "4") || strings.Contains(pitch, "5") {
		pitch = strings.TrimRight(pitch, "45")
	}

	s := stats[note.Pitch]
	s.Attempts++
	outlierAdded := false
	if strings.ToUpper(answer) == pitch {
		fmt.Printf("%sCorrect!%s\n", Green, Reset)
		if elapsedMs <= 13000 { // 13 seconds threshold
			s.TotalCorrectMs += elapsedMs
			s.CorrectCount++
			s.AvgCorrectSec = float64(s.TotalCorrectMs) / float64(s.CorrectCount) / 1000.0
		} else {
			*outliers = append(*outliers, Outlier{Pitch: note.Pitch, WasCorrect: true, TimeSec: elapsedSec})
			outlierAdded = true
		}
		stats[note.Pitch] = s
		return true, false, outlierAdded
	}

	// Check for fast miss outlier
	if elapsedMs < 2100 { // 2.1 seconds
		fmt.Printf("%sWrong. It was %s. (Too fast, not counted)%s\n", Red, pitch, Reset)
		*outliers = append(*outliers, Outlier{Pitch: note.Pitch, WasCorrect: false, TimeSec: elapsedSec})
		outlierAdded = true
	} else {
		fmt.Printf("%sWrong. It was %s.%s\n", Red, pitch, Reset)
		s.Misses++
	}
	stats[note.Pitch] = s
	return false, false, outlierAdded
}

// PrintStats shows per-note performance in light blue
func PrintStats(stats map[string]NoteStats) {
	orderedPitches := []string{"A5", "G5", "F5", "E5", "D5", "C5", "B4", "A4", "G4", "F4", "E4", "D4", "C4"}
	// 13 pitches
	fmt.Printf("%s--- Note Stats ---%s\n", colorYellow, Reset)
	for _, pitch := range orderedPitches {
		s := stats[pitch]
		avgTime := "N/A"
		if s.CorrectCount > 0 {
			avgTime = fmt.Sprintf("%.3f seconds", s.AvgCorrectSec)
		}
		fmt.Printf("%s%s: Attempts=%d, Misses=%d, Avg Correct Time=%s%s\n", LightBlue, pitch, s.Attempts, s.Misses, avgTime, Reset)
	}
}

// printOutliers displays the outlier ledger in light blue
func printOutliers(outliers []Outlier) {
	fmt.Printf("%s--- Outlier Ledger ---%s\n", colorYellow, Reset)
	if len(outliers) == 0 {
		fmt.Printf("%sNo outliers recorded.%s\n", LightBlue, Reset)
		return
	}
	for i, o := range outliers {
		result := "Correct (too slow, over 13s)"
		if !o.WasCorrect {
			result = "Missed (too fast, under 2.1s)"
		}
		fmt.Printf("%s%d: %s - %s, %.3f seconds%s\n", LightBlue, i+1, o.Pitch, result, o.TimeSec, Reset)
	}
}

// DrawStaff generates a five-line ASCII staff with the note placed
func DrawStaff(note Note) string {
	staff := []string{
		"        ------ ",      // 15 chars (A5)
		"                    ", // 20 chars (G5)
		"  ------------------", // 20 chars (F5)
		"                    ", // 20 chars (E5)
		"  ------------------", // 20 chars (D5)
		"                    ", // 20 chars (C5)
		"  ------------------", // 20 chars (B4)
		"                    ", // 20 chars (A4)
		"  ------------------", // 20 chars (G4)
		"                    ", // 20 chars (F4)
		"  ------------------", // 20 chars (E4)
		"                    ", // 20 chars (D4)
		"        ------ ",      // 15 chars (C4)
	}

	pitchMap := map[string]int{
		"A5": 0, "G5": 1, "F5": 2, "E5": 3, "D5": 4, "C5": 5,
		"B4": 6, "A4": 7, "G4": 8, "F4": 9, "E4": 10, "D4": 11, "C4": 12,
	}

	lineIndex, exists := pitchMap[note.Pitch]
	if !exists {
		fmt.Printf("\nlineIndex: %d, exists: %t\n", lineIndex, exists)
		return "Invalid pitch"
	}

	// Use position 8—works for all (before dashes, mid-spaces)
	staff[lineIndex] = staff[lineIndex][:10] + "●" + staff[lineIndex][11:]
	return strings.Join(staff, "\n")
}
