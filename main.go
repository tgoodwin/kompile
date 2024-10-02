package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
)

func findGoroutines(node ast.Node) {
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.GoStmt:
			fmt.Printf("Found a goroutine at line %d\n", x.Pos())
		}
		return true
	})
}

func main() {
	fmt.Println("vim-go")
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <filename>")
		return
	}
	filename := os.Args[1]
	fset := token.NewFileSet()

	// parse the source file into an AST
	node, err := parser.ParseFile(fset, filename, nil, 0)
	if err != nil {
		fmt.Println(err)
		return
	}

	findGoroutines(node)
}
