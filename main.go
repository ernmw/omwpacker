package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ernmw/omwpacker/omwpack"
)

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

	switch ext {
	case ".omwscripts":
		// Convert text → addon
		if outPath == "" {
			outPath = strings.TrimSuffix(inPath, ext) + ".omwaddon"
		}
		fmt.Printf("Packing %s → %s\n", inPath, outPath)

		// If you have a default ESP template, specify it here, e.g.:
		template := "" // or "S3maphore.esp"
		if err := omwpack.PackageOmwScripts(inPath, outPath, template); err != nil {
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
