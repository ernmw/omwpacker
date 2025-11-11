package main

import (
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/ernmw/omwpacker/esm"
	"github.com/spf13/pflag"
	"go.coder.com/cli"
	"golang.org/x/term"
)

// readCmd implements read; supports extra positional args for later filtering etc.
type readCmd struct {
	record    string // -r record
	subrecord string // -s subrecord
}

func (cmd *readCmd) Spec() cli.CommandSpec {
	return cli.CommandSpec{
		Name:    "read",
		Usage:   "<input> [-r record] [-s subrecord]",
		Aliases: []string{"r"},
		Desc:    "Read and display contents of an .omwaddon/.esp/.esp.",
	}
}

func (cmd *readCmd) RegisterFlags(fl *pflag.FlagSet) {
	fl.StringVarP(&cmd.record, "record", "r", "", "Filter to records of the given type.")
	fl.StringVarP(&cmd.subrecord, "subrecord", "s", "", "Filter to subrecords of the given type.")
}

func (cmd *readCmd) Run(fl *pflag.FlagSet) {
	if fl.NArg() < 1 {
		fl.Usage()
		fmt.Fprintln(os.Stderr, "input file required")
		os.Exit(2)
	}
	inPath := fl.Arg(0)

	if !fileExists(inPath) {
		fmt.Printf("💀 Failed: File %q not found\n", inPath)
		os.Exit(1)
	}

	// set up record filter
	var recFilter func(rec *esm.Record) bool
	if len(cmd.record) > 0 {
		expected := esm.RecordTag(strings.ToUpper(cmd.record))
		recFilter = func(rec *esm.Record) bool {
			return rec.Tag == expected
		}
	} else {
		recFilter = func(rec *esm.Record) bool { return true }
	}

	// set up subrecord filter
	var subrecFilter func(sub *esm.Subrecord) bool
	if len(cmd.subrecord) > 0 {
		expected := esm.SubrecordTag(strings.ToUpper(cmd.subrecord))
		subrecFilter = func(subrec *esm.Subrecord) bool {
			return subrec.Tag == expected
		}
	} else {
		subrecFilter = func(sub *esm.Subrecord) bool { return true }
	}

	if err := cmd.readCommand(inPath, recFilter, subrecFilter); err != nil {
		fmt.Printf("💀 Failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("🩷 Done reading %q\n", inPath)
}

func (cmd *readCmd) readCommand(
	inPath string,
	recordFilter func(rec *esm.Record) bool,
	subrecordFilter func(sub *esm.Subrecord) bool,
) error {

	inRecords, err := esm.ParsePluginFile(inPath)
	if err != nil {
		return fmt.Errorf("failed to parse %q: %w", inPath, err)
	}

	width := 120
	if fd := int(os.Stdout.Fd()); term.IsTerminal(fd) {
		width, _, err = term.GetSize(fd)
		if err != nil {
			return fmt.Errorf("get terminal size: %w", err)
		}
	}

	for _, rec := range inRecords {
		if !recordFilter(rec) {
			continue
		}
		fmt.Printf("%s:\n", rec.Tag)
		for _, subRec := range rec.Subrecords {
			if !subrecordFilter(subRec) {
				continue
			}
			fmt.Printf("  %s:\n", subRec.Tag)
			_ = printHex(width, subRec.Data)
		}
	}
	return nil
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
