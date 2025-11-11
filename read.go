package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"unicode"

	"github.com/ernmw/omwpacker/cfg"
	"github.com/ernmw/omwpacker/esm"
	"github.com/spf13/pflag"
	"go.coder.com/cli"
	"golang.org/x/term"
)

// readCmd implements read; supports extra positional args for later filtering etc.
type readCmd struct {
	record    string // -r record
	subrecord string // -s subrecord
	filter    string // -f subrecordtag=string
}

func (cmd *readCmd) Spec() cli.CommandSpec {
	return cli.CommandSpec{
		Name:    "read",
		Usage:   "<input> [-r record] [-s subrecord] [-f subrecordtag=string]",
		Aliases: []string{"r"},
		Desc:    "Read and display contents of an .omwaddon/.esp/.esp/openmw.cfg/morrowind.ini.",
	}
}

func (cmd *readCmd) RegisterFlags(fl *pflag.FlagSet) {
	fl.StringVarP(&cmd.record, "record", "r", "", "Display records of the given type. Specify multiples by delimiting with a comma.")
	fl.StringVarP(&cmd.subrecord, "subrecord", "s", "", "Display subrecords of the given type. Specify multiples by delimiting with a comma.")
	fl.StringVarP(&cmd.filter, "filter", "f", "", "Filter records to those that contain the given subrecord, and that subrecord contains the provided string. Example: 'NAME=Balmora'. Prefix the string with '0x' to interpret it as hex-encoded.")
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
		expectedTokens := []esm.RecordTag{}
		for tok := range strings.SplitSeq(cmd.record, ",") {
			expectedTokens = append(expectedTokens, esm.RecordTag(strings.ToUpper(tok)))
		}
		recFilter = func(rec *esm.Record) bool {
			return slices.Contains(expectedTokens, rec.Tag)
		}
	} else {
		recFilter = func(_ *esm.Record) bool { return true }
	}

	// set up subrecord filter
	var subrecFilter func(sub *esm.Subrecord) bool
	if len(cmd.subrecord) > 0 {
		expectedTokens := []esm.SubrecordTag{}
		for tok := range strings.SplitSeq(cmd.subrecord, ",") {
			expectedTokens = append(expectedTokens, esm.SubrecordTag(strings.ToUpper(tok)))
		}
		subrecFilter = func(subrec *esm.Subrecord) bool {
			return slices.Contains(expectedTokens, subrec.Tag)
		}
	} else {
		subrecFilter = func(_ *esm.Subrecord) bool { return true }
	}

	// set up filter
	var filter func(rec *esm.Record) bool
	tokens := strings.SplitN(cmd.filter, "=", 2)
	if len(cmd.filter) > 0 && len(tokens) == 2 {
		name := esm.SubrecordTag(strings.ToUpper(tokens[0]))
		sub := []byte(tokens[1])
		if strings.HasPrefix(tokens[1], "0x") {
			var err error
			sub, err = hex.DecodeString(strings.TrimPrefix(tokens[1], "0x"))
			if err != nil {
				fmt.Printf("💀 Failed: String %q is not hex.\n", tokens[1])
				os.Exit(1)
			}
		}
		filter = func(rec *esm.Record) bool {
			return slices.ContainsFunc(rec.Subrecords, func(s *esm.Subrecord) bool {
				if s.Tag != name {
					return false
				}
				if len(sub) > 0 {
					return bytes.Contains(s.Data, sub)
				}
				return true
			})
		}
	} else {
		filter = func(_ *esm.Record) bool { return true }
	}

	combinedRecordFilter := func(rec *esm.Record) bool {
		return recFilter(rec) && filter(rec)
	}

	var plugins []string

	if strings.EqualFold(filepath.Ext(inPath), ".cfg") {
		var err error
		plugins, _, err = cfg.OpenMWPlugins(inPath)
		if err != nil {
			fmt.Printf("💀 Failed: %q couldn't be parsed: %v\n", inPath, err)
			os.Exit(1)
		}
	} else {
		plugins = []string{inPath}
	}
	for _, plugin := range plugins {
		if err := cmd.readCommand(
			plugin,
			combinedRecordFilter,
			subrecFilter); err != nil {
			fmt.Printf("💀 Failed parsing %s: %v\n", plugin, err)
			os.Exit(1)
		}
	}

	fmt.Printf("🩷 Done reading %q\n", inPath)
}

func (cmd *readCmd) readCommand(
	in string,
	recordFilter func(rec *esm.Record) bool,
	subrecordFilter func(sub *esm.Subrecord) bool,
) error {

	inRecords, err := esm.ParsePluginFile(in)
	if err != nil {
		return fmt.Errorf("failed to parse %q: %w", in, err)
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
		headerPrinted := false
		for _, subRec := range rec.Subrecords {
			if !subrecordFilter(subRec) {
				continue
			}
			if !headerPrinted {
				fmt.Printf("\n%s: (%s)\n", rec.Tag, filepath.Base(in))
				headerPrinted = true
			}
			fmt.Printf("  %s:\n", subRec.Tag)
			if err = printHex(width, subRec.Data); err != nil {
				return fmt.Errorf("printing %s/%s from %q", rec.Tag, subRec.Tag, in)
			}
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
