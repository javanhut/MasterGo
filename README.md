# golings

A rustlings-style tour of Go. Each exercise is a small broken or incomplete
program. Fix it until it compiles / passes its test, then remove the
`// I AM NOT DONE` marker to advance.

## Setup

```fish
go build -o golings ./cmd/golings
```

(Or just `go run ./cmd/golings <command>`.)

## Usage

```fish
./golings watch        # main loop: re-runs the next pending exercise on save
./golings list         # show all exercises and their status
./golings run intro1      # run a single exercise
./golings hint intro1     # print the hint
./golings solution intro1 # print the reference solution from solutions/
./golings reset intro1    # restore the I-AM-NOT-DONE marker
./golings verify          # run every exercise once; non-zero exit if any fail
```

The `info.json` at the project root is the registry — name, path, mode
(`compile` / `run` / `test`), and hint for each exercise. Edit it to add
your own.

## How an exercise is "done"

1. It builds / runs / tests successfully (depending on its `mode`).
2. The line `// I AM NOT DONE` has been removed from the source file.

`golings watch` will keep re-running the same exercise until both are true,
then move on to the next.

## Layout

```
cmd/golings/        the CLI
exercises/          start here — each is broken/incomplete
  00_intro/
  01_variables/
  ...
solutions/          reference solutions, mirroring exercises/ — peek when stuck
info.json           exercise registry
```
