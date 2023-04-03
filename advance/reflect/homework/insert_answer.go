package homework

import (
	"errors"
	"reflect"
	"strings"
)

var errInvalidEntity = errors.New("invalid entity")

// InsertStmt 会为实例生成一个 INSERT 语句。
// INSERT 语句只需要考虑 MySQL 的语法
// 只接收非 nil，一级指针，结构体实例
// 结构体字段只能是基本类型，或者实现了 driver.Valuer 接口
// 但是我们只做最简单的校验，不会全部情况都校验
func InsertStmt(entity interface{}) (string, []interface{}, error) {
	if entity == nil {
		return "", nil, errInvalidEntity
	}
	val := reflect.ValueOf(entity)
	typ := val.Type()
	if typ.Kind() == reflect.Pointer {
		val = val.Elem()
		typ = val.Type()
	}

	if typ.Kind() != reflect.Struct {
		return "", nil, errInvalidEntity
	}

	bd := strings.Builder{}
	_, _ = bd.WriteString("INSERT INTO `")
	bd.WriteString(typ.Name())
	bd.WriteString("`(")
	fields, values := fieldNameAndValues(val)
	for i, name := range fields {
		if i > 0 {
			bd.WriteRune(',')
		}
		bd.WriteRune('`')
		bd.WriteString(name)
		bd.WriteRune('`')
	}
	bd.WriteString(") VALUES(")
	args := make([]interface{}, 0, len(values))
	for i, fd := range fields {
		if i > 0 {
			bd.WriteRune(',')
		}
		bd.WriteRune('?')
		args = append(args, values[fd])
	}
	if len(args) == 0 {
		return "", nil, errInvalidEntity
	}
	bd.WriteString(");")
	return bd.String(), args, nil
}

// 我们这种写法会导致在出现组合的时候会有额外的内存分配
// 第一个数组来保证顺序，第二个map保存结果，并且用于去重
// 重复的时候，第一个胜出
func fieldNameAndValues(val reflect.Value) ([]string, map[string]interface{}) {
	typ := val.Type()
	fieldNum := typ.NumField()
	fields := make([]string, 0, fieldNum)
	values := make(map[string]interface{}, fieldNum)
	for i := 0; i < fieldNum; i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i)

		// Anonymous 只处理真正的组合，这是区别我们测试用例里面 Buyer 和 Seller 不同声明方式的差异点
		if field.Type.Kind() == reflect.Struct && field.Anonymous {
			subFields, subValues := fieldNameAndValues(fieldVal)
			for _, k := range subFields {
				if _, ok := values[k]; ok {
					// 重复字段，只会出现在组合情况下。我们直接忽略重复字段。
					continue
				}
				fields = append(fields, k)
				values[k] = subValues[k]
			}
			continue
		}
		fields = append(fields, field.Name)
		values[field.Name] = fieldVal.Interface()
	}
	return fields, values
}
