package valuer

import (
	"database/sql"
	"gitee.com/geektime-geekbang/geektime-go/orm/demo3/model"
)

// 先来一个反射和 unsafe 的抽象

// Value 是对结构体实例的内部抽象
type Value interface {
	// SetColumns 设置新值
	SetColumns(rows *sql.Rows) error

	Field(name string) (any, error)
}

// 本质上也可以看所是 factory 模式，极其简单的 factory 模式
type Creator func(t any, model *model.Model) Value