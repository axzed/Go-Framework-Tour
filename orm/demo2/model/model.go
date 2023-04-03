package model

import (
	"gitee.com/geektime-geekbang/geektime-go/orm/demo2/internal/errs"
	"reflect"
)

// 带校验的 option 模式
type Option func(m *Model) error

type Model struct {
	// 结果体对应的表名
	TableName string
	Fields    []*Field
	// 字段名到字段的元数据
	FieldMap map[string]*Field

	// 列名到字段的映射
	ColumnMap map[string]*Field
}

func WithTableName(name string) Option {
	return func(m *Model) error {
		// if name == "" {
		// 	return errs.ErrEmptyTableName
		// }
		m.TableName = name
		return nil
	}
}

func ModelWithColumnName(field string, colName string) Option {
	return func(m *Model) error {
		fd, ok := m.FieldMap[field]
		if !ok {
			return errs.NewErrUnknownField(field)
		}
		fd.ColName = colName
		return nil
	}
}

func ModelWithColumn(field string, col *Field) Option {
	return func(m *Model) error {
		m.FieldMap[field] = col
		return nil
	}
}

// func ModelWithColumnAutoIncrement(field string) Option {
// 	return func(m *Model) error {
// 		m.FieldMap[field].autoIncrement = true
// 		return nil
// 	}
// }

type Field struct {
	// 字段名
	GoName string
	// 字段对应的列名
	ColName string

	Type reflect.Type

	// 表达相对量的概念是可以
	// 字段偏移量
	Offset uintptr

	Index []int
}

type TableName interface {
	TableName() string
}


