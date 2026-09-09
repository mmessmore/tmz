package main

import "strings"

// translate maps a tmz argument list to the argument list for tmux.
//
// If the first argument is one of screen's attach flags (-r, -x, -dr, -rd, or
// -d and -r as separate tokens) it is rewritten into a "tmux attach" invocation.
// A following token that does not start with "-" is treated as the session name
// and becomes "-t <name>". Any remaining arguments are appended untouched.
//
// Anything that is not a recognized screen attach flag is returned unchanged and
// passed straight through to tmux.
//
// As a special case, a leading "--" is dropped and everything after it is
// returned verbatim, with no translation. This is the escape hatch for handing
// tmux an argument that tmz would otherwise claim, e.g. "tmz -- -V".
func translate(args []string) []string {
	if len(args) == 0 {
		return []string{}
	}

	if args[0] == "--" {
		return args[1:]
	}

	first := args[0]
	rest := args[1:]

	recognized := false
	detach := false

	switch first {
	case "-r", "-x":
		recognized = true
	case "-dr", "-rd":
		recognized = true
		detach = true
	case "-d":
		if len(rest) > 0 && rest[0] == "-r" {
			recognized = true
			detach = true
			rest = rest[1:]
		}
	}

	if !recognized {
		return args
	}

	// Handle "screen -r -d" (detach given as a separate trailing token).
	if first == "-r" && len(rest) > 0 && rest[0] == "-d" {
		detach = true
		rest = rest[1:]
	}

	out := []string{"attach"}
	if detach {
		out = append(out, "-d")
	}
	if len(rest) > 0 && !strings.HasPrefix(rest[0], "-") {
		out = append(out, "-t", rest[0])
		rest = rest[1:]
	}
	out = append(out, rest...)
	return out
}
