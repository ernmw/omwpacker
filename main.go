package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/ernmw/omwpacker/esm"
	"github.com/ernmw/omwpacker/esm/tags"
	"github.com/ernmw/omwpacker/omwscripts"
)

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <input> [output]\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	inPath := os.Args[1]
	var outPath string
	if len(os.Args) >= 3 {
		outPath = os.Args[2]
	}

	ext := strings.ToLower(filepath.Ext(inPath))

	if !fileExists(inPath) {
		fmt.Printf("File %q not found\n", inPath)
		os.Exit(1)
	}

	switch ext {
	case ".omwscripts":
		// Convert text → addon
		if outPath == "" {
			outPath = strings.TrimSuffix(inPath, ext) + ".omwaddon"
		}
		fmt.Printf("Packing %s → %s\n", inPath, outPath)

		var outRecords []*esm.Record

		if fileExists(outPath) {
			// file exists, so load it
			var err error
			outRecords, err = esm.ParsePluginFile(outPath)
			if err != nil {
				fmt.Printf("Failed to parse %q: %v", outPath, err)
				os.Exit(1)
			}
			// delete existing luaf/luas subrecords
			for _, rec := range outRecords {
				if rec.Tag == tags.LUAL {
					rec.Subrecords = slices.DeleteFunc(rec.Subrecords, func(e *esm.Subrecord) bool {
						return e.Tag == tags.LUAF || e.Tag == tags.LUAS
					})
				}
			}
		} else {
			// make new empty records
			firstRec, err := esm.NewTES3Record("", "Made with https://github.com/ernmw/omwpacker/")
			if err != nil {
				fmt.Printf("Failed to make empty recs: %v", err)
				os.Exit(1)
			}
			outRecords = []*esm.Record{firstRec}
		}

		inContents, err := os.ReadFile(inPath)
		if err != nil {
			fmt.Printf("Failed to read input file: %v", err)
			os.Exit(1)
		}
		subRecs, err := omwscripts.Package(string(inContents))
		if err != nil {
			fmt.Printf("Failed to read file %q: %v", inPath, err)
			os.Exit(1)
		}

		found := false
		for _, rec := range outRecords {
			if rec.Tag == tags.LUAL {
				found = true
				rec.Subrecords = append(rec.Subrecords, subRecs...)
			}
		}
		if !found {
			// make new lual
			outRecords = append(outRecords, &esm.Record{
				Tag:        tags.LUAL,
				Subrecords: subRecs,
			})
		}
		writeOut, err := os.Create(outPath)
		if err != nil {
			fmt.Printf("Failed to read file %q: %v", inPath, err)
			os.Exit(1)
		}
		if err := esm.WriteRecords(writeOut, slices.Values(outRecords)); err != nil {
			fmt.Printf("Failed to write file %q: %v", outPath, err)
			os.Exit(1)
		}
		fmt.Println("✓ Created", outPath)

	case ".omwaddon", ".esp":
		fmt.Println("not implemented yet")
		os.Exit(1)

	default:
		fmt.Fprintf(os.Stderr, "Unsupported file extension: %s\n", ext)
		os.Exit(1)
	}
}
