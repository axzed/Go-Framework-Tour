//go:build v16
package orm

import (
	"gitee.com/geektime-geekbang/geektime-go/orm/internal/errs"
)

var (
	MySQL Dialect = &mysqlDialect{}
	SQLite3 Dialect = &sqlite3Dialect{}
)

type Dialect interface {
	// quoter 返回一个引号，引用列名，表名的引号
	quoter() byte
	// buildUpsert 构造插入冲突部分
	buildUpsert(b *builder, odk *Upsert) error
}

type standardSQL struct {

}

func (s *standardSQL) quoter() byte {
	// TODO implement me
	panic("implement me")
}

func (s *standardSQL) buildUpsert(b *builder,
	odk *Upsert) error {
	panic("implement me")
}

type mysqlDialect struct {
	standardSQL
}

func (m *mysqlDialect)  quoter() byte {
	return '`'
}

func (m *mysqlDialect) buildUpsert(b *builder,
	odk *Upsert) error {
	b.sb.WriteString(" ON DUPLICATE KEY UPDATE ")
	for idx, a := range odk.assigns {
		if idx > 0 {
			b.sb.WriteByte(',')
		}
		switch assign := a.(type) {
		case Column:
			colName, err := b.colName(assign.table, assign.name)
			if err != nil {
				return err
			}
			b.quote(colName)
			b.sb.WriteString("=VALUES(")
			b.quote(colName)
			b.sb.WriteByte(')')
		case Assignment:
			err := b.buildColumn(nil, assign.column)
			if err != nil {
				return err
			}
			b.sb.WriteString("=")
			return b.buildExpression(assign.val)
		default:
			return errs.NewErrUnsupportedAssignableType(a)
		}
	}
	return nil
}

type sqlite3Dialect struct {
	standardSQL
}


func (s *sqlite3Dialect)  quoter() byte {
	return '`'
}

func (s *sqlite3Dialect) buildUpsert(b *builder,
	odk *Upsert) error {
	b.sb.WriteString(" ON CONFLICT")
	if len(odk.conflictColumns) > 0 {
		b.sb.WriteByte('(')
		for i, col := range odk.conflictColumns {
			if i > 0 {
				b.sb.WriteByte(',')
			}
			err := b.buildColumn(nil, col)
			if err != nil {
				return err
			}
		}
		b.sb.WriteByte(')')
	}
	b.sb.WriteString(" DO UPDATE SET ")

	for idx, a := range odk.assigns {
		if idx > 0 {
			b.sb.WriteByte(',')
		}
		switch assign := a.(type) {
		case Column:
			colName, err := b.colName(assign.table, assign.name)
			if err != nil {
				return err
			}
			b.quote(colName)
			b.sb.WriteString("=excluded.")
			b.quote(colName)
		case Assignment:
			err := b.buildColumn(nil, assign.column)
			if err != nil {
				return err
			}
			b.sb.WriteString("=")
			return b.buildExpression(assign.val)
		default:
			return errs.NewErrUnsupportedAssignableType(a)
		}
	}
	return nil
}