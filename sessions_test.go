package main

import (
	"reflect"
	"testing"
)

func TestQuoteTarget(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"plain", "work", "'work'"},
		{"single quote", "o'brien", `'o'\''brien'`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := quoteTarget(tt.in); got != tt.want {
				t.Errorf("quoteTarget(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestChooserArgs(t *testing.T) {
	tests := []struct {
		name     string
		sessions []string
		want     []string
	}{
		{
			"no sessions",
			nil,
			[]string{
				"attach-session", ";",
				"display-menu", "-T", "Sessions",
				"New session", "n", "new-session",
			},
		},
		{
			"one session",
			[]string{"work"},
			[]string{
				"attach-session", ";",
				"display-menu", "-T", "Sessions",
				"New session", "n", "new-session",
				"work", "1", "switch-client -t 'work'",
			},
		},
		{
			"multiple sessions keep given order and number sequentially",
			[]string{"alpha", "beta"},
			[]string{
				"attach-session", ";",
				"display-menu", "-T", "Sessions",
				"New session", "n", "new-session",
				"alpha", "1", "switch-client -t 'alpha'",
				"beta", "2", "switch-client -t 'beta'",
			},
		},
		{
			"session name needing quote escaping",
			[]string{"o'brien"},
			[]string{
				"attach-session", ";",
				"display-menu", "-T", "Sessions",
				"New session", "n", "new-session",
				"o'brien", "1", `switch-client -t 'o'\''brien'`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := chooserArgs(tt.sessions); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("chooserArgs(%q) = %q, want %q", tt.sessions, got, tt.want)
			}
		})
	}
}

func TestChooserArgsKeyNumberingBeyondNine(t *testing.T) {
	sessions := make([]string, 10)
	for i := range sessions {
		sessions[i] = string(rune('a' + i))
	}
	args := chooserArgs(sessions)

	// Each session contributes 3 args, after the 8 fixed header args.
	const header = 8
	for i, name := range sessions {
		base := header + i*3
		if got := args[base]; got != name {
			t.Fatalf("session %d: name = %q, want %q", i, got, name)
		}
		wantKey := ""
		if i < 9 {
			wantKey = string(rune('1' + i))
		}
		if got := args[base+1]; got != wantKey {
			t.Errorf("session %d (%q): key = %q, want %q", i, name, got, wantKey)
		}
	}
}
