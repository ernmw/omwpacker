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
	"github.com/ernmw/omwpacker/omwpack"
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
			}
			// delete existing records
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
			}
			outRecords = []*esm.Record{firstRec}
		}

		if err := omwpack.PackageOmwScripts(inPath, outPath); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
		fmt.Println("✓ Created", outPath)

	case ".omwaddon", ".esp":
		// Convert addon → text
		if outPath == "" {
			outPath = strings.TrimSuffix(inPath, ext) + ".omwscripts"
		}
		fmt.Printf("Extracting %s → %s\n", inPath, outPath)

		if err := omwpack.ExtractOmwScripts(inPath, outPath); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
		fmt.Println("✓ Created", outPath)

	default:
		fmt.Fprintf(os.Stderr, "Unsupported file extension: %s\n", ext)
		os.Exit(1)
	}
}
