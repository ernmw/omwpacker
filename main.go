// main.go
package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"unicode"

	"github.com/spf13/pflag"
	"go.coder.com/cli"
)

// ---------- Utilities (unchanged logic) ----------

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

func backup(path string) (string, error) {
	if fileExists(path) {
		source, err := os.Open(path)
		if err != nil {
			return "", fmt.Errorf("failed to open file %q: %w", path, err)
		}
		defer source.Close()

		tmp, err := os.CreateTemp("", filepath.Base(path))
		if err != nil {
			return "", fmt.Errorf("failed to create temporary file %q: %w", path, err)
		}
		defer tmp.Close()

		if _, err := io.Copy(tmp, source); err != nil {
			return "", fmt.Errorf("failed to backup file %q: %w", path, err)
		}
		return tmp.Name(), nil
	}
	return "", nil
}

// printHex prints binary data with ASCII row above hex row (terminal-friendly).
func printHex(width int, dump []byte) error {
	// Each byte = "xx " -> 3 columns
	bytesPerLine := width / 3
	if bytesPerLine > 32 {
		bytesPerLine = 32
	} else if bytesPerLine < 4 {
		bytesPerLine = 4
	}

	for i := 0; i < len(dump); i += bytesPerLine {
		end := i + bytesPerLine
		if end > len(dump) {
			end = len(dump)
		}
		line := dump[i:end]

		// ASCII row
		for _, b := range line {
			if unicode.IsPrint(rune(b)) {
				fmt.Printf(" %c ", b)
			} else {
				fmt.Printf(" . ")
			}
		}
		fmt.Println()

		// Hex row
		for _, b := range line {
			fmt.Printf("%02x ", b)
		}
		fmt.Println()
	}
	return nil
}

// ---------- Core command implementations (pure functions) ----------

// ---------- CLI wiring using go.coder.com/cli + spf13/pflag ----------

// root command; provides top-level metadata and subcommands.
type root struct{}

// Spec satisfies cli.Command
func (r *root) Spec() cli.CommandSpec {
	return cli.CommandSpec{
		Name:  "omwpacker",
		Usage: "[command] [flags]",
		Desc:  "Utilities for packing/unpacking Morrowind omw scripts and addons.",
	}
}

// Run for root just shows usage
func (r *root) Run(fl *pflag.FlagSet) {
	fl.Usage()
}

// Subcommands implements ParentCommand, attaching child commands.
func (r *root) Subcommands() []cli.Command {
	return []cli.Command{
		new(packCmd),
		new(extractCmd),
		new(readCmd),
	}
}

// ---------- Entry point: wire into cdr/cli runner ----------

func main() {
	cli.RunRoot(&root{})
}
