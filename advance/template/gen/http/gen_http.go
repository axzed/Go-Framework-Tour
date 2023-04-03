package http

import (
	"io"
	"text/template"
)

// 这部分和课堂的很像，但是有一些地方被我改掉了
const serviceTpl = `
`

func Gen(writer io.Writer, def ServiceDefinition) error {
	tpl := template.New("service")
	tpl, err := tpl.Parse(serviceTpl)
	if err != nil {
		return err
	}
	// 还可以进一步调用 format.Source 来格式化生成代码
	return tpl.Execute(writer, def)
}
