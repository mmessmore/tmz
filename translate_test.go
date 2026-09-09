package main

import (
	"reflect"
	"testing"
)

func TestTranslate(t *testing.T) {
	tests := []struct {
		name string
		in   []string
		want []string
	}{
		{"no args", nil, []string{}},
		{"empty slice", []string{}, []string{}},

		{"screen -r", []string{"-r"}, []string{"attach"}},
		{"screen -r with session", []string{"-r", "work"}, []string{"attach", "-t", "work"}},
		{"screen -x", []string{"-x"}, []string{"attach"}},
		{"screen -x with session", []string{"-x", "work"}, []string{"attach", "-t", "work"}},

		{"screen -dr", []string{"-dr"}, []string{"attach", "-d"}},
		{"screen -rd", []string{"-rd"}, []string{"attach", "-d"}},
		{"screen -dr with session", []string{"-dr", "main"}, []string{"attach", "-d", "-t", "main"}},
		{"screen -d -r", []string{"-d", "-r"}, []string{"attach", "-d"}},
		{"screen -r -d", []string{"-r", "-d"}, []string{"attach", "-d"}},
		{"screen -d -r with session", []string{"-d", "-r", "main"}, []string{"attach", "-d", "-t", "main"}},

		{"trailing flag is not a session name", []string{"-r", "-c", "cmd"}, []string{"attach", "-c", "cmd"}},
		{"extra args after session preserved", []string{"-r", "work", "-E"}, []string{"attach", "-t", "work", "-E"}},

		{"-- passes the rest through untouched", []string{"--", "-V"}, []string{"-V"}},
		{"-- alone", []string{"--"}, []string{}},
		{"-- shields a screen-style flag", []string{"--", "-r", "work"}, []string{"-r", "work"}},
		{"-- only consumed in first position", []string{"attach", "--", "-x"}, []string{"attach", "--", "-x"}},

		{"raw tmux subcommand", []string{"new", "-s", "foo"}, []string{"new", "-s", "foo"}},
		{"raw tmux attach", []string{"attach", "-t", "foo"}, []string{"attach", "-t", "foo"}},
		{"raw tmux global flag", []string{"-L", "sock", "ls"}, []string{"-L", "sock", "ls"}},
		{"unrelated leading flag passed through", []string{"-2"}, []string{"-2"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := translate(tt.in)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("translate(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
