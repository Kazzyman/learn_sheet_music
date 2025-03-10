package main

import (
	"fmt"
)

// DrawStaff generates a twelve-line, twelve-space ASCII staff with a note placed in Red
func DrawStaff(note Note, prompting bool, correct bool) { // accepts a simple structure of notes ::: - -

	// First we make an slice of strings, one string for each staff line
	staff := []string{
		"        ------ ",      // 1 (A5)
		"                    ", // 2 (G5)
		"  ------------------", // 3 (F5)
		"                    ", // 4 (E5)
		"  ------------------", // 5 (D5)
		"                    ", // 6 (C5) --------- 5th octave starts on C and runs through B
		"  ------------------", // 7 (B4)
		"                    ", // 8 (A4)
		"  ------------------", // 9 (G4)
		"                    ", // 10(F4)
		"  ------------------", // 11(E4)
		"                    ", // 12(D4)
		"        ------ ",      // 13(C4) ======  4th octave starts on C and runs through B
		"                    ", // 14(B3)  A & B are ::: pattern breakers
		"  ------------------", // 15(A3)           ::: pattern breaker A
		"                    ", // 16 (G3)
		"  ------------------", // (F3)
		"                    ", // (E3)
		"  ------------------", // (D3)
		"                    ", // (C3) --------- 3erd octave starts on C and runs through B
		"  ------------------", // (B2)
		"                    ", // (A2)
		"  ------------------", // (G2)
		"                    ", // (F2)
	}
	pitchMap := map[string]int{
		"A5": 0, "G5": 1, "F5": 2, "E5": 3, "D5": 4, "C5": 5,
		"B4": 6, "A4": 7, "G4": 8, "F4": 9, "E4": 10, "D4": 11, "C4": 12,
		"B3": 13, "A3": 14, "G3": 15, "F3": 16, "E3": 17, "D3": 18, "C3": 19,
		"B2": 20, "A2": 21, "G2": 22, "F2": 23,
	}
	lineMap := map[int]string{
		0: "A5", 1: "G5", 2: "F5", 3: "E5", 4: "D5", 5: "C5",
		6: "B4", 7: "A4", 8: "G4", 9: "F4", 10: "E4", 11: "D4", 12: "C4",
		13: "B3", 14: "A3", 15: "G3", 16: "F3", 17: "E3", 18: "D3", 19: "C3",
		20: "B2", 21: "A2", 22: "G2", 23: "F2",
	}

	// ::: obtain an index into the lines of the staff (counting begins at 0)
	lineIndex, exists := pitchMap[note.Pitch] // use a random note as an index|key into pitchMap ...
	// ... this finds the line of the staff that the note.Pitch should go on
	// handle error
	// fmt.Printf("\nlineIndex: %d, exists: %t\n", lineIndex, exists)
	if !exists { // if the value of exists is false
		fmt.Printf("\nlineIndex: %d, exists-not: %t\n", lineIndex, exists)
	}

	if prompting {
		if !trainingWheels {
			fmt.Println(showStaff)
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
				} else if lineCounter == 13 { // When displaying the staff as a prompt, add a blank space between treble and bass clefts.
					fmt.Println()
				}
			}
			// When correct, immediately reprint the staff with the note in green such that it appears that the note has changed color...
			// ... and includes note.Pitch -- the player should adjust the terminal view scale to allow full view of prior staff
		} else {
			fmt.Println(showStaff)
			/*
					When prompting, we need to print the lines of the staff one at a time so that we can
				insert an extra line between the treble and bass clefts.
			*/
			// place the note on the staff -- ::: insert the note within the indexed line
			staff[lineIndex] = staff[lineIndex][:10] + Red + "●" + Reset + staff[lineIndex][11:]
			lineCounter := 0
			for _, oneOfThe25Lines := range staff {
				indexToPitch, _ := lineMap[lineCounter]
				fmt.Printf("%s%s\n", oneOfThe25Lines, indexToPitch) // Print one line, indexed
				lineCounter++
				if lineCounter > 24 {
					return
				} else if lineCounter == 13 { // When displaying the staff as a prompt, add a blank space between treble and bass clefts.
					fmt.Println()
				}
			}
		}
	} else if correct { // ::: Correct case !!
		staff[lineIndex] = staff[lineIndex][:10] + Green + "●" + note.Pitch + Reset + staff[lineIndex][11:]
		lineCounter := 0
		for _, oneOfThe25Lines := range staff {
			fmt.Printf("%s\n", oneOfThe25Lines) // Print one line, indexed
			lineCounter++
			if lineCounter > 24 {
				return
			} else if lineCounter == 13 { // When displaying the staff as a prompt, add a blank space between treble and bass clefts.
				fmt.Println()
			}
		}
		tryThatAgain = false
	} else { // ::: Wrong case !!
		// When wrong, reprint the staff with the note in yellow such that it appears that the note has changed color ...
		// ... and, add the correct note, e.g., E5, next to the yellow note to inform the player of the right answer.
		// ::: -------- First, we setup the staff slice ... ----------------------------------------------------------------------------
		staff[lineIndex] = staff[lineIndex][:10] + colorYellow + "● " + note.Pitch + " " + Reset + staff[lineIndex][11:]

		// ::: ... then, we print all the lines, one at a time
		lineCounter := 0
		for _, oneOfThe25Lines := range staff {
			fmt.Printf("%s\n", oneOfThe25Lines) // Print one line, per indexed
			lineCounter++
			if lineCounter >= 24 {
				tryThatAgain = true
				fmt.Printf("%sYour guess was: %s  %s", Red, answer, Reset) // ::: fix to answer
				return
			} else if lineCounter == 13 { // When displaying the staff as a prompt, add a blank space between treble and bass clefts.
				fmt.Println()
			}
		}
	}
	return
}
