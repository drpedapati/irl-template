// Package style provides consistent terminal styling for the IRL CLI.
package style

import "fmt"

// ANSI color codes
const (
	Reset     = "\033[0m"
	Bold      = "\033[1m"
	Dim       = "\033[2m"
	Cyan      = "\033[36m"
	Green     = "\033[32m"
	Yellow    = "\033[33m"
	Red       = "\033[31m"
	BoldCyan  = "\033[1;36m"
	BoldGreen = "\033[1;32m"
)

// Semantic styles
func Header(s string) string  { return BoldCyan + s + Reset }
func Success(s string) string { return Green + s + Reset }
func Warn(s string) string    { return Yellow + s + Reset }
func Err(s string) string     { return Red + s + Reset }
func Cmd(s string) string     { return Cyan + s + Reset }
func Faint(s string) string   { return Dim + s + Reset }
func B(s string) string       { return Bold + s + Reset }

// Common symbols
const (
	Check = "✓"
	Cross = "✗"
	Dot   = "•"
	Arrow = "→"
)

// Prefixed output helpers
func OK(msg string) string   { return Green + Check + Reset + " " + msg }
func Fail(msg string) string { return Red + Cross + Reset + " " + msg }
func Item(msg string) string { return Cyan + Dot + Reset + " " + msg }

// Print helpers
func PrintHeader(s string)  { fmt.Println(Header(s)) }
func PrintOK(s string)      { fmt.Println(OK(s)) }
func PrintFail(s string)    { fmt.Println(Fail(s)) }
func PrintItem(s string)    { fmt.Println(Item(s)) }
func PrintCmd(label, c string) { fmt.Printf("  %s %s\n", Faint(label), Cmd(c)) }
