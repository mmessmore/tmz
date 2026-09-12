package main

import "testing"

func TestShortHostname(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"bare hostname", "myhost", "myhost"},
		{"fqdn", "myhost.example.com", "myhost"},
		{"single trailing dot", "myhost.", "myhost"},
		{"empty", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shortHostname(tt.in); got != tt.want {
				t.Errorf("shortHostname(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestTitleEscape(t *testing.T) {
	want := "\033]2;myhost\007"
	if got := titleEscape("myhost"); got != want {
		t.Errorf("titleEscape(%q) = %q, want %q", "myhost", got, want)
	}
}
