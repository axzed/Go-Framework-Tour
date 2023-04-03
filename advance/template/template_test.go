package template

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"text/template"
)

func TestIfElseBlock(t *testing.T)  {
	// 用一点小技巧来实现 for i 循环
	tpl := template.New("hello-world")
	tpl, err := tpl.Parse(`
{{- if and (gt .Age 0) (le .Age 6) }}
儿童 0<age<6
{{ else if and (gt .Age 6) (le .Age 18) }}
少年 6<age<=18
{{ else }}
成人 > 18
{{ end -}}
`)
	assert.Nil(t, err)
	bs := &bytes.Buffer{}
	err = tpl.Execute(bs, map[string]any{"Age": 5})
	assert.Nil(t, err)
	assert.Equal(t, `
儿童 0<age<6
`, bs.String())
}

func TestForILoop(t *testing.T) {
	// 用一点小技巧来实现 for i 循环
	tpl := template.New("hello-world")
	tpl, err := tpl.Parse(`
{{ range $idx,$elem := . -}}
下标：{{$idx -}},
{{- end}}
`)
	assert.Nil(t, err)
	bs := &bytes.Buffer{}
	// 假设我们要从 0 迭代到 100，即 [0, 100)
	// 这里的切片可以是任意类型，[]bool, []byte 都可以
	// 因为我们本身并不关心里面元素，只是借用一下下标而已
	data := make([]bool, 100)
	err = tpl.Execute(bs, data)
	assert.Nil(t, err)
	assert.Equal(t, `
下标：0,下标：1,下标：2,下标：3,下标：4,下标：5,下标：6,下标：7,下标：8,下标：9,下标：10,下标：11,下标：12,下标：13,下标：14,下标：15,下标：16,下标：17,下标：18,下标：19,下标：20,下标：21,下标：22,下标：23,下标：24,下标：25,下标：26,下标：27,下标：28,下标：29,下标：30,下标：31,下标：32,下标：33,下标：34,下标：35,下标：36,下标：37,下标：38,下标：39,下标：40,下标：41,下标：42,下标：43,下标：44,下标：45,下标：46,下标：47,下标：48,下标：49,下标：50,下标：51,下标：52,下标：53,下标：54,下标：55,下标：56,下标：57,下标：58,下标：59,下标：60,下标：61,下标：62,下标：63,下标：64,下标：65,下标：66,下标：67,下标：68,下标：69,下标：70,下标：71,下标：72,下标：73,下标：74,下标：75,下标：76,下标：77,下标：78,下标：79,下标：80,下标：81,下标：82,下标：83,下标：84,下标：85,下标：86,下标：87,下标：88,下标：89,下标：90,下标：91,下标：92,下标：93,下标：94,下标：95,下标：96,下标：97,下标：98,下标：99,
`, bs.String())
}

type FuncCall struct {
	Slice []string
}

func (f FuncCall) Hello(firstName string, lastName string) string {
	return fmt.Sprintf("Hello, %s·%s", firstName, lastName)
}

func TestFuncCall(t *testing.T) {
	tpl := template.New("hello-world")
	tpl, err := tpl.Parse(`
切片长度: {{len .Slice}}
say hello: {{.Hello "Tom" "Jerry"}}
打印数字: {{printf "%.2f" 1.234}}
`)
	assert.Nil(t, err)
	bs := &bytes.Buffer{}
	err = tpl.Execute(bs,
		&FuncCall{Slice: []string{"Tom", "Jerry"}})
	assert.Nil(t, err)
	assert.Equal(t, `
切片长度: 2
say hello: Hello, Tom·Jerry
打印数字: 1.23
`, bs.String())
}

func TestHelloWorld(t *testing.T) {
	tpl := template.New("hello-world")
	tpl, err := tpl.Parse(`Hello, {{.Name}}`)
	if err != nil {
		t.Fatal(err)
	}
	bs := &bytes.Buffer{}
	err = tpl.Execute(bs, &User{Name: "Tom"})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, `Hello, Tom`, bs.String())
}

func TestMapData(t *testing.T) {
	tpl := template.New("map-data")
	tpl, err := tpl.Parse(`Hello, {{.Name}}`)
	if err != nil {
		t.Fatal(err)
	}
	bs := &bytes.Buffer{}
	err = tpl.Execute(bs, map[string]string{"Name": "Jerry"})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, `Hello, Jerry`, bs.String())
}

func TestSliceData(t *testing.T) {
	tpl := template.New("map-data")
	tpl, err := tpl.Parse(`Hello, {{index . 0}}`)
	if err != nil {
		t.Fatal(err)
	}
	bs := &bytes.Buffer{}
	err = tpl.Execute(bs, []string{"John"})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, `Hello, John`, bs.String())
}

func TestPipeline(t *testing.T) {
	testCases := []struct{
		name string

		tpl  string
		data any

		want string
	} {
		// 这些例子来自官方文档
		// https://pkg.go.dev/text/template#hdr-Pipelines
		{
			name: "string constant",
			tpl:`{{"\"output\""}}`,
			want: `"output"`,
		},
		{
			name: "raw string constant",
			tpl: "{{`\"output\"`}}",
			want: `"output"`,
		},
		{
			name: "function call",
			tpl: `{{printf "%q" "output"}}`,
			want: `"output"`,
		},
		{
			name: "take argument from pipeline",
			tpl: `{{"output" | printf "%q"}}`,
			want: `"output"`,
		},
		{
			name: "parenthesized argument",
			tpl: `{{printf "%q" (print "out" "put")}}`,
			want: `"output"`,
		},
		{
			name: "elaborate call",
			tpl: `{{"put" | printf "%s%s" "out" | printf "%q"}}`,
			want: `"output"`,
		},
		{
			name: "longer chain",
			tpl: `{{"output" | printf "%s" | printf "%q"}}`,
			want: `"output"`,
		},
		{
			name: "with action using dot",
			tpl: `{{with "output"}}{{printf "%q" .}}{{end}}`,
			want: `"output"`,
		},
		{
			name: "with action that creates and uses a variable",
			tpl: `{{with $x := "output" | printf "%q"}}{{$x}}{{end}}`,
			want: `"output"`,
		},
		{
			name: "with action that uses the variable in another action",
			tpl: `{{with $x := "output"}}{{printf "%q" $x}}{{end}}`,
			want: `"output"`,
		},
		{
			name: "pipeline with action that uses the variable in another action",
			tpl: `{{with $x := "output"}}{{$x | printf "%q"}}{{end}}`,
			want: `"output"`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tpl := template.New(tc.name)
			tpl, err := tpl.Parse(tc.tpl)
			if err != nil {
				t.Fatal(err)
			}
			bs := &bytes.Buffer{}
			err = tpl.Execute(bs, tc.data)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tc.want, bs.String())
		})
	}
}
