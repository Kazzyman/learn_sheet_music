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
	fmt.Println("\nRick's first Sheet Music Learning App")
	// Rick knows about options such as: s, o, and q, so no need to bother printing then to the terminal

	for {
		// Get the next random note
		note := NewRandomNote() // NewRandomNote() returns a simple struct, i.e. custom type Note.
		// Refer to the type definition in globals.go for why it is done this way.

		// Quiz :: process a single question, track time, and return result: isCorrect, shouldQuit, outlierAdded
		isCorrect, shouldQuit, outlierAdded := Quiz(note, stats, &outliers) // "&outliers" gets a pure and simple pointer
		// ... the above is a good example of passing by pointer.
		if shouldQuit {
			fmt.Printf("Final Score: %d/%d (%.1f%%)\n", correct, total, float64(correct)/float64(total)*100)
			PrintStats(stats)
			printOutliers(outliers)
			break
		}

		// Tally, calculate, and print the current score
		if isCorrect {
			correct++
		}
		total++
		percent := float64(correct) / float64(total) * 100
		fmt.Printf("Score: %d/%d (%.1f%%)\n", correct, total, percent)

		// Before looping to obtain the next random note, conditionally notify the player if an outlier occurred
		if outlierAdded {
			fmt.Printf("%s That answer will be added to the outlier ledger.%s\n", Green, Reset)
		}

		// Pause for a bit before re-prompting the player
		// time.Sleep(1 * time.Second) // optionally set to a fraction or multiple of one second 
		time.Sleep(time.Second) // one second 
	}
}

// NewRandomNote generates a random note
func NewRandomNote() Note { // Returns a simple struct; refer to Note's definition for details ::: - -
	pitches := []string{"A6",
		"G5", "F5", "E5", "D5", "C5", "B5", "A5",
		"G4", "F4", "E4", "D4", "C4", "B4", "A4",
		"G3", "F3", "E3", "D3", "C3", "B3", "A3",
		"G2", "F2"}
	//  24 pitches in the slice
	r := rand.Intn(24)             // so, also 24 random indexes for the pitches slice
	return Note{Pitch: pitches[r]} // Pitch is a member of the Note struct ...
	// ... this says: make a Note var with its Pitch field set to pitches[r], and return that to the caller.
}

// Quiz process the question, track time, return results: (isCorrect, shouldQuit, outlierAdded)
func Quiz(note Note, stats map[string]NoteStats, outliers *[]Outlier) (bool, bool, bool) { // ::: - -
	/*
		Usage of the Note type instead of a simple string is explained in the globals.go file.
		&outliers makes a pointer to the slice, so Quiz gets *[]Outlier as a parameter (a pointer).
		*[]Outlier is a pointer to that slice; the * de-references, while pointing, resulting in the slice, not just ...
		... the pointer itself. Why do this? Because passing outliers []Outlier would give us a copy rather than the original.
		Here outliers *[]Outlier gives us the outliers "out there", whereas if we had used outliers []Outlier we'd get a local
		copy of that outside world.
		"go" always passes by value, "exception": maps are reference types, so maps are kinda pseudo pointers by default.
	*/

	reader := bufio.NewReader(os.Stdin) // Create local "reader" which is an object of type bufio.NewReader

	fmt.Println("What note is this? (or give a directive: s, o, q etc.)\n")
	fmt.Println(DrawStaff(note)) // DrawStaff is passed Pitch via a type Note
	fmt.Print("Guess: ")

	// Obtain player's answer, on the clock
	start := time.Now()                           // start the clock
	answer, _ := reader.ReadString('\n')          // obtain the player's guess
	elapsedMs := time.Since(start).Milliseconds() // stop the clock
	elapsedSec := float64(elapsedMs) / 1000.0     // recast time to a float, and convert Ms to sec
	answer = strings.TrimSpace(answer)            // trim the answer (essential)

	// Three ways to return to the main loop
	if strings.ToLower(answer) == "q" {
		return false, true, false // (isCorrect, shouldQuit, outlierAdded)
	}
	if strings.ToLower(answer) == "o" {
		printOutliers(*outliers)
		return false, false, false // Continue, no score change, no outlier
	}
	if strings.ToLower(answer) == "s" {
		PrintStats(stats)
		return false, false, false // Continue, no score change, no outlier
	}

	pitch := note.Pitch // the note struct was
	pitch = strings.TrimRight(pitch, "23456")

	// "stats" was passed-in as map[string]NoteStats // a map of key:string, NoteStats pairs
	s := stats[note.Pitch] // "Pitch" is a field/member of the Note struct (internally: note) ...
	// ... stats[note.Pitch] obtains the NoteStats struct "indexed" by the Key string: note.Pitch
	// therefore, "s" is created as a NoteStats object.
	// Here "note" is actually a Note struct (because it was passed into this func as such)
	// "stats" is a map, a correspondence between pitch (a string) and a NoteStats struct

	s.Attempts++ // increment the Attempts field/member of the current (specific) NoteStats struct
	outlierAdded := false
	if strings.ToUpper(answer) == pitch {
		fmt.Printf("%sCorrect!%s\n", Green, Reset) // todo::: move this into a reprint of staff
		if elapsedMs <= 13000 {                    // 13 seconds threshold
			s.TotalCorrectMs += elapsedMs
			s.CorrectCount++
			s.AvgCorrectSec = float64(s.TotalCorrectMs) / float64(s.CorrectCount) / 1000.0
		} else {
			//
			*outliers = append(*outliers, Outlier{Pitch: note.Pitch, WasCorrect: true, TimeSec: elapsedSec}) // add some literals to a slice
			outlierAdded = true
		}
		stats[note.Pitch] = s
		return true, false, outlierAdded
	}

	// Check for fast miss outlier
	if elapsedMs < 2100 { // 2.1 seconds
		fmt.Printf("%sActually it was %s. (Too fast, not counted)%s\n", Red, pitch, Reset)                // todo::: move this into a reprint of staff
		*outliers = append(*outliers, Outlier{Pitch: note.Pitch, WasCorrect: false, TimeSec: elapsedSec}) // essential pointer magic here...
		// Without pointer magic: outliers := append(*outliers, Outlier{Pitch: note.Pitch, WasCorrect: false, TimeSec: elapsedSec})
		// ... append by itself (without any pointer magic) would just make an appended copy
		// ... outliers naked is just a memory address which contains another memory address (to a value) 
		// outliers is a global var of type []Outlier
		outlierAdded = true
	} else {
		fmt.Printf("%sActually it was %s, specifically %s %s\n", Red, pitch, note.Pitch, Reset) // todo::: move this into a reprint of staff
		s.Misses++
	}
	stats[note.Pitch] = s // update the stats map at Key = note.Pitch; update it with the "s" structure, remembering that s := stats[note.Pitch] ... 
	// ... and stats is a map of the form map[string]NoteStats
	return false, false, outlierAdded // return three bools 
}

// PrintStats shows per-note performance in light blue
func PrintStats(stats map[string]NoteStats) {
	orderedPitches := []string{"A6",
		"G5", "F5", "E5", "D5", "C5", "B5", "A5",
		"G4", "F4", "E4", "D4", "C4", "B4", "A4",
		"G3", "F3", "E3", "D3", "C3", "B3", "A3",
		"G2", "F2",
	}
	// 13 pitches
	fmt.Printf("%s--- Note Stats ---%s\n", colorYellow, Reset)
	for _, pitch := range orderedPitches {
		s := stats[pitch] // 
		avgTime := "N/A"
		if s.CorrectCount > 0 {
			avgTime = fmt.Sprintf("%.3f seconds", s.AvgCorrectSec)
		}
		if s.Misses > 0 {
			fmt.Printf("%s%s: Attempts=%d, %sMisses=%d%s, Avg Correct Time=%s%s\n", LightBlue, pitch, s.Attempts, Red, s.Misses, LightBlue, avgTime, Reset)
		} else {
			fmt.Printf("%s%s: Attempts=%d, Misses=%d, Avg Correct Time=%s%s\n", LightBlue, pitch, s.Attempts, s.Misses, avgTime, Reset)
		}
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
func DrawStaff(note Note) string { // accepts a simple structure of notes ::: - -
	staff := []string{
		"        ------ ",      // (A6)  6th octave starts on A and runs through G
		"                    ", // (G5)
		"  ------------------", // (F5)
		"                    ", // (E5)
		"  ------------------", // (D5)
		"                    ", // (C5)
		"  ------------------", // (B5)
		"                    ", // (A5)  5th octave starts on A and runs through G
		"  ------------------", // (G4)
		"                    ", // (F4)
		"  ------------------", // (E4)
		"                    ", // (D4)
		"        ------ ",      // (C4)
		// somehow we'd like to need/want to skip this space ???
		"                    ", // (B4)
		"  ------------------", // (A4)  4th octave starts on A and runs through G
		"                    ", // (G3)
		"  ------------------", // (F3)
		"                    ", // (E3)
		"  ------------------", // (D3)
		"                    ", // (C3)
		"  ------------------", // (B3)
		"                    ", // (A3) 3erd octave starts on A and runs through G
		"  ------------------", // (G2)
		"                    ", // (F2)
	}

	pitchMap := map[string]int{
		"A6": 0,
		"G5": 1, "F5": 2, "E5": 3, "D5": 4, "C5": 5, "B5": 6, "A5": 7,
		"G4": 8, "F4": 9, "E4": 10, "D4": 11, "C4": 12, "B4": 13, "A4": 14,
		"G3": 15, "F3": 16, "E3": 17, "D3": 18, "C3": 19, "B3": 20, "A3": 21,
		"G2": 22, "F2": 23,
	}

	// obtain an index into the lines of the staff (counting begins at 0)
	lineIndex, exists := pitchMap[note.Pitch] // use a random note as an index|key into pitchMap

	// handle error
	// fmt.Printf("\nlineIndex: %d, exists: %t\n", lineIndex, exists)
	if !exists { // if the value of exists is false
		fmt.Printf("\nlineIndex: %d, exists-not: %t\n", lineIndex, exists)
		return "Invalid pitch" // return early without a staff line
	}

	// place the note on the staff
	staff[lineIndex] = staff[lineIndex][:10] + Red + "●" + Reset + staff[lineIndex][11:]
	return strings.Join(staff, "\n") // returns a line with an embedded note
}
