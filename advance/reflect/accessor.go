package reflect

import (
	"errors"
	"reflect"
)

type ReflectAccessor struct {
	val reflect.Value
	typ reflect.Type
}

func NewReflectAccessor(val any) (*ReflectAccessor, error) {
	typ := reflect.TypeOf(val)
	if typ.Kind() != reflect.Pointer || typ.Elem().Kind() != reflect.Struct {
		return nil, errors.New("invalid entity")
	}
	return &ReflectAccessor{
		val: reflect.ValueOf(val).Elem(),
		typ: typ.Elem(),
	}, nil
}

func (r *ReflectAccessor) Field(field string) (int, error) {
	if _, ok := r.typ.FieldByName(field); !ok {
		return 0, errors.New("非法字段")
	}
	return r.val.FieldByName(field).Interface().(int), nil
}

func (r *ReflectAccessor) SetField(field string, val int) error {
	if _, ok := r.typ.FieldByName(field); !ok {
		return errors.New("非法字段")
	}
	fdVal := r.val.FieldByName(field)
	if !fdVal.CanSet() {
		return errors.New("无法设置新值的字段")
	}
	fdVal.Set(reflect.ValueOf(val))
	return nil
}
