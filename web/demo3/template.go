package web

import "context"

type TemplateEngine interface {
	// 第一个返回值：渲染好的数据
	Render(ctx context.Context, tplName string, data any) ([]byte, error)
}
