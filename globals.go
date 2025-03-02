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

// Note represents a musical note
type Note struct {
	Pitch string // e.g., "C5"
}

// NoteStats is a structure that tracks performance per note
type NoteStats struct {
	Attempts       int     // Total tries
	Misses         int     // Wrong answers
	TotalCorrectMs int64   // Sum of correct answer times (milliseconds)
	AvgCorrectSec  float64 // Average time for correct answers (seconds)
	CorrectCount   int     // Number of correct answers (for averaging)
}

// Outlier is a structure that will be used to records each outlier event
type Outlier struct {
	Pitch      string  // its pitch
	WasCorrect bool    // and if it was correct or not
	TimeSec    float64 // Time taken by the user to respond, in seconds
}

// making a map initializes it well-enough
// create stats map, which will be a correspondence between pitch (the string) and the NoteStats struct
var stats = make(map[string]NoteStats)

// making an array of structures (called outliers) initializes it well-enough
var outliers []Outlier

// creating two global int vars, initializes them well-enough
var correct, total int
