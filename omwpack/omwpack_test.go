package omwpack

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPackageOmwScripts(t *testing.T) {
	tmpDir := t.TempDir()

	// Create sample .omwscripts text
	scriptContent := `
PLAYER: scripts/s3/music/core.lua
GLOBAL: scripts/s3/music/staticCollection.lua
`
	inFile := filepath.Join(tmpDir, "test.omwscripts")
	if err := os.WriteFile(inFile, []byte(scriptContent), 0644); err != nil {
		t.Fatal(err)
	}

	template := "testdata/expected.esp"
	if _, err := os.Stat(template); err != nil {
		t.Skip("template S3maphore.esp not found — skipping integration test")
	}

	outFile := filepath.Join(tmpDir, "out.omwaddon")

	err := PackageOmwScripts(inFile, outFile, template)
	if err != nil {
		t.Fatalf("PackageOmwScripts failed: %v", err)
	}

	// Validate output
	data, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatal(err)
	}

	if len(data) < 256 {
		t.Errorf("output too small: %d bytes", len(data))
	}

	mustContain := []string{"TES3", "LUAL", "LUAS", "LUAF", "core.lua", "staticCollection.lua"}
	for _, s := range mustContain {
		if !bytes.Contains(data, []byte(s)) {
			t.Errorf("output missing substring: %q", s)
		}
	}

	// Optional: sanity check order (LUAS before LUAF)
	iLuas := bytes.Index(data, []byte("LUAS"))
	iLuaf := bytes.Index(data, []byte("LUAF"))
	if iLuas < 0 || iLuaf < 0 || iLuas > iLuaf {
		t.Errorf("LUAS/LUAF ordering unexpected (LUAS=%d, LUAF=%d)", iLuas, iLuaf)
	}

	t.Logf("Output .omwaddon written to %s (%d bytes)", outFile, len(data))
}

func TestExtractOmwScripts(t *testing.T) {
	tmp := t.TempDir()
	addon := filepath.Join(tmp, "input.omwaddon")
	scripts := filepath.Join(tmp, "out.omwscripts")

	// First pack some data
	in := filepath.Join(tmp, "input.omwscripts")
	os.WriteFile(in, []byte("PLAYER: scripts/foo.lua\nGLOBAL: scripts/bar.lua\n"), 0644)
	if err := PackageOmwScripts(in, addon, "testdata/expected.esp"); err != nil {
		t.Fatal(err)
	}

	if err := ExtractOmwScripts(addon, scripts); err != nil {
		t.Fatal(err)
	}

	out, _ := os.ReadFile(scripts)
	txt := string(out)
	if !strings.Contains(txt, "PLAYER:") || !strings.Contains(txt, "GLOBAL:") {
		t.Errorf("unexpected output:\n%s", txt)
	}
}
