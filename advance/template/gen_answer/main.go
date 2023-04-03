package main

import (
	"errors"
	"fmt"
	"gitee.com/geektime-geekbang/geektime-go/advance/template/gen_answer/annotation"
	"gitee.com/geektime-geekbang/geektime-go/advance/template/gen_answer/http"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

// 实际上 main 函数这里要考虑接收参数
// src 源目标
// dst 目标目录
// type src 里面可能有很多类型，那么用户可能需要指定具体的类型
// 这里我们简化操作，只读取当前目录下的数据，并且扫描下面的所有源文件，然后生成代码
// 在当前目录下运行 go install 就将 main 安装成功了，
// 可以在命令行中运行 gen
// 在 testdata 里面运行 gen，则会生成能够通过所有测试的代码
func main() {
	err := gen(".")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("success")
}

func gen(src string) error {
	srcFiles, err := scanFiles(src)
	if err != nil {
		return err
	}
	defs, err := parseFiles(srcFiles)
	if err != nil {
		return err
	}
	return genFiles(src, defs)
}

func genFiles(src string, defs []http.ServiceDefinition) error {
	for _, def := range defs {
		name := def.GenName()
		fileName := underscoreName(name) + ".go"
		f, err := os.Create(filepath.Join(src, fileName))
		if err != nil {
			return err
		}
		err = http.Gen(f, def)
		if err != nil {
			return err
		}
		f.Close()
	}
	return nil
}

func parseFiles(srcFiles []string) ([]http.ServiceDefinition, error) {
	defs := make([]http.ServiceDefinition, 0, 20)
	for _, src := range srcFiles {
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, src, nil, parser.ParseComments)
		if err != nil {
			return nil, err
		}
		tv := &annotation.SingleFileEntryVisitor{}
		ast.Walk(tv, f)
		file := tv.Get()

		for _, typ := range file.Types {
			_, ok := typ.Annotations.Get("HttpClient")
			if !ok {
				continue
			}
			def, err := parseServiceDefinition(file.Node.Name.Name, typ)
			if err != nil {
				return nil, err
			}
			defs = append(defs, def)
		}
	}
	return defs, nil
}

func parseServiceDefinition(pkg string, typ annotation.Type) (http.ServiceDefinition, error) {
	methods := make([]http.ServiceMethod, 0, len(typ.Fields))
	for _, fd := range typ.Fields {
		ft, ok := fd.Node.Type.(*ast.FuncType)
		if !ok {
			continue
		}
		// 不符合要求的可以报错，也可以调过，看设计者
		if ft.Params == nil || len(ft.Params.List) != 2 {
			return http.ServiceDefinition{},
				errors.New("gen: 方法必须接收两个参数，其中第一个参数是 context.Context，第二个参数请求")
		}
		if ft.Results == nil || len(ft.Results.List) != 2 {
			return http.ServiceDefinition{},
				errors.New("gen: 方法必须返回两个参数，其中第一个返回值是响应，第二个返回值是error")
		}
		pType, ok := ft.Params.List[1].Type.(*ast.StarExpr)
		// 必须使用指针
		if !ok {
			return http.ServiceDefinition{}, errors.New("gen：第二个参数必须是指针")
		}
		rType, ok := ft.Results.List[0].Type.(*ast.StarExpr)
		if !ok {
			return http.ServiceDefinition{}, errors.New("gen：第一个返回值必须是指针")
		}
		name := toName(fd.Node.Names)
		path := string(append([]byte{filepath.Separator}, name...))
		pathAns, ok := fd.Annotations.Get("Path")
		if ok {
			// 没有必要检测这个 value，我们防呆不防傻
			path = pathAns.Value
		}
		methods = append(methods, http.ServiceMethod{
			Name:         name,
			Path:         path,
			ReqTypeName:  pType.X.(*ast.Ident).Name,
			RespTypeName: rType.X.(*ast.Ident).Name,
		})
	}

	if len(methods) == 0 {
		return http.ServiceDefinition{}, errors.New("gen：没有符合条件的方法")
	}

	name := typ.Node.Name.String()
	snAn, ok := typ.Annotations.Get("ServiceName")
	if ok {
		// 防呆不防傻
		name = snAn.Value
	}

	return http.ServiceDefinition{
		Package: pkg,
		Name:    name,
		Methods: methods,
	}, nil
}

func toName(names []*ast.Ident) string {
	var sb strings.Builder
	for _, n := range names {
		sb.WriteString(n.String())
	}
	return sb.String()
}

// 返回符合条件的 Go 源代码文件
func scanFiles(src string) ([]string, error) {
	srcFiles := make([]string, 0, 10)
	files, err := os.ReadDir(src)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".go") &&
			!strings.HasSuffix(file.Name(), "_test.go") &&
			!strings.HasSuffix(file.Name(), "gen.go") {
			src, err = filepath.Abs(src)
			if err != nil {
				return nil, err
			}
			srcFiles = append(srcFiles, filepath.Join(src, file.Name()))
		}
	}
	return srcFiles, nil
}

// underscoreName 驼峰转字符串命名
// 可以用正则表达式，然而我写不出来，我是正则渣
func underscoreName(name string) string {
	var buf []byte
	for i, v := range name {
		if unicode.IsUpper(v) {
			if i != 0 {
				buf = append(buf, '_')
			}
			buf = append(buf, byte(unicode.ToLower(v)))
		} else {
			buf = append(buf, byte(v))
		}

	}
	return string(buf)
}
