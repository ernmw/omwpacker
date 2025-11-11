package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/pflag"
	"go.coder.com/cli"
)

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

func main() {
	cli.RunRoot(&root{})
}
