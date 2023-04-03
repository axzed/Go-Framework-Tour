package demo

import (
	"github.com/stretchr/testify/assert"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestFileVisitor(t *testing.T) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "src.go",
		`
// annotation go through the source code and extra the annotation
// @author Deng Ming
// @date 2022/04/02
// @
package annotation

type (
	// Interface is a test interface
	// @author Deng Ming
	/* @multiple first line
	   second line
	*/
	// @date 2022/04/04
	Interface interface {
		// MyFunc is a test func
		// @parameter arg1 int
		// @parameter arg2 int32
		// @return string
		MyFunc(arg1 int, arg2 int32) string

		// second is a test func
		// @return string
		second() string
	}
)
`, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}
	fv := &FileVisitor{
		ans:   map[string]string{},
		types: []*TypeSpecVisitor{},
	}
	ast.Walk(fv, f)
	res := map[string]string{
		"date":   "2022/04/02",
		"author": "Deng Ming",
	}
	assert.Equal(t, res, fv.ans)
}
