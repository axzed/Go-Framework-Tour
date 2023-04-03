//go:build v6

package orm

import "gitee.com/geektime-geekbang/geektime-go/orm/internal/errs"

type ModelOpt func(model *Model) error

type Model struct {
	// tableName 结构体对应的表名
	tableName string
	fieldMap  map[string]*field
}

func ModelWithTableName(tableName string) ModelOpt {
	return func(model *Model) error {
		model.tableName = tableName
		return nil
	}
}

func ModelWithColumnName(field string, columnName string) ModelOpt {
	return func(model *Model) error {
		fd, ok := model.fieldMap[field]
		if !ok {
			return errs.NewErrUnknownField(field)
		}
		// 注意，这里我们根本没有检测 colName 会不会是空字符串
		// 因为正常情况下，用户都不会写错
		// 即便写错了，也很容易在测试中发现
		fd.colName = columnName
		return nil
	}
}

// field 字段
type field struct {
	colName string
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
