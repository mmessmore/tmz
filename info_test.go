package main

import (
	"strings"
	"testing"
)

func TestVersionText(t *testing.T) {
	want := "tmz " + version + "\n"
	if got := versionText(); got != want {
		t.Errorf("versionText() = %q, want %q", got, want)
	}
	if version != "1.0.0" {
		t.Errorf("version = %q, want %q", version, "1.0.0")
	}
}

func TestUsageText(t *testing.T) {
	got := usageText()
	for _, want := range []string{
		"Usage: tmz [OPTION]... [TMUX-ARGUMENT]...",
		"-h, --help",
		"-V, --version",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("usageText() missing %q\n---\n%s", want, got)
		}
	}
	if !strings.HasSuffix(got, "\n") {
		t.Errorf("usageText() should end with a newline")
	}
}

func TestInfoFlag(t *testing.T) {
	tests := []struct {
		name string
		in   []string
		want infoKind
	}{
		{"no args", nil, infoNone},
		{"passthrough", []string{"ls"}, infoNone},
		{"-h", []string{"-h"}, infoHelp},
		{"--help", []string{"--help"}, infoHelp},
		{"-V", []string{"-V"}, infoVersion},
		{"--version", []string{"--version"}, infoVersion},
		{"lowercase -v is not version", []string{"-v"}, infoNone},
		{"only first position", []string{"attach", "-h"}, infoNone},
		{"help wins with trailing args", []string{"-h", "foo"}, infoHelp},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := infoFlag(tt.in); got != tt.want {
				t.Errorf("infoFlag(%q) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}
