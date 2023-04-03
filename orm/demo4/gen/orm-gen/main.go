package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"text/template"
)

// Go 会读取 tpl.gohtml 里面的内容填充到变量 tpl 里面
//go:embed template.gohtml
var genOrm string

func main() {
	// 用户必须输入一个 src，限制为文件
	// 然后我们会在同目录下生成代码
	// src := os.Args[1]
	// 源代码目录，也是目标文件目录
	// dir := filepath.Dir(src)
	// 输入的文件名
	// srcFileName := filepath.Base(src)
	src := "testdata/user.go"
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, src, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	tv := &SingleFileEntryVisitor{}
	ast.Walk(tv, f)
	file := tv.Get()

	tpl := template.New("orm-gen")
	fmt.Println(genOrm)
	tpl, err = tpl.Parse(genOrm)
	if err != nil {
		panic(err)
	}
	bs := &bytes.Buffer{}
	err = tpl.Execute(bs, file)
	if err != nil {
		panic(err)
	}
}
