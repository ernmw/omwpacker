package omwpack

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

// ExtractOmwScripts reads an .omwaddon (or .esp) file and reconstructs the text-based .omwscripts list.
// It finds the LUAL record, parses LUAS/LUAF pairs, and writes "ATTACH: path" lines to outScriptsPath.
func ExtractOmwScripts(inAddonPath, outScriptsPath string) error {
	data, err := os.ReadFile(inAddonPath)
	if err != nil {
		return fmt.Errorf("read input: %w", err)
	}

	pos := bytes.Index(data, []byte("LUAL"))
	if pos < 0 {
		return fmt.Errorf("LUAL record not found in %s", inAddonPath)
	}
	if pos+8 > len(data) {
		return fmt.Errorf("invalid LUAL header")
	}
	size := binary.LittleEndian.Uint32(data[pos+4 : pos+8])
	end := int(pos) + 8 + int(size)
	if end > len(data) {
		return fmt.Errorf("truncated LUAL record (size=%d, file=%d)", size, len(data))
	}
	record := data[pos+8 : end]

	// Skip first 8 bytes (flags/padding)
	if len(record) < 8 {
		return fmt.Errorf("LUAL record too short")
	}
	record = record[8:]

	var (
		offset   int
		pairs    [][2]string
		lastLUAS string
	)

	for offset < len(record) {
		if offset+8 > len(record) {
			break
		}
		tag := string(record[offset : offset+4])
		sz := int(binary.LittleEndian.Uint32(record[offset+4 : offset+8]))
		offset += 8
		if offset+sz > len(record) {
			break
		}
		val := string(record[offset : offset+sz])
		offset += sz

		switch tag {
		case "LUAS":
			lastLUAS = val
		case "LUAF":
			if lastLUAS == "" {
				return fmt.Errorf("LUAF found before LUAS")
			}
			pairs = append(pairs, [2]string{lastLUAS, val})
			lastLUAS = ""
		default:
			// ignore unknown subrecord
		}
	}

	if len(pairs) == 0 {
		return fmt.Errorf("no LUAS/LUAF pairs found")
	}

	var out bytes.Buffer
	for _, p := range pairs {
		fmt.Fprintf(&out, "%s: %s\n", p[1], p[0])
	}

	return os.WriteFile(outScriptsPath, out.Bytes(), 0644)
}
