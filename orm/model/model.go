package model

import (
	"reflect"
)

type Model struct {
	// TableName 结构体对应的表名
	TableName string
	Fields []*Field
	FieldMap  map[string]*Field
	ColumnMap map[string]*Field
}

// Field 字段
type Field struct {
	NotNull bool
	ColName string
	GoName string
	Type   reflect.Type
	Index int
	// Offset 相对于对象起始地址的字段偏移量
	Offset uintptr
}

// 我们支持的全部标签上的 key 都放在这里
// 方便用户查找，和我们后期维护
const (
	tagKeyColumn = "column"
)

// 用户自定义一些模型信息的接口，集中放在这里
// 方便用户查找和我们后期维护

// TableName 用户实现这个接口来返回自定义的表名
type TableName interface {
	TableName() string
}

// type Order struct {
// 	BuyerId int64
// 	Region string
// }
//
// func (o Order) DBName() string {
// 	return fmt.Sprintf("%s_order_db_%04d", o.Region, o.BuyerId % 1000)
// }
//
// func (o Order) TableName() string {
//
// 	// 最简单的分表
// 	return fmt.Sprintf("order_tab_%04d", o.BuyerId % 1000)
// }
