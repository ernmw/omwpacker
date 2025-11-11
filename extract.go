package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/pflag"
	"go.coder.com/cli"
)

// extractCmd implements extract
type extractCmd struct {
	out string
}

func (cmd *extractCmd) Spec() cli.CommandSpec {
	return cli.CommandSpec{
		Name:    "extract",
		Usage:   "<input> [-o output]",
		Aliases: []string{"x"},
		Desc:    "Extract an .omwaddon/.esp into .omwscripts (not implemented).",
	}
}

func (cmd *extractCmd) RegisterFlags(fl *pflag.FlagSet) {
	fl.StringVarP(&cmd.out, "output", "o", "", "Output file path (defaults to <input>.omwscripts)")
}

func (cmd *extractCmd) Run(fl *pflag.FlagSet) {
	if fl.NArg() < 1 {
		fl.Usage()
		fmt.Fprintln(os.Stderr, "input file required")
		os.Exit(2)
	}
	inPath := fl.Arg(0)
	outPath := cmd.out
	if outPath == "" {
		ext := filepath.Ext(inPath)
		outPath = strings.TrimSuffix(inPath, ext) + ".omwscripts"
	}

	if err := cmd.extractCommand(inPath, outPath); err != nil {
		fmt.Printf("💀 Failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("🩷 Extracted to %q\n", outPath)
}

func (cmd *extractCmd) extractCommand(inPath, outPath string) error {
	// placeholder — keep behaviour as before: not implemented
	return fmt.Errorf("extract not implemented")
}
