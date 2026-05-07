package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

const (
	ansiReset = "\033[0m"
	ansiBold  = "\033[1m"
	ansiGreen = "\033[32m"
	ansiRed   = "\033[31m"
	ansiDim   = "\033[2m"
)

const notDoneMarker = "I AM NOT DONE"

type Exercise struct {
	Name string `json:"name"`
	Path string `json:"path"` // package directory, e.g. exercises/00_intro/intro1
	File string `json:"file"` // file the user edits
	Mode string `json:"mode"` // compile | run | test
	Hint string `json:"hint"`
}

type info struct {
	Welcome   string     `json:"welcome"`
	Exercises []Exercise `json:"exercises"`
}

func loadInfo() (*info, error) {
	data, err := os.ReadFile("info.json")
	if err != nil {
		return nil, fmt.Errorf("read info.json (run golings from the project root): %w", err)
	}
	var i info
	if err := json.Unmarshal(data, &i); err != nil {
		return nil, fmt.Errorf("parse info.json: %w", err)
	}
	return &i, nil
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	inf, err := loadInfo()
	if err != nil {
		die(err)
	}
	switch os.Args[1] {
	case "list":
		cmdList(inf)
	case "run":
		mustArg("run")
		cmdRun(inf, os.Args[2])
	case "hint":
		mustArg("hint")
		cmdHint(inf, os.Args[2])
	case "solution":
		mustArg("solution")
		cmdSolution(inf, os.Args[2])
	case "reset":
		mustArg("reset")
		cmdReset(inf, os.Args[2])
	case "verify":
		cmdVerify(inf)
	case "watch":
		cmdWatch(inf)
	case "help", "-h", "--help":
		usage()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", os.Args[1])
		usage()
		os.Exit(2)
	}
}

func usage() {
	fmt.Println(`golings — a rustlings-style tour of Go.

Commands:
  watch              Run exercises in order, re-running when files change.
  run    <name>      Run a single exercise by name.
  hint   <name>      Print the hint for an exercise.
  solution <name>    Print the path and contents of the reference solution.
  list               List all exercises and their status.
  verify             Run every exercise once; exit non-zero if any fail.
  reset  <name>      Restore the "I AM NOT DONE" marker on an exercise.
  help               Show this message.

An exercise is "done" once it builds/passes AND the line containing
"I AM NOT DONE" has been removed.`)
}

func mustArg(cmd string) {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "usage: golings %s <name>\n", cmd)
		os.Exit(2)
	}
}

func die(err error) {
	fmt.Fprintln(os.Stderr, "error:", err)
	os.Exit(1)
}

func find(inf *info, name string) (*Exercise, error) {
	for i := range inf.Exercises {
		if inf.Exercises[i].Name == name {
			return &inf.Exercises[i], nil
		}
	}
	return nil, fmt.Errorf("no exercise named %q", name)
}

// hasMarker is true iff the file contains a line that, after trimming spaces,
// is exactly `// I AM NOT DONE`. Substring matches inside other comments
// (e.g. an explanation that mentions the marker) do not count.
func hasMarker(file string) (bool, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return false, err
	}
	want := "// " + notDoneMarker
	for _, line := range strings.Split(string(data), "\n") {
		if strings.TrimSpace(line) == want {
			return true, nil
		}
	}
	return false, nil
}

// runExercise compiles/tests/runs the exercise and returns (output, success, error).
func runExercise(ex *Exercise) (string, bool, error) {
	pkg := "./" + filepath.ToSlash(ex.Path)
	var cmd *exec.Cmd
	switch ex.Mode {
	case "compile":
		cmd = exec.Command("go", "build", "-o", os.DevNull, pkg)
	case "run":
		cmd = exec.Command("go", "run", pkg)
	case "test":
		cmd = exec.Command("go", "test", "-count=1", pkg)
	default:
		return "", false, fmt.Errorf("exercise %s: unknown mode %q", ex.Name, ex.Mode)
	}
	out, err := cmd.CombinedOutput()
	return string(out), err == nil, nil
}

func cmdList(inf *info) {
	for _, ex := range inf.Exercises {
		status := "done"
		marker, err := hasMarker(ex.File)
		switch {
		case err != nil:
			status = "missing"
		case marker:
			status = "pending"
		}
		fmt.Printf("  [%-7s] %-14s %s (%s)\n", status, ex.Name, ex.File, ex.Mode)
	}
}

func cmdRun(inf *info, name string) {
	ex, err := find(inf, name)
	if err != nil {
		die(err)
	}
	runOne(ex, true)
}

func runOne(ex *Exercise, verbose bool) bool {
	if verbose {
		fmt.Printf("▶ Running %s (%s) ...\n", ex.Name, ex.Mode)
	}
	out, ok, err := runExercise(ex)
	if err != nil {
		die(err)
	}
	if out != "" && (verbose || !ok) {
		fmt.Print(out)
		if !strings.HasSuffix(out, "\n") {
			fmt.Println()
		}
	}
	marker, err := hasMarker(ex.File)
	if err != nil {
		die(err)
	}
	switch {
	case !ok:
		fmt.Printf("✗ %s failed to %s.\n", ex.Name, ex.Mode)
		return false
	case marker:
		fmt.Printf("✓ %s %ss successfully — now remove the `// %s` line.\n", ex.Name, ex.Mode, notDoneMarker)
		return false
	default:
		fmt.Printf("✓ %s done!\n", ex.Name)
		return true
	}
}

func cmdHint(inf *info, name string) {
	ex, err := find(inf, name)
	if err != nil {
		die(err)
	}
	if ex.Hint == "" {
		fmt.Println("(no hint available)")
		return
	}
	fmt.Println(ex.Hint)
}

func cmdSolution(inf *info, name string) {
	ex, err := find(inf, name)
	if err != nil {
		die(err)
	}
	// Mirror the exercise file path under solutions/.
	rel := strings.TrimPrefix(ex.File, "exercises/")
	path := filepath.Join("solutions", rel)
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "no reference solution at %s: %v\n", path, err)
		os.Exit(1)
	}
	fmt.Printf("// %s\n\n%s", path, data)
}

func cmdReset(inf *info, name string) {
	ex, err := find(inf, name)
	if err != nil {
		die(err)
	}
	data, err := os.ReadFile(ex.File)
	if err != nil {
		die(err)
	}
	s := string(data)
	if strings.Contains(s, notDoneMarker) {
		fmt.Printf("%s already has the marker.\n", ex.Name)
		return
	}
	// Insert marker on a new comment line just after the package declaration.
	lines := strings.Split(s, "\n")
	for i, l := range lines {
		if strings.HasPrefix(strings.TrimSpace(l), "package ") {
			lines = append(lines[:i+1], append([]string{"", "// " + notDoneMarker}, lines[i+1:]...)...)
			break
		}
	}
	if err := os.WriteFile(ex.File, []byte(strings.Join(lines, "\n")), 0644); err != nil {
		die(err)
	}
	fmt.Printf("Reset %s — marker restored.\n", ex.Name)
}

func cmdVerify(inf *info) {
	failed := 0
	for i := range inf.Exercises {
		ex := &inf.Exercises[i]
		marker, err := hasMarker(ex.File)
		if err != nil {
			fmt.Printf("✗ %s: %v\n", ex.Name, err)
			failed++
			continue
		}
		out, ok, err := runExercise(ex)
		if err != nil {
			die(err)
		}
		switch {
		case !ok:
			fmt.Printf("✗ %s: %s failed\n", ex.Name, ex.Mode)
			fmt.Print(out)
			failed++
		case marker:
			fmt.Printf("… %s: pending (marker still present)\n", ex.Name)
			failed++
		default:
			fmt.Printf("✓ %s\n", ex.Name)
		}
	}
	if failed > 0 {
		fmt.Printf("\n%d exercise(s) not yet done.\n", failed)
		os.Exit(1)
	}
	fmt.Println("\nAll exercises pass — nice work!")
}

// cmdWatch is the interactive loop: render progress + current exercise,
// run it, then wait for either a file change OR a single keypress
// (h: hint, l: list, c: check all, x: reset, q: quit).
func cmdWatch(inf *info) {
	saved, raw := enableRawMode()
	if raw {
		defer restoreTerminal(saved)
		// Restore terminal even on Ctrl-C.
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			<-sigCh
			restoreTerminal(saved)
			fmt.Println()
			os.Exit(0)
		}()
	}

	fsCh := startFsWatcher("exercises")
	keyCh := make(chan byte, 8)
	if raw {
		go readKeys(keyCh)
	}

	first := true
	for {
		clearScreen()
		if first && inf.Welcome != "" {
			fmt.Println(ansiDim + inf.Welcome + ansiReset)
			fmt.Println()
			first = false
		}
		ex := nextPending(inf)
		render(inf, ex)
		if ex == nil {
			fmt.Println("\n" + ansiGreen + "All exercises complete!" + ansiReset + " Run `golings verify` to double-check.")
			return
		}
		runOne(ex, false)
		printPrompt()

		// Race: file change OR keystroke.
		select {
		case <-fsCh:
			// loop -> redraw
		case k := <-keyCh:
			if !handleKey(k, inf, ex, keyCh) {
				return // quit
			}
		}
	}
}

// handleKey returns false if the loop should exit.
func handleKey(k byte, inf *info, ex *Exercise, keyCh <-chan byte) bool {
	switch k {
	case 'q', 3 /* Ctrl-C */ :
		fmt.Println()
		return false
	case 'h':
		fmt.Println()
		fmt.Println(ansiBold + "Hint:" + ansiReset + " " + ex.Hint)
		waitAnyKey(keyCh)
	case 'l':
		fmt.Println()
		cmdList(inf)
		waitAnyKey(keyCh)
	case 'c':
		fmt.Println()
		fmt.Println(ansiBold + "Checking all exercises..." + ansiReset)
		runVerifyInline(inf)
		waitAnyKey(keyCh)
	case 'x':
		fmt.Println()
		cmdReset(inf, ex.Name)
		waitAnyKey(keyCh)
	}
	return true
}

func waitAnyKey(keyCh <-chan byte) {
	fmt.Print(ansiDim + "\n(press any key to continue)" + ansiReset)
	<-keyCh
}

// runVerifyInline is like cmdVerify but does not call os.Exit.
func runVerifyInline(inf *info) {
	failed := 0
	for i := range inf.Exercises {
		ex := &inf.Exercises[i]
		marker, err := hasMarker(ex.File)
		if err != nil {
			fmt.Printf("%s✗%s %s: %v\n", ansiRed, ansiReset, ex.Name, err)
			failed++
			continue
		}
		_, ok, err := runExercise(ex)
		if err != nil {
			fmt.Println(err)
			failed++
			continue
		}
		switch {
		case !ok:
			fmt.Printf("%s✗%s %s: %s failed\n", ansiRed, ansiReset, ex.Name, ex.Mode)
			failed++
		case marker:
			fmt.Printf("%s…%s %s: marker still present\n", ansiDim, ansiReset, ex.Name)
			failed++
		default:
			fmt.Printf("%s✓%s %s\n", ansiGreen, ansiReset, ex.Name)
		}
	}
	if failed == 0 {
		fmt.Println("\n" + ansiGreen + "All exercises pass!" + ansiReset)
	} else {
		fmt.Printf("\n%d exercise(s) not yet done.\n", failed)
	}
}

// render prints the green progress bar plus the current-exercise header.
func render(inf *info, current *Exercise) {
	done, total := progress(inf)
	const width = 40
	filled := 0
	if total > 0 {
		filled = done * width / total
	}
	bar := ansiGreen + strings.Repeat("█", filled) + ansiReset + strings.Repeat("░", width-filled)
	pct := 0
	if total > 0 {
		pct = done * 100 / total
	}
	fmt.Printf("Progress: [%s] %s%d/%d%s  %d%%\n", bar, ansiBold, done, total, ansiReset, pct)
	if current != nil {
		fmt.Printf("Current exercise: %s%s%s\n\n", ansiBold, current.File, ansiReset)
	} else {
		fmt.Println()
	}
}

func printPrompt() {
	k := func(c, label string) string {
		return ansiBold + c + ansiReset + ":" + label
	}
	fmt.Printf("\n%s / %s / %s / %s / %s ? ",
		k("h", "hint"), k("l", "list"), k("c", "check all"), k("x", "reset"), k("q", "quit"))
}

// --- terminal raw mode (via stty; no third-party deps) -------------------

func sttyCmd(args ...string) *exec.Cmd {
	c := exec.Command("stty", args...)
	c.Stdin = os.Stdin
	return c
}

// enableRawMode switches the terminal to cbreak (no line buffering, no echo).
// Returns the saved state and whether raw mode was actually enabled
// (false when stdin isn't a TTY — the loop then falls back to file-watch only).
func enableRawMode() (string, bool) {
	out, err := sttyCmd("-g").Output()
	if err != nil {
		return "", false
	}
	saved := strings.TrimSpace(string(out))
	if err := sttyCmd("-icanon", "-echo", "min", "1", "time", "0").Run(); err != nil {
		return "", false
	}
	return saved, true
}

func restoreTerminal(saved string) {
	if saved == "" {
		return
	}
	_ = sttyCmd(saved).Run()
}

func readKeys(out chan<- byte) {
	buf := make([]byte, 1)
	for {
		n, err := os.Stdin.Read(buf)
		if err != nil {
			return
		}
		if n > 0 {
			out <- buf[0]
		}
	}
}

// --- filesystem watcher (polling) ----------------------------------------

func startFsWatcher(root string) <-chan struct{} {
	ch := make(chan struct{}, 1)
	go func() {
		last, _ := snapshot(root)
		for {
			time.Sleep(500 * time.Millisecond)
			cur, err := snapshot(root)
			if err != nil {
				continue
			}
			if changed(last, cur) {
				last = cur
				select {
				case ch <- struct{}{}:
				default:
				}
			} else {
				last = cur
			}
		}
	}()
	return ch
}

// progress counts exercises whose I-AM-NOT-DONE marker has been removed.
// We trust the marker as the source of truth so the bar updates instantly
// on save without having to recompile every "done" exercise.
func progress(inf *info) (done, total int) {
	for _, ex := range inf.Exercises {
		total++
		marker, err := hasMarker(ex.File)
		if err == nil && !marker {
			done++
		}
	}
	return
}

func clearScreen() {
	// ANSI: clear screen + move cursor home. Harmless if terminal ignores it.
	fmt.Print("\033[2J\033[H")
}

func nextPending(inf *info) *Exercise {
	for i := range inf.Exercises {
		ex := &inf.Exercises[i]
		marker, err := hasMarker(ex.File)
		if err != nil {
			// File missing — surface it so the user notices.
			fmt.Fprintf(os.Stderr, "warning: cannot read %s: %v\n", ex.File, err)
			continue
		}
		if marker {
			return ex
		}
		// Marker removed — make sure it actually passes; if not, treat as pending.
		_, ok, err := runExercise(ex)
		if err != nil {
			die(err)
		}
		if !ok {
			return ex
		}
	}
	return nil
}

func snapshot(root string) (map[string]time.Time, error) {
	m := make(map[string]time.Time)
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		m[path] = info.ModTime()
		return nil
	})
	if err != nil {
		return nil, err
	}
	return m, nil
}

func changed(a, b map[string]time.Time) bool {
	if len(a) != len(b) {
		return true
	}
	for k, v := range b {
		if old, ok := a[k]; !ok || !old.Equal(v) {
			return true
		}
	}
	return false
}
