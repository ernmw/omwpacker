package main

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/ernmw/omwpacker/esm"
	"github.com/ernmw/omwpacker/esm/record/lua"
	"github.com/ernmw/omwpacker/esm/record/tes3"
	"github.com/ernmw/omwpacker/omwscripts"
	"github.com/spf13/pflag"
	"go.coder.com/cli"
)

// packCmd implements the pack subcommand.
type packCmd struct {
	out string // -o output
}

func (cmd *packCmd) Spec() cli.CommandSpec {
	return cli.CommandSpec{
		Name:    "pack",
		Usage:   "<input> [-o output]",
		Aliases: []string{"p"},
		Desc:    "Package a .omwscripts file into an .omwaddon (or inject into existing addon).",
	}
}

func (cmd *packCmd) RegisterFlags(fl *pflag.FlagSet) {
	fl.StringVarP(&cmd.out, "output", "o", "", "Output file path (defaults to <input>.omwaddon)")
}

func (cmd *packCmd) Run(fl *pflag.FlagSet) {
	// positional args are fl.Args()
	if fl.NArg() < 1 {
		fl.Usage()
		fmt.Fprintln(os.Stderr, "input file required")
		os.Exit(2)
	}

	inPath := fl.Arg(0)
	outPath := cmd.out
	ext := strings.ToLower(filepath.Ext(inPath))
	if outPath == "" {
		outPath = strings.TrimSuffix(inPath, ext) + ".omwaddon"
	}

	if !fileExists(inPath) {
		fmt.Printf("💀 Failed: File %q not found\n", inPath)
		os.Exit(1)
	}

	// backup output if exists
	if backupFile, err := backup(outPath); err != nil {
		fmt.Printf("💀 Failed: Couldn't back up %q: %v\n", outPath, err)
		os.Exit(1)
	} else if backupFile != "" {
		fmt.Printf("Backed up %q → %q\n", outPath, backupFile)
	}

	fmt.Printf("Packing %q → %q\n", inPath, outPath)
	if err := cmd.packCommand(inPath, outPath); err != nil {
		fmt.Printf("💀 Failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("🩵 Done: %q\n", outPath)
}

func (cmd *packCmd) packCommand(inPath, outPath string) error {
	var outRecords []*esm.Record

	if fileExists(outPath) {
		var err error
		outRecords, err = esm.ParsePluginFile(outPath)
		if err != nil {
			return fmt.Errorf("failed to parse %q: %v", outPath, err)
		}
		// remove existing LUAF/LUAS entries under LUAL
		for _, rec := range outRecords {
			if rec.Tag == lua.LUAL {
				rec.Subrecords = slices.DeleteFunc(rec.Subrecords, func(e *esm.Subrecord) bool {
					return e.Tag == lua.LUAF || e.Tag == lua.LUAS
				})
			}
		}
	} else {
		firstRec, err := tes3.NewTES3Record("", "Made with https://github.com/ernmw/omwpacker/")
		if err != nil {
			return fmt.Errorf("failed to make TES3 record: %v", err)
		}
		outRecords = []*esm.Record{firstRec}
	}

	inContents, err := os.ReadFile(inPath)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}
	subRecs, err := omwscripts.Package(string(inContents))
	if err != nil {
		return fmt.Errorf("failed to package file %q: %w", inPath, err)
	}

	found := false
	for _, rec := range outRecords {
		if rec.Tag == lua.LUAL {
			found = true
			rec.Subrecords = append(rec.Subrecords, subRecs...)
		}
	}
	if !found {
		outRecords = append(outRecords, &esm.Record{
			Tag:        lua.LUAL,
			Subrecords: subRecs,
		})
	}

	writeOut, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("failed to create output file %q: %w", outPath, err)
	}
	defer writeOut.Close()

	if err := esm.WriteRecords(writeOut, slices.Values(outRecords)); err != nil {
		return fmt.Errorf("failed to write file %q: %w", outPath, err)
	}
	return nil
}
