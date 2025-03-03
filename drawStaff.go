package main

import (
	"fmt"
	"strings"
)

// DrawStaff generates a twelve-line, twelve-space ASCII staff with a note placed in Red
func DrawStaff(note Note, prompting bool, correct bool) { // accepts a simple structure of notes ::: - -
	// First we make an slice of strings, one string for each staff line
	staff := []string{
		"        ------ ",      // 1 (A6)  6th octave starts on A and runs through G
		"                    ", // 2 (G5)
		"  ------------------", // 3 (F5)
		"                    ", // 4 (E5)
		"  ------------------", // 5 (D5)
		"                    ", // 6 (C5)
		"  ------------------", // 7 (B5)
		"                    ", // 8 (A5)  5th octave starts on A and runs through G
		"  ------------------", // 9 (G4)
		"                    ", // 10(F4)
		"  ------------------", // 11(E4)
		"                    ", // 12(D4)
		"        ------ ",      // 13(C4)
		"                    ", // 14(B4)                                  A & B are ::: pattern breakers
		"  ------------------", // 15(A4)  4th octave starts on A and runs through G ::: pattern breaker A
		"                    ", // 16 (G3)
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

	// ::: obtain an index into the lines of the staff (counting begins at 0)
	lineIndex, exists := pitchMap[note.Pitch] // use a random note as an index|key into pitchMap
	// handle error
	// fmt.Printf("\nlineIndex: %d, exists: %t\n", lineIndex, exists)
	if !exists { // if the value of exists is false
		fmt.Printf("\nlineIndex: %d, exists-not: %t\n", lineIndex, exists)
	}

	if prompting {
		/*
				When prompting, we need to print the lines of the staff one at a time so that we can
			insert an extra line between the treble and bass clefts.
		*/
		// place the note on the staff -- ::: insert the note within the indexed line
		staff[lineIndex] = staff[lineIndex][:10] + Red + "●" + Reset + staff[lineIndex][11:]
		lineCounter := 0
		for _, oneOfThe25Lines := range staff {
			fmt.Printf("%s\n", oneOfThe25Lines) // Print one line, indexed
			lineCounter++
			if lineCounter > 24 {
				return
			} else if lineCounter == 14 {
				fmt.Println()
			}
		}
		// When correct, immediately reprint the staff with the note in green such that it appears that the note has changed color...
		// ... and includes note.Pitch -- the player should adjust the terminal view scale to allow full view of prior staff
	} else if correct { // ::: Correct case
		staff[lineIndex] = staff[lineIndex][:10] + Green + "●" + note.Pitch + Reset + staff[lineIndex][11:]
		fmt.Println(strings.Join(staff, "\n")) // Print all the lines in one go
		tryThatAgain = false
	} else { // ::: Wrong case
		// When wrong, reprint the staff with the note in yellow such that it appears that the note has changed color ...
		// ... and, add the correct note, e.g., E5, next to the yellow note to inform the player of the right answer.
		// ::: first we setup the staff slice ...
		if lineIndex == 14 {
			staff[lineIndex] = staff[lineIndex][:10] + colorYellow + "● " + note.Pitch + " A breaks the pattern" + Reset + staff[lineIndex][11:]
		} else if lineIndex == 13 {
			staff[lineIndex] = staff[lineIndex][:10] + colorYellow + "● " + note.Pitch + " B breaks the pattern" + Reset + staff[lineIndex][11:]
		} else {
			staff[lineIndex] = staff[lineIndex][:10] + colorYellow + "● " + note.Pitch + " " + Reset + staff[lineIndex][11:]
		}
		// ::: ... then we print all the lines, one at a time
		lineCounter := 0
		for _, oneOfThe25Lines := range staff {
			fmt.Printf("%s\n", oneOfThe25Lines) // Print one line, indexed
			lineCounter++
			if lineCounter > 24 {
				tryThatAgain = true
				return
			}
		}
	}
	return
}
