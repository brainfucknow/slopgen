// Package check validates generated source in-process with the Go frontend.
package check

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
)

func Source(filename string, src []byte) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, src, 0)
	if err != nil {
		return err
	}
	conf := types.Config{Error: func(error) {}}
	if _, err = conf.Check("generated", fset, []*ast.File{f}, nil); err != nil {
		return fmt.Errorf("type check: %w", err)
	}
	return nil
}
