package ast

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"testing"
)

func TestHello(t *testing.T) {

	testCases := []struct{
		name string
		src string
		wantRes string
	}{
		{
			name:"one input and no import",
			src: `
package ast

func Hello(name string) string {
	return "hello, " + name
}
`,
wantRes: `package ast

import "log"

func Hello(name string) string {
	log.
		Println(name,
		)
	return "hello, " + name
}
`,
		},
		{
			name: "no input",
			src: `
package ast

func Hello() string {
	return "hello, world"
}
`,
wantRes: `package ast

func Hello() string {
	return "hello, world"
}
`,
		},
		{
			name: "multiple input",
			src: `
package ast

import "fmt"

func Hello(firstName, lastName string, age int) string {
	return fmt.Sprintf("hello, %s %s, %d", firstName, lastName, age)
}
`,
wantRes: `package ast

import (
	"fmt"
	"log"
)

func Hello(firstName, lastName string, age int) string {
	log.
		Println(firstName,

			lastName,

			age)
	return fmt.Sprintf("hello, %s %s, %d", firstName, lastName, age)
}
`,
		},
		{
			name: "multiple func",
			src: `
package ast

import (
	"fmt"
)

func Hello(name string) string {
	return "hello, " + name
}

func Hello1(firstName, lastName string, age int) string {
	return fmt.Sprintf("hello, %s %s, %d", firstName, lastName, age)
}
`,
wantRes: `package ast

import (
	"fmt"
	"log"
)

func Hello(name string) string {
	log.
		Println(name)

	return "hello, " + name
}

func Hello1(firstName, lastName string, age int) string {
	log.
		Println(firstName,

			lastName,

			age,
		)

	return fmt.Sprintf("hello, %s %s, %d", firstName, lastName, age)
}
`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, "src.go", tc.src, parser.ParseComments)
			if err != nil {
				t.Fatal(err)
			}
			ast.Walk(&loggingVisitor{}, f)
			sb := &bytes.Buffer{}

			err = format.Node(sb, fset, f)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tc.wantRes, sb.String())
		})
	}
}
