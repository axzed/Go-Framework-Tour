package reflect

import (
	"errors"
	"reflect"
)

// IterateFuncs 输出方法信息，并执行调用
func IterateFuncs(val any) (map[string]*FuncInfo, error) {
	typ := reflect.TypeOf(val)
	if typ.Kind() != reflect.Struct && typ.Kind() != reflect.Ptr {
		return nil, errors.New("非法类型")
	}
	num := typ.NumMethod()
	result := make(map[string]*FuncInfo, num)
	for i := 0; i < num; i++ {
		f := typ.Method(i)
		numIn := f.Type.NumIn()
		ps := make([]reflect.Value, 0, f.Type.NumIn())
		// 第一个参数永远都是接收器，类似于 java 的 this 概念
		ps = append(ps, reflect.ValueOf(val))
		in := make([]reflect.Type, 0, f.Type.NumIn())
		for j := 0; j < numIn; j++ {
			p := f.Type.In(j)
			in = append(in, p)
			if j > 0 {
				ps = append(ps, reflect.Zero(p))
			}
		}
		// 调用结果
		ret := f.Func.Call(ps)
		outNum := f.Type.NumOut()
		out := make([]reflect.Type, 0, outNum)
		res := make([]any, 0, outNum)
		for k := 0; k < outNum; k++ {
			out = append(out, f.Type.Out(k))
			res = append(res, ret[k].Interface())
		}
		result[f.Name] = &FuncInfo{
			Name:   f.Name,
			In:     in,
			Out:    out,
			Result: res,
		}
	}
	return result, nil
}

type FuncInfo struct {
	Name string
	In   []reflect.Type
	Out  []reflect.Type

	// 反射调用得到的结果
	Result []any
}
