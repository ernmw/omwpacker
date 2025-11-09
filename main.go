package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"unicode"

	"golang.org/x/term"

	"github.com/ernmw/omwpacker/esm"
	"github.com/ernmw/omwpacker/esm/record/lua"
	"github.com/ernmw/omwpacker/esm/record/tes3"
	"github.com/ernmw/omwpacker/omwscripts"
)

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

func pack(inPath, outPath string) error {
	var outRecords []*esm.Record

	if fileExists(outPath) {
		// file exists, so load it
		var err error
		outRecords, err = esm.ParsePluginFile(outPath)
		if err != nil {
			return fmt.Errorf("Failed to parse %q: %v", outPath, err)
		}
		// delete existing luaf/luas subrecords
		for _, rec := range outRecords {
			if rec.Tag == lua.LUAL {
				rec.Subrecords = slices.DeleteFunc(rec.Subrecords, func(e *esm.Subrecord) bool {
					return e.Tag == lua.LUAF || e.Tag == lua.LUAS
				})
			}
		}
	} else {
		// make new empty records
		firstRec, err := tes3.NewTES3Record("", "Made with https://github.com/ernmw/omwpacker/")
		if err != nil {
			return fmt.Errorf("Failed to make empty recs: %v", err)
		}
		outRecords = []*esm.Record{firstRec}
	}

	inContents, err := os.ReadFile(inPath)
	if err != nil {
		return fmt.Errorf("Failed to read input file: %v", err)
	}
	subRecs, err := omwscripts.Package(string(inContents))
	if err != nil {
		return fmt.Errorf("Failed to read file %q: %v", inPath, err)
	}

	found := false
	for _, rec := range outRecords {
		if rec.Tag == lua.LUAL {
			found = true
			rec.Subrecords = append(rec.Subrecords, subRecs...)
		}
	}
	if !found {
		// make new lual
		outRecords = append(outRecords, &esm.Record{
			Tag:        lua.LUAL,
			Subrecords: subRecs,
		})
	}
	writeOut, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("Failed to read file %q: %v", inPath, err)
	}
	if err := esm.WriteRecords(writeOut, slices.Values(outRecords)); err != nil {
		return fmt.Errorf("Failed to write file %q: %v", outPath, err)
	}
	return nil
}

// printHex prints binary data in a readable, aligned format.
// For terminals: each ASCII char appears *above* its corresponding byte’s hex.
func printHex(width int, dump []byte) error {

	// Each byte = 3 columns ("xx ").
	// Determine how many bytes per line.
	bytesPerLine := width / 3
	if bytesPerLine > 32 {
		bytesPerLine = 32
	} else if bytesPerLine < 4 {
		bytesPerLine = 4
	}

	for i := 0; i < len(dump); i += bytesPerLine {
		end := min(i+bytesPerLine, len(dump))
		line := dump[i:end]

		// Build top row: printable chars, padded to same column positions as hex
		for _, b := range line {
			if unicode.IsPrint(rune(b)) {
				fmt.Printf(" %c ", b)
			} else {
				fmt.Printf(" . ")
			}
		}
		fmt.Println()

		// Build bottom row: hex values aligned under the chars
		for _, b := range line {
			fmt.Printf("%02x ", b)
		}
		fmt.Println()
	}
	return nil
}

func read(inPath string) error {
	inRecords, err := esm.ParsePluginFile(inPath)
	if err != nil {
		return fmt.Errorf("Failed to parse %q: %v", inPath, err)
	}
	width := 120

	if fd := int(os.Stdout.Fd()); term.IsTerminal(fd) {
		width, _, err = term.GetSize(fd)
		if err != nil {
			return fmt.Errorf("get terminal size: %w", err)
		}
	}
	// delete existing luaf/luas subrecords
	for _, rec := range inRecords {
		fmt.Printf("%s: \n", rec.Tag)
		for _, subRec := range rec.Subrecords {
			fmt.Printf("  %s: \n", subRec.Tag)
			printHex(width, subRec.Data)
		}
	}
	return nil
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

func main() {
	var usage = fmt.Sprintf("Usage: %s [pack|extract|read] <input> [output]\n", filepath.Base(os.Args[0]))
	if len(os.Args) < 3 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(1)
	}

	verb := strings.ToLower(strings.TrimSpace(os.Args[1]))
	inPath := os.Args[2]
	var outPath string
	if len(os.Args) >= 4 {
		outPath = os.Args[3]
	}

	ext := strings.ToLower(filepath.Ext(inPath))

	if !fileExists(inPath) {
		fmt.Printf("💀 Failed: File %q not found\n", inPath)
		os.Exit(1)
	}

	switch verb {
	case "pack":
		// Convert text → addon
		if outPath == "" {
			outPath = strings.TrimSuffix(inPath, ext) + ".omwaddon"
		}

		modify := fileExists(outPath)
		if backupFile, err := backup(outPath); err != nil {
			fmt.Printf("💀 Failed: Couldn't back up %q: %v\n", outPath, err)
			os.Exit(1)
		} else if len(backupFile) > 0 {
			fmt.Printf("Backed up %q → %q\n", outPath, backupFile)
		}

		fmt.Printf("Packing %q → %q\n", inPath, outPath)
		err := pack(inPath, outPath)
		if err != nil {
			fmt.Printf("💀 Failed: %v\n", err)
			os.Exit(1)
		}

		if modify {
			fmt.Printf("🤍 Modified %q\n", outPath)
		} else {
			fmt.Printf("🩵 Created %q\n", outPath)
		}
	case "extract":
		fmt.Printf("💀 Failed: %v\n", "not implemented yet\n")
		os.Exit(1)
	case "read":
		err := read(inPath)
		if err != nil {
			fmt.Printf("💀 Failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("🩷 Done reading %q\n", inPath)
	default:
		fmt.Fprint(os.Stderr, usage)
		os.Exit(1)
	}
}
