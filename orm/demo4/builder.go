package orm

import (
	"gitee.com/geektime-geekbang/geektime-go/orm/demo4/model"
	"strings"
)

type builder struct {
	model *model.Model
	sb strings.Builder
	args []any
	dialect Dialect
}

func (b *builder) quote(name string) {
	b.sb.WriteByte(b.dialect.quoter())
	b.sb.WriteString(name)
	b.sb.WriteByte(b.dialect.quoter())
}
