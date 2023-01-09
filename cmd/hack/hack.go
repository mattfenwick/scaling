package main

import (
	"debug/buildinfo"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"

	"github.com/mattfenwick/collections/pkg/file"
	"github.com/mattfenwick/collections/pkg/json"
	"github.com/mattfenwick/scaling/pkg/utils"
	"github.com/pkg/errors"
)

func PrintBuildInfo() {
	name := os.Args[1]
	info, err := buildinfo.ReadFile(name)
	utils.Die(err)
	fmt.Printf("name %s\ninfo %s\n", name, json.MustMarshalToString(info))
}

func main() {
	filename := "cmd/hack/hack.go"
	src, err := file.ReadString(filename)
	utils.Die(err)

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, src, 0)
	utils.Die(errors.Wrapf(err, "unable to ParseFile %s", filename))

	//fmt.Printf("%s\n\n", json.MustMarshalToString(f)) // hits cycle

	// Inspect the AST and print all identifiers and literals.
	ast.Inspect(f, func(n ast.Node) bool {
		//fmt.Printf("node: %T; %s\n\n", n, json.MustMarshalToString(n))
		fmt.Printf("node: %T\n", n)
		var s string
		switch x := n.(type) {
		case *ast.BasicLit:
			s = x.Value
		case *ast.Ident:
			s = x.Name
		}
		if s != "" {
			fmt.Printf("%s:\t%s\n", fset.Position(n.Pos()), s)
		}
		return true
	})

	utils.Die(ast.Print(fset, f))
}
