package annotation

import (
	"github.com/stretchr/testify/assert"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestFileVisitor_Get(t *testing.T) {
	testCases := []struct {
		src  string
		want File
	}{
		{
			src: `
// annotation go through the source code and extra the annotation
// @author Deng Ming
/* @multiple first line
second line
*/
// @date 2022/04/02
package annotation

type (
	// FuncType is a type
	// @author Deng Ming
	/* @multiple first line
	   second line
	*/
	// @date 2022/04/02
	FuncType func()
)

type (
	// StructType is a test struct
	//
	// @author Deng Ming
	/* @multiple first line
	   second line
	*/
	// @date 2022/04/02
	StructType struct {
		// Public is a field
		// @type string
		Public string
	}

	// SecondType is a test struct
	//
	// @author Deng Ming
	/* @multiple first line
	   second line
	*/
	// @date 2022/04/03
	SecondType struct {
	}
)

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
`,
			want: File{
				Annotations: Annotations[*ast.File]{
					Ans: []Annotation{
						{
							Key:   "author",
							Value: "Deng Ming",
						},
						{
							Key:   "multiple",
							Value: "first line\nsecond line\n",
						},
						{
							Key:   "date",
							Value: "2022/04/02",
						},
					},
				},
				Types: []Type{
					{
						Annotations: Annotations[*ast.TypeSpec]{
							Ans: []Annotation{
								{
									Key:   "author",
									Value: "Deng Ming",
								},
								{
									Key:   "multiple",
									Value: "first line\n\t   second line\n\t",
								},
								{
									Key:   "date",
									Value: "2022/04/02",
								},
							},
						},
					},
					{
						Annotations: Annotations[*ast.TypeSpec]{
							Ans: []Annotation{
								{
									Key:   "author",
									Value: "Deng Ming",
								},
								{
									Key:   "multiple",
									Value: "first line\n\t   second line\n\t",
								},
								{
									Key:   "date",
									Value: "2022/04/02",
								},
							},
						},
						Fields: []Field{
							{
								Annotations: Annotations[*ast.Field]{
									Ans: []Annotation{
										{
											Key:   "type",
											Value: "string",
										},
									},
								},
							},
						},
					},
					{
						Annotations: Annotations[*ast.TypeSpec]{
							Ans: []Annotation{
								{
									Key:   "author",
									Value: "Deng Ming",
								},
								{
									Key:   "multiple",
									Value: "first line\n\t   second line\n\t",
								},
								{
									Key:   "date",
									Value: "2022/04/03",
								},
							},
						},
					},
					{
						Annotations: Annotations[*ast.TypeSpec]{
							Ans: []Annotation{
								{
									Key:   "author",
									Value: "Deng Ming",
								},
								{
									Key:   "multiple",
									Value: "first line\n\t   second line\n\t",
								},
								{
									Key:   "date",
									Value: "2022/04/04",
								},
							},
						},
						Fields: []Field{
							{
								Annotations: Annotations[*ast.Field]{
									Ans: []Annotation{
										{
											Key:   "parameter",
											Value: "arg1 int",
										},
										{
											Key:   "parameter",
											Value: "arg2 int32",
										},
										{
											Key:   "return",
											Value: "string",
										},
									},
								},
							},
							{
								Annotations: Annotations[*ast.Field]{
									Ans: []Annotation{
										{
											Key:   "return",
											Value: "string",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tc := range testCases {
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, "src.go", tc.src, parser.ParseComments)
		if err != nil {
			t.Fatal(err)
		}
		tv := &SingleFileEntryVisitor{}
		ast.Walk(tv, f)
		file := tv.Get()
		assertAnnotations(t, tc.want.Annotations, file.Annotations)
		if len(tc.want.Types) != len(file.Types) {
			t.Fatal()
		}
		for i, typ := range file.Types {
			wantType := tc.want.Types[i]
			assertAnnotations(t, wantType.Annotations, typ.Annotations)
			if len(wantType.Fields) != len(typ.Fields) {
				t.Fatal()
			}
			for j, fd := range typ.Fields {
				wantFd := wantType.Fields[j]
				assertAnnotations(t, wantFd.Annotations, fd.Annotations)
			}
		}
	}
}

func assertAnnotations[N ast.Node](t *testing.T, wantAns Annotations[N], dst Annotations[N]) {
	want := wantAns.Ans
	if len(want) != len(dst.Ans) {
		t.Fatal()
	}
	for i, an := range want {
		val := dst.Ans[i]
		assert.Equal(t, an.Value, val.Value)
	}
}
