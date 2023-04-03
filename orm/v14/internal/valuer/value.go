//go:build v14

package valuer

import (
	"database/sql"
	"gitee.com/geektime-geekbang/geektime-go/orm/v14/model"
)

// Value 是对结构体实例的内部抽象
type Value interface {
	// Field 返回字段对应的值
	Field(name string) (any, error)
	// SetColumns 设置新值
	SetColumns(rows *sql.Rows) error
}

type Creator func(val interface{}, meta *model.Model) Value

// ResultSetHandler 这是另外一种可行的设计方案
// type ResultSetHandler interface {
// 	// SetColumns 设置新值，column 是列名
// 	SetColumns(val any, rows *sql.Rows) error
// }
