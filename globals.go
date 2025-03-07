package main

var pitchesAll = []string{
	"A5", "G5", "F5", "E5", "D5", "C5",
	"B4", "A4", "G4", "F4", "E4", "D4", "C4",
	"B3", "A3", "G3", "F3", "E3", "D3", "C3",
	"B2", "A2", "G2", "F2",
}
var pitchesLeft = []string{
	"B3", "A3", "G3", "F3", "E3", "D3", "C3",
	"B2", "A2", "G2", "F2",
}
var pitchesRight = []string{
	"A5", "G5", "F5", "E5", "D5", "C5",
	"B4", "A4", "G4", "F4", "E4", "D4", "C4",
}

// First we make an slice of strings, one string for each staff line
// staff := []string{
var showStaff = []string{
	"        ---a--- \n",      // 1 (A5)
	"              G      \n", // 2 (G5)
	"  ---------f---------\n", // 3 (F5)
	"              E      \n", // 4 (E5)
	"  ---------d---------\n", // 5 (D5)
	"              C      \n", // 6 (C5) --------- 5th octave starts on C and runs through B
	"  ---------b---------\n", // 7 (B4)
	"              A      \n", // 8 (A4)
	"  ---------g---------\n", // 9 (G4)
	"              F      \n", // 10(F4)
	"  ---------e---------\n", // 11(E4)
	"              D      \n", // 12(D4)
	"        ---c--- \n\n",    // 13(C4) ======  4th octave starts on C and runs through B
	"              B      \n", // 14(B3)  A & B are ::: two steps up from G would be B
	"  ---------a---------\n", // 15(A3)           ::: two steps up from F would be A
	"              G      \n", // 16 (G3)
	"  ---------f---------\n", // (F3)
	"              E      \n", // (E3)
	"  ---------d---------\n", // (D3)
	"              C      \n", // (C3) --------- 3erd octave starts on C and runs through B
	"  ---------b---------\n", // (B2)
	"              A      \n", // (A2)
	"  ---------g---------\n", // (G2)
	"              F      \n", // (F2)
}

// Create two global int vars, initialized to 0
var correct, total, rememberLastPickL, rememberLastPickR, rememberLastPickAll int

var tryThatAgain, left, right, givenCreditForCorrectAnswer bool

var answer string

// ANSI color codes being assigned to constant string-like user-defined reserved words.
const (
	Green     = "\033[32m"
	Red       = "\033[31m"
	LightBlue = "\033[94m"
	Reset     = "\033[0m"
)
const colorCyan = "\033[36m"   // ::: - -
const colorPurple = "\033[35m" // ::: - -
const colorYellow = "\033[33m"

/*
Three reasons that justify deploying ::: the following Note struct:
	1. later we may want to add fields to Note (e.g., duration like "quarter" or "half"), a struct is easier to extend.
Strings would require bigger changes.
    2. Type Safety: properly used, Note signals “this is a musical note,” not just any string. A string is more generic and could lead
to mistakes in a larger app (e.g., passing a random string like "hello"). e.g. note.Pitch signals thusly.
    3.  Most important is Readability: note.Pitch is self-explanatory; pitch alone might be less obvious in a bigger project.
*/

// Note represents a musical note; see below for why it is "useful" to deploy this simple struct in place of a sting
type Note struct { // ::: - -
	Pitch string // e.g., "C5"
}

/*
.
The following two statements provide for the logging of statistics relative to the player's performance with staff notations.
*/

// NoteStats is a structure that tracks performance per note
type NoteStats struct { // ::: - -
	Attempts       int     // Total tries tally
	Misses         int     // Wrong answers, tally
	TotalCorrectMs int64   // sum of Correct answer times (in milliseconds)
	AvgCorrectSec  float64 // calculated Average time for correct answers (in seconds)
	CorrectCount   int     // Number of correct answers (for calculating the average)
}

// A map is a table of  Key:value pairs (and is always a reference type). Here the Key will be like C4 or F5
// create stats map, which will be a correspondence between a pitch (a string var used as a Key) and NoteStats structures
var stats = make(map[string]NoteStats) // ::: - -

/*
.
These last two statements merely allow for the logging and reporting of particularly untimely player interactions.
*/

// Outlier is a structure that will be used to record each outlier event ...
type Outlier struct { // used directly only in two append statements, else below ::: - -
	Pitch      string  // its pitch
	WasCorrect bool    // if it was correct or not
	TimeSec    float64 // Time taken by the user to respond, in seconds
}

// ... outliers is a slice (array) of those structures; and is the primary way to access Outlier data
var outliers []Outlier // made an empty initialized array (slice) of Outlier structures // ::: - -
