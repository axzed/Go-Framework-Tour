package main

import (
	"fmt"
	"gitee.com/geektime-geekbang/geektime-go/advance/template/gen/annotation"
	"gitee.com/geektime-geekbang/geektime-go/advance/template/gen/http"
	"os"
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
	// 第一步找出符合条件的文件
	srcFiles, err := scanFiles(src)
	if err != nil {
		return err
	}
	// 第二步，AST 解析源代码文件，拿到 service definition 定义
	defs, err := parseFiles(srcFiles)
	if err != nil {
		return err
	}
	// 生成代码
	return genFiles(src, defs)
}

// 根据 defs 来生成代码
// src 是源代码所在目录，在测试里面它是 ./testdata
func genFiles(src string, defs []http.ServiceDefinition) error {
	panic("implement me")
}

func parseFiles(srcFiles []string) ([]http.ServiceDefinition, error) {
	defs := make([]http.ServiceDefinition, 0, 20)
	for _, src := range srcFiles {
		fmt.Println(src)
		// 你需要利用 annotation 里面的东西来扫描 src，然后生成 file
		panic("implement me")
		var file annotation.File

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

// 你需要利用 typ 来构造一个 http.ServiceDefinition
// 注意你可能需要检测用户的定义是否符合你的预期
func parseServiceDefinition(pkg string, typ annotation.Type) (http.ServiceDefinition, error) {
	panic("implement me")
}

// 返回符合条件的 Go 源代码文件，也就是你要用 AST 来分析这些文件的代码
func scanFiles(src string) ([]string, error) {
	panic("implement me")
}

// underscoreName 驼峰转字符串命名，在决定生成的文件名的时候需要这个方法
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
