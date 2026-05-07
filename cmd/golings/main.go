package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
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

func hasMarker(file string) (bool, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return false, err
	}
	return strings.Contains(string(data), notDoneMarker), nil
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

// cmdWatch finds the next pending exercise and re-runs it whenever any file
// under exercises/ changes. Polling-based to avoid third-party deps.
func cmdWatch(inf *info) {
	clearScreen()
	if inf.Welcome != "" {
		fmt.Println(inf.Welcome)
		fmt.Println()
	}
	lastSnap, _ := snapshot("exercises")
	for {
		ex := nextPending(inf)
		render(inf, ex)
		if ex == nil {
			fmt.Println("\n🎉 All exercises complete! Run `golings verify` to double-check.")
			return
		}
		runOne(ex, false)
		fmt.Println("\n(watching for changes — Ctrl-C to quit)")
		// Poll for any change under exercises/
		for {
			time.Sleep(500 * time.Millisecond)
			snap, err := snapshot("exercises")
			if err != nil {
				die(err)
			}
			if changed(lastSnap, snap) {
				lastSnap = snap
				break
			}
			lastSnap = snap
		}
		clearScreen()
	}
}

// render prints a progress bar plus a header for the current exercise.
func render(inf *info, current *Exercise) {
	done, total := progress(inf)
	const width = 40
	filled := 0
	if total > 0 {
		filled = done * width / total
	}
	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	pct := 0
	if total > 0 {
		pct = done * 100 / total
	}
	fmt.Printf("Progress: [%s] %d/%d  %d%%\n\n", bar, done, total, pct)
	if current != nil {
		fmt.Printf("▶ %s  (%s)\n  %s\n\n", current.Name, current.Mode, current.File)
	}
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
