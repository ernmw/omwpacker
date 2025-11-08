package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type ConstInfo struct {
	Name string
	Type string
	Tag  string
}

// readConstants parses a Go source file and extracts constants with a comment containing "Type: X"
func readConstants(filename string) ([]ConstInfo, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var consts []ConstInfo

	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.CONST {
			continue
		}

		for _, spec := range genDecl.Specs {
			valueSpec := spec.(*ast.ValueSpec)
			if len(valueSpec.Names) == 0 || len(valueSpec.Values) == 0 {
				continue
			}

			name := valueSpec.Names[0].Name
			value := valueSpec.Values[0]

			basicLit, ok := value.(*ast.BasicLit)
			if !ok {
				continue
			}
			tag := strings.Trim(basicLit.Value, `"`)

			// Look for comment with "Type: X"
			var typ string
			comments := append(valueSpec.Comment.List, genDecl.Doc.List...)
			for _, c := range comments {
				if strings.Contains(c.Text, "Type:") {
					parts := strings.Split(c.Text, "Type:")
					if len(parts) == 2 {
						typ = strings.TrimSpace(parts[1])
						break
					}
				}
			}
			if typ == "" {
				continue
			}

			consts = append(consts, ConstInfo{
				Name: name,
				Type: typ,
				Tag:  tag,
			})
		}
	}

	return consts, nil
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: go run main.go <source-file.go>")
	}

	sourceFile := os.Args[1]
	consts, err := readConstants(sourceFile)
	if err != nil {
		log.Fatalf("failed to read constants: %v", err)
	}

	outDir := "genout"
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		log.Fatalf("failed to create output dir: %v", err)
	}

	for _, c := range consts {
		tmplFile := filepath.Join("templates", c.Type+".tmpl")
		tmplContent, err := os.ReadFile(tmplFile)
		if err != nil {
			log.Printf("warning: template for type %q not found: %v", c.Type, err)
			continue
		}

		tmpl, err := template.New(c.Type).Parse(string(tmplContent))
		if err != nil {
			log.Fatalf("failed to parse template %s: %v", tmplFile, err)
		}

		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, c); err != nil {
			log.Fatalf("failed to execute template for %s: %v", c.Name, err)
		}

		outFile := filepath.Join(outDir, strings.ToLower(c.Name)+"_gen.go")
		if err := os.WriteFile(outFile, buf.Bytes(), 0o644); err != nil {
			log.Fatalf("failed to write output file %s: %v", outFile, err)
		}

		fmt.Printf("generated %s\n", outFile)
	}
}
