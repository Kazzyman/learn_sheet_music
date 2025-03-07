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
	fmt.Println("Identify the note below (or give a directive: L, R, all, S, O, Q, or dir.)\n")
	/*
		trainingWheels = false
		didADirective = false
		left = false  // When both are false it signifies we are being prompted on the entire Grand staff ...
		right = false // ... left: Lower; right: upper, or right-hand notes

		tryThatAgain = false // When the player commits an error, the player is forced to try that note again

	*/

	for {
		// Print the current score
		if total > 0 {
			percent := float64(correct) / float64(total) * 100
			fmt.Printf("%sScore: %d/%d (%.1f%%)\n%s\n", colorCyan, correct, total, percent, Reset)
		}

		// Get the next random note
		note :=
			NewRandomNote() // NewRandomNote() returns a simple struct, i.e. custom type Note ...
		// ... refer to the type definition in globals.go for why it is done this way.

		// Quiz :: process a single question, track time, and return 3 bool results
		givenCreditForCorrectAnswer, shouldQuit, outlierAdded =
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

		// Before looping to obtain the next random note, conditionally notify the player if an outlier occurred
		if outlierAdded {
			fmt.Printf("%s That answer was added to the outlier ledger.%s\n", Green, Reset)
		}

		if givenCreditForCorrectAnswer {
			// Pause for a bit before re-prompting the player
			time.Sleep(time.Second / 4) // one second / 4 ; or one quarter second
		} else if didADirective {
			didADirective = false
			// Pause for a bit before re-prompting the player
			time.Sleep(time.Second / 4) // one second / 4 ; or one quarter second
		} else {
			time.Sleep(time.Second * 2)
		}
	}
	// end of main loop
}

/*
.
*/

// NewRandomNote generates a random note
func NewRandomNote() Note { // Returns a simple struct; refer to Note's definition for details ::: - -
	// ::: force player to answer correctly prior to being given a new novel note to guess
	if tryThatAgain {
		tryThatAgain = false
		if left {
			return Note{Pitch: pitchesLeft[rememberLastPickL]}
		} else if right {
			return Note{Pitch: pitchesRight[rememberLastPickR]}
		} else {
			return Note{Pitch: pitchesAll[rememberLastPickAll]}
		}
	} else {
		if left { // 11
			r := rand.Intn(11) // ... 11 random indexes for the pitches-left const
			rememberLastPickL = r
			return Note{Pitch: pitchesLeft[r]} // ::: Pitch is a member of the Note struct ...
		} else if right { // 13
			r := rand.Intn(13) // ... 13 random indexes for the pitches-right const
			rememberLastPickR = r
			return Note{Pitch: pitchesRight[r]}
		} else {
			r := rand.Intn(24) // ... 24 random indexes for the pitches-all const
			rememberLastPickAll = r
			return Note{Pitch: pitchesAll[r]}
		}
		// ... those say: make a Note type with its Pitch field set to pitches__[r], and return that to the caller.
	}
}

/*
.
*/

// Quiz prompts, process the player's response, track time, return results: (givenCreditForCorrectAnswer, shouldQuit, outlierAdded)
func Quiz(note Note, mapOfNoteStats map[string]NoteStats, outliers *[]Outlier) (givenCredit bool, shouldQuit bool, outlierAdded bool) { // ::: - -
	/*
		Usage of the Note type instead of a simple string is explained in the globals.go file.
		&outliers makes a pointer to a slice of structures, so Quiz gets *[]Outlier as a parameter (a pointer).
		*[]Outlier is a pointer to that slice; the * de-references, while pointing, resulting in the slice, not just ...
		... the pointer itself. Why do this? Because passing outliers []Outlier would give us a copy rather than the original.
		Here outliers *[]Outlier gives us the outliers "out there", whereas if we had used outliers []Outlier we'd get a local
		copy of the object (in this case a slice of structures) in the outside world and persistence would be lost.
		"go" always passes by value, "exception": maps are reference types, so maps are kinda pseudo pointers by default.
	*/
	// ::: givenCreditForCorrectAnswer := false
	shouldQuit = false

	// DrawStaff is passed Pitch via a Note type struct
	// DrawStaff is used to prompt ::: ONLY here!
	DrawStaff(note, true, true) // prompting true causes a normal display of the staff

	if !givenCreditForCorrectAnswer {
		// Prompt three ways: Lower cleft, Right-hand cleft, or the entire Grand staff:
		if left {
			fmt.Print("Again; Guess-L: ") // Prompt the player for a guess in the lower staff.
		} else if right {
			fmt.Print("Again; Guess-R: ") // Prompt the player for a guess in the upper staff.
		} else {
			fmt.Print("Again; Guess: ") // Prompt the player for a guess throughout the entire staff.
		}
	} else {
		// Prompt three ways: Lower cleft, Right-hand cleft, or the entire Grand staff:
		if left {
			fmt.Print("Guess-L: ") // Prompt the player for a guess in the lower staff.
		} else if right {
			fmt.Print("Guess-R: ") // Prompt the player for a guess in the upper staff.
		} else {
			fmt.Print("Guess: ") // Prompt the player for a guess throughout the entire staff.
		}
	}

	givenCreditForCorrectAnswer = false

	reader := bufio.NewReader(os.Stdin) // Create local "reader" which is an object of type bufio.NewReader
	// ::: Obtain player's answer, "on the clock"
	start := time.Now()                           // start the clock
	answer, _ = reader.ReadString('\n')           // obtain the player's guess
	elapsedMs := time.Since(start).Milliseconds() // stop the clock
	elapsedSec := float64(elapsedMs) / 1000.0     // recast time to a float, and convert Ms to sec
	answer = strings.TrimSpace(answer)            // trim the answer (essential)

	// All the ways to return early to the main loop
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
	if strings.ToLower(answer) == "l" {
		left = true
		right = false
		didADirective = true
		return givenCreditForCorrectAnswer, shouldQuit, false // Continue, no score change, no outlier
	}
	if strings.ToLower(answer) == "r" {
		right = true
		left = false
		didADirective = true
		return givenCreditForCorrectAnswer, shouldQuit, false // Continue, no score change, no outlier
	}
	if strings.ToLower(answer) == "all" {
		right = false
		left = false
		didADirective = true
		return givenCreditForCorrectAnswer, shouldQuit, false // Continue, no score change, no outlier
	}
	//
	if strings.ToLower(answer) == "dir" {
		fmt.Println("directives: L, R, all, S, O, Q, tw, two)\n")
		total--
		return givenCreditForCorrectAnswer, shouldQuit, false // Continue, no score change, no outlier
	}

	// trainingWheels
	if strings.ToLower(answer) == "tw" {
		trainingWheels = true
		didADirective = true
		total--
		return givenCreditForCorrectAnswer, shouldQuit, false // Continue, no score change, no outlier
	}
	if strings.ToLower(answer) == "two" {
		trainingWheels = false
		didADirective = true
		total--
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
	outlierAdded = false

	// Test the player's guess
	if strings.ToUpper(answer) == pitch { // ::: Correct
		DrawStaff(note, false, true) // redraw the staff with the note shown in green to signify a correct guess.
		CurentNoteStatsObject.TotalCorrectMs += elapsedMs
		CurentNoteStatsObject.CorrectCount++
		CurentNoteStatsObject.AvgCorrectSec = float64(CurentNoteStatsObject.TotalCorrectMs) / float64(CurentNoteStatsObject.CorrectCount) / 1000.0
		if elapsedMs < 250 { // 1000=1.00s ::: because time.Sleep 1/4s
			fmt.Printf("%sIt was %s. (but too fast, answer given prior to query being shown, not counted)%s\n", Red, pitch, Reset)
			outlierAdded = true // We use this flag to inform the player of the disposition of this super-fast screw-up
			tryThatAgain = true // because even though it was correct, it was pure luck (answered prior to the query)
			*outliers = append(*outliers, Outlier{Pitch: note.Pitch, WasCorrect: true, TimeSec: elapsedSec})
		}
		mapOfNoteStats[note.Pitch] = CurentNoteStatsObject
		if elapsedMs > 250 {
			givenCreditForCorrectAnswer = true
			return givenCreditForCorrectAnswer, shouldQuit, outlierAdded // ::: in any case we return if answer == pitch
		}
	} else { // ::: Wrong
		// Check for fast miss outlier
		if elapsedMs < 250 { // 2100 = 2.1 seconds, 700 = 0.7s, 1090 = 1.09s ::: because time.Sleep is 1/4s
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
	orderedPitches := []string{"A5",
		"G5", "F5", "E5", "D5", "C5",
		"B4", "A4", "G4", "F4", "E4", "D4", "C4",
		"B3", "A3", "G3", "F3", "E3", "D3", "C3",
		"B2", "A2", "G2", "F2",
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
		result := "(too fast, under 250 Ms)"
		fmt.Printf("%s%d: %s - %s, %.3f seconds%s\n", LightBlue, i+1, o.Pitch, result, o.TimeSec, Reset)
	}
}
