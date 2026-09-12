package main

// version is the current release of tmz.
const version = "1.0.0"

// infoKind identifies an informational flag that tmz handles itself instead of
// passing on to tmux.
type infoKind int

const (
	infoNone infoKind = iota
	infoHelp
	infoVersion
)

// infoFlag reports whether the first argument is one of tmz's own
// informational flags (-h/--help or -V/--version). Like the screen attach
// flags, these are recognized only in the first position. -V (not -v) is used
// for the version so it shadows tmux's own version flag rather than its
// verbose-logging flag.
func infoFlag(args []string) infoKind {
	if len(args) == 0 {
		return infoNone
	}
	switch args[0] {
	case "-h", "--help":
		return infoHelp
	case "-V", "--version":
		return infoVersion
	default:
		return infoNone
	}
}

// versionText is the output of "tmz -V", GNU style.
func versionText() string {
	return "tmz " + version + "\n"
}

// usageText is the output of "tmz -h", GNU style.
func usageText() string {
	return `Usage: tmz [OPTION]... [TMUX-ARGUMENT]...
Thin tmux wrapper that also understands screen's attach flags.

Arguments are passed to tmux unchanged, except that a screen-style attach
flag in the first position is rewritten to a tmux attach invocation:

  -r  [session]       reattach              ->  tmux attach [-t session]
  -x  [session]       shared attach         ->  tmux attach [-t session]
  -dr, -rd [session]  detach then attach    ->  tmux attach -d [-t session]
  -d -r [session]     same, flags separate  ->  tmux attach -d [-t session]

Options:
  -h, --help     display this help and exit
  -V, --version  output version information and exit

Anything else is run as 'tmux <arguments>'.  See tmux(1) for tmux's own
options and commands.
`
}
