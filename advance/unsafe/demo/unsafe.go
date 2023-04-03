package demo

import (
	"errors"
	"reflect"
	"unsafe"
)

type FieldAccessor interface {
	Field(field string) (int, error)
	SetField(field string, val int) error
}

type UnsafeAccessor struct {
	fields     map[string]FieldMeta
	entityAddr unsafe.Pointer
}

func NewUnsafeAccessor(entity interface{}) (*UnsafeAccessor, error) {
	if entity == nil {
		return nil, errors.New("invalid entity")
	}

	typ := reflect.TypeOf(entity)
	if typ.Kind() != reflect.Pointer || typ.Elem().Kind() != reflect.Struct {
		return nil, errors.New("invalid entity")
	}
	fields := make(map[string]FieldMeta, typ.Elem().NumField())
	elemType := typ.Elem()
	for i := 0; i < elemType.NumField(); i++ {
		fd := elemType.Field(i)
		fields[fd.Name] = FieldMeta{offset: fd.Offset, typ: fd.Type}
	}

	val := reflect.ValueOf(entity)
	return &UnsafeAccessor{entityAddr: val.UnsafePointer(), fields: fields}, nil
}

func (u *UnsafeAccessor) Field(field string) (int, error) {
	meta, ok := u.fields[field]
	if !ok {
		return 0, errors.New("不存在字段")
	}

	// 计算地址了
	res := *(*int)(unsafe.Pointer(uintptr(u.entityAddr) + meta.offset))
	return res, nil
}

func (u *UnsafeAccessor) FieldAny(field string) (any, error) {
	meta, ok := u.fields[field]
	if !ok {
		return 0, errors.New("不存在字段")
	}

	res := reflect.NewAt(meta.typ, unsafe.Pointer(uintptr(u.entityAddr)+meta.offset)).Elem()
	// 计算地址了
	return res.Interface(), nil
}

func (u *UnsafeAccessor) SetField(field string, val int) error {
	meta, ok := u.fields[field]
	if !ok {
		return errors.New("不存在字段")
	}

	// 计算地址了
	*(*int)(unsafe.Pointer(uintptr(u.entityAddr) + meta.offset)) = val
	return nil
}

func (u *UnsafeAccessor) SetFieldAny(field string, val any) error {
	meta, ok := u.fields[field]
	if !ok {
		return errors.New("不存在字段")
	}

	// 计算地址了
	res := reflect.NewAt(meta.typ, unsafe.Pointer(uintptr(u.entityAddr)+meta.offset))
	if res.CanSet() {
		res.Set(reflect.ValueOf(val))
	}

	return nil
}

type FieldMeta struct {
	typ reflect.Type
	// offset 后期在我们考虑组合，或者复杂类型字段的时候，它的含义衍生为表达相当于最外层的结构体的偏移量
	offset uintptr
}
