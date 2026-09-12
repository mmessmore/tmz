package main

import (
	"os"
	"strings"
)

// shortHostname trims a fully qualified hostname down to just the leading
// label, mirroring `uname -n | cut -d. -f1`.
func shortHostname(full string) string {
	name, _, _ := strings.Cut(full, ".")
	return name
}

// titleEscape returns the OSC escape sequence that sets the terminal window
// title to name.
func titleEscape(name string) string {
	return "\033]2;" + name + "\007"
}

// isTerminal reports whether f is attached to a terminal.
func isTerminal(f *os.File) bool {
	fi, err := f.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}

// setTerminalTitle sets the terminal's title to the local host's short name,
// writing to out. It is a no-op if out is not attached to a terminal or the
// hostname cannot be determined.
func setTerminalTitle(out *os.File) {
	if !isTerminal(out) {
		return
	}
	hostname, err := os.Hostname()
	if err != nil {
		return
	}
	out.WriteString(titleEscape(shortHostname(hostname)))
}
