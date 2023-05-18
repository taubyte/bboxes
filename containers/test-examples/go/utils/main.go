package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
)

func main() {
	fs := token.NewFileSet()
	pkgs, err := parser.ParseDir(fs, os.Args[1], nil, 0)
	if err != nil {
		log.Fatal("parsing dir:", err)
	}

	fmt.Printf("Pakages %#v", pkgs)

	pkg, ok := pkgs[os.Args[2]]
	if !ok {
		log.Fatal("package " + os.Args[2] + " not found")
	}

	ast.Walk(VisitorFunc(FindTypes), pkg)
}

type VisitorFunc func(n ast.Node) ast.Visitor

func (f VisitorFunc) Visit(n ast.Node) ast.Visitor {
	return f(n)
}

func FindTypes(n ast.Node) ast.Visitor {
	switch n := n.(type) {
	case *ast.Package:
		return VisitorFunc(FindTypes)
	case *ast.File:
		return VisitorFunc(FindTypes)
	case *ast.FuncDecl:
		if n.Body == nil {
			os.WriteFile(os.Args[1]+"/"+n.Name.Name+".s", []byte{}, 0640)
		}
	}
	return nil
}
