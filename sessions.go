package main

import (
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

// listSessions returns the names of existing tmux sessions, sorted
// alphabetically. If no tmux server is running (or it has no sessions), it
// returns an empty slice and a nil error — that is treated as "nothing to
// choose from", not a failure.
func listSessions(tmuxPath string) ([]string, error) {
	out, err := exec.Command(tmuxPath, "list-sessions", "-F", "#{session_name}").Output()
	if err != nil {
		return nil, nil
	}

	trimmed := strings.TrimSpace(string(out))
	if trimmed == "" {
		return []string{}, nil
	}

	names := strings.Split(trimmed, "\n")
	sort.Strings(names)
	return names, nil
}

// quoteTarget wraps name in single quotes for use inside a tmux command
// string, escaping any single quotes it contains.
func quoteTarget(name string) string {
	return "'" + strings.ReplaceAll(name, "'", `'\''`) + "'"
}

// chooserArgs builds the tmux argument list for a single chained invocation
// that attaches to the most-recently-used session and immediately overlays a
// menu offering each of sessions plus "New session". Sessions are shown in
// the order given, each keyed 1-9 (unlabeled beyond the ninth, but still
// reachable with the arrow keys); "New session" is pinned at the top under
// key 'n' and starts a plain new session, attaching the current client to it.
func chooserArgs(sessions []string) []string {
	args := []string{
		"attach-session", ";",
		"display-menu", "-T", "Sessions",
		"New session", "n", "new-session",
	}
	for i, name := range sessions {
		key := ""
		if i < 9 {
			key = strconv.Itoa(i + 1)
		}
		args = append(args, name, key, "switch-client -t "+quoteTarget(name))
	}
	return args
}
