package unsafe

import (
	"errors"
	"fmt"
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
	val := reflect.ValueOf(entity)
	typ := reflect.TypeOf(entity)
	val.UnsafeAddr()
	if typ.Kind() != reflect.Pointer || typ.Elem().Kind() != reflect.Struct {
		return nil, errors.New("invalid entity")
	}
	fields := make(map[string]FieldMeta, typ.Elem().NumField())
	elemType := typ.Elem()
	for i := 0; i < elemType.NumField(); i++ {
		fd := elemType.Field(i)
		fields[fd.Name] = FieldMeta{offset: fd.Offset}
	}
	return &UnsafeAccessor{entityAddr: val.UnsafePointer(), fields: fields}, nil
}

func (u *UnsafeAccessor) Field(field string) (int, error) {
	fdMeta, ok := u.fields[field]
	if !ok {
		return 0, fmt.Errorf("invalid field %s", field)
	}
	ptr := unsafe.Pointer(uintptr(u.entityAddr) + fdMeta.offset)
	if ptr == nil {
		return 0, fmt.Errorf("invalid address of the field: %s", field)
	}
	res := *(*int)(ptr)
	return res, nil
}

func (u *UnsafeAccessor) SetField(field string, val int) error {
	fdMeta, ok := u.fields[field]
	if !ok {
		return fmt.Errorf("invalid field %s", field)
	}
	ptr := unsafe.Pointer(uintptr(u.entityAddr) + fdMeta.offset)
	if ptr == nil {
		return fmt.Errorf("invalid address of the field: %s", field)
	}
	*(*int)(ptr) = val
	return nil
}

type FieldMeta struct {
	// offset 后期在我们考虑组合，或者复杂类型字段的时候，它的含义衍生为表达相当于最外层的结构体的偏移量
	offset uintptr
}
