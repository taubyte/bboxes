// rewrite_package rewrites non-main package declarations to "package main"
// in all .go files under the given directory (for compat with old-style layouts).
package main

import (
	"fmt"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
)

func main() {
	dir := "."
	if len(os.Args) >= 2 {
		dir = os.Args[1]
	}
	if err := rewriteDir(dir); err != nil {
		fmt.Fprintf(os.Stderr, "rewrite-package: %v\n", err)
		os.Exit(1)
	}
}

func rewriteDir(dir string) error {
	dir, err := filepath.Abs(dir)
	if err != nil {
		return err
	}
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if info.Name() == "vendor" {
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(path) != ".go" {
			return nil
		}
		return rewriteFile(path)
	})
}

func rewriteFile(path string) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("parse %s: %w", path, err)
	}
	if f.Name.Name == "main" {
		return nil
	}
	f.Name.Name = "main"
	out, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("open %s: %w", path, err)
	}
	defer out.Close()
	if err := printer.Fprint(out, fset, f); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}
