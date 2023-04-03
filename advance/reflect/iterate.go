package reflect

import (
	"errors"
	"reflect"
)

// Iterate 迭代数组，切片，或者字符串
func Iterate(input any) ([]any, error) {
	val := reflect.ValueOf(input)
	typ := val.Type()
	kind := typ.Kind()
	if kind != reflect.Array && kind != reflect.Slice && kind != reflect.String {
		return nil, errors.New("非法类型")
	}
	res := make([]any, 0, val.Len())
	for i := 0; i < val.Len(); i++ {
		ele := val.Index(i)
		res = append(res, ele.Interface())
	}
	return res, nil
}

// IterateMapV1 返回键，值
func IterateMapV1(input any) ([]any, []any, error) {
	val := reflect.ValueOf(input)
	if val.Kind() != reflect.Map {
		return nil, nil, errors.New("非法类型")
	}
	l := val.Len()
	keys := make([]any, 0, l)
	values := make([]any, 0, l)
	for _, k := range val.MapKeys() {
		keys = append(keys, k.Interface())
		v := val.MapIndex(k)
		values = append(values, v.Interface())
	}
	return keys, values, nil
}

// IterateMapV2 返回键，值
func IterateMapV2(input any) ([]any, []any, error) {
	val := reflect.ValueOf(input)
	if val.Kind() != reflect.Map {
		return nil, nil, errors.New("非法类型")
	}
	l := val.Len()
	keys := make([]any, 0, l)
	values := make([]any, 0, l)
	itr := val.MapRange()
	for itr.Next() {
		keys = append(keys, itr.Key().Interface())
		values = append(values, itr.Value().Interface())
	}
	return keys, values, nil
}
