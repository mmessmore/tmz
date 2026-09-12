# tmz

`tmz` is a thin wrapper around `tmux` for people whose fingers still type
`screen` commands. It forwards every argument to `tmux` unchanged, except that
it recognizes `screen`'s attach flags as the first argument and rewrites them
into the equivalent `tmux attach` invocation. Once the arguments are settled it
`exec`s `tmux`, replacing itself, so `tmz` leaves nothing behind in the process
tree.

## Usage

```
tmz [screen-style attach flags | tmux arguments ...]
```

Run it exactly as you would run `tmux`:

```sh
tmz                     # exec: tmux
tmz new -s work         # exec: tmux new -s work
tmz ls                  # exec: tmux ls
tmz -L mysock attach    # exec: tmux -L mysock attach
```

Running `tmz` bare (no arguments) when one or more tmux sessions already
exist shows a menu of those sessions instead of blindly starting a new one,
with an option to start a new session too. Pick a session to attach to it, or
"New session" to create one; `Escape`/`q` cancels. If no tmux server is
running yet, `tmz` behaves as before and just starts tmux normally.

If the **first** argument is one of `screen`'s attach flags, it is translated:

| You type            | `tmz` runs                  | Meaning                              |
| ------------------- | --------------------------- | ------------------------------------ |
| `tmz -r`            | `tmux attach`               | reattach to a session               |
| `tmz -r work`       | `tmux attach -t work`       | reattach to session `work`          |
| `tmz -x`            | `tmux attach`               | attach (shared; tmux's default)      |
| `tmz -x work`       | `tmux attach -t work`       | attach to session `work` (shared)   |
| `tmz -dr`           | `tmux attach -d`            | detach elsewhere, then attach here   |
| `tmz -dr main`      | `tmux attach -d -t main`    | same, for session `main`            |
| `tmz -d -r main`    | `tmux attach -d -t main`    | same, flags given separately         |

The token after the flag is treated as a session name only when it does not
start with `-`; anything else (including further flags and trailing arguments)
is passed through to `tmux attach` untouched. Screen-style flags are recognized
only in the first position — anywhere else they are left alone and handed to
`tmux` as-is.

`tmz` handles two options itself instead of passing them to `tmux`, and only
when they appear as the **first** argument:

| You type         | `tmz` does                          |
| ---------------- | ----------------------------------- |
| `tmz -h`, `--help`    | print usage and exit            |
| `tmz -V`, `--version` | print `tmz <version>` and exit  |

`-V` shadows `tmux`'s own version flag; `tmux`'s other options, including `-v`
(verbose logging), are untouched.

To hand `tmux` an argument that `tmz` would otherwise claim, put `--` first:
everything after it is passed through verbatim with no translation.

```sh
tmz -- -V        # exec: tmux -V   (tmux's version)
tmz -- -h        # exec: tmux -h
tmz --           # exec: tmux
```

`tmz` exits with status 1 if `tmux` cannot be found in `PATH`.

When run from a terminal, `tmz` also sets the terminal's title to the local
host's short name (the hostname up to the first `.`, like `uname -n | cut -d.
-f1`), before handing off to `tmux`.

## Building

`tmz` is a single Go module with no dependencies. You need Go 1.26 or newer.

### The Go toolchain

Nothing special is required — use `go` in the usual way.

```sh
go build -o tmz .                        # produces ./tmz
go install github.com/mmessmore/tmz@latest   # build and drop it in $GOBIN
go test ./...                            # run the tests
```

### Using `make`

The `Makefile` wraps the same commands and adds a few conveniences.

| Target        | What it does                                                                  |
| ------------- | ---------------------------------------------------------------------------- |
| `make`        | Build `./tmz` (alias for `make build`).                                     |
| `make test`   | Run `go test ./...`.                                                        |
| `make lint`   | Check formatting with `gofmt -l` and run `go vet ./...`.                    |
| `make fmt`    | Reformat the source with `go fmt ./...`.                                    |
| `make install`| Build, then `install -m 755` the binary into `DEST`.                        |
| `make dist`   | Cross-compile static release binaries into `dist/` for macOS, Linux, and Windows. |
| `make clean`  | Remove `./tmz` and `dist/`.                                                 |

`make install` picks `DEST` automatically: the first of `$HOME/bin`,
`$HOME/.local/bin`, or `/usr/local/bin` that exists. Override it explicitly
when needed:

```sh
make install DEST=~/.local/bin
```
