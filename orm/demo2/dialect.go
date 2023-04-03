package orm

import (
	"gitee.com/geektime-geekbang/geektime-go/orm/demo2/internal/errs"
)

// Dialect 方言，构造个性部分
type Dialect interface {
	// 引号
	quoter() byte
	buildDuplicateKey(b *builder, odk *OnDuplicateKey) error
}

// SQL 标准的方言实现
type standardSQL struct {

}

type mysqlDialect struct {
	standardSQL
}

func (dialect *mysqlDialect) quoter() byte {
	return '`'
}

func (dialect *mysqlDialect) buildDuplicateKey(b *builder, odk *OnDuplicateKey) error {
	// 你在这里可以进一步检测 assigns 是不是为空
	// 构造 ON DUPLICATE KEY 部分
	b.sb.WriteString(" ON DUPLICATE KEY UPDATE ")
	for idx, assign := range odk.assigns {
		if idx > 0 {
			b.sb.WriteByte(',')
		}
		switch expr := assign.(type) {
		case Assignment:
			fd, ok := b.model.FieldMap[expr.column]
			if !ok {
				return errs.NewErrUnknownField(expr.column)
			}
			b.quote(fd.ColName)
			b.sb.WriteString("=?")
			b.args = append(b.args, expr.val)
		case Column:
			fd, ok := b.model.FieldMap[expr.name]
			if !ok {
				return errs.NewErrUnknownField(expr.name)
			}
			b.quote(fd.ColName)
			b.sb.WriteString("=VALUES(")
			b.quote(fd.ColName)
			b.sb.WriteByte(')')
		}
	}
	return nil
}

type sqliteDialect struct {
	standardSQL
}

// func (dialect *sqliteDialect) buildDuplicateKey(b *builder, odk *Update) error {
// 	// 你在这里可以进一步检测 assigns 是不是为空
// 	// 构造 ON DUPLICATE KEY 部分
// 	b.sb.WriteString(" ON CONFLICT xxxx DO UPDATE SET ")
// 	for idx, assign := range odk.assigns {
// 		if idx > 0 {
// 			b.sb.WriteByte(',')
// 		}
// 		switch expr := assign.(type) {
// 		case Assignment:
// 			fd, ok := b.model.FieldMap[expr.column]
// 			if !ok {
// 				return errs.NewErrUnknownField(expr.column)
// 			}
// 			b.quote(fd.ColName)
// 			b.sb.WriteString("=?")
// 			b.args = append(b.args, expr.val)
// 		case Column:
// 			fd, ok := b.model.FieldMap[expr.name]
// 			if !ok {
// 				return errs.NewErrUnknownField(expr.name)
// 			}
// 			b.quote(fd.ColName)
// 			b.sb.WriteString("=VALUES(")
// 			b.quote(fd.ColName)
// 			b.sb.WriteByte(')')
// 		}
// 	}
// 	return nil
// }
