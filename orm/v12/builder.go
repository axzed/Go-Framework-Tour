//go:build v12

package orm

import (
	"gitee.com/geektime-geekbang/geektime-go/orm/v12/internal/errs"
	"gitee.com/geektime-geekbang/geektime-go/orm/v12/model"
	"strings"
)

type builder struct {
	sb      strings.Builder
	args    []any
	model   *model.Model
	dialect Dialect
	quoter  byte
}

// buildColumn 构造列
func (b *builder) buildColumn(fd string) error {
	meta, ok := b.model.FieldMap[fd]
	if !ok {
		return errs.NewErrUnknownField(fd)
	}
	b.quote(meta.ColName)
	return nil
}

func (b *builder) quote(name string) {
	b.sb.WriteByte(b.quoter)
	b.sb.WriteString(name)
	b.sb.WriteByte(b.quoter)
}

func (b *builder) addArgs(args ...any) {
	if b.args == nil {
		// 很少有查询能够超过八个参数
		// INSERT 除外
		b.args = make([]any, 0, 8)
	}
	b.args = append(b.args, args...)
}
