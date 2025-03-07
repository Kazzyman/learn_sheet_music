package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

/*
// ::: About:
This project is spread across 3 go files and is formatted for JetBrains GoLand. The 3 files are globals.go, drawStaff.go
and this one. The author is strictly a hobbyist who specializes in learning apps of personal interest, sometimes with a
game-style option. The comments herein reflect the fact that it is highly-unlikely that anyone other than the author will
ever read them. Three sequential colons: ::: causes the rest of a comment line to be highlighted.
*/

/*
// ::: Features:
Notes show in Red for added visibility.
Correct answers cause the note to turn green.
Wrong answers turn the note yellow, and the correct response is shown alongside it.
Wrong answers demand a repeat of the failed query, until correctly answered.
All statistics are logged and reported during app shutdown; or upon request via directives: "s", or "o"
Progress and scoring is shown during play.
*/
/*
 */
func main() {
	fmt.Println("\nRick's first Sheet Music Learning App")
	tryThatAgain = false // When the player commits an error, the player is forced to try that note again
	for {
		// Get the next random note
		note :=
			NewRandomNote() // NewRandomNote() returns a simple struct, i.e. custom type Note ...
		// ... refer to the type definition in globals.go for why it is done this way.

		// Quiz :: process a single question, track time, and return 3 bool results
		givenCreditForCorrectAnswer, shouldQuit, outlierAdded :=
			Quiz(note, stats, &outliers) // "&outliers" gets a pure and simple pointer
		// ... the above is a good example of passing by pointer, &outliers evaluates|resolves to a pointer.
		// Notice that since stats is a map, and therefore of reference type, there is no need to pass it by pointer.

		// Calculate the current score
		if givenCreditForCorrectAnswer {
			correct++ // Global
		}
		total++ // Global

		if shouldQuit {
			fmt.Printf("Final Score: %d/%d (%.1f%%)\n", correct, total, float64(correct)/float64(total)*100)
			PrintStats(stats)
			printOutliers(outliers)
			break // The only exit point for the app
		}

		// Print the current score
		percent := float64(correct) / float64(total) * 100
		fmt.Printf("Score: %d/%d (%.1f%%)\n", correct, total, percent)

		// Before looping to obtain the next random note, conditionally notify the player if an outlier occurred
		if outlierAdded {
			fmt.Printf("%s That answer was added to the outlier ledger.%s\n", Green, Reset)
		}

		// Pause for a bit before re-prompting the player
		// time.Sleep(1.5 * time.Second) // optionally set to a fraction or multiple of one second
		time.Sleep(time.Second) // one second
	}
}

/*
.
*/

// NewRandomNote generates a random note
func NewRandomNote() Note { // Returns a simple struct; refer to Note's definition for details ::: - -
	pitches := []string{"A6",
		"G5", "F5", "E5", "D5", "C5", "B5", "A5",
		"G4", "F4", "E4", "D4", "C4", "B4", "A4",
		"G3", "F3", "E3", "D3", "C3", "B3", "A3",
		"G2", "F2"}
	//  24 pitches in the slice ...
	if tryThatAgain {
		tryThatAgain = false
		return Note{Pitch: pitches[rememberLastPick]} // ::: force player to answer correctly
	} else {
		r := rand.Intn(24) // ... so, also 24 random indexes for the pitches slice
		rememberLastPick = r
		return Note{Pitch: pitches[r]} // Pitch is a member of the Note struct ...
		// ... this says: make a Note type with its Pitch field set to pitches[r], and return that to the caller.
	}
}

/*
.
*/

// Quiz process the player's response, track time, return results: (givenCreditForCorrectAnswer, shouldQuit, outlierAdded)
func Quiz(note Note, mapOfNoteStats map[string]NoteStats, outliers *[]Outlier) (bool, bool, bool) { // ::: - -
	/*
		Usage of the Note type instead of a simple string is explained in the globals.go file.
		&outliers makes a pointer to a slice of structures, so Quiz gets *[]Outlier as a parameter (a pointer).
		*[]Outlier is a pointer to that slice; the * de-references, while pointing, resulting in the slice, not just ...
		... the pointer itself. Why do this? Because passing outliers []Outlier would give us a copy rather than the original.
		Here outliers *[]Outlier gives us the outliers "out there", whereas if we had used outliers []Outlier we'd get a local
		copy of the object (in this case a slice of structures) in the outside world and persistence would be lost.
		"go" always passes by value, "exception": maps are reference types, so maps are kinda pseudo pointers by default.
	*/
	givenCreditForCorrectAnswer := false
	shouldQuit := false

	reader := bufio.NewReader(os.Stdin) // Create local "reader" which is an object of type bufio.NewReader

	fmt.Println("Identify the note below (or give a directive: s, o, q etc.)\n")
	// DrawStaff is passed Pitch via a Note type struct
	DrawStaff(note, true, true) // prompting true causes a normal display of the staff
	fmt.Print("Guess: ")        // Prompt the player for a guess.

	// Obtain player's answer, "on the clock"
	start := time.Now()                           // start the clock
	answer, _ := reader.ReadString('\n')          // obtain the player's guess
	elapsedMs := time.Since(start).Milliseconds() // stop the clock
	elapsedSec := float64(elapsedMs) / 1000.0     // recast time to a float, and convert Ms to sec
	answer = strings.TrimSpace(answer)            // trim the answer (essential)

	// Three ways to return early to the main loop
	if strings.ToLower(answer) == "q" {
		shouldQuit = true
		return givenCreditForCorrectAnswer, shouldQuit, false // (isCorrect, shouldQuit, outlierAdded)
	}
	if strings.ToLower(answer) == "o" {
		printOutliers(*outliers)                              // printOutliers expects a slice of structures : type []Outlier
		return givenCreditForCorrectAnswer, shouldQuit, false // Continue, no score change, no outlier
	}
	if strings.ToLower(answer) == "s" {
		PrintStats(mapOfNoteStats)
		return givenCreditForCorrectAnswer, shouldQuit, false // Continue, no score change, no outlier
	}

	pitch := note.Pitch                       // note.Pitch will eval to a note+octave couplet, such as C4, which needs trimming ...
	pitch = strings.TrimRight(pitch, "23456") // 2-6 are octave suffixes of the various staff notes ...
	// ... and comprise here a "cut-set" which strings.TrimRight uses to know where and when to do its trimming.

	// "mapOfNoteStats" was passed-in as map[string]NoteStats -- a map of key-As-String+NoteStats pairs
	CurentNoteStatsObject := mapOfNoteStats[note.Pitch] // "Pitch" is a field/member of the Note struct (locally: note) ...
	// ... mapOfNoteStats[note.Pitch] obtains the NoteStats struct "indexed" by the Key string: note.Pitch
	// ... it says: give me the NoteStats in the map of NoteStats at Pitch in the map
	// ::: therefore, "CurentNoteStatsObject" is created as a NoteStats object.
	// Here "note" is actually a Note struct (because it was passed into this func as such)
	// "mapOfNoteStats" is a map, a correspondence between pitch (a string) and a NoteStats struct

	CurentNoteStatsObject.Attempts++ // increment the Attempts field/member of the current (specific) NoteStats struct
	outlierAdded := false

	// Test the player's guess
	// ::: Correct
	if strings.ToUpper(answer) == pitch {
		DrawStaff(note, false, true) // redraw the staff with the note shown in green to signify a correct guess.
		if elapsedMs <= 13000 {      // 13 seconds threshold
			CurentNoteStatsObject.TotalCorrectMs += elapsedMs
			CurentNoteStatsObject.CorrectCount++
			CurentNoteStatsObject.AvgCorrectSec = float64(CurentNoteStatsObject.TotalCorrectMs) / float64(CurentNoteStatsObject.CorrectCount) / 1000.0
		} else if elapsedMs < 1090 { // 1000=1.00s ::: 1.00s because time.Sleep is 1.09s
			fmt.Printf("%sIt was %s. (but too fast, answer given prior to query being shown, not counted)%s\n", Red, pitch, Reset)
			outlierAdded = true // We use this flag to inform the player of the disposition of this super-fast screw-up
			tryThatAgain = true // because even though it was correct, it was pure luck (answered prior to the query)
		} else {
			// The player took a long time to get it right, so log it as an outlier and set a flag to nudge the player for being pokey.
			fmt.Printf("%sToo slow%s, therefore ... ", colorYellow, Reset)
			*outliers = append(*outliers, Outlier{Pitch: note.Pitch, WasCorrect: true, TimeSec: elapsedSec}) // add some literals to a slice
			outlierAdded = true
		}
		mapOfNoteStats[note.Pitch] = CurentNoteStatsObject
		if elapsedMs > 1089 && elapsedMs <= 13000 { // is good, was not too fast, and not too slow
			givenCreditForCorrectAnswer = true
			return givenCreditForCorrectAnswer, shouldQuit, outlierAdded // ::: in any case we return if answer == pitch
		} else {
			givenCreditForCorrectAnswer = false
			return givenCreditForCorrectAnswer, shouldQuit, outlierAdded // ::: in any case we return if answer == pitch
		}
	} else { // ::: Wrong
		// Check for fast miss outlier
		if elapsedMs < 1090 { // 2100 = 2.1 seconds, 700 = 0.7s, 1090 = 1.09s ::: 1.09s because time.Sleep is 1.0s
			fmt.Printf("%sActually it was %s. (Too fast, answer given prior to query being shown, not counted)%s\n", Red, pitch, Reset)
			*outliers = append(*outliers, Outlier{Pitch: note.Pitch, WasCorrect: false, TimeSec: elapsedSec}) // essential pointer magic here...
			// Without pointer magic: outliers := append(*outliers, Outlier{Pitch: note.Pitch, WasCorrect: false, TimeSec: elapsedSec})
			// ... append by itself (without any pointer magic) would just make an appended copy
			// ... outliers naked is just a memory address which contains another memory address (to a value)
			// outliers is a global var of type []Outlier
			outlierAdded = true // We use this flag to inform the player of the disposition of this super-fast screw-up
			tryThatAgain = true
		} else {
			DrawStaff(note, false, false) // Re-draw the staff with the note highlighted in Yellow + correction !
			CurentNoteStatsObject.Misses++
			tryThatAgain = true
		}
	}
	mapOfNoteStats[note.Pitch] = CurentNoteStatsObject // update the mapOfNoteStats map at Key = note.Pitch; update it with the "s" structure, remembering that s := mapOfNoteStats[note.Pitch] ...
	// ... and mapOfNoteStats is a map of the form map[string]NoteStats
	return givenCreditForCorrectAnswer, shouldQuit, outlierAdded // return three bools: givenCreditForCorrectAnswer, shouldQuit, outlierAdded
}

// PrintStats display per-note performance stats
func PrintStats(stats map[string]NoteStats) {
	orderedPitches := []string{"A6",
		"G5", "F5", "E5", "D5", "C5", "B5", "A5",
		"G4", "F4", "E4", "D4", "C4", "B4", "A4",
		"G3", "F3", "E3", "D3", "C3", "B3", "A3",
		"G2", "F2",
	}
	fmt.Printf("%s--- Note Stats ---%s\n", colorYellow, Reset)
	for _, pitch := range orderedPitches { // ::: each loop prints the stats for one note
		s := stats[pitch] // get the NoteStats for one pitch, from the orderedPitches slice
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

// printOutliers displays the outliers ledger
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
