package main

// ANSI color codes
const (
	Green     = "\033[32m"
	Red       = "\033[31m"
	LightBlue = "\033[94m"
	Reset     = "\033[0m"
)
const colorCyan = "\033[36m"
const colorPurple = "\033[35m"
const colorYellow = "\033[33m"

// Note represents a musical note; see below for why it is "useful" to deploy this simple struct in place of a sting
type Note struct {
	Pitch string // e.g., "C5"
}

/*
     Three reasons that justify deploying the above Note struct: later we may want to add fields to Note (e.g., duration
like "quarter" or "half"), a struct is easier to extend. Strings would require bigger changes.
     Type Safety: properly used, Note signals “this is a musical note,” not just any string. A string is more generic and could lead to
mistakes in a larger app (e.g., passing a random string like "hello"). e.g. note.Pitch signals thusly.
     Most important is Readability: note.Pitch is self-explanatory; pitch alone might be less obvious in a bigger project.
*/

// NoteStats is a structure that tracks performance per note
type NoteStats struct {
	Attempts       int     // Total tries tally
	Misses         int     // Wrong answers, tally
	TotalCorrectMs int64   // sum of Correct answer times (in milliseconds)
	AvgCorrectSec  float64 // calculated Average time for correct answers (in seconds)
	CorrectCount   int     // Number of correct answers (for calculating the average)
}

// Outlier is a structure that will be used to records each outlier event
type Outlier struct {
	Pitch      string  // its pitch
	WasCorrect bool    // if it was correct or not
	TimeSec    float64 // Time taken by the user to respond, in seconds
}

// making a map initializes it, perfectly-well
// create stats map, which will be a correspondence between pitch (the string) and NoteStats structures
var stats = make(map[string]NoteStats)

// making an empty initialized array (slice) of Outlier structures (called outliers)
var outliers []Outlier

// creating two global int vars, initialized to 0
var correct, total int
