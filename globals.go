package main

// Create two global int vars, initialized to 0
var correct, total, rememberLastPick int

var tryThatAgain bool

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
